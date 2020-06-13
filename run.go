package main

import (
	"newdocker/alidocker/cgroups"
	"newdocker/alidocker/cgroups/subsystems"
	"newdocker/alidocker/container"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	cgroupManager := cgroups.NewCgroupManager("donkey-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkSpace(rootURL, mntURL)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
