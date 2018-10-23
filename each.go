package workerpool

import (
	"reflect"
	"sync"
)

// Each takes a slice or map, asynchronously iterates it, and returns a collection of return
// values or nil
func (wp *WorkerPool) Each(collection interface{}, iterator interface{}) Result {
	c := reflect.ValueOf(collection)
	ct := c.Type()
	i := reflect.ValueOf(iterator)
	it := i.Type()

	if it.Kind() != reflect.Func {
		panic("iterator is not a function")
	}

	if it.NumOut() > 2 {
		panic("iterator can only return two, one or zero values")
	}

	switch ct.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		r, err := wp.eachArray(c, i)
		if r.IsValid() {
			return Result{r.Interface(), nil}
		}
		return Result{}
	case reflect.Map:
		r := wp.eachMap(c, i)
		if r.IsValid() {
			return Result{r.Interface(), nil}
		}
		return Result{}
	}

	panic("collection is not a map, array or slice")
}

func (wp *WorkerPool) eachArray(c reflect.Value, i reflect.Value) []reflect.Value {
	var ret []reflect.Value
	cl := c.Len()
	it := i.Type()
	numReturn := it.NumOut()
	hasReturn := numReturn != 0
	if hasReturn {
		ret = make([]reflect.Value, 0, numReturn)
	}
	for i := 0; i < numReturn; i++ {
		ret = append(ret,
			reflect.MakeSlice(reflect.SliceOf(it.Out(i)), cl, cl))
	}
	wg := sync.WaitGroup{}
	wg.Add(cl)
	m := sync.Mutex{}
	for ci := 0; ci < cl; ci++ {
		idx := ci
		err := wp.Push(func() {
			args := []reflect.Value{reflect.ValueOf(idx), c.Index(idx)}
			retValue := i.Call(args)
			if hasReturn {
				m.Lock()
				for i := 0; i < numReturn; i++ {
					ret[i].Index(idx).Set(retValue[i])
				}
				m.Unlock()
			}
			wg.Done()
		})
		if err != nil {
			panic(err)
		}
	}
	wg.Wait()
	return ret
}

func (wp *WorkerPool) eachMap(c reflect.Value, i reflect.Value) []reflect.Value {
	var ret []reflect.Value
	ct := c.Type()
	cl := c.Len()
	it := i.Type()
	numReturn := it.NumOut()
	hasReturn := numReturn != 0
	for i := 0; i < numReturn; i++ {
		ret = append(ret,
			reflect.MakeMap(reflect.MapOf(ct.Key(), it.Out(i))))
	}
	wg := sync.WaitGroup{}
	wg.Add(cl)
	m := sync.Mutex{}
	for _, key := range c.MapKeys() {
		idx := key
		err := wp.Push(func() {
			args := []reflect.Value{idx, c.MapIndex(idx)}
			retValue := i.Call(args)
			if hasReturn {
				m.Lock()
				for i := 0; i < numReturn; i++ {
					ret[i].SetMapIndex(idx, retValue[i])
				}
				m.Unlock()
			}
			wg.Done()
		})
		if err != nil {
			panic(err)
		}
	}
	wg.Wait()
	return ret
}
