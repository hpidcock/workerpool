package workerpool

import (
	"reflect"
	"sync"
)

// Each takes a slice or map, asynchronously iterates it, and returns a collection of return
// values or nil
func (wp *WorkerPool) Each(collection interface{}, iterator interface{}) interface{} {
	c := reflect.ValueOf(collection)
	ct := c.Type()
	i := reflect.ValueOf(iterator)
	it := i.Type()

	if it.Kind() != reflect.Func {
		panic("iterator is not a function")
	}

	if it.NumOut() > 1 {
		panic("iterator can only return one or zero values")
	}

	switch ct.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		r := wp.eachArray(c, i)
		if r.IsValid() {
			return r.Interface()
		}
		return nil
	case reflect.Map:
		r := wp.eachMap(c, i)
		if r.IsValid() {
			return r.Interface()
		}
		return nil
	}

	panic("collection is not a map, array or slice")
}

func (wp *WorkerPool) eachArray(c reflect.Value, i reflect.Value) reflect.Value {
	cl := c.Len()
	it := i.Type()
	ret := reflect.Value{}
	hasReturn := it.NumOut() != 0
	if hasReturn {
		ret = reflect.MakeSlice(reflect.SliceOf(it.Out(0)), cl, cl)
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
				ret.Index(idx).Set(retValue[0])
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

func (wp *WorkerPool) eachMap(c reflect.Value, i reflect.Value) reflect.Value {
	ct := c.Type()
	cl := c.Len()
	it := i.Type()
	ret := reflect.Value{}
	hasReturn := it.NumOut() != 0
	if hasReturn {
		ret = reflect.MakeMap(reflect.MapOf(ct.Key(), it.Out(0)))
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
				ret.SetMapIndex(idx, retValue[0])
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
