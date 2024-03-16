package platform

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 虚拟化(云)平台模型
type Platform struct {
	Id         uint              `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	CreatedAt  util.Time         `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt  util.Time         `gorm:"column:updated_time; type:datetime; not null"`
	DeletedAt  *util.Time        `gorm:"column:deleted_time; type:datetime"`
	Name       string            `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc       string            `gorm:"column:desc; type:varchar(512)"                     json:"desc"`
	Type       contants.AppType  `gorm:"column:type; type:int; not null; default:-1"        json:"type"`
	Ip         string            `gorm:"column:ip; type:varchar(255); not null"             json:"ip"`
	Port       uint              `gorm:"column:port; type:int; not null"                    json:"port"`
	UserName   string            `gorm:"column:username; type:varchar(255); not null"       json:"username"`
	PassWord   string            `gorm:"column:password; type:varchar(255); not null"`
	Version    string            `gorm:"column:version; type:varchar(255); not null"        json:"version"`
	CreateUser string            `gorm:"column:create_user; type:varchar(50); not null"     json:"create_user"`
	AuthType   contants.AuthType `gorm:"column:auth_type; type:int; not null; default:111"  json:"auth_type"`
}

func (p *Platform) TableName() string {
	return "tangula_platform"
}

func (p *Platform) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
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

// 租户模型
type Tenant struct {
	Id         uint     `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	Name       string   `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc       string   `gorm:"column:desc; type:varchar(512)"                     json:"desc"`
	DomainId   string   `gorm:"column:domain_id; type:varchar(255)"                json:"domain_id"`
	UserName   string   `gorm:"column:username; type:varchar(255); not null"       json:"username"`
	PassWord   string   `gorm:"column:password; type:varchar(255); not null"`
	Platform   Platform `gorm:"ForeignKey:PlatformId;AssociationForeignKey:Id"     json:"platform"`
	PlatformId uint     `gorm:"column:platform_id; type:int(11); not null"         json:"platform_id"`
}

func (te *Tenant) TableName() string {
	return "tangula_tenant"
}

func (te *Tenant) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(te)
	v := reflect.ValueOf(te)
	data := make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		k := t.Elem().Field(i).Tag.Get("json")
		if len(k) == 0 {
			continue
		}
		if k == "platform" {
			if v.Elem().Field(i).Interface() != nil {
				temp, ok := v.Elem().Field(i).Interface().(Platform)
				if ok {
					platform := temp.TransformMap()
					data[k] = platform
				}
			}
			continue
		}
		data[k] = v.Elem().Field(i).Interface()
	}
	return data
}
