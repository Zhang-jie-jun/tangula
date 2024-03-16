package resource

import (
	"context"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/host"
	"github.com/Zhang-jie-jun/tangula/internal/dao/platform"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/Zhang-jie-jun/tangula/service/app"
	"github.com/Zhang-jie-jun/tangula/service/app/vmware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func CreateHost(createParam *view.CreateResource, user *auth.User) (result map[string]interface{}, err error) {
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_HOST_FAILED, createParam.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_HOST, createParam.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.CREATE_HOST_SUCCESS, createParam.Name)
			service.CreateLogRecord(msg.CREATE_HOST, createParam.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	if host.HostMgm.CheckIsExistByIp(createParam.Ip) {
		err = errors.New(msg.ERROR_HOST_IP_IS_EXIST, msg.GetMsg(msg.ERROR_HOST_IP_IS_EXIST))
		return
	}
	if host.HostMgm.CheckIsExistByNameAndCreateUser(createParam.Name, user.Account) {
		err = errors.New(msg.ERROR_HOST_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_HOST_NAME_IS_EXIST))
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
		err = errors.New(msg.ERROR_CREATE_HOST, msg.GetMsg(msg.ERROR_CREATE_HOST, err.Error()))
		logrus.Error(err)
		return
	}
	// 密码入库加密
	password, err := util.AesEncrypt(createParam.PassWord)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_HOST, msg.GetMsg(msg.ERROR_CREATE_HOST, err.Error()))
		logrus.Error(err)
		return
	}

	var obj host.Host
	obj.Name = createParam.Name
	obj.Desc = createParam.Desc
	obj.HostName = response.HostName
	obj.Type = createParam.Type
	obj.Status = 1
	obj.Os = response.OSType
	obj.Arch = response.Arch
	obj.Ip = createParam.Ip
	obj.Port = createParam.Port
	obj.UserName = createParam.UserName
	obj.PassWord = password
	obj.CreateUser = user.Account
	obj.AuthType = contants.PRIVATE
	obj, err = host.HostMgm.CreateHost(obj)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_HOST, msg.GetMsg(msg.ERROR_CREATE_HOST, err.Error()))
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func GetHostList(queryParam *view.QueryResourceParam, user *auth.User) (totalNum int64, result []map[string]interface{}, err error) {
	var objs []host.Host
	if queryParam.Auth == contants.PRIVATE {

		totalNum, objs, err = host.HostMgm.
			GetPrivateHostList(queryParam.Index, queryParam.Count, queryParam.Type, user.Account, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
			logrus.Error(err)
			return
		}

	} else if queryParam.Auth == contants.PUBLIC {
		totalNum, objs, err = host.HostMgm.
			GetPublicHostList(queryParam.Index, queryParam.Count, queryParam.Type, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
			logrus.Error(err)
			return
		}
	} else {
		totalNum, objs, err = host.HostMgm.
			GetHostList(queryParam.Index, queryParam.Count, queryParam.Type, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
			logrus.Error(err)
			return
		}
		// 所有用户只能获取公共主机与自身创建的私有主机资源
		for i := 0; i < len(objs); {
			if e := service.CheckResource(service.QUERY_RESOURCE, user, objs[i].AuthType, objs[i].CreateUser); e != nil {
				objs = append(objs[:i], objs[i+1:]...)
			} else {
				i++
			}
		}
	}

	for index, item := range objs {
		//检查最后一次部署状态
		lastRes, lastErr := host.HostMgm.FindLastDeployByHost(item.Id)
		if lastErr != nil {
			logrus.Error("查询部署记录异常")
			objs[index].DeployStatus = 1
		} else {
			if lastRes == (host.DeployApp{}) { //没有部署记录
				objs[index].DeployStatus = 1
			} else {
				objs[index].DeployStatus = lastRes.Status
			}
		}

	}
	fmt.Println(objs)

	for _, iter := range objs {
		result = append(result, iter.TransformMap())
	}

	return
}

