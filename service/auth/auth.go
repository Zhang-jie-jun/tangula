package auth

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUserInfo(context *gin.Context) (user *auth.User, err error) {
	claims := jwt.ExtractClaims(context)
	mail, ok := claims["mail"].(string)
	if !ok {
		err = errors.New(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, "mail is nil!"))
		return
	}
	var userInfo auth.User
	SuperAdminMail := contants.AppCfg.System.SuperAdminName
	ATMail := contants.AppCfg.System.ATName
	if mail == SuperAdminMail {
		role := auth.Role{Name: contants.SUPERADMIN, Desc: string(contants.SUPERADMIN)}
		userInfo = auth.User{
			Name:    fmt.Sprintf("超级管理员(%s)", contants.AppCfg.System.SuperAdminName),
			Account: contants.AppCfg.System.SuperAdminName,
			Mail:    SuperAdminMail,
			Phone:   "",
			Role:    role,
			Status:  contants.ENABLE,
		}
	} else if mail == ATMail {
		role := auth.Role{Name: contants.SUPERADMIN, Desc: string(contants.SUPERADMIN)}
		userInfo = auth.User{
			Name:    fmt.Sprintf("自动化测试(%s)", contants.AppCfg.System.ATName),
			Account: contants.AppCfg.System.ATName,
			Mail:    ATMail,
			Phone:   "",
			Role:    role,
			Status:  contants.ENABLE,
		}
	} else {
		userInfo, err = auth.AuthMgm.FindByMail(mail)
		if err != nil {
			err = errors.New(msg.ERROR_GET_USER_INFO, msg.GetMsg(msg.ERROR_GET_USER_INFO, err.Error()))
			logrus.Error(err)
			return nil, err
		}
	}
	return &userInfo, nil
}

func GetAssignUserInfo(id uint) (user auth.User, err error) {
	user, err = auth.AuthMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_USER_NOT_EXIST, msg.GetMsg(msg.ERROR_USER_NOT_EXIST))
		logrus.Error(err)
		return
	}
	return
}

func GetUserList(queryParam *view.UserQueryParam) (totalNum int64, users []auth.User, err error) {
	totalNum, users, err = auth.AuthMgm.GetUserList(queryParam.Index, queryParam.Count, queryParam.RoleName, queryParam.Filter)
	if err != nil {
		err = errors.New(msg.ERROR_GET_USER_INFO, msg.GetMsg(msg.ERROR_GET_USER_INFO, err.Error()))
		logrus.Error(err)
		return totalNum, users, err
	}
	return totalNum, users, err
}

func SetUserRole(id uint, roleType contants.RoleType, userInfo *auth.User) (user auth.User, err error) {
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.EDIT_USER_ROLE_FAILED, roleType, err.Error())
			service.CreateLogRecord(msg.EDIT_USER_ROLE, user.Account, detail, userInfo.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.EDIT_USER_ROLE_SUCCESS, roleType)
			service.CreateLogRecord(msg.EDIT_USER_ROLE, user.Account, detail, userInfo.Account, contants.LOG_SUCCESS)
		}
	}()
	if roleType == contants.SUPERADMIN {
		err = errors.New(msg.ERROR_SET_USER_ROLE, msg.GetMsg(msg.ERROR_SET_USER_ROLE,
			"The user role cannot be set to super administrator."))
		logrus.Error(err)
		return
	}
	user, err = auth.AuthMgm.SetUserRole(id, roleType)
	if err != nil {
		err = errors.New(msg.ERROR_SET_USER_ROLE, msg.GetMsg(msg.ERROR_SET_USER_ROLE, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func GetRoleList(queryParam *view.UserQueryParam) (totalNum int64, role []auth.Role, err error) {
	totalNum, role, err = auth.AuthMgm.GetRoleList(queryParam.Index, queryParam.Count, queryParam.Filter)
	if err != nil {
		err = errors.New(msg.ERROR_GET_ROLE_INFO, msg.GetMsg(msg.ERROR_GET_ROLE_INFO, err.Error()))
		logrus.Error(err)
		return totalNum, role, err
	}
	return totalNum, role, err
}

func EnableUser(id uint, userInfo *auth.User) (user auth.User, err error) {
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.ENABLE_USER_FAILED, user.Account, err.Error())
			service.CreateLogRecord(msg.ENABLE_USER, user.Account, detail, userInfo.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.ENABLE_USER_SUCCESS, user.Account)
			service.CreateLogRecord(msg.ENABLE_USER, user.Account, detail, userInfo.Account, contants.LOG_SUCCESS)
		}
	}()
	user, err = auth.AuthMgm.EnableUser(id)
	if err != nil {
		err = errors.New(msg.ERROR_ENABLE_USER, msg.GetMsg(msg.ERROR_ENABLE_USER, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func DisableUser(id uint, userInfo *auth.User) (user auth.User, err error) {
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DISABLE_USER_FAILED, user.Account, err.Error())
			service.CreateLogRecord(msg.DISABLE_USER, user.Account, detail, userInfo.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DISABLE_USER_SUCCESS, user.Account)
			service.CreateLogRecord(msg.DISABLE_USER, user.Account, detail, userInfo.Account, contants.LOG_SUCCESS)
		}
	}()
	user, err = auth.AuthMgm.DisableUser(id)
	if err != nil {
		err = errors.New(msg.ERROR_DISABLE_USER, msg.GetMsg(msg.ERROR_DISABLE_USER, err.Error()))
		logrus.Error(err)
		return
	}
	return
}
