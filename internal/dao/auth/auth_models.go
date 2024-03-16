package auth

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"reflect"
)

// 角色类型模型
type Role struct {
	Id        uint              `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	Name      contants.RoleType `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Desc      string            `gorm:"column:desc; type:varchar(50); not null"            json:"desc"`
	CreatedAt util.Time         `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
}

func (r *Role) TableName() string {
	return "tangula_role"
}

// 类型转换方法
func (r *Role) Transform() map[string]interface{} {
	t := reflect.TypeOf(r)
	v := reflect.ValueOf(r)
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

// 用户模型
type User struct {
	Id         uint                `gorm:"column:id; primaryKey; type:int(11); not null"      json:"id"`
	Name       string              `gorm:"column:name; type:varchar(50); not null"            json:"name"`
	Account    string              `gorm:"column:account; type:varchar(50); not null"         json:"account"`
	Mail       string              `gorm:"column:mail; type:varchar(50);"                     json:"mail"`
	Phone      string              `gorm:"column:phone; type:varchar(11);"                    json:"phone"`
	Role       Role                `gorm:"ForeignKey:RoleId;AssociationForeignKey:Id"`
	RoleId     uint                `gorm:"column:role_id; type:int(11); not null"             json:"role_id"`
	Status     contants.UserStatus `gorm:"column:status; type:int; not null; default:1"       json:"status"`
	UsageCount int                 `gorm:"column:usage_count; type:int; default:0"            json:"usage_count"`
	CreatedAt  util.Time           `gorm:"column:created_time; type:datetime; not null"       json:"created_time"`
	UpdatedAt  util.Time           `gorm:"column:updated_time; type:datetime; not null"       json:"updated_time"`
}

func (u *User) TableName() string {
	return "tangula_user"
}

// 类型转换方法
func (u *User) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)
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