func GetHostById(id uint, user *auth.User) (result map[string]interface{}, err error) {
	obj, err := host.HostMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
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

func EditHost(id uint, updateParam *view.UpdateResource, user *auth.User) (result map[string]interface{}, err error) {
	var obj host.Host
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.EDIT_HOST_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.EDIT_HOST, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.EDIT_HOST_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.EDIT_HOST, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	if host.HostMgm.CheckIsExistByIp(updateParam.Ip) {
		err = errors.New(msg.ERROR_HOST_IP_IS_EXIST, msg.GetMsg(msg.ERROR_HOST_IP_IS_EXIST))
		return
	}
	obj, err = host.HostMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
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
		err = errors.New(msg.ERROR_UPDATE_HOST, msg.GetMsg(msg.ERROR_UPDATE_HOST, err.Error()))
		logrus.Error(err)
		return
	}
	// 密码入库加密
	password, err := util.AesEncrypt(updateParam.PassWord)
	if err != nil {
		err = errors.New(msg.ERROR_UPDATE_HOST, msg.GetMsg(msg.ERROR_UPDATE_HOST, err.Error()))
		logrus.Error(err)
		return
	}

	obj.Desc = updateParam.Desc
	obj.HostName = response.HostName
	obj.Os = response.OSType
	obj.Arch = response.Arch
	obj.Ip = updateParam.Ip
	obj.Port = updateParam.Port
	obj.UserName = updateParam.UserName
	obj.PassWord = password
	obj, err = host.HostMgm.UpdateHost(obj)
	if err != nil {
		err = errors.New(msg.ERROR_UPDATE_HOST, msg.GetMsg(msg.ERROR_UPDATE_HOST, err.Error()))
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func DeleteHost(id uint, user *auth.User) (err error) {
	var obj host.Host
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DELETE_HOST_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.DELETE_HOST, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DELETE_HOST_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.DELETE_HOST, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = host.HostMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
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
		err = errors.New(msg.ERROR_DELETE_HOST, msg.GetMsg(msg.ERROR_DELETE_HOST, err.Error()))
		logrus.Error(err)
		return
	}
	if len(mountInfos) != 0 {
		err = errors.New(msg.ERROR_DELETE_HOST,
			msg.GetMsg(msg.ERROR_DELETE_HOST, msg.GetMsg(msg.ERROR_PLATFORM_EXIST_MOUNT_STATUS_REPLICA)))
		logrus.Error(err)
		return
	}
	err = host.HostMgm.DeleteHost(id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_HOST, msg.GetMsg(msg.ERROR_DELETE_HOST, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func PublishHost(id uint, user *auth.User) (result map[string]interface{}, err error) {
	var obj host.Host
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.PUBLISH_HOST_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.PUBLISH_HOST, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.PUBLISH_HOST_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.PUBLISH_HOST, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = host.HostMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.PUBLISH_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	obj.AuthType = contants.PUBLIC
	obj, err = host.HostMgm.UpdateHost(obj)
	if err != nil {
		err = errors.New(msg.ERROR_PUBLISH_HOST, msg.GetMsg(msg.ERROR_PUBLISH_HOST, err.Error()))
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func UpdHost(id uint, user *auth.User) (err error) {
	var obj host.Host
	// 闭包的方式记录操作日志
	defer func() {
		_, updErr := host.HostMgm.UpdateHost(obj)
		if updErr != nil {
			err = errors.New(msg.ERROR_UPD_HOST_INFO, msg.GetMsg(msg.ERROR_UPD_HOST_INFO, updErr.Error()))
			logrus.Error("================更新主机信息发生错误================", updErr)
		}
		if err != nil {
			detail := msg.GetOperation(msg.UPD_HOST_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.UPD_HOST, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.UPD_HOST_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.UPD_HOST, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = host.HostMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info("开始更新主机信息")
	// 校验每台机器
	realPassWord, _ := util.AesDecrypt(obj.PassWord)
	loginInfo := contants.LoginInfo{
		Ip:       obj.Ip,
		Port:     obj.Port,
		UserName: obj.UserName,
		PassWord: realPassWord,
	}

	adapter := app.NewAdapter(obj.Type, &loginInfo)
	res, verifyErr := adapter.VerifyConfigInfo()
	if verifyErr != nil {
		logrus.Error("================校验主机，获取不到信息================", verifyErr.Error())
		err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, verifyErr.Error()))
		obj.Status = 2

	} else {
		//更新主机状态为已连通
		obj.Status = 1
		obj.Os = res.OSType
		obj.Arch = res.Arch

		if obj.PlatformId != 0 {
			platformRes, findPlatformErr := platform.PlatformMgm.FindById(obj.PlatformId)
			if findPlatformErr != nil {
				logrus.Error(fmt.Sprintf("查询虚拟化平台信息报错:%s", findPlatformErr))
			} else {

				realPlatformPassword, _ := util.AesDecrypt(platformRes.PassWord)
				ctx := context.Background()
				vms := vmware.NewVmWare(platformRes.Ip, platformRes.UserName, realPlatformPassword)
				machineName, findVmByIpErr := vms.FindVMByIP(ctx, vms.Client.Client, obj.Ip)
				if findVmByIpErr != nil {
					logrus.Error(fmt.Sprintf("根据ip查询虚拟机报错:%s", findVmByIpErr))
				} else {
					obj.Name = machineName
				}
			}
		}
	}
	return
}

func DeployApps(hostId uint, deployParam *view.DeployParam, user *auth.User, context *gin.Context) (result map[string]interface{}, err error) {
	var obj host.Host
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DEPLOY_CLIENT_FAILED, obj.Ip, err.Error())
			service.CreateLogRecord(msg.DEPLOY_CLIENT, obj.Ip, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DEPLOY_CLIENT_START, obj.Ip)
			service.CreateLogRecord(msg.DEPLOY_CLIENT, obj.Ip, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()

	obj, err = host.HostMgm.FindById(hostId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
		logrus.Error(err)
		return
	}

	//检查最后一次部署状态
	lastRes, lastErr := host.HostMgm.FindLastDeployByHost(hostId)
	if lastErr != nil {
		logrus.Error("查询部署记录异常")
	} else {
		if lastRes != (host.DeployApp{}) {
			if lastRes.Status == 2 {
				err = errors.New(msg.ERROR_HOST_DEPLOYING, msg.GetMsg(msg.ERROR_HOST_DEPLOYING, obj.Ip))
				logrus.Error(err)
				return
			}
		}
	}

	//先创建记录，防止丢失记录
	var deployApp host.DeployApp
	deployApp.HostId = hostId
	deployApp.Status = 2
	deployApp.Log = "================开始部署客户端================"
	deployApp.ClientIp = obj.Ip
	deployApp.BaseDir = deployParam.Basedir
	deployApp.ServerIp = deployParam.ServerIp
	deployApp.AbTpye = deployParam.Type
	deployApp.Apps = deployParam.Apps
	deployApp.Create_user = user.Account
	crtRes, crtDeployAppErr := host.HostMgm.CreateDeployApp(deployApp)
	if crtDeployAppErr != nil {
		err = errors.New(msg.ERROR_UPD_HOST_INFO, msg.GetMsg(msg.ERROR_UPD_HOST_INFO, crtDeployAppErr.Error()))
		logrus.Error("================创建部署记录发生错误================", crtDeployAppErr)
		return
	}

	basicFtpPath := "defalut"
	if deployParam.BasicFtpPath != "" {
		basicFtpPath = deployParam.BasicFtpPath
	}

	appFtpPath := "default"
	if deployParam.AppFtpPath != "" {
		basicFtpPath = deployParam.AppFtpPath
	}

	realPassWord, _ := util.AesDecrypt(obj.PassWord)
	paramMap := map[string]string{
		"baseDir":      deployParam.Basedir,
		"serverIp":     deployParam.ServerIp,
		"clientIp":     obj.Ip,
		"user":         obj.UserName,
		"userpass":     realPassWord,
		"apps":         deployParam.Apps,
		"abType":       deployParam.Type,
		"deployId":     cast.ToString(crtRes.Id),
		"basicFtpPath": basicFtpPath,
		"appFtpPath":   appFtpPath,
	}

	logrus.Info(paramMap)
	buildId, buildErr := util.BuildJenkins(context, contants.AppCfg.JENKINS.INSTALL_CLIENT_NAME, paramMap)
	if buildErr != nil {
		err = errors.New(msg.ERROR_HOST_START_JENKINSJOB, msg.GetMsg(msg.ERROR_HOST_START_JENKINSJOB, buildErr.Error()))
		logrus.Error(buildErr)
		//更新部署状态为失败
		crtRes.Status = 4
		crtRes.Log = cast.ToString(buildErr)
		_, updErr := host.HostMgm.UpdateDeployApp(crtRes)
		if updErr != nil {
			err = errors.New(msg.ERROR_UPD_DEPLOYAPP, msg.GetMsg(msg.ERROR_UPD_DEPLOYAPP, updErr.Error()))
			logrus.Error("================更新部署信息发生错误================", updErr)
			return
		}
		return
	}

	res := map[string]interface{}{
		"buildId": buildId,
	}
	return res, err

}

func UpdDeploy(id uint, deployStatus uint, log string) (err error) {
	//检查最后一次部署状态
	lastRes, lastErr := host.HostMgm.FindLastDeployById(id)
	if lastErr != nil {
		err = errors.New(msg.ERROR_GET_DEPLOYAPP, msg.GetMsg(msg.ERROR_GET_DEPLOYAPP, lastRes.ClientIp))
		logrus.Error(err)
		return
	}

	currentLog := lastRes.Log
	currentLog += "\r\n"
	currentLog += log
	lastRes.Status = deployStatus
	lastRes.Log = currentLog
	_, updErr := host.HostMgm.UpdateDeployApp(lastRes)
	if updErr != nil {
		err = errors.New(msg.ERROR_UPD_DEPLOYAPP, msg.GetMsg(msg.ERROR_UPD_DEPLOYAPP, updErr.Error()))
		logrus.Error("================更新部署信息发生错误================", updErr)
		return
	}

	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DEPLOY_CLIENT_FAILED, id, err.Error())
			service.CreateLogRecord(msg.DEPLOY_CLIENT, cast.ToString(id), detail, "system", contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DEPLOY_CLIENT_SUCCESS, id)
			service.CreateLogRecord(msg.DEPLOY_CLIENT, cast.ToString(id), detail, "system", contants.LOG_SUCCESS)
		}
	}()
	return
}

func GetDeployRecord(id uint, queryParam *view.PageQueryParam) (int64, []map[string]interface{}, error) {
	totalNum, res, err := host.HostMgm.GetDeployAppByHostId(id, queryParam.Index, queryParam.Count)
	if err != nil {
		err = errors.New(msg.ERROR_GET_DEPLOYAPP, msg.GetMsg(msg.ERROR_GET_DEPLOYAPP, err.Error()))
		logrus.Error(err)
		return 0, nil, err
	}
	var result []map[string]interface{}
	for _, iter := range res {
		instanceMap := iter.TransformMap()
		result = append(result, instanceMap)
	}
	return totalNum, result, nil
}
