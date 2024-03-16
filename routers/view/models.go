package view

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"time"
)

type IdParam struct {
	Id uint `uri:"id" binding:"required"` // URL绑定id
}

type PageQueryParam struct {
	Index int `form:"index"    json:"index"    binding:"gte=0"`   // 起始索引
	Count int `form:"count"    json:"count"    binding:"lte=100"` // 请求个数
}

type NameParam struct {
	Name string `form:"name"  json:"name"   binding:"required"` // 名称
	Desc string `form:"desc"  json:"desc"`                      // 名称
}

type QueryParam struct {
	PageQueryParam
	Filter string `form:"filter"    json:"filter"` // 名称参数
}

type UserQueryParam struct {
	PageQueryParam
	RoleName string `form:"roleName"  json:"roleName"` // 角色名称类型
	Filter   string `form:"filter"    json:"filter"`   // 用户名称过滤
}

type UserParam struct {
	UserName string `form:"username"    json:"username"    binding:"required"` // 用户名称
	PassWord string `form:"password"    json:"password"    binding:"required"` // 用户密码
}

type RoleParam struct {
	RoleType contants.RoleType `form:"roleName" json:"roleName"` // 角色类型
}

type QueryResourceParam struct {
	PageQueryParam
	Auth   contants.AuthType `form:"auth"    json:"auth"`   // 权限类型
	Type   contants.AppType  `form:"type"    json:"type"`   // 资源类型
	Filter string            `form:"filter"  json:"filter"` // 过滤参数
}

type CreateResource struct {
	Name     string           `form:"name"         json:"name"        binding:"required"` // 资源名称
	Desc     string           `form:"desc"         json:"desc"`                           // 资源描述
	Type     contants.AppType `form:"type"         json:"type"        binding:"required"` // 资源类型
	Ip       string           `form:"ip"           json:"ip"          binding:"required"` // 资源IP
	Port     uint             `form:"port"         json:"port"        binding:"required"` // 资源端口
	UserName string           `form:"username"     json:"username"    binding:"required"` // 资源登录用户名
	PassWord string           `form:"password"     json:"password"    binding:"required"` // 资源登录密码
}

type UpdateResource struct {
	Desc     string `form:"desc"         json:"desc"`                        // 资源描述
	Ip       string `form:"ip"           json:"ip"       binding:"required"` // 资源IP
	Port     uint   `form:"port"         json:"port"     binding:"required"` // 资源端口
	UserName string `form:"username"     json:"username" binding:"required"` // 资源登录用户名
	PassWord string `form:"password"     json:"password" binding:"required"` // 资源登录密码
}

type DeployParam struct {
	HostId       uint   `form:"hostId"           json:"hostId"          binding:"required"` //主机id
	Basedir      string `form:"basedir"         json:"basedir"          binding:"required"` // 安装目录
	ServerIp     string `form:"serverIp"         json:"serverIp"        binding:"required"` // 控制台IP
	Apps         string `form:"apps"         json:"apps"        binding:"required"`         // 应用列表
	Type         string `form:"type"         json:"type"        binding:"required"`         // 安装类型AB or HW
	BasicFtpPath string `form:"basicFtpPath" json:"basicFtpPath"`
	AppFtpPath   string `form:"appFtpPath"  json:"appFtpPath"`
}

type UpdDeployParam struct {
	Id     uint   `form:"id"           json:"id"          binding:"required"`     //id
	Status uint   `form:"status"       json:"status"          binding:"required"` //status
	Log    string `form:"log"          json:"log"`                                // 日志
}

type LogRecordQueryParam struct {
	PageQueryParam
	User string `form:"user"    json:"user"` // 用户名称过滤
}

type ImageQueryParam struct {
	PageQueryParam
	Auth   contants.AuthType  `form:"auth"    json:"auth"`   // 镜像权限类型
	Type   contants.ImageType `form:"type"    json:"type"`   // 镜像类型
	Filter string             `form:"filter"  json:"filter"` // 过滤参数
}

type ImageCreateParam struct {
	Name string `form:"name"         json:"name"        binding:"required"` // 名称
	Desc string `form:"desc"         json:"desc"`                           // 描述
}

type ReplicaCreateParam struct {
	Name        string             `form:"name"         json:"name"        binding:"required"` // 名称
	Desc        string             `form:"desc"         json:"desc"`                           // 描述
	Type        contants.ImageType `form:"type"         json:"type"        binding:"required"` // 类型
	Size        uint64             `form:"size"         json:"size"        binding:"required"` // 大小
	StorePoolId string             `form:"storePoolId"  json:"storePoolId" binding:"required"` // 存储池UUID
}

