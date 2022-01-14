package config

import (
	"github.com/spf13/viper"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

var exports = map[string]lua.LGFunction{
	"read": func(state *lua.LState) int {
		cfgstr := state.Get(1).String()
		if err := viper.ReadConfig(strings.NewReader(cfgstr)); err != nil {
			state.Push(lua.LString(err.Error()))
			return 1
		}
		return 0
	},
	"setType": func(state *lua.LState) int {
		cfgType := state.Get(1).String()

		viper.SetConfigType(cfgType)
		return 0
	},
	"reset": func(state *lua.LState) int {
		viper.Reset()
		return 0
	},
	"getString": func(state *lua.LState) int {
		key := state.Get(1).String()

		value := viper.GetString(key)
		state.Push(lua.LString(value))
		return 1
	},
	"getInt": func(state *lua.LState) int {
		key := state.Get(1).String()

		value := viper.GetInt64(key)
		state.Push(lua.LNumber(value))
		return 1
	},
	"getFloat": func(state *lua.LState) int {
		key := state.Get(1).String()

		value := viper.GetFloat64(key)
		state.Push(lua.LNumber(value))
		return 1
	},
	"getBool": func(state *lua.LState) int {
		key := state.Get(1).String()

		value := viper.GetBool(key)
		state.Push(lua.LBool(value))
		return 1
	},
}

func Load(L *lua.LState) int {
	tb := L.NewTable()
	L.SetFuncs(tb, exports)

	L.Push(tb)
	return 1
}
