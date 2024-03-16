package image

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/pool"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 镜像模型
type Image struct {
	Id         uint                 `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	CreatedAt  util.Time            `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt  util.Time            `gorm:"column:updated_time; type:datetime; not null"`
	DeletedAt  *util.Time           `gorm:"column:deleted_time; type:datetime"`
	Name       string               `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc       string               `gorm:"column:desc; type:varchar(512); not null"           json:"desc"`
	Uuid       string               `gorm:"column:uuid; type:varchar(32); not null"            json:"uuid"`
	Size       uint64               `gorm:"column:size; type:bigint; not null"                 json:"size"`
	Type       contants.ImageType   `gorm:"column:type; type:int; not null; default:1003"      json:"type"`
	Status     contants.ImageStatus `gorm:"column:status; type:int; not null; default:1024"    json:"status"`
	Pool       pool.StorePool       `gorm:"ForeignKey:PoolId;AssociationForeignKey:Id"         json:"pool"`
	PoolId     uint                 `gorm:"column:pool_id; type:int(11); not null"             json:"pool_id"`
	CreateUser string               `gorm:"column:create_user; type:varchar(50); not null"     json:"create_user"`
	AuthType   contants.AuthType    `gorm:"column:auth_type; type:int; not null; default:111"  json:"auth_type"`
}

func (i *Image) TableName() string {
	return "tangula_image"
}

// 类型转换方法
func (i *Image) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
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
