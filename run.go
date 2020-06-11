package main

import (
	"newdocker/mydocker/container"
	"os"

	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, command string) {
	//版本1：创建容器进程
	parent := container.NewParentProcess(tty, command)
	//start先执行init 后执行传递的参数比如bash ls等命令
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
