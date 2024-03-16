package resource

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/platform"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/Zhang-jie-jun/tangula/service/app"
	"github.com/sirupsen/logrus"
)

func CreatePlatform(createParam *view.CreateResource, user *auth.User) (result map[string]interface{}, err error) {
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_PLATFORM_FAILED, createParam.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_PLATFORM, createParam.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.CREATE_PLATFORM_SUCCESS, createParam.Name)
			service.CreateLogRecord(msg.CREATE_PLATFORM, createParam.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	if platform.PlatformMgm.CheckIsExistByIp(createParam.Ip) {
		err = errors.New(msg.ERROR_PLATFORM_IP_IS_EXIST, msg.GetMsg(msg.ERROR_PLATFORM_IP_IS_EXIST))
		return
	}
	if platform.PlatformMgm.CheckIsExistByNameAndCreateUser(createParam.Name, user.Account) {
		err = errors.New(msg.ERROR_PLATFORM_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_PLATFORM_NAME_IS_EXIST))
		return
	}
	// 应用认证处理
	var login contants.LoginInfo
	login.Ip = createParam.Ip
	login.Port = createParam.Port
	login.UserName = createParam.UserName
	login.PassWord = createParam.PassWord
	adapter := app.NewAdapter(createParam.Type, &login)
	response, err := adapter.VerifyConfigInfo()
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_PLATFORM, msg.GetMsg(msg.ERROR_CREATE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	// 密码入库加密
	password, err := util.AesEncrypt(createParam.PassWord)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_PLATFORM, msg.GetMsg(msg.ERROR_CREATE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	var obj platform.Platform
	obj.Name = createParam.Name
	obj.Desc = createParam.Desc
	obj.Type = createParam.Type
	obj.Ip = createParam.Ip
	obj.Port = createParam.Port
	obj.UserName = createParam.UserName
	obj.PassWord = password
	obj.Version = response.Version
	obj.CreateUser = user.Account
	obj.AuthType = contants.PRIVATE
	obj, err = platform.PlatformMgm.CreatePlatform(obj)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_PLATFORM, msg.GetMsg(msg.ERROR_CREATE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func GetPlatform(queryParam *view.QueryResourceParam, user *auth.User) (totalNum int64,
	result []map[string]interface{}, err error) {
	var platforms []platform.Platform
	if queryParam.Auth == contants.PRIVATE {
		totalNum, platforms, err = platform.PlatformMgm.
			GetPrivatePlatformList(queryParam.Index, queryParam.Count, queryParam.Type, user.Account, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
			logrus.Error(err)
			return
		}
	} else if queryParam.Auth == contants.PUBLIC {
		totalNum, platforms, err = platform.PlatformMgm.
			GetPublicPlatformList(queryParam.Index, queryParam.Count, queryParam.Type, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
			logrus.Error(err)
			return
		}
	} else {
		totalNum, platforms, err = platform.PlatformMgm.
			GetPlatformList(queryParam.Index, queryParam.Count, queryParam.Type, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
			logrus.Error(err)
			return
		}
		// 所有用户只能获取公共平台与自身创建的私有平台资源
		for i := 0; i < len(platforms); {
			if e := service.CheckResource(service.QUERY_RESOURCE, user, platforms[i].AuthType, platforms[i].CreateUser); e != nil {
				platforms = append(platforms[:i], platforms[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, iter := range platforms {
		result = append(result, iter.TransformMap())
	}
	return
}

func GetPlatformById(id uint, user *auth.User) (result map[string]interface{}, err error) {
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func EditPlatform(id uint, updateParam *view.UpdateResource, user *auth.User) (result map[string]interface{}, err error) {
	var obj platform.Platform
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.EDIT_PLATFORM_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.EDIT_PLATFORM, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.EDIT_PLATFORM_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.EDIT_PLATFORM, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	if platform.PlatformMgm.CheckIsExistByIp(updateParam.Ip) {
		err = errors.New(msg.ERROR_PLATFORM_IP_IS_EXIST, msg.GetMsg(msg.ERROR_PLATFORM_IP_IS_EXIST))
		return
	}
	obj, err = platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.EDIT_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 应用认证处理
	var login contants.LoginInfo
	login.Ip = updateParam.Ip
	login.Port = updateParam.Port
	login.UserName = updateParam.UserName
	login.PassWord = updateParam.PassWord
	adapter := app.NewAdapter(obj.Type, &login)
	response, err := adapter.VerifyConfigInfo()
	if err != nil {
		err = errors.New(msg.ERROR_UPDATE_PLATFORM, msg.GetMsg(msg.ERROR_UPDATE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	// 密码入库加密
	password, err := util.AesEncrypt(updateParam.PassWord)
	if err != nil {
		err = errors.New(msg.ERROR_UPDATE_PLATFORM, msg.GetMsg(msg.ERROR_UPDATE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}

	obj.Desc = updateParam.Desc
	obj.Ip = updateParam.Ip
	obj.Port = updateParam.Port
	obj.UserName = updateParam.UserName
	obj.PassWord = password
	obj.Version = response.Version
	obj, err = platform.PlatformMgm.UpdatePlatform(obj)
	if err != nil {
		err = errors.New(msg.ERROR_UPDATE_PLATFORM, msg.GetMsg(msg.ERROR_UPDATE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func DeletePlatform(id uint, user *auth.User) (err error) {
	var obj platform.Platform
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DELETE_PLATFORM_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.DELETE_PLATFORM, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DELETE_PLATFORM_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.DELETE_PLATFORM, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 检查资源权限
	err = service.CheckResource(service.DELETE_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 检查平台上是否存在挂载的副本
	mountInfos, err := replica.MountInfoMgm.GetMountInfoByTarget(obj.Type, obj.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_PLATFORM, msg.GetMsg(msg.ERROR_DELETE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	if len(mountInfos) != 0 {
		err = errors.New(msg.ERROR_DELETE_PLATFORM,
			msg.GetMsg(msg.ERROR_DELETE_PLATFORM, msg.GetMsg(msg.ERROR_PLATFORM_EXIST_MOUNT_STATUS_REPLICA)))
		logrus.Error(err)
		return
	}
	err = platform.PlatformMgm.DeletePlatform(id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_PLATFORM, msg.GetMsg(msg.ERROR_DELETE_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func PublishPlatform(id uint, user *auth.User) (result map[string]interface{}, err error) {
	var obj platform.Platform
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.PUBLISH_PLATFORM_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.PUBLISH_PLATFORM, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.PUBLISH_PLATFORM_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.PUBLISH_PLATFORM, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.PUBLISH_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	obj.AuthType = contants.PUBLIC
	obj, err = platform.PlatformMgm.UpdatePlatform(obj)
	if err != nil {
		err = errors.New(msg.ERROR_PUBLISH_PLATFORM, msg.GetMsg(msg.ERROR_PUBLISH_PLATFORM, err.Error()))
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}
