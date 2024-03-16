package system

import (
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/host"
	"github.com/Zhang-jie-jun/tangula/internal/dao/image"
	"github.com/Zhang-jie-jun/tangula/internal/dao/platform"
	"github.com/Zhang-jie-jun/tangula/internal/dao/pool"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/internal/dao/script"
)

type UserStat struct {
	TotalNum  int `json:"totalNum"`  // 用户总数
	ActiveNum int `json:"activeNum"` // 活跃用户数
	TodayNum  int `json:"todayNum"`  // 今日访问用户数
}

type ResourceStat struct {
	Public  int `json:"public"`
	Private int `json:"private"`
}

type Dashboard struct {
	UserStat     UserStat     `json:"userStat"`
	PoolStat     int          `json:"poolStat"`
	PlatformStat ResourceStat `json:"platformStat"`
	HostStat     ResourceStat `json:"hostStat"`
	ReplicaStat  int          `json:"replicaStat"`
	ImageStat    ResourceStat `json:"imageStat"`
	ScriptStat   int          `json:"scriptStat"`
}

//@author: [jick.zhang](zhang.jiejun@outlook.com)
//@function: GetDashboardInfo
//@description: 获取统计信息
//@return: c Dashboard, err error

func GetDashboardInfo(userName string) *Dashboard {
	var d Dashboard
	d.UserStat = InitUserStat()
	d.PoolStat = InitPoolStat()
	d.PlatformStat = InitPlatformStat(userName)
	d.HostStat = InitHostStat(userName)
	d.ReplicaStat = InitReplicaStat(userName)
	d.ImageStat = InitImageStat(userName)
	d.ScriptStat = InitScriptStat(userName)
	return &d
}

func InitUserStat() (u UserStat) {
	u.TotalNum, u.ActiveNum, u.TodayNum = auth.AuthMgm.GetUserStat()
	return u
}

func InitPoolStat() int {
	return pool.StorePoolMgm.GetPoolStat()
}

func InitPlatformStat(userName string) ResourceStat {
	var r ResourceStat
	r.Public, r.Private = platform.PlatformMgm.GetPlatformStat(userName)
	return r
}

func InitHostStat(userName string) ResourceStat {
	var r ResourceStat
	r.Public, r.Private = host.HostMgm.GetHostStat(userName)
	return r
}

func InitReplicaStat(userName string) int {
	return replica.ReplicaMgm.GetReplicaStat(userName)
}

func InitImageStat(userName string) ResourceStat {
	var r ResourceStat
	r.Public, r.Private = image.ImageMgm.GetImageStat(userName)
	return r
}

func InitScriptStat(userName string) int {
	return script.ScriptMgm.GetScriptStat(userName)
}
