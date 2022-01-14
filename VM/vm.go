package VM

import (
	lua "github.com/yuin/gopher-lua"
	"log"
	"sync"
	"yottachain/ytfs-daemon/luamod/cmd"
	"yottachain/ytfs-daemon/luamod/config"
	lmhttp "yottachain/ytfs-daemon/luamod/http"
	md5 "yottachain/ytfs-daemon/luamod/md5"
	lmos "yottachain/ytfs-daemon/luamod/process"
	tm "yottachain/ytfs-daemon/luamod/time"
)

var pool sync.Pool

func init() {
	pool = sync.Pool{New: getVM}
}

func GetVM() *lua.LState {
	ls := pool.Get().(*lua.LState)
	return ls
}

func PutVM(L *lua.LState) {
	pool.Put(L)
}

func getVM() interface{} {
	ls := lua.NewState()
	ls.PreloadModule("cmd", cmd.Load)
	ls.PreloadModule("time", tm.Load)
	ls.PreloadModule("http", lmhttp.Load)
	ls.PreloadModule("process", lmos.Load)
	ls.PreloadModule("config", config.Load)
	ls.PreloadModule("md5", md5.Load)

	load(ls)

	return ls
}

func Run(shells ...string) {
	wg := sync.WaitGroup{}
	for _, v := range shells {
		wg.Add(1)
		go func(shell string) {
			defer wg.Done()

			vm := GetVM()
			defer PutVM(vm)

			if err := vm.DoFile(shell); err != nil {
				log.Println(err.Error())
			}
		}(v)
	}
	wg.Wait()
}
