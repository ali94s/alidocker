package subsystems

/*
	用于传递资源限制配置的结构体
	MemoryLimit：内存限制
	CpuShare：cpu时间片权重
	CpuSet：cpu核心数
*/
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}


/*
	subsystem接口，每一个subsystem实现四个接口
	Name：返回subsystem的名字
	Set：设置某个cgroup在这个subsystem中的资源限制
	Apply：将进程添加到某个cgroup中
	remove：移除某个cgroup
*/
type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}


/*
	通过不同的subsystem初始化实例创建资源限制处理链数组
*/
var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
