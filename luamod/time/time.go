package time

import (
	lua "github.com/yuin/gopher-lua"
	"strconv"
	"time"
)

var exports = map[string]lua.LGFunction{
	"sleep": func(state *lua.LState) int {
		var d time.Duration
		arg1 := state.Get(1).String()
		i, err := strconv.ParseInt(arg1, 10, 64)
		if err != nil {
			return 0
		}

		d = time.Duration(i)
		time.Sleep(d * time.Millisecond)
		return 0
	},
}

func Load(L *lua.LState) int {
	tb := L.NewTable()
	tb.RawSetString("second", lua.LNumber(1000))
	tb.RawSetString("minute", lua.LNumber(1000*60))
	L.SetFuncs(tb, exports)

	L.Push(tb)
	return 1
}
