package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

/*
	在/proc/self/mountinfo文件中查找subsystem的挂在目录
*/
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		fmt.Println("打不开/proc/self/mountinfi")
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f) //按行检索
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				//37 28 0:31 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:19 - cgroup cgroup rw,memory
				//返回挂载目录
				// fmt.Println("++++++++++++++++++++++", fields[4])
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

//将系统中的cgroup根目录和现在操作的cgroup目录拼接
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	fmt.Println("==========", cgroupRoot)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
