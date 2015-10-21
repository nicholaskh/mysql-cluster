package core

import (
	"database/sql"

	"github.com/nicholaskh/golib/breaker"
	"github.com/nicholaskh/golib/cache"
	log "github.com/nicholaskh/log4go"
	_ "github.com/nicholaskh/mysql"
	"github.com/nicholaskh/mysql-cluster/config"
)

// A mysql conn to a single mysql instance
// Conn pool is natively supported by golang
type mysql struct {
	dsn        string
	db         *sql.DB         // a pool of connections to a single db instance
	stmtsStore *cache.LruCache // {query: stmt}
	breaker    *breaker.Consecutive

	maxIdleConns int
	maxOpenConns int

	config *config.MysqlConfig
}

func newMysql(dsn string, config *config.MysqlConfig) *mysql {
	this := new(mysql)
	this.dsn = dsn
	this.config = config
	this.breaker = &breaker.Consecutive{
		FailureAllowance: config.FailureAllowance,
		RetryTimeout:     config.RetryTimeout,
	}
	this.maxIdleConns = config.MaxIdleConns
	this.maxOpenConns = config.MaxOpenConns
	if config.MaxStmtCache > 0 {
		this.stmtsStore = cache.NewLruCache(config.MaxStmtCache)
		this.stmtsStore.OnEvicted = func(key cache.Key, value interface{}) {
			query := key.(string)
			stmt := value.(*sql.Stmt)
			stmt.Close()

			log.Debug("[%s] stmt[%s] closed", this.dsn, query)
		}
	}

	return this
}

func (this *mysql) Open() (err error) {
	this.db, err = sql.Open("mysql", this.dsn)
	this.db.SetMaxIdleConns(this.maxIdleConns)
	this.db.SetMaxOpenConns(this.maxOpenConns)

	return
}

func (this *mysql) Ping() error {
	if this.db == nil {
		return ErrNotOpen
	}

	return this.db.Ping()
}

func (this *mysql) String() string {
	return this.dsn
}

func (this *mysql) Query(query string, args ...interface{}) (rows *sql.Rows,
	err error) {
	if this.db == nil {
		return nil, ErrNotOpen
	}
	if this.breaker.Open() {
		return nil, ErrCircuitOpen
	}

	var stmt *sql.Stmt = nil
	if this.stmtsStore != nil {
		if stmtc, present := this.stmtsStore.Get(query); present {
			stmt = stmtc.(*sql.Stmt)
		} else {
			// FIXME thundering hurd
			stmt, err = this.db.Prepare(query)
			if err != nil {
				if this.isSystemError(err) {
					log.Warn("mysql prepare breaks: %s", err.Error())
					this.breaker.Fail()
				}

				return nil, err
			}

			this.stmtsStore.Set(query, stmt)
			log.Debug("[%s] stmt[%s] open", this.dsn, query)
		}
	}

	// Under the hood, db.Query() actually prepares, executes, and closes
	// a prepared statement. That's three round-trips to the database.
	if stmt != nil {
		rows, err = stmt.Query(args...)
	} else {
		rows, err = this.db.Query(query, args...)
	}
	if err != nil {
		if this.isSystemError(err) {
			log.Warn("mysql query breaks: %s", err.Error())
			this.breaker.Fail()
		}
	} else {
		this.breaker.Succeed()
	}

	return
}

func (this *mysql) Exec(query string, args ...interface{}) (afftectedRows int64,
	lastInsertId int64, err error) {
	if this.db == nil {
		return 0, 0, ErrNotOpen
	}
	if this.breaker.Open() {
		return 0, 0, ErrCircuitOpen
	}

	var result sql.Result
	result, err = this.db.Exec(query, args...)
	if err != nil {
		if this.isSystemError(err) {
			log.Warn("mysql exec breaks: %s", err.Error())
			this.breaker.Fail()
		}

		return 0, 0, err
	}

	afftectedRows, err = result.RowsAffected()
	if err != nil {
		if this.isSystemError(err) {
			log.Warn("mysql exec2 breaks: %s", err.Error())
			this.breaker.Fail()
		}
	} else {
		this.breaker.Succeed()
	}

	lastInsertId, _ = result.LastInsertId()
	return
}

func (this *mysql) isSystemError(err error) bool {
	// "Error %d:" skip the leading 6 chars: "Error "
	var errcode = err.Error()[6:10] // TODO confirm mysql err always "Error %d: %s"
	if _, present := mysqlNonSystemErrors[errcode]; present {
		return false
	}

	return true
}
