package md5

import (
	"crypto/md5"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"os"
)

func CalcMd5(l *lua.LState) int {
	path := l.Get(1).String()

	f, err := os.OpenFile(path, os.O_RDONLY, 660)
	if err != nil {
		fmt.Printf("open file %s err %s\n", err.Error())
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ReadAll file", err)
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}

	fMd5 := md5.Sum(body)
	fmt.Printf("%s md5 is %x\n", path, fMd5)
	str := fmt.Sprintf("%x", fMd5)
	l.Push(lua.LString(str))
	l.Push(lua.LNil)

	return 2
}

var exports = map[string]lua.LGFunction{"CalcFileMd5": CalcMd5}


func Load(L *lua.LState) int {
	tb := L.NewTable()
	L.SetFuncs(tb, exports)

	L.Push(tb)
	return 1
}