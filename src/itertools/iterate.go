package itertools

import (
	"container/list"
	"container/ring"
	"fmt"
	"reflect"
	"strings"
)

func Iterate(l interface{}) (out chan Pair) {
	out = make(chan Pair, GetIterBuffer())

	go func() {
		defer close(out)

		if isBuiltIn := builtInIterate(l, out); !isBuiltIn {
			if isContainer := containerIterate(l, out); !isContainer {
				panic(fmt.Sprintf("Iterate function does not support the type: %s", l))
			}
		}
	}()
	return
}

func builtInIterate(l interface{}, out chan Pair) (isBuiltInType bool) {
	isBuiltInType = true
	valueOfIter := reflect.ValueOf(l)
	k := valueOfIter.Kind()

	if k == reflect.Ptr {
		valueOfIter = valueOfIter.Elem()
		k = valueOfIter.Kind()
	}

	switch k {
	case reflect.Map:
		iterateMap(valueOfIter, out)

	case reflect.Array, reflect.Slice:
		iterateArraySlice(valueOfIter, out)

	case reflect.Chan:
		iterateChan(valueOfIter, out)

	case reflect.String:
		iterateString(l.(string), out)

	default:
		isBuiltInType = false
	}

	return
}

func containerIterate(l interface{}, out chan Pair) (isContainerType bool) {
	isContainerType = true

	switch l.(type) {
	case *list.List:
		iterateList(l.(*list.List), out)

	case *ring.Ring:
		iterateRing(l.(*ring.Ring), out)

	default:
		isContainerType = false
	}

	return
}

func iterateList(lst *list.List, out chan Pair) {
    fmt.Println("we will have List support shortly", lst, out)
}

func iterateRing(rng *ring.Ring, out chan Pair) {
    fmt.Println("we will have Ring support shortly", rng, out)
}

func iterateString(s string, out chan Pair) {
	for i, v := range strings.Split(s, "") {
		out <- Pair{i, v}
	}
}

func iterateChan(valueOfIter reflect.Value, out chan Pair) {
	i := 0

	for v, ok := valueOfIter.Recv(); ok; {
		out <- Pair{i, v.Interface()}
		i++
	}
}

func iterateArraySlice(valueOfIter reflect.Value, out chan Pair) {
	for i := 0; i < valueOfIter.Len(); i++ {
		out <- Pair{i, valueOfIter.Index(i).Interface()}
	}
}

func iterateMap(valueOfIter reflect.Value, out chan Pair) {
	for _, v := range valueOfIter.MapKeys() {
		out <- Pair{v.Interface(), valueOfIter.MapIndex(v).Interface()}
	}
}
