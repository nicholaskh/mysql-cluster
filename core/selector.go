package core

type Selector interface {
	PickServer(pool string, hintId int, sql string) (*mysql, error)
}
