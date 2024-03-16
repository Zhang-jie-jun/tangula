package host

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 主机模型
type Host struct {
	Id           uint              `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	CreatedAt    util.Time         `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt    util.Time         `gorm:"column:updated_time; type:datetime; not null"       json:"updated_time"`
	DeletedAt    *util.Time        `gorm:"column:deleted_time; type:datetime"`
	Name         string            `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc         string            `gorm:"column:desc; type:varchar(512)"                     json:"desc"`
	HostName     string            `gorm:"column:hostname; type:varchar(128); not null"       json:"hostname"`
	Type         contants.AppType  `gorm:"column:type; type:int; not null"                    json:"type"`
	Status       uint              `gorm:"column:status; type:int; not null"                  json:"status"`
	DeployStatus uint              `gorm:"column:deployStatus; type:int; not null"           json:"deployStatus"`
	Os           string            `gorm:"column:os; type:varchar(128); not null"             json:"os"`
	Arch         string            `gorm:"column:arch; type:varchar(50); not null"            json:"arch"`
	Ip           string            `gorm:"column:ip; type:varchar(255); not null"             json:"ip"`
	Port         uint              `gorm:"column:port; type:int; not null"                    json:"port"`
	PlatformId   uint              `gorm:"column:platformId; type:int; not null"              json:"PlatformId"`
	UserName     string            `gorm:"column:username; type:varchar(255); not null"       json:"username"`
	PassWord     string            `gorm:"column:password; type:varchar(255); not null"`
	CreateUser   string            `gorm:"column:create_user; type:varchar(50); not null"     json:"create_user"`
	AuthType     contants.AuthType `gorm:"column:auth_type; type:int; not null; default:111"  json:"auth_type"`
}

type DeployApp struct {
	Id          uint       `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	CreatedAt   util.Time  `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt   util.Time  `gorm:"column:updated_time; type:datetime; not null"       json:"updated_time"`
	DeletedAt   *util.Time `gorm:"column:deleted_time; type:datetime"`
	HostId      uint       `gorm:"column:hostId; type:int; not null"                  json:"hostId"`
	Status      uint       `gorm:"column:status; type:int; not null"                  json:"status"`
	ServerIp    string     `gorm:"column:serverIp; type:varchar(128); not null"       json:"serverIp"`
	ClientIp    string     `gorm:"column:clientIp; type:varchar(128); not null"       json:"clientIp"`
	BaseDir     string     `gorm:"column:baseDir; type:varchar(200); not null"        json:"baseDir"`
	Apps        string     `gorm:"column:apps; type:varchar(500); not null"           json:"apps"`
	AbTpye      string     `gorm:"column:abTpye; type:varchar(128); not null"         json:"abTpye"`
	Log         string     `gorm:"column:log; type:longtext; not null"                json:"log"`
	Create_user string     `gorm:"column:create_user; type:varchar(128); not null"         json:"create_user"`
}

func (h *Host) TableName() string {
	return "tangula_host"
}

func (d *DeployApp) TableName() string {
	return "tangula_deploy_app"
}

//func (h *Host) AfterFind(tx *gorm.DB) (err error) {
//	h.PassWord = "******"
//	return
//}

func (h *Host) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	data := make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		k := t.Elem().Field(i).Tag.Get("json")
		if len(k) == 0 {
			continue
		}
		data[k] = v.Elem().Field(i).Interface()
	}
	return data
}

func (h *DeployApp) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	data := make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		k := t.Elem().Field(i).Tag.Get("json")
		if len(k) == 0 {
			continue
		}
		data[k] = v.Elem().Field(i).Interface()
	}
	return data
}
