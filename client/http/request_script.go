package http

import (
	"fmt"
	"strings"
	"time"

	json "github.com/json-iterator/go"

	"github.com/taodev/apiman/storage"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func (h *ApiHttp) processBeforeScript() (err error) {
	for _, v := range h.BeforeScript {
		if err = h.doScript(v); err != nil {
			return
		}
	}

	return
}

func (h *ApiHttp) processAfterScript() (err error) {
	for _, v := range h.AfterScript {
		if err = h.doScript(v); err != nil {
			return
		}
	}
	return
}

func (h *ApiHttp) doScript(codeField LineField[string]) (err error) {
	l := lua.NewState()

	l.SetGlobal("request", luar.New(l, h.Request))
	l.SetGlobal("response", luar.New(l, h.Response))

	variables := storage.NewFromMemory()
	variables.SetData(h.Variables)
	l.SetGlobal("variables", luar.New(l, variables))
	l.SetGlobal("session", luar.New(l, h.sessionDB))
	l.SetGlobal("global", luar.New(l, storage.GlobalDB))
	l.SetGlobal("environment", luar.New(l, storage.EnvironmentDB))

	l.SetGlobal("printf", luar.New(l, fmt.Printf))
	l.SetGlobal("tojson", luar.New(l, func(v any) string {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(jsonBytes)
	}))
	l.SetGlobal("json", luar.New(l, func(v string) any {
		var out any
		if err = json.Unmarshal([]byte(v), &out); err != nil {
			return nil
		}

		return out
	}))
	l.SetGlobal("wait", luar.New(l, func(t int64) {
		<-time.After(time.Duration(t) * time.Millisecond)
	}))

	l.SetGlobal("test", l.NewFunction(h.Test.Test))

	// 设置包搜索路径
	l.SetGlobal("workDir", lua.LString(h.workDir))
	l.SetGlobal("fileDir", lua.LString(h.fileDir))

	packagePath := l.GetGlobal("package")
	if tbl, ok := packagePath.(*lua.LTable); ok {
		p1 := tbl.RawGetString("path").String()
		p2 := fmt.Sprintf("%s;%s/?.lua;%s/?.lua", p1, h.workDir, h.fileDir)
		tbl.RawSetString("path", lua.LString(p2))
	}

	// 通过在lua脚本中，插入对应的行数，用取巧的方式获取正确行数，方便调试时定位!_!
	var code string
	for i := 0; i < codeField.Line; i++ {
		code += "\n"
	}
	code += codeField.Value

	var fn *lua.LFunction
	if fn, err = l.Load(strings.NewReader(code), h.filepath); err != nil {
		return
	} else {
		l.Push(fn)
		err = l.PCall(0, lua.MultRet, nil)
	}

	return
}
