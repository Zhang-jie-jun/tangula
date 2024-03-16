package cas

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/instance"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"
)

type Adapter struct {
	LoginInfo *contants.LoginInfo
}

func NewAdapter(login *contants.LoginInfo) (adapter *Adapter) {
	return &Adapter{LoginInfo: login}
}

func (a *Adapter) VerifyConfigInfo() (*contants.VerifyResult, error) {
	//err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	var res contants.VerifyResult
	var err error
	if a.LoginInfo == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
	}
	loginRes, err := a.login()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	res.Version = loginRes.Version
	res.HostpoolInfo = loginRes.HostpoolInfo
	return &res, err
}

func setParam(sharePath *string, mountParam *view.MountInfo) (*MountDetailInfo, error) {
	var detailInfo MountDetailInfo
	detailInfo.HostId = cast.ToString(mountParam.AppConfig["hostId"])
	detailInfo.HostName = cast.ToString(mountParam.AppConfig["hostName"])
	detailInfo.VsId = cast.ToUint(mountParam.AppConfig["vsId"])
	detailInfo.HostpoolId = cast.ToUint(mountParam.AppConfig["hostPoolId"])
	detailInfo.ClusterId = cast.ToUint(mountParam.AppConfig["clusterId"])
	detailInfo.StorePath = cast.ToString(mountParam.AppConfig["storePath"])
	detailInfo.IsRegisterVM = cast.ToBool(mountParam.AppConfig["isRegisterVM"])
	detailInfo.IsCrtByJson = cast.ToBool(mountParam.AppConfig["isCrtCasByJson"])
	detailInfo.Title = cast.ToString(mountParam.AppConfig["title"])
	detailInfo.PoolName = cast.ToString(mountParam.AppConfig["poolName"])
	start1 := strings.LastIndex(*sharePath, ":")
	remotePath := string([]byte(*sharePath)[start1+1:])
	detailInfo.RemoteHost = string([]byte(*sharePath)[0:start1])  // 截取IP
	detailInfo.RemotePath = string([]byte(*sharePath)[start1+1:]) // 截取路径：/tangula/mnt/pool_uuid/image_uuid
	start2 := strings.LastIndex(remotePath, "/")
	if mountParam.AppConfig["storeName"] != "" {
		detailInfo.StoreName = cast.ToString(mountParam.AppConfig["storeName"])
	} else {
		detailInfo.StoreName = fmt.Sprintf("%s", string([]byte(remotePath)[start2+1:])) // 截取image_uuid作为CAS存储池名称
	}

	return &detailInfo, nil

}

