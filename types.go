package workerpool

import "context"

// Task defines a asynchronous operation that can be canceled and return an error.
type Task func(context.Context) error

// Result returned from batch operation.
type Result struct {
	value interface{}
	err   error
}

func (r Result) Value() interface{} {
	return r.value
}

func (r Result) Error() error {
	return r.err
}
