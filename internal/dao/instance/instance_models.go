package instance

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 执行实例模型
type Instance struct {
	Id         uint                    `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	Type       contants.InstanceType   `gorm:"column:type; type:int(11); not null;"               json:"type"`
	Status     contants.InstanceStatus `gorm:"column:status; type:int; not null; default:1"       json:"status"`
	ReplicaId  uint                    `gorm:"column:replica_id; type:int(11); not null"          json:"replica_id"`
	TargetType contants.AppType        `gorm:"column:target_type; type:int(11); not null"         json:"target_type"`
	TargetId   uint                    `gorm:"column:target_id; type:int(11); not null"           json:"target_id"`
	MountPoint string                  `gorm:"column:mount_point; type:varchar(512);"             json:"mount_point"`
	CreatedAt  util.Time               `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt  util.Time               `gorm:"column:updated_time; type:datetime; not null"`
	DeletedAt  *util.Time              `gorm:"column:deleted_time; type:datetime"`
}

func (i *Instance) TableName() string {
	return "tangula_instance"
}

// 类型转换方法
func (i *Instance) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
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

// 执行日志模型
type InstanceLog struct {
	Id         uint           `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	Level      contants.Level `gorm:"column:level; type:int; not null"                   json:"level"`
	Info       string         `gorm:"column:info; type:varchar(512); not null"           json:"info"`
	Detail     string         `gorm:"column:detail; type:varchar(512); not null"         json:"detail"`
	InstanceId uint           `gorm:"column:instance_id; type:int(11); not null"         json:"instance_id"`
	CreatedAt  util.Time      `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt  util.Time      `gorm:"column:updated_time; type:datetime; not null"`
	DeletedAt  *util.Time     `gorm:"column:deleted_time; type:datetime"`
}

func (i *InstanceLog) TableName() string {
	return "tangula_instance_log"
}

func (i *InstanceLog) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
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
