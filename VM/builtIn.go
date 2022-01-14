package VM

import (
	"crypto/md5"
	"encoding/hex"
	lua "github.com/yuin/gopher-lua"
	"io"
	"os"
	"os/user"
	"runtime"
)

func run(L *lua.LState) int {
	shellPath := L.Get(1).String()
	if shellPath == "" {
		L.Push(lua.LString("shellPath is nil"))
		return 1
	}

	vm := GetVM()
	err := vm.DoFile(shellPath)
	PutVM(vm)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}

	return 0
}

func md5sum(L *lua.LState) int {
	filename := L.Get(1).String()

	fl, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	m5 := md5.New()

	io.Copy(m5, fl)

	resbuf := m5.Sum(nil)
	resstr := hex.EncodeToString(resbuf)

	L.Push(lua.LString(resstr))
	return 1
}

func mv(L *lua.LState) int {
	oname := L.Get(1).String()
	nname := L.Get(2).String()

	err := os.Rename(oname, nname)

	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// GetCurrentUserHome 获取当前用户主目录
func GetCurrentUserHome() string {
	var userDir string
	if u, err := user.Current(); err == nil {
		userDir = u.HomeDir
	}
	return userDir
}

// GetYTFSPath 获取YTFS文件存放路径
//
// 如果存在环境变量ytfs_path则使用环境变量ytfs_path
func GetYTFSPath() string {
	ps, ok := os.LookupEnv("ytfs_path")
	if ok {
		return ps
	}
	return GetCurrentUserHome() + "/YTFS"
}

func load(L *lua.LState) int {
	// 运行新的脚本
	L.SetGlobal("run", L.NewFunction(run))
	L.SetGlobal("md5sum", L.NewFunction(md5sum))
	L.SetGlobal("mv", L.NewFunction(mv))
	L.SetGlobal("YTFS_PATH", L.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(GetYTFSPath()))
		return 1
	}))
	// 常量---------------------------------
	L.SetGlobal("L_OS", lua.LString(runtime.GOOS))
	L.SetGlobal("L_ARCH", lua.LString(runtime.GOARCH))
	// 常量---------------------------------

	return 0
}
