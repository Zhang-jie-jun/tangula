package fusioncompute

// 挂载详细信息
type MountFcInfo struct {
	HostName      string
	HostUrn       string
	HostIp        string
	Os            string
	VsId          uint
	ClusterId     uint
	HostpoolId    uint
	StorePath     string
	StoreName     string //存储资源名称
	DataStoreName string //数据存储名称
	RemoteHost    string
	RemotePath    string
	IsRegisterVM  bool
	IsRefresh     bool //是否在FusionCompute立即扫描存储设备
	Title         string
	PoolName      string
}
