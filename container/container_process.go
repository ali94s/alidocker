package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	NewWorkSpace(rootURL, mntURL)
	cmd.Dir = mntURL
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

//Create a AUFS filesystem as container root workspace
/*
	CreateReadOnlyLayer：新建busybox文件夹将busybox.tar解压到busybox目录，作为只读层
	CreateWriteLayer：创建一个名字为writeLayer的文件夹作为可写层
	CreateMountPoint：创建mnt文件夹，作为挂载点，然后把writeLayer和busybox目录mount到mnt目录
*/
func NewWorkSpace(rootURL string, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
}

func CreateReadOnlyLayer(rootURL string) {
	fmt.Println("=====CreateReadOnlyLayer")
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf("Mkdir  %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	fmt.Println("=====CreateWriteLayer")
	writeURL := rootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	fmt.Println("=====CreateMountPoint")
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntURL, err)
	}
	dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
	fmt.Println("=====dirs:", dirs)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

//Delete the AUFS filesystem while container exit
func DeleteWorkSpace(rootURL string, mntURL string) {
	DeleteMountPoint(rootURL, mntURL)
	DeleteWriteLayer(rootURL)
}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove dir %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("Remove dir %s error %v", writeURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
