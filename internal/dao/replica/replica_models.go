package replica

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/pool"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 副本模型
type Replica struct {
	Id         uint                 `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	CreatedAt  util.Time            `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt  util.Time            `gorm:"column:updated_time; type:datetime; not null"`
	DeletedAt  *util.Time           `gorm:"column:deleted_time; type:datetime"`
	Name       string               `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc       string               `gorm:"column:desc; type:varchar(512)"                     json:"desc"`
	Uuid       string               `gorm:"column:uuid; type:varchar(32); not null"            json:"uuid"`
	Size       uint64               `gorm:"column:size; type:bigint; not null"                 json:"size"`
	Type       contants.ImageType   `gorm:"column:type; type:int; not null; default:1003"      json:"type"`
	Status     contants.ImageStatus `gorm:"column:status; type:int; not null; default:1024"    json:"status"`
	Export     string               `gorm:"column:export; type:varchar(512)"                   json:"export"`
	Pool       pool.StorePool       `gorm:"ForeignKey:PoolId;AssociationForeignKey:Id"         json:"pool"`
	PoolId     uint                 `gorm:"column:pool_id; type:int(11); not null"             json:"pool_id"`
	CreateUser string               `gorm:"column:create_user; type:varchar(50); not null"     json:"create_user"`
}

func (p *Replica) TableName() string {
	return "tangula_replica"
}

// 类型转换方法
func (p *Replica) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	data := make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		k := t.Elem().Field(i).Tag.Get("json")
		if len(k) == 0 {
			continue
		}
		if k == "pool" {
			if v.Elem().Field(i).Interface() != nil {
				temp, ok := v.Elem().Field(i).Interface().(pool.StorePool)
				if ok {
					pol := temp.TransformMap()
					data[k] = pol
				}
			}
			continue
		}
		data[k] = v.Elem().Field(i).Interface()
	}
	return data
}

// 副本挂载信息模型[挂载成功后生成，卸载后删除]
type MountInfo struct {
	Id         uint             `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	ReplicaId  uint             `gorm:"column:replica_id; type:int(11); not null"          json:"replica_id"`
	TargetType contants.AppType `gorm:"column:target_type; type:int(11); not null"         json:"target_type"`
	TargetId   uint             `gorm:"column:target_id; type:int(11); not null"           json:"target_id"`
	MountParam string           `gorm:"column:mount_param; type:varchar(8192);"            json:"mount_param"`
}

func (m *MountInfo) TableName() string {
	return "tangula_mount_info"
}

// 类型转换方法
func (m *MountInfo) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
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