func (a *Adapter) Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error) {

	var err error
	// 设置挂载详细参数
	param, err := setParam(sharePath, mountParam)

	if err != nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, err.Error()))
		logrus.Error(err)
		return "", err
	}
	if param == nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER,
			msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return "", err
	}

	//登录校验
	if a.LoginInfo == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return "", err
	}

	_, err = a.login()
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	time.Sleep(time.Second * 5)
	// 创建存储池
	reqData := CrtStoragePoolParams(param)
	instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_CREATE_STORAGE), "")
	reqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/storage/add", a.LoginInfo.Ip)
	storageRes, storageErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "POST", reqUrl, reqData)
	if storageErr != nil {
		err = errors.New(msg.ERRIR_CREATE_POOL, msg.GetMsg(msg.ERRIR_CREATE_POOL, param.StoreName, storageErr.Error()))
		return "", err
	}
	if cast.ToInt(storageRes.StatusMap["errorCode"]) != 0 {
		err = errors.New(msg.ERRIR_CREATE_POOL, msg.GetMsg(msg.ERRIR_CREATE_POOL, param.StoreName, cast.ToString(storageRes.StatusMap["message"])))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERRIR_CREATE_POOL), err.Error())
		return "", err
	}
	time.Sleep(time.Second * 3)

	//启动存储池
	logrus.Info("启动存储池=========>", param.StoreName)
	reqUrlStart := fmt.Sprintf("http://%s:8080/cas/casrs/storage/start?id=%d&poolName=%s&hostName=%s", a.LoginInfo.Ip, cast.ToUint(param.HostId), param.StoreName, param.HostName)
	_, startErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrlStart, nil)
	if startErr != nil {
		err = errors.New(msg.ERRIR_START_POOL, msg.GetMsg(msg.ERRIR_START_POOL, param.StoreName, err.Error()))
		return "", err
	}
	time.Sleep(time.Second * 3)

	//先查询存储池,如果没启动尝试重启
	var count = 0
	for {
		if count > 3 {
			break
		}
		reqUrlQueryStorage := fmt.Sprintf("http://%s:8080/cas/casrs/storage/info?hostId=%d&poolName=%s", a.LoginInfo.Ip, cast.ToUint(param.HostId), param.StoreName)
		resQueryStorage, storageErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrlQueryStorage, nil)
		if storageErr != nil {
			logrus.Error(storageErr, "查询不到存储池，跳过")
		} else {
			if cast.ToInt(resQueryStorage.StatusMap["status"]) == 1 { //已启动
				logrus.Info("存储池已启动=========>", param.StoreName)
				break
			} else {
				//启动存储池
				logrus.Info("尝试再次启动存储池=========>", param.StoreName)
				util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrlStart, nil)
				time.Sleep(time.Second * 3)
			}
		}
		count++
	}

	//有报错时，要卸载存储池
	defer func() {
		if err != nil {
			logrus.Info("开始回滚，卸载存储池:", mountParam.ReplicaId)
			var loginInfo contants.LoginInfo
			loginInfo.Ip = a.LoginInfo.Ip
			loginInfo.Port = a.LoginInfo.Port
			loginInfo.UserName = a.LoginInfo.UserName
			loginInfo.PassWord = a.LoginInfo.PassWord
			DelStorage(loginInfo, cast.ToInt(param.HostId), param.StoreName)
		}
	}()

	//创建虚拟机
	if param.IsRegisterVM {

		crtCasReqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/vm/add", a.LoginInfo.Ip)

		replicaObj, findErr := replica.ReplicaMgm.FindById(mountParam.ReplicaId)
		if findErr != nil {
			err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, "查询副本失败："+findErr.Error()))
			logrus.Error(err)
			return "", err
		}
		if param.IsCrtByJson { //通过副本内json文件创建虚拟机
			cacheDirName := contants.AppCfg.System.MountPath + "/" + replicaObj.Pool.Uuid + "/" + replicaObj.Uuid
			configFleName := fmt.Sprintf("%s/%s", cacheDirName, "cas.json")
			if !util.IsFileExists(configFleName) {
				err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, "无法获取脚本:", configFleName))

				logrus.Error(err)
				return "", err
			}
			//解析json文件
			jsonList, jsonErr := util.ReadCasJson(configFleName)
			if jsonErr != nil {
				err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, "读取json配置文件失败:"+jsonErr.Error()))
				logrus.Error(err)
				return "", err
			}

			for _, jsonMap := range jsonList {
				crtParam, crtParamErr := SetCasVmParamsByJson(param, jsonMap)
				if crtParamErr != nil {
					err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, crtParamErr.Error()))
					return "", err
				}
				_, crtErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "POST", crtCasReqUrl, crtParam)
				if crtErr != nil {
					err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, err.Error()))
					return "", err
				}
				logrus.Info("创建虚拟机=====================================>", jsonMap)
			}
		} else { //默认配置

			//解析磁盘文件
			diskPath := param.RemotePath
			cmdRes, cmdErr := util.GetMountDiskInfo(diskPath)
			if cmdErr != nil {
				err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, err.Error()))
				return "", err
			}
			var diskInfoList []view.DiskInfo
			for _, line := range cmdRes {
				if !strings.Contains(line, ".json") {
					logrus.Info(fmt.Sprintf("获取磁盘信息:%s", line))
					splitRet := strings.Split(line, " ")
					if len(splitRet) != 2 {
						logrus.Error("split line ret error.")
						continue
					}
					diskInfo := view.DiskInfo{}
					diskInfo.Name = splitRet[1]
					diskInfo.PoolName = replicaObj.Uuid
					diskInfo.SizeByte, _ = strconv.ParseInt(splitRet[0], 10, 64)
					diskInfoList = append(diskInfoList, diskInfo)
				}

			}
			logrus.Info("-----------diskInfoList----------", diskInfoList)

			vsId, vsErr := GetVswitchId(*a.LoginInfo, param.HostName)
			if vsErr != nil || vsId == 0 {
				err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, vsErr.Error()))
				return "", err
			}

			crtParam, crtParamErr := SetCasVmParams(replicaObj.Name, param.HostId, vsId, diskInfoList)
			if crtParamErr != nil {
				err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, crtParamErr.Error()))
				return "", err
			}
			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_CREATE_VM), "")

			crtRes, crtErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "POST", crtCasReqUrl, crtParam)
			if crtErr != nil {
				err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, crtErr.Error()))
				logrus.Error(err)
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERRIR_CREATE_CAS, replicaObj.Name, crtErr.Error()), cast.ToString(crtRes.StatusMap["message"]))
				return "", err
			}

			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_CREATE_VM_SUCCESS, replicaObj.Name), "")
		}
	}
	return "", err

}

