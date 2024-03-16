package script

import (
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 脚本模型
type Script struct {
	Id         uint       `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	Uuid       string     `gorm:"column:uuid; type:varchar(32); not null"            json:"uuid"`
	Name       string     `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc       string     `gorm:"column:desc; type:varchar(512)"                     json:"desc"`
	Label      string     `gorm:"column:label; type:varchar(50)"                     json:"label"`
	CreateUser string     `gorm:"column:create_user; type:varchar(50); not null"     json:"create_user"`
	CreatedAt  util.Time  `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt  util.Time  `gorm:"column:updated_time; type:datetime; not null"`
	DeletedAt  *util.Time `gorm:"column:deleted_time; type:datetime"`
}

func (s *Script) TableName() string {
	return "tangula_script"
}

// 类型转换方法
func (s *Script) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
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
