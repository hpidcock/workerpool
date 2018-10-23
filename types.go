package workerpool

import "context"

// Task defines a asynchronous operation that can be canceled and return an error.
type Task func(context.Context) error
