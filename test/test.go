package test

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type TestManager struct {
	Status string       `yaml:"status,omitempty"`
	Cases  []*TestClass `yaml:"cases,omitempty"`
}

func (tm *TestManager) Test(l *lua.LState) int {
	t := new(TestClass)
	tm.Cases = append(tm.Cases, t)

	t.Name = l.CheckString(1)
	testFunc := l.CheckFunction(2)

	l.SetGlobal("expect", luar.New(l, t.Expect))

	// 执行testFunc，原型为function()，无返回值
	l.Push(testFunc)
	l.Call(0, 0)

	t.Done()

	return 0
}

func (tm *TestManager) Pass() bool {
	return tm.Status == "PASS"
}

func (tm *TestManager) Done() {
	for _, v := range tm.Cases {
		if !v.Pass() {
			tm.Status = "FAIL"
			return
		}
	}

	tm.Status = "PASS"
}

type TestClass struct {
	Name    string `yaml:"name,omitempty"`
	Status  string `yaml:"status,omitempty"`
	expects []*Expect
	Error   []string `yaml:"expects,omitempty"`
}

func (t *TestClass) Expect(v any, desc string) *Expect {
	e := new(Expect)
	e.Value = v
	e.desc = desc

	t.expects = append(t.expects, e)
	return e
}

func (t *TestClass) Pass() bool {
	for _, v := range t.expects {
		if !v.Pass() {
			return false
		}
	}

	return true
}

func (t *TestClass) Done() {
	t.Status = "PASS"

	for _, v := range t.expects {
		if !v.Pass() {
			t.Status = "FAIL"
			t.Error = append(t.Error, v.Error)
			continue
		}
	}
}
