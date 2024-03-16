package service

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/service/logmgm"
	"github.com/sirupsen/logrus"
)

// 资源操作方式
type Operation int

const (
	QUERY_RESOURCE   Operation = 1 // 查询
	EDIT_RESOURCE    Operation = 2 // 编辑
	DELETE_RESOURCE  Operation = 3 // 删除
	PUBLISH_RESOURCE Operation = 4 // 发布

)

// 资源权限检查
func CheckResource(operation Operation, user *auth.User, authType contants.AuthType, resourceUser string) error {
	switch operation {
	case QUERY_RESOURCE:
		// 所有用户不可查看不属于自身的私有资源
		if user.Account != resourceUser && authType == contants.PRIVATE {
			err := errors.New(msg.ERROR_QUERY_RESOURCE_AUTH,
				msg.GetMsg(msg.ERROR_QUERY_RESOURCE_AUTH, msg.GetMsg(msg.AUTH_PERMISSION_DENIED)))
			return err
		}
	case EDIT_RESOURCE:
		// 普通用户只能编辑属于自身的私有资源
		// 系统管理员可以编辑公共资源及属于自身的私有资源
		if user.Role.Name == contants.USER {
			if user.Account != resourceUser || authType == contants.PUBLIC {
				err := errors.New(msg.ERROR_EDIT_RESOURCE_AUTH,
					msg.GetMsg(msg.ERROR_EDIT_RESOURCE_AUTH, msg.GetMsg(msg.AUTH_PERMISSION_DENIED)))
				return err
			}
		} else if user.Role.Name == contants.SYSADMIN {
			if user.Account != resourceUser && authType == contants.PRIVATE {
				err := errors.New(msg.ERROR_EDIT_RESOURCE_AUTH,
					msg.GetMsg(msg.ERROR_EDIT_RESOURCE_AUTH, msg.GetMsg(msg.AUTH_PERMISSION_DENIED)))
				return err
			}
		}
	case DELETE_RESOURCE:
		// 普通用户只能删除属于自身的私有资源
		// 系统管理员可以删除公共资源及属于自身的私有资源
		if user.Role.Name == contants.USER {
			if user.Account != resourceUser || authType == contants.PUBLIC {
				err := errors.New(msg.ERROR_DELETE_RESOURCE_AUTH,
					msg.GetMsg(msg.ERROR_DELETE_RESOURCE_AUTH, msg.GetMsg(msg.AUTH_PERMISSION_DENIED)))
				return err
			}
		} else if user.Role.Name == contants.SYSADMIN {
			if user.Account != resourceUser && authType == contants.PRIVATE {
				err := errors.New(msg.ERROR_DELETE_RESOURCE_AUTH,
					msg.GetMsg(msg.ERROR_DELETE_RESOURCE_AUTH, msg.GetMsg(msg.AUTH_PERMISSION_DENIED)))
				return err
			}
		}
	case PUBLISH_RESOURCE:
		// 只能发布属于用户自身的私有资源
		if authType != contants.PRIVATE || user.Account != resourceUser {
			err := errors.New(msg.ERROR_PUBLISH_RESOURCE_AUTH,
				msg.GetMsg(msg.ERROR_PUBLISH_RESOURCE_AUTH, msg.GetMsg(msg.AUTH_PERMISSION_DENIED)))
			return err
		}
	default:
		err := errors.New(msg.ERROR_PUBLISH_RESOURCE_AUTH,
			msg.GetMsg(msg.ERROR_PUBLISH_RESOURCE_AUTH, msg.GetMsg(msg.AUTH_PERMISSION_DENIED)))
		return err
	}
	return nil
}

func CheckType(imageType contants.ImageType, targetType contants.AppType) bool {
	var res bool
	logrus.Info("imageType", imageType)
	switch imageType {
	case contants.VMWAREVM:
		if targetType == contants.VMWARE {
			res = true
		} else {
			res = false
		}
	case contants.CUSTOME:
		if targetType == contants.LINUX || targetType == contants.WINDOWS {
			res = true
		} else {
			res = false
		}
	case contants.CASVM:
		if targetType == contants.CAS {
			res = true
		} else {
			res = false
		}
	case contants.FUSIONCOMPUTEVM:
		if targetType == contants.FUSIONCOMPUTE {
			res = true
		} else {
			res = false
		}
	default:
		res = false
	}
	return res
}

// 创建日志记录
func CreateLogRecord(opt msg.Operation, object, detail, userName string, status contants.LogStatus) {
	operation := msg.GetOperation(opt)
	record, err := logmgm.CreateLogRecord(operation, object, detail, userName, status)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(record)
}
