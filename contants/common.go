package contants

// 服务依赖全局变量
var (
	ConfPath *string    // 配置文件路径
	AppCfg   *AppConfig // 服务配置器
)

// 运行模式
const (
	DEBUG   = "debug"
	RELEASE = "release"
)

// 基础角色类型
type RoleType string

const (
	USER       RoleType = "普通用户"
	SYSADMIN   RoleType = "系统管理员"
	SUPERADMIN RoleType = "超级管理员"
)

// 用户状态
type UserStatus int

const (
	DISABLE UserStatus = -1 // 禁用
	ENABLE  UserStatus = 1  // 启用
)

// 存储资源类型
type CephType int

const (
	STOREPOOL CephType = 101 // 存储池
	IMAGE     CephType = 102 // 数据镜像
	REPLICA   CephType = 103 // 数据副本
	SNAPSHOT  CephType = 104 // 副本快照
)

// 副本镜像类型
type ImageType int

const (
	CUSTOME         ImageType = 1003 // 自定义类型
	VMWAREVM        ImageType = 1004 // VMware
	CASVM           ImageType = 1005 // CAS
	FUSIONCOMPUTEVM ImageType = 1006 // FusionCompute
)

var ImageFlags = map[ImageType]string{
	CUSTOME:         "Custom",
	VMWAREVM:        "VMware",
	CASVM:           "CAS",
	FUSIONCOMPUTEVM: "FusionCompute",
}

// 副本镜像状态
type ImageStatus int

const (
	NOT_MOUNT  ImageStatus = 1024  // 未挂载
	MOUNTING   ImageStatus = 2048  // 挂载中
	MOUNTED    ImageStatus = 4096  // 已挂载
	UNMOUNTING ImageStatus = 8192  // 卸载中
	CLONEING   ImageStatus = 16384 // 克隆中
)

// 实例类型
type InstanceType int

const (
	INSTANCE_MOUNT   InstanceType = 1 // 挂载
	INSTANCE_UNMOUNT InstanceType = 2 // 卸载
)

// 实例状态
type InstanceStatus int

const (
	NONE           InstanceStatus = 1   // 未启动
	READY          InstanceStatus = 2   // 准备中
	MOUNT_RUNNING  InstanceStatus = 4   // 挂载中
	MOUNT_SUCCESS  InstanceStatus = 8   // 挂载成功
	MOUNT_FAILED   InstanceStatus = 16  // 挂载失败
	UMOUNT_RUNNING InstanceStatus = 32  // 卸载中
	UMOUNT_SUCCESS InstanceStatus = 64  // 卸载成功
	UMOUNT_FAILED  InstanceStatus = 128 // 卸载失败
	ABNORMAL       InstanceStatus = 256 // 异常
)

// 挂载类型
type MountType int

const (
	APP_DEFAULT MountType = 1
	EXPORT_PATH MountType = 2
)

var MountTypeFlags = map[MountType]string{
	APP_DEFAULT: "应用默认挂载",
	EXPORT_PATH: "仅导出共享路径",
}

// 应用类型
type AppType int

const (
	VMWARE        AppType = 11
	CAS           AppType = 12
	FUSIONCOMPUTE AppType = 13
	LINUX         AppType = 51
	WINDOWS       AppType = 52
)

var AppFlags = map[AppType]string{
	VMWARE:        "VMware",
	CAS:           "CAS",
	FUSIONCOMPUTE: "FusionCompute",
	LINUX:         "Linux",
	WINDOWS:       "Windows",
}

var ReAppFlags = map[string]AppType{
	"VMware":  VMWARE,
	"Cas":     CAS,
	"Linux":   LINUX,
	"Windows": WINDOWS,
}

// 资源权限类型
type AuthType int

const (
	PRIVATE AuthType = 111 // 私有
	PUBLIC  AuthType = 222 // 公有
)

// 日志级别
type Level int

const (
	LOG_INFO     Level = 1 // 普通信息
	LOG_WARNNING Level = 2 // 警告信息
	LOG_ERROR    Level = 3 // 错误信息
	LOG_FATAL    Level = 4 // 严重错误
)

// 日志状态
type LogStatus int

const (
	LOG_SUCCESS LogStatus = 1 // 成功
	LOG_FAILED  LogStatus = 2 // 失败
)

// 认证登录信息
type LoginInfo struct {
	Ip       string
	Port     uint
	UserName string
	PassWord string
}

// 认证响应信息
type VerifyResult struct {
	// 虚拟化平台使用
	Version      string
	HostpoolInfo map[string]interface{} //主机池，cas平台使用
	// 租户使用
	DomainId string
	// 主机使用
	HostName string
	OSType   string
	Arch     string
	// 自定义响应，建议使用JSON字符串
	Custom string
}

type MountInfo struct {
	UserId      uint   // 用户id
	ReplicaId   uint   // 副本id
	TargetId    uint   // 目标平台id
	RecordId    uint   // 任务记录id
	MountPoint  string // 挂载目录
	CustomParam string // 自定义挂载参数, json字符串
}