func (a *Adapter) UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) (err error) {
	//err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	param, err := setParam(sharePath, mountParam)

	if err != nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, err.Error()))
		logrus.Error(err)
		return err
	}
	if param == nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return err
	}
	//登录校验
	if a.LoginInfo == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return err
	}

	_, err = a.login()
	if err != nil {
		logrus.Error(err)
		return err
	}

	var isDelete = true
	//根据配置文件查询虚拟机id
	replicaObj, findErr := replica.ReplicaMgm.FindById(mountParam.ReplicaId)
	if findErr != nil {
		logrus.Error(findErr)
		isDelete = false
	}

	cacheDirName := contants.AppCfg.System.MountPath + "/" + replicaObj.Pool.Uuid + "/" + replicaObj.Uuid
	configFleName := fmt.Sprintf("%s/%s", cacheDirName, "cas.json")
	if !util.IsFileExists(configFleName) {
		logrus.Error("无法获取脚本:", cacheDirName+"/cas.json")
		isDelete = false
	}
	jsonList, jsonErr := util.ReadCasJson(configFleName)
	if jsonErr != nil {
		logrus.Error(jsonErr)
		isDelete = false
	}

	if isDelete {
		//删除虚拟机
		for _, jsonMap := range jsonList {
			var vmName = jsonMap.Name
			logrus.Info("通过名字查询虚拟机id============>", vmName)
			reqUrlQueryName := fmt.Sprintf("http://%s:8080/cas/casrs/vm/basicInfo/%s", a.LoginInfo.Ip, vmName)
			//reqUrlDel := fmt.Sprintf("http://%s:8080/cas/casrs/vm/deleteVm?id=%d&type=0&isWipeVolume=false", a.LoginInfo.Ip, cast.ToInt(vmId))
			resQueryName, resQueryErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrlQueryName, nil)
			if resQueryErr != nil {
				logrus.Error("查询虚拟机id失败:" + resQueryErr.Error())
			} else {
				var vmId = resQueryName.StatusMap["id"]
				logrus.Info("虚拟机id===========>", vmId)
				//查询是否开机状态
				reqUrlQueryStatus := fmt.Sprintf("http://%s:8080/cas/casrs/vm/detail/%d", a.LoginInfo.Ip, cast.ToInt(vmId))
				queryStatusRes, queryStatusErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrlQueryStatus, nil)
				if queryStatusErr != nil {
					logrus.Error(queryStatusErr)
				} else {
					var vmStatus = cast.ToString(queryStatusRes.StatusMap["status"])
					if vmStatus == "running" {
						logrus.Info("关机=======================》", vmName)
						reqUrlShutdown := fmt.Sprintf("http://%s:8080/cas/casrs/vm/stop/%d", a.LoginInfo.Ip, cast.ToInt(vmId))
						_, shutdownErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "PUT", reqUrlShutdown, nil)
						if shutdownErr != nil {
							logrus.Error(shutdownErr)
						}
						time.Sleep(time.Second * 5)
					}
				}

				//删除虚拟机
				reqUrlDel := fmt.Sprintf("http://%s:8080/cas/casrs/vm/deleteVm?id=%d&type=0&isWipeVolume=false", a.LoginInfo.Ip, cast.ToInt(vmId))
				_, err = util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "DELETE", reqUrlDel, nil)
				if err != nil {
					err = errors.New(msg.ERRIR_DELETE_CAS, msg.GetMsg(msg.ERRIR_DELETE_CAS, "param.VmName", err.Error()))
					logrus.Error(err)
				}
			}
		}
		time.Sleep(time.Second * 10)
	}

	//删除存储池
	hostId := cast.ToInt(param.HostId)
	poolName := param.StoreName
	if poolName != "" {
		//先查询存储池
		reqUrl5 := fmt.Sprintf("http://%s:8080/cas/casrs/storage/info?hostId=%d&poolName=%s", a.LoginInfo.Ip, hostId, poolName)
		res, storageErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrl5, nil)
		if storageErr != nil {
			logrus.Error(storageErr, "查询不到存储池，跳过")
		} else {
			if res.StatusMap["path"] != nil {
				logrus.Info("开始删除存储池:", poolName)
				reqUrl3 := fmt.Sprintf("http://%s:8080/cas/casrs/storage/delete?id=%d&hostName=%s&poolName=%s", a.LoginInfo.Ip, hostId, "", poolName)
				delRes, delErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "DELETE", reqUrl3, nil)
				time.Sleep(time.Second * 5)
				if delErr != nil {
					err = errors.New(msg.ERRIR_DELETE_POOL, msg.GetMsg(msg.ERRIR_DELETE_POOL, poolName, delErr.Error()))
					logrus.Error(err)
					return err
				}
				if cast.ToInt(delRes.StatusMap["errorCode"]) != 0 {
					err = errors.New(msg.ERRIR_DELETE_POOL, msg.GetMsg(msg.ERRIR_DELETE_POOL, poolName, delRes.StatusMap["message"]))
					logrus.Error(err)
					return err
				}
			}
		}
	}
	return err
}

