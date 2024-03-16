package vmware

import "github.com/vmware/govmomi/vim25/types"

// 挂载参数
type MountParam struct {
	LocationPath    string `json:"LOCATION_PATH"`     // 虚拟机位置路径
	ComputeResource string `json:"COMPUTE_RESOURCE"`  // 虚拟机计算资源路径
	IsRegisterVm    string `json:"IS_REGISTER_VM"`    // 是否注册虚拟机
	VmConfig        string `json:"CUSTOMIZE_VM_INFO"` // 虚拟机定制配置
}

// 虚拟机自定义参数
type CustomizeVmInfo struct {
	VmName   string `json:"vmname"`   // 虚拟机名称
	PowerOn  string `json:"poweron"`  // 是否打开虚拟机电源
	HostName string `json:"hostname"` // 虚拟机系统主机名
	Addr     string `json:"addr"`     // 虚拟机系统IP
	Netmask  string `json:"netmask"`  // 虚拟机系统网关
	Gateway  string `json:"gateway"`  // 虚拟机系统掩码
}

type vmBaseInfo struct {
	VmxPath string
	VmName  string
	VmUuid  string
}

// 挂载详细信息
type MountDetailInfo struct {
	LocationPath     string
	HostPath         string
	HostName         string
	ResourcePoolPath string

	IsRegisterVm  bool
	IsAutoPowerOn bool
	IsSetIp       bool

	StoreName  string
	StoreType  string
	StoreMode  string
	RemoteHost string
	RemotePath string

	VmxPath string
	VmName  string
	VmUuid  string

	VmUsername string
	VmPassword string
	VmHostName string
	VmAddr     string
	VmNetmask  string
	VmGateway  string
	Os         string
}

type VirtualMachines struct {
	Uuid       string
	Name       string
	System     string
	Ip         string
	Hostname   string
	PowerState string
	Self       Self
	VM         types.ManagedObjectReference
}

type TemplateInfo struct {
	Name   string
	System string
	Self   Self
	VM     types.ManagedObjectReference
}

type DatastoreSummary struct {
	Datastore          Datastore `json:"Datastore"`
	Name               string    `json:"Name"`
	URL                string    `json:"Url"`
	Capacity           int64     `json:"Capacity"`
	FreeSpace          int64     `json:"FreeSpace"`
	Uncommitted        int64     `json:"Uncommitted"`
	Accessible         bool      `json:"Accessible"`
	MultipleHostAccess bool      `json:"MultipleHostAccess"`
	Type               string    `json:"Type"`
	MaintenanceMode    string    `json:"MaintenanceMode"`
	DatastoreSelf      types.ManagedObjectReference
}

type Datastore struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type HostSummary struct {
	Host        Host   `json:"Host"`
	Name        string `json:"Name"`
	UsedCPU     int64  `json:"UsedCPU"`
	TotalCPU    int64  `json:"TotalCPU"`
	FreeCPU     int64  `json:"FreeCPU"`
	UsedMemory  int64  `json:"UsedMemory"`
	TotalMemory int64  `json:"TotalMemory"`
	FreeMemory  int64  `json:"FreeMemory"`
	HostSelf    types.ManagedObjectReference
}

type Host struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type HostVM struct {
	Host map[string][]VMS
}

type VMS struct {
	Name  string
	Value string
}

type DataCenter struct {
	Datacenter      Self
	Name            string
	VmFolder        Self
	HostFolder      Self
	DatastoreFolder Self
}

type ClusterInfo struct {
	Cluster      Self
	Name         string
	Parent       Self
	ResourcePool Self
	Hosts        []types.ManagedObjectReference
	Datastore    []types.ManagedObjectReference
}

type ResourcePoolInfo struct {
	ResourcePool     Self
	Name             string
	Parent           Self
	ResourcePoolList []types.ManagedObjectReference
	Resource         types.ManagedObjectReference
}

type FolderInfo struct {
	Folder      Self
	Name        string
	ChildEntity []types.ManagedObjectReference
	Parent      Self
	FolderSelf  types.ManagedObjectReference
}

type Self struct {
	Type  string
	Value string
}

type CreateMap struct {
	TempName    string
	Datacenter  string
	Cluster     string
	Host        string
	Resources   string
	Storage     string
	VmName      string
	SysHostName string
	Network     string
}