type ReplicaEditParam struct {
	Name string             `form:"name"         json:"name"` // 描述
	Desc string             `form:"desc"         json:"desc"` // 描述
	Size uint64             `form:"size"         json:"size"` // 大小
	Type contants.ImageType `form:"type"         json:"type"` // 类型
}

type ReplicaQueryParam struct {
	PageQueryParam
	Status contants.ImageStatus `form:"status"  json:"status"` // 副本状态
	Type   contants.ImageType   `form:"type"    json:"type"`   // 副本类型
	Filter string               `form:"filter"  json:"filter"` // 过滤参数
	Uuid   string               `form:"uuid"  json:"uuid"`     // uuid
}

type MountInfo struct {
	ReplicaId    uint                   `form:"replicaId"          json:"replicaId"     binding:"required"` // 目标副本id
	TargetType   contants.AppType       `form:"targetType"         json:"targetType"    binding:"required"` // 目标平台类型
	TargetId     uint                   `form:"targetId"           json:"targetId"      binding:"required"` // 目标平台id
	MountType    contants.MountType     `form:"mountType"          json:"mountType"     binding:"required"` // 挂载类型
	IsExecScript bool                   `form:"isExecScript"       json:"isExecScript"`                     // 是否执行脚本
	ScriptId     uint                   `form:"scriptId"           json:"scriptId"`                         // 所选脚本ID
	AppConfig    map[string]interface{} `form:"appConfig"          json:"appConfig"`                        // 应用配置
}

type MountParam struct {
	MountInfo MountInfo `form:"mountInfo"    json:"mountInfo"      binding:"required"` // 挂载信息
}

type BatchMountParam struct {
	MountInfos []MountInfo `form:"mountInfos"       json:"mountInfos"   binding:"required"` // 挂载信息列表
}

type BatchUnmountParam struct {
	UnMountInfo []uint `form:"unMountInfo"    json:"unMountInfo"   binding:"required"` // 卸载载信息列表
}
type BatchDeleteReplicaParam struct {
	ReplicaIdList []uint `form:"replicaIdList"    json:"replicaIdList"   binding:"required"`
}

type ScriptCreateParam struct {
	Desc string `form:"desc"         json:"desc"` // 描述
}

type VMwareDataSources struct {
	FullPath string `form:"fullPath"    json:"fullPath"` // 数据源路径
	HostName string `form:"hostName"    json:"hostName"` // 主机名称
	ShowType int    `form:"showType"    json:"showType"` // 浏览类型
	ISGetAll bool   `form:"isGetAll"    json:"isGetAll"` // 是否获取所有
	IsGetVm  bool   `form:"isGetVm"     json:"isGetVm"`  // 是否获取虚拟机
}

type VMwareHostParam struct {
	Path string `form:"path"         json:"path" binding:"required"`
}

type CasHostParam struct {
	Id     uint `form:"id"         json:"id" binding:"required"`
	HostId uint `form:"hostId"      json:"hostId" binding:"required"`
}

type UpdIdParam struct {
	Id uint `form:"id"         json:"id" binding:"required"`
}

type DoCasJobParam struct {
	Id        uint   `form:"id"         json:"id" binding:"required"`
	ConsoleIp string `form:"consoleIp"         json:"consoleIp" binding:"required"`
	Tag       string `form:"tag"         json:"tag" binding:"required"`
}

type MetaParam struct {
	ConsoleIp string `form:"consoleIp"        json:"consoleIp" binding:"required"`
	Port      string `form:"port"             json:"port" binding:"required"`
	Username  string `form:"username"         json:"username" binding:"required"`
	Password  string `form:"password"         json:"password" binding:"required"`
	Jobname   string `form:"jobname"          json:"jobname" binding:"required"`
}

type ImageEditParam struct {
	Name string `form:"name"         json:"name"` // 描述
	Desc string `form:"desc"         json:"desc"` // 描述
}

type VmFileRecoveryResult struct {
	Sequence  string    `form:"sequence" json:"sequence"`
	Pass      uint      `form:"pass"         json:"pass"`
	Failed    uint      `form:"failed"         json:"failed"`
	Duration  uint      `form:"duration"         json:"duration"`
	Report    string    `form:"report"         json:"report"`
	StartTime time.Time `gorm:"column:startTime; type:datetime; not null"     json:"startTime"`
}

type DiskInfo struct {
	Name     string
	PoolName string
	SizeByte int64
}

type ImageDiskInfo struct {
	DiskList []DiskInfo
}
