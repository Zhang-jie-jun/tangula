package snapshot

import (
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 快照模型
type Snapshot struct {
	Id        uint            `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	CreatedAt util.Time       `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt util.Time       `gorm:"column:updated_time; type:datetime; not null"`
	DeletedAt *util.Time      `gorm:"column:deleted_time; type:datetime"`
	Name      string          `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc      string          `gorm:"column:desc; type:varchar(512); not null"           json:"desc"`
	Uuid      string          `gorm:"column:uuid; type:varchar(32); not null"            json:"uuid"`
	Replica   replica.Replica `gorm:"ForeignKey:ReplicaId;AssociationForeignKey:Id"      json:"replica"`
	ReplicaId uint            `gorm:"column:replica_id; type:int(11); not null"          json:"replica_id"`
}

func (s *Snapshot) TableName() string {
	return "tangula_snapshot"
}

// 类型转换方法
func (s *Snapshot) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	data := make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		k := t.Elem().Field(i).Tag.Get("json")
		if len(k) == 0 {
			continue
		}
		if k == "replica" {
			// todo: 涉及到多层外键关系 不在此处转换
			continue
		}
		data[k] = v.Elem().Field(i).Interface()
	}
	return data
}