func (a *Adapter) login() (result *contants.VerifyResult, err error) {
	var vr contants.VerifyResult

	//先查询版本
	reqUrl1 := "http://" + a.LoginInfo.Ip + ":8080/cas/casrs/version"
	verResp, verResErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrl1, nil)
	if verResErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_CAS, a.LoginInfo.Ip, verResErr.Error()))
		logrus.Info(verResErr.Error())
		return
	}
	if cast.ToInt(verResp.StatusCode) == 200 {
		vr.Version = cast.ToString(verResp.StatusMap["casVersion"])
	}

	//查询主机池
	reqUrl := "http://" + a.LoginInfo.Ip + ":8080/cas/casrs/hostpool/all"
	resp, resErr := util.CasRequest(a.LoginInfo.UserName, a.LoginInfo.PassWord, "GET", reqUrl, nil)
	if resErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_CAS, a.LoginInfo.Ip, resErr.Error()))
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_CAS, a.LoginInfo.Ip, fmt.Sprintf("%d", resp.StatusCode)+":虚拟化平台校验未通过"))
		return
	}
	if resErr != nil {
		logrus.Error("查询主机池报错:", resErr)
	} else {
		vr.HostpoolInfo = resp.StatusMap
	}
	return &vr, err
}

func GetCasHosts(loginInfo contants.LoginInfo) (result map[string]interface{}, err error) {
	//查询主机列表
	reqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/host", loginInfo.Ip)
	queryRes, queryErr := util.CasRequest(loginInfo.UserName, loginInfo.PassWord, "GET", reqUrl, nil)
	if queryErr != nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, loginInfo.Ip, queryErr.Error()))
		logrus.Error(err)
		return nil, err
	}

	return queryRes.StatusMap, err
}

func GetHostInfo(loginInfo contants.LoginInfo, hostId uint) (result map[string]interface{}, err error) {
	//查询主机详细信息
	reqUrl := fmt.Sprintf("http://%s:%d/cas/casrs/host/id/%d", loginInfo.Ip, loginInfo.Port, hostId)
	queryRes, queryErr := util.CasRequest(loginInfo.UserName, loginInfo.PassWord, "GET", reqUrl, nil)
	if err != nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, loginInfo.Ip, queryErr.Error()))
		logrus.Error(err)
		return nil, err
	}
	return queryRes.StatusMap, err
}

func GetCasVms(loginInfo contants.LoginInfo, hostId uint) (result map[string]interface{}, err error) {
	//查询虚拟机列表
	reqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/host/id/%d/vm", loginInfo.Ip, hostId)
	queryRes, queryErr := util.CasRequest(loginInfo.UserName, loginInfo.PassWord, "GET", reqUrl, nil)
	if err != nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, loginInfo.Ip, queryErr.Error()))
		logrus.Error(err)
		return nil, err
	}

	return queryRes.StatusMap, err
}

func GetCasStorages(loginInfo contants.LoginInfo, hostId uint) (result map[string]interface{}, err error) {
	//查询存储
	reqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/host/id/%d/storage", loginInfo.Ip, hostId)
	queryRes, queryErr := util.CasRequest(loginInfo.UserName, loginInfo.PassWord, "GET", reqUrl, nil)
	if err != nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, loginInfo.Ip, queryErr.Error()))
		logrus.Error(err)
		return nil, err
	}
	return queryRes.StatusMap, err
}

func DelStorage(loginInfo contants.LoginInfo, hostId int, poolName string) {
	reqUrl3 := fmt.Sprintf("http://%s:8080/cas/casrs/storage/delete?id=%d&hostName=%s&poolName=%s", loginInfo.Ip, hostId, "", poolName)
	_, err := util.CasRequest(loginInfo.UserName, loginInfo.PassWord, "DELETE", reqUrl3, nil)
	if err != nil {
		logrus.Error("回滚时删除存储池失败:=====>", err)
	}
	logrus.Info("回滚时删除存储池:===========>", poolName)

}

func GetVswitchId(loginInfo contants.LoginInfo, hostName string) (vsId uint, err error) {
	//获取CAS主机的虚拟机交换机ID
	reqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/vswitch/info/%s/vswitch0", loginInfo.Ip, hostName)
	response, resErr := util.CasRequest(loginInfo.UserName, loginInfo.PassWord, "GET", reqUrl, nil)
	if resErr != nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, hostName, resErr.Error()))
	}
	if response.StatusMap["id"] == nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, hostName, "查询不到虚拟交换机vswitch0"))
	}
	return cast.ToUint(response.StatusMap["id"]), err

}
