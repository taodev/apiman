package test

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type TestClass struct {
	Name    string
	Expects []*Expect

	Error string
}

func (t *TestClass) Test(l *lua.LState) int {
	t.Name = l.CheckString(1)
	testFunc := l.CheckFunction(2)
	t.Error = ""
	t.Expects = make([]*Expect, 0)

	l.SetGlobal("expect", luar.New(l, t.Expect))

	// 执行testFunc，原型为function()，无返回值
	l.Push(testFunc)
	l.Call(0, 0)

	for _, v := range t.Expects {
		if v.Pass() {
			continue
		}

		t.Error += fmt.Sprintf("%s\n", v.Error)
	}

	if len(t.Error) > 0 {
		fmt.Printf("test %s: faild\n  %s", t.Name, t.Error)
		return 0
	}

	fmt.Printf("test %s: pass\n", t.Name)

	return 0
}

func (t *TestClass) Expect(v any) *Expect {
	e := new(Expect)
	e.Value = v

	t.Expects = append(t.Expects, e)
	return e
}
