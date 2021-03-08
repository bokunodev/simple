package iterator

type Iterator interface {
	Err() error
	Next() bool
	Value() interface{}
}
