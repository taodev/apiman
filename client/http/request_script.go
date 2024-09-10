package http

import (
	"fmt"

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

func (h *ApiHttp) doScript(code string) (err error) {
	l := lua.NewState()

	l.SetGlobal("request", luar.New(l, h.Request))
	l.SetGlobal("response", luar.New(l, h.Response))

	variables := storage.NewFromMemory()
	variables.SetData(h.Variables)
	l.SetGlobal("variables", luar.New(l, variables))
	l.SetGlobal("session", luar.New(l, h.sessionDB))

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

	if err = l.DoString(code); err != nil {
		return
	}

	return
}
