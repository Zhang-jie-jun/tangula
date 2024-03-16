package system

// 磁盘文件系统信息
type DevInfo struct {
	DevPath string // 设备名称
	Uuid    string // 文件系统UUID
	Type    string // 文件系统类型
}

// 系统服务状态
type ServerStatus int

const (
	NONE     ServerStatus = 0
	ACTIVE   ServerStatus = 1
	INACTIVE ServerStatus = 2
	FAILED   ServerStatus = 3
)
