package test

import (
	"fmt"
	"reflect"
)

type Expect struct {
	Value any

	Error string
}

// 数值比较
func (e *Expect) Equal(v1, v2 any) *Expect {
	if reflect.DeepEqual(v1, v2) {
		return e
	}

	e.Error = fmt.Sprintf("expect %v equal %v", v1, v2)

	return e
}

// ToEqual 数值比较
func (e *Expect) ToEqual(v any) *Expect {
	if reflect.DeepEqual(e.Value, v) {
		return e
	}

	e.Error = fmt.Sprintf("expect %v toEqual %v", e.Value, v)
	return e
}

func (e *Expect) Pass() bool {
	return len(e.Error) == 0
}
