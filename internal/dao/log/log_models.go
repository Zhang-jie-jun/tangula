package log

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

type LogRecord struct {
	Id         uint               `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	Operation  string             `gorm:"column:operation; type:varchar(128); not null"      json:"operation"`
	Object     string             `gorm:"column:object; type:varchar(512)"                   json:"object"`
	Detail     string             `gorm:"column:detail; type:varchar(4096)"                  json:"detail"`
	Status     contants.LogStatus `gorm:"column:status; type:int; not null"                  json:"status"`
	User       string             `gorm:"column:user; type:varchar(50); not null"            json:"user"`
	CreateTime util.Time          `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
}

func (l *LogRecord) TableName() string {
	return "tangula_record"
}

func (h *LogRecord) TransformMap() map[string]interface{} {
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
