package lmhttp

import (
	lua "github.com/yuin/gopher-lua"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var exports = map[string]lua.LGFunction{
	"download": func(state *lua.LState) int {
		var url, name string

		url = state.Get(1).String()
		name = state.Get(2).String()

		resp, err := http.Get(url)
		if err != nil {
			state.Push(lua.LString(err.Error()))
			return 1
		}

		if resp.StatusCode != 200 {
			state.Push(lua.LString(resp.Status))
			return 1
		}

		fl, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0744)
		defer fl.Close()
		if err != nil {
			state.Push(lua.LString(err.Error()))
			return 1
		}

		if _, err := io.Copy(fl, resp.Body); err != nil {
			state.Push(lua.LString(err.Error()))
			return 1
		}
		if err := fl.Close(); err != nil {
			state.Push(lua.LString(err.Error()))
			return 1
		}
		if err := resp.Body.Close(); err != nil {
			state.Push(lua.LString(err.Error()))
			return 1
		}

		state.Push(lua.LNil)
		return 1
	},
	"get": func(state *lua.LState) int {
		url := state.Get(1).String()

		resp, err := http.Get(url)

		if err != nil {
			state.Push(lua.LNil)
			state.Push(lua.LString(err.Error()))
			return 2
		}

		if resp.StatusCode != 200 {
			state.Push(lua.LNil)
			state.Push(lua.LString(resp.Status))
			return 2
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			state.Push(lua.LNil)
			state.Push(lua.LString(err.Error()))
			return 2
		}

		state.Push(lua.LString(buf))
		state.Push(lua.LNil)
		return 2
	},
}

func Load(L *lua.LState) int {
	tb := L.NewTable()
	L.SetFuncs(tb, exports)
	L.Push(tb)
	return 1
}
