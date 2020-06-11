package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

//memory subsystem的实现
type MemorySubSystem struct {
}

//设置cgroupPath对应的cgroup的内存限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	//GetcgroupPath()是获取当前subsystem在虚拟文件系统中的路径
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}

}

//移除cgroupPath对应的cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		//删除cgroup对应的目录即可
		fmt.Println(subsysCgroupPath)
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

//讲一个进程加入到cgroupPath对应的cgroup
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		//讲进程id写入到cgroup的虚拟文件系统对应目录下的task文件中
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

//返回cgroup的名字
func (s *MemorySubSystem) Name() string {
	return "memory"
}
