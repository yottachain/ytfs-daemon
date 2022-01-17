package main

import (
	"flag"
	"fmt"
	"github.com/natefinch/lumberjack"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"
	"time"
	"yottachain/ytfs-daemon/VM"
)

var isDaemon bool

var FileLogger = &lumberjack.Logger{
	Filename:   path.Join("output.log"),
	MaxSize:    128,
	Compress:   false,
	MaxAge:     7,
	MaxBackups: 30,
}

func main() {
	flag.BoolVar(&isDaemon, "d", false, "是否以守护进程启动")
	flag.Parse()

	if isDaemon {
		log.Println("日志文件:output.log")
		log.SetOutput(FileLogger)
		log.Println("守护进程已启动")
		//for {
		boot()
		//}
		log.Println("守护进程退出")
	} else {
		VM.Run("update.lua", "boot.lua")
	}
}

func boot() {
	log.Println("主脚本启动")
	defer func() {
		log.Println("主脚本退出")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	var daemonC *exec.Cmd
	go func() {
		var yOrN byte
		<-sigs
		fmt.Println("Are you sure you want to quit ?（y/n）")
		fmt.Scanf("%c\n", &yOrN)
		if yOrN == 'y' {
			daemonC.Process.Signal(syscall.SIGQUIT)
			os.Exit(0)
		}
	}()

	for {
		cmd := exec.Command(os.Args[0])
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
		cmd.Env = os.Environ()

		fl, err := os.OpenFile(".pid", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

		cmd.Start()
		if err != nil {
			log.Println(err.Error())
		} else {
			fl.WriteString(fmt.Sprintf("%d", cmd.Process.Pid))
			fl.Close()
		}
		cmd.Wait()

		time.Sleep(10 * time.Second)
	}
}


