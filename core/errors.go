package core

import (
	"errors"
)

var (
	ErrNotOpen             = errors.New("mysql not open")
	ErrServerNotFound      = errors.New("mysql server not found")
	ErrCircuitOpen         = errors.New("mysql circuit open")
	ErrShardLookupNotFound = errors.New("shard lookup fails")
	ErrInvalidHintId       = errors.New("hintId=0?")
	ErrEntityLocked        = errors.New("entity being locked")
	ErrLookupTableNotFound = errors.New("mysql lookup table not configured")

	ErrInvalidSql          = errors.New("invalid sql")
	ErrMysqlServerNotFound = errors.New("mysql server not found")

	ErrInvalidShardingPrivacy = errors.New("invalid sharding privacy")
	ErrInvalidSqlType         = errors.New("invalid sql type")
	ErrInvalidPrivacyBehavior = errors.New("invalid privacy behavior")
)

// http://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html
// mysql error code is always 4 digits
var mysqlNonSystemErrors = map[string]bool{
	"1054": true, // Error 1054: Unknown column 'curve_internal_id' in 'field list'
	"1062": true, // Error 1062: Duplicate entry '1' for key 'PRIMARY'
}
