package lmos

import (
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

func kill(L *lua.LState) int {
	pidstr := L.Get(1).String()
	if runtime.GOOS == "linux" {
		exec.Command("kill", "-9", pidstr).Output()
	} else {
		L.Push(lua.LString("windows no kill cmd"))
		return 1
	}
	return 0
}

var exports = map[string]lua.LGFunction{
	"kill": kill,
	"killSelf": func(state *lua.LState) int {
		fl, err := os.OpenFile(".pid", os.O_RDONLY, 0644)
		if err == nil {
			pidbuf, err := ioutil.ReadAll(fl)
			defer fl.Close()
			if err == nil {
				if runtime.GOOS == "linux" {
					exec.Command("kill", "-9", string(pidbuf)).Output()
				}
			}
		}
		return 0
	},
	"pkill": func(state *lua.LState) int {
		name := state.Get(1).String()
		switch runtime.GOOS {
		case "linux", "darwin":
			exec.Command("pkill", "-9", name)
		case "windows":
			exec.Command("taskkill.exe", "/f", "/im", name)
		}
		return 0
	},
}

func Load(L *lua.LState) int {
	tb := L.NewTable()
	L.SetFuncs(tb, exports)

	L.Push(tb)
	return 1
}
