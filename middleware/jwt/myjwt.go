package jwt

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/gin-gonic/gin"
)

// JWT Authorizator授权规则接口
type Authorizator interface {
	HandleAuthorizator(data interface{}, c *gin.Context) bool
}

type SuperAdminAuthorizator struct {
}

// 超级管理员实现授权规则接口
func (*SuperAdminAuthorizator) HandleAuthorizator(data interface{}, c *gin.Context) bool {
	user, ok := data.(*auth.User)
	Role, err := auth.AuthMgm.FindRoleById(user.RoleId)
	if err != nil {
		return false
	}
	if ok {
		if Role.Name == contants.SUPERADMIN {
			return true
		}
	}
	return false
}

type AdminAuthorizator struct {
}

// 系统管理员实现授权规则接口
func (*AdminAuthorizator) HandleAuthorizator(data interface{}, c *gin.Context) bool {
	if user, ok := data.(*auth.User); ok {
		Role, err := auth.AuthMgm.FindRoleById(user.RoleId)
		if err != nil {
			return false
		}
		if Role.Name == contants.SUPERADMIN || Role.Name == contants.SYSADMIN {
			return true
		}
	}
	return false
}

type UserAuthorizator struct {
}

// 普通用户实现授权规则接口
func (*UserAuthorizator) HandleAuthorizator(data interface{}, c *gin.Context) bool {
	return true
}
