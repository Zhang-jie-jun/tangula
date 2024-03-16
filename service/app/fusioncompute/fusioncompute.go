package fusioncompute

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/instance"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/internal/system"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Adapter struct {
	LoginInfo *contants.LoginInfo
}

type FcResult struct {
	StatusCode int
	ResBody    map[string]interface{}
}

func NewAdapter(login *contants.LoginInfo) (adapter *Adapter) {
	return &Adapter{LoginInfo: login}
}

func (a *Adapter) login() (result *contants.VerifyResult, err error) {
	var vr contants.VerifyResult
	//获取版本号
	versionRes, versionErr := get_version(a.LoginInfo.Ip)
	if versionErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, a.LoginInfo.Ip, versionErr.Error()))
		return
	}
	if versionRes.StatusCode == 200 {
		resbody := versionRes.ResBody
		versonList := cast.ToSlice(resbody["versions"])
		lastVersionMap := versonList[len(versonList)-1]
		lastVersion := cast.ToStringMap(lastVersionMap)["version"]
		vr.Version = cast.ToString(lastVersion)
	} else {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, a.LoginInfo.Ip,
			cast.ToString(fmt.Sprintf("获取fusioncompute版本号异常:%s", versionRes.ResBody))))
		return
	}
	return &vr, err
}
func (a Adapter) VerifyConfigInfo() (*contants.VerifyResult, error) {
	var res contants.VerifyResult
	var err error
	if a.LoginInfo == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
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

func setParam(sharePath *string, mountParam *view.MountInfo) (*MountFcInfo, error) {
	var detailInfo MountFcInfo
	detailInfo.HostName = cast.ToString(mountParam.AppConfig["fcHostName"])
	detailInfo.HostUrn = cast.ToString(mountParam.AppConfig["fcHostUrn"])
	detailInfo.HostIp = cast.ToString(mountParam.AppConfig["fcHostIp"])
	detailInfo.IsRegisterVM = cast.ToBool(mountParam.AppConfig["isRegisterVM"])
	detailInfo.Os = cast.ToString(mountParam.AppConfig["os"])
	start1 := strings.LastIndex(*sharePath, ":")
	detailInfo.RemoteHost = string([]byte(*sharePath)[0:start1])  // 截取IP
	detailInfo.RemotePath = string([]byte(*sharePath)[start1+1:]) // 截取路径：/tangula/mnt/pool_uuid/image_uuid
	remotePath := string([]byte(*sharePath)[start1+1:])
	start2 := strings.LastIndex(remotePath, "/")
	detailInfo.DataStoreName = fmt.Sprintf("tangula_%s", string([]byte(remotePath)[start2+1:])) // 截取image_uuid作为CAS存储池名称
	detailInfo.StoreName = "tangula"                                                            //存储池名称
	detailInfo.IsRefresh = cast.ToBool(mountParam.AppConfig["isRefresh"])
	logrus.Info("【setParam】", util.SerializeToJson(detailInfo))
	return &detailInfo, nil

}

func (a Adapter) Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error) {
	var err error
	// 设置挂载详细参数
	param, err := setParam(sharePath, mountParam)

	if err != nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, err.Error()))
		logrus.Error(err)
		return "", err
	}
	if param == nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return "", err
	}

	//添加存储资源
	var loginInfo contants.LoginInfo
	loginInfo.Ip = a.LoginInfo.Ip
	loginInfo.UserName = a.LoginInfo.UserName
	loginInfo.PassWord = a.LoginInfo.PassWord

	//查询是否已添加存储
	storageFlag := false
	var storageId uint
	storageRes, stroageErr := GetStorageResource(loginInfo)
	if stroageErr != nil {
		logrus.Error(fmt.Sprintf("查询存储资源异常:%s", stroageErr.Error()))
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, stroageErr.Error()))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, fmt.Sprintf("查询存储资源异常:%s", stroageErr.Error()), "")
		return "", err
	}
	if cast.ToUint(storageRes.ResBody["total"]) > 0 {
		logrus.Info("存储资源tangula已存在")
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, "存储资源tangula已存在", "")
		storageFlag = true
		storegeResInfo := cast.ToSlice(storageRes.ResBody["storeResInfoList"])[0]
		storageUrn := cast.ToStringMap(storegeResInfo)["urn"]
		storageUrnList := strings.Split(cast.ToString(storageUrn), ":")
		storageId = cast.ToUint(storageUrnList[len(storageUrnList)-1])
	}

	//添加存储资源
	if storageFlag == false {
		logrus.Info(fmt.Sprintf("开始在FusionCompute平台添加存储资源:%s", param.StoreName))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, loginInfo.Ip, param.StoreName), "")
		addrRes, addErr := AddStorageResource(param.StoreName, param.RemoteHost, loginInfo)
		if addErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, addErr.Error()))
			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, loginInfo.Ip, addErr.Error()), "")
			return "", err
		}
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, loginInfo.Ip, param.StoreName), "")
		storageUrn := addrRes.ResBody["urn"]
		storageUrnList := strings.Split(cast.ToString(storageUrn), ":")
		storageId = cast.ToUint(storageUrnList[len(storageUrnList)-1])
		time.Sleep(time.Second * 5)
	}
	//查询是否已关联主机
	hostFlag := false
	hostsRes, hostsErr := QueryStorageresourcehosts(storageId, param.HostIp, loginInfo)
	if hostsErr != nil {
		logrus.Error(fmt.Sprintf("查询关联主机异常:%s", hostsErr.Error()))
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, hostsErr.Error()))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, loginInfo.Ip, stroageErr.Error()), "")
		return "", err
	}
	if cast.ToUint(hostsRes.ResBody["total"]) > 0 {
		logrus.Info(fmt.Sprintf("存储资源已关联主机:%s", param.HostIp))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, fmt.Sprintf("存储资源已关联主机:%s", param.HostIp), "")
		hostFlag = true

	}

	//关联主机
	if hostFlag == false {
		logrus.Info(fmt.Sprintf("开始在FusionCompute平台为存储资源关联主机:%s", param.HostIp))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_ADD_FUSIONCOMPUTE_RESOURCE_HOST, loginInfo.Ip, param.HostIp), "")
		_, connectErr := ConnectHost(storageId, param.HostUrn, loginInfo)
		if connectErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, connectErr.Error()))
			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, loginInfo.Ip, connectErr.Error()), "")
			return "", err
		}
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_SUCESS_ADD_FUSIONCOMPUTE_RESOURCE_HOST, loginInfo.Ip, param.HostIp), "")
		time.Sleep(time.Second * 5)
	}
	if param.IsRefresh {
		//先查询当前是不是有扫描任务正在进行
		var count1 = 0
		refreshFlag := false
		for {
			//查询任务信息
			taskRes, _ := GetTask("RefreshHostStorageUnitTask", loginInfo)
			taskList := cast.ToSlice(taskRes.ResBody["tasks"])
			if len(taskList) > 0 {
				taskStatus := cast.ToStringMap(taskList[0])["status"]
				if taskStatus == "running" {
					instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, "=======等待其他扫描任务结束=======次数:"+cast.ToString(count1+1), "")
				} else {
					logrus.Info(fmt.Sprintf("开始在FusionCompute平台扫描存储设备:%s", *sharePath))
					instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_FUSIONCOMPUTE_REFRESH, loginInfo.Ip), "")
					refreshFlag = true
				}
			} else {
				refreshFlag = true
			}

			if count1 > 20 {
				logrus.Error("等待超时")
				break
			}
			if refreshFlag {
				break
			}
			time.Sleep(time.Second * 5)
			count1++
		}
		//扫描存储设备
		refreshErr := Refresh(param.HostUrn, loginInfo)
		if refreshErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, refreshErr.Error()))
			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_REFRESH, loginInfo.Ip, refreshErr.Error()), "")
			return "", err
		}

		//等待刷新完成
		var count = 0
		queryFlag := false
		//获取设备urn
		var storageUnitUrn string
		for {
			logrus.Info(fmt.Sprintf("第%d次查询存储设备……", count+1))
			//查询任务信息
			taskRes, _ := GetTask("RefreshHostStorageUnitTask", loginInfo)
			taskList := cast.ToSlice(taskRes.ResBody["tasks"])
			if len(taskList) > 0 {
				process := cast.ToStringMap(taskList[0])["progress"]
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, "当前扫描进度-----"+cast.ToString(process)+"%", "")
			}

			if count > 20 {
				logrus.Error("扫描超时")
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_REFRESH, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_REFRESH, loginInfo.Ip, "扫描超时，仍未识别到存储设备"))
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, "扫描超时，仍未识别到存储设备", "")
				break
			}
			queryRes, queryErr := Queryallstorageunit(loginInfo)
			if queryErr != nil {
				logrus.Error(fmt.Sprintf("查询存储设备报错：%s", queryErr.Error()))
			} else {
				suList := cast.ToSlice(queryRes.ResBody["suList"])
				if len(suList) > 0 {
					for _, iter := range suList {
						if cast.ToString(cast.ToStringMap(iter)["name"]) == *sharePath {
							logrus.Error(fmt.Sprintf("已识别到存储设备：%s", cast.ToStringMap(iter)["name"]))
							storageUnitUrn = cast.ToString(cast.ToStringMap(iter)["urn"])
							queryFlag = true
							break
						}
					}
				}

			}
			if queryFlag {
				break
			}
			time.Sleep(time.Second * 5)
			count++
		}
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_SUCESS_FUSIONCOMPUTE_REFRESH, loginInfo.Ip, *sharePath), "")
		logrus.Info(fmt.Sprintf("存储设备urn：%s", storageUnitUrn))
		time.Sleep(time.Second * 5)

		//添加数据存储
		addDatastoreRes, addDatastoreErr := AddDatastore(param.DataStoreName, param.HostUrn, storageUnitUrn, loginInfo)
		if addDatastoreErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, addDatastoreErr.Error()))
			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, loginInfo.Ip, addDatastoreErr.Error()), "")
			return "", err
		}
		if addDatastoreRes.StatusCode != 200 {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, util.SerializeToJson(addDatastoreRes.ResBody)))
			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, loginInfo.Ip, util.SerializeToJson(addDatastoreRes.ResBody)), "")
			return "", err
		}
		datastoreUrn := addDatastoreRes.ResBody["urn"]
		logrus.Info(fmt.Sprintf("添加数据存储成功:%s", param.DataStoreName))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_SUCESS_FUSIONCOMPUTE_ADD_DATASTORE, loginInfo.Ip, param.DataStoreName), "")

		//创建虚拟机
		if param.IsRegisterVM {
			//等待添加数据存储完成
			var addDatastoreCount = 0
			addDatastoreFlag := false
			for {
				//查询任务信息
				taskRes, _ := GetTask("AddDataStoreTask", loginInfo)
				taskList := cast.ToSlice(taskRes.ResBody["tasks"])
				if len(taskList) > 0 {
					taskStatus := cast.ToStringMap(taskList[0])["status"]
					if taskStatus == "running" {
						instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, "=======等待添加数据存储完成======="+cast.ToString(addDatastoreCount+1), "")
					} else if taskStatus == "success" {
						addDatastoreFlag = true
					}
				} else {
					addDatastoreFlag = true
				}

				if addDatastoreCount > 20 {
					logrus.Error("等待超时")
					break
				}
				if addDatastoreFlag {
					break
				}
				time.Sleep(time.Second * 5)
				addDatastoreCount++
			}

			replicaObj, findErr := replica.ReplicaMgm.FindById(mountParam.ReplicaId)
			if findErr != nil {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, findErr.Error()))
				logrus.Error(err)
				return "", err
			}

			//检查副本中磁盘文件，要求是裸qcow2磁盘文件
			diskPath := param.RemotePath
			dirSub, readErr := ioutil.ReadDir(diskPath)
			if readErr != nil {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, readErr.Error()))
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, readErr.Error()), "")
				return "", err
			}
			for _, iter := range dirSub {
				if iter.IsDir() {
					err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, "副本中包含文件夹，请检查磁盘文件"))
					return "", err
				}
			}

			//解析磁盘文件
			cmdRes, cmdErr := util.GetMountDiskInfo(diskPath)
			if cmdErr != nil {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, cmdErr.Error()))
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, cmdErr.Error()), "")
				return "", err
			}
			var diskInfoList []view.DiskInfo
			for _, line := range cmdRes {

				logrus.Info(fmt.Sprintf("获取磁盘信息:%s", line))
				splitRet := strings.Split(line, " ")
				if len(splitRet) != 2 {
					logrus.Error("split line ret error.")
					continue
				}
				//检查是否qcow类型文件
				fileCmd := fmt.Sprintf("file %s", diskPath+"/"+splitRet[1])
				fileRes, fileErr := system.SysManage.RunCommand(fileCmd)
				logrus.Info(fmt.Sprintf("file执行结果:%s", cast.ToString(fileRes[0])))
				if fileErr != nil {
					logrus.Error(fileErr)
				} else {
					if strings.Contains(fileRes[0], "QCOW") {
						diskInfo := view.DiskInfo{}
						diskInfo.Name = splitRet[1]
						diskInfo.PoolName = replicaObj.Uuid
						diskInfo.SizeByte, _ = strconv.ParseInt(splitRet[0], 10, 64)
						diskInfoList = append(diskInfoList, diskInfo)
					}
				}
			}
			logrus.Info("-----------diskInfoList----------", diskInfoList)

			var crtDiskList []map[string]interface{}
			for index, _ := range diskInfoList {
				crtDiskInfo := map[string]interface{}{
					"datastoreUrn":   datastoreUrn,
					"quantityGB":     1,
					"sequenceNum":    index + 1,
					"type":           "normal",
					"indepDisk":      false,
					"persistentDisk": true,
					"isThin":         false,
					"volType":        0,
					"pciType":        "VIRTIO",
					"bootOrder":      -1,
					"encrypted":      "0",
				}
				crtDiskList = append(crtDiskList, crtDiskInfo)
			}

			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_FUSIONCOMPUTE_CRT_VM, loginInfo.Ip, replicaObj.Name), "")
			crtRes, crtErr := CrtVm(replicaObj.Name, param.Os, param.HostUrn, crtDiskList, loginInfo)
			if crtErr != nil {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, crtErr.Error()))
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, crtErr.Error()), "")
				return "", err
			}
			if crtRes.StatusCode != 200 {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, util.SerializeToJson(crtRes.ResBody)))
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, util.SerializeToJson(crtRes.ResBody)), "")
				return "", err
			}
			vmUrn := cast.ToString(crtRes.ResBody["urn"])
			vmSplitList := strings.Split(cast.ToString(vmUrn), ":")
			vmId := vmSplitList[len(vmSplitList)-1]
			logrus.Info(fmt.Sprintf("虚拟机urn：%s", crtRes.ResBody["urn"]))
			time.Sleep(time.Second * 5)
			//等待虚拟机创建完成
			//先查询当前是不是有扫描任务正在进行
			var vmCount = 0
			vmFlag := false
			for {
				//查询任务信息
				taskRes, _ := GetTask("CreateVmTask", loginInfo)
				taskList := cast.ToSlice(taskRes.ResBody["tasks"])
				if len(taskList) > 0 {
					taskStatus := cast.ToStringMap(taskList[0])["status"]
					if taskStatus == "running" {
						instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, "=======等待虚拟机创建完成======="+cast.ToString(vmCount+1), "")
					} else if taskStatus == "success" {
						vmFlag = true
					}
				} else {
					vmFlag = true
				}

				if vmCount > 20 {
					logrus.Error("等待超时")
					break
				}
				if vmFlag {
					break
				}
				time.Sleep(time.Second * 5)
				vmCount++
			}

			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_SUCESS_FUSIONCOMPUTE_CRT_VM, loginInfo.Ip, replicaObj.Name), "")

			//替换磁盘文件
			getVmRes, getVmErr := GetVmInfo(vmId, loginInfo)
			if getVmErr != nil {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, getVmErr.Error()))
				return "", err
			}

			if getVmRes.StatusCode != 200 {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, util.SerializeToJson(getVmRes.ResBody)))
				return "", err
			}
			vmDiskList := cast.ToSlice(cast.ToStringMap(getVmRes.ResBody["vmConfig"])["disks"])
			logrus.Info("-----------虚拟机磁盘信息----------", vmDiskList)

			if len(vmDiskList) != len(diskInfoList) {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, "磁盘文件数量与磁盘数量不一致!"))
				return "", err
			}

			for index, vmDisk := range vmDiskList {
				volumeUuid := cast.ToStringMap(vmDisk)["volumeUuid"]
				diskFileName := fmt.Sprintf("vol_%s.img", volumeUuid)
				bakDiskFileName := fmt.Sprintf("vol_%s_bak.img", volumeUuid)
				diskFilePath := fmt.Sprintf("%s/vol/vol_%s", param.RemotePath, volumeUuid)
				bakDiskCmd := fmt.Sprintf("mv %s/%s %s/%s", diskFilePath, diskFileName, diskFilePath, bakDiskFileName)

				_, bakDiskErr := system.SysManage.RunCommand(bakDiskCmd)
				if bakDiskErr != nil {
					logrus.Error(bakDiskErr)
					err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, bakDiskErr.Error()))
					return "", err
				}
				mvDiskCmd := fmt.Sprintf("cp %s/%s %s/%s", param.RemotePath, diskInfoList[index].Name, diskFilePath, diskFileName)
				_, mvDiskErr := system.SysManage.RunCommand(mvDiskCmd)
				if mvDiskErr != nil {
					logrus.Error(mvDiskErr)
					err = errors.New(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_MOUNT_FAILED, loginInfo.Ip, mvDiskErr.Error()))
					return "", err
				}

				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_SUCESS_FUSIONCOMPUTE_MV_DISK, loginInfo.Ip, replicaObj.Name), "")
			}

		}
	}

	return "", err
}

func (a Adapter) UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) error {
	var err error
	logrus.Infof("卸载FusionCompute……, %s", *sharePath)
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
	var loginInfo contants.LoginInfo
	loginInfo.Ip = a.LoginInfo.Ip
	loginInfo.UserName = a.LoginInfo.UserName
	loginInfo.PassWord = a.LoginInfo.PassWord

	//查询数据存储
	queryRes, queryErr := QueryDatastore(param.DataStoreName, loginInfo)
	if queryErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, loginInfo.Ip, queryErr.Error()))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, loginInfo.Ip, queryErr.Error()), "")
		return err
	}
	if queryRes.StatusCode != 200 {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, loginInfo.Ip, queryRes.ResBody))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, loginInfo.Ip, queryRes.ResBody), "")
		return err
	}

	datastores := queryRes.ResBody["datastores"]
	if len(cast.ToSlice(datastores)) == 0 {
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, loginInfo.Ip, "查询不到数据存储信息,跳过此步骤"), "")
	} else {
		datastoreUrn := cast.ToStringMap(cast.ToSlice(datastores)[0])["urn"]
		datastoreUrnList := strings.Split(cast.ToString(datastoreUrn), ":")
		datastoreId := cast.ToUint(datastoreUrnList[len(datastoreUrnList)-1])

		//卸载数据存储
		delRes, delErr := DelDatastore(datastoreId, param.HostUrn, loginInfo)
		if delErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, loginInfo.Ip, delErr.Error()))
			instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, loginInfo.Ip, delErr.Error()), "")
			return err
		}
		if delRes.StatusCode != 200 {
			if delRes.ResBody["errorCode"] != "10410006" { //数据存储已经不存在，忽略
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED, loginInfo.Ip, delRes.ResBody))
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, loginInfo.Ip, delRes.ResBody), "")
				return err
			}
		}
		time.Sleep(time.Second * 10)
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_SUCESS_FUSIONCOMPUTE_DEL_DATASTORE, loginInfo.Ip, cast.ToStringMap(cast.ToSlice(datastores)[0])["name"]), "")
	}

	return err
}

func get_version(platformIp string) (response *FcResult, err error) {
	var fcRes FcResult
	reqUrl := fmt.Sprintf("https://%s:7443/service/versions", platformIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	//增加header选项
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", platformIp))

	if reqErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, reqErr.Error()))
		return
	}
	//处理返回结果
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, resErr := client.Do(req)
	if resErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, resErr.Error()))
		return
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, readErr.Error()))
		return
	}
	resBody := make(map[string]interface{})
	jsonErr := json.Unmarshal(body, &resBody)
	if jsonErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, jsonErr.Error()))
		return
	}
	fcRes.StatusCode = resp.StatusCode
	fcRes.ResBody = resBody
	return &fcRes, err
}

func login_platform(platformIp string, username string, password string) (token string, err error) {
	//获取版本号
	var version string
	versionRes, versionErr := get_version(platformIp)
	if versionErr != nil {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, versionErr.Error()))
		logrus.Error(err)
		return
	}
	if versionRes.StatusCode == 200 {
		resbody := versionRes.ResBody
		versonList := cast.ToSlice(resbody["versions"])
		lastVersionMap := versonList[len(versonList)-1]
		lastVersion := cast.ToStringMap(lastVersionMap)["version"]
		version = cast.ToString(lastVersion)
	} else {
		err = errors.New(msg.ERROR_LOGIN_CAS, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp,
			cast.ToString(fmt.Sprintf("获取fusioncompute版本号异常:%s", versionRes.ResBody))))
		logrus.Error(err)
		return
	}

	reqUrl := fmt.Sprintf("https://%s:7443/service/session", platformIp)
	req, reqErr := http.NewRequest("POST", reqUrl, nil)
	//增加header选项
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Accept", fmt.Sprintf("application/json;charset=UTF-8;version=%s", version))
	req.Header.Set("Content-Type", "appilcation/json;charset=UTF-8")
	req.Header.Set("X-Auth-User", username)
	req.Header.Set("X-Auth-Key", password)
	req.Header.Set("X-Auth-UserType", cast.ToString(2))
	req.Header.Set("X-ENCRIPT-ALGORITHM", cast.ToString(1))

	if reqErr != nil {
		err = errors.New(msg.ERROR_LOGIN_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, reqErr.Error()))
		logrus.Error(err)
		return
	}
	//处理返回结果
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, resErr := client.Do(req)
	if resErr != nil {
		err = errors.New(msg.ERROR_LOGIN_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, resErr.Error()))
		logrus.Error(err)
		return
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		err = errors.New(msg.ERROR_LOGIN_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, readErr.Error()))
		logrus.Error(err)
		return
	}
	resMap := make(map[string]interface{})
	jsonErr := json.Unmarshal(body, &resMap)
	if jsonErr != nil {
		err = errors.New(msg.ERROR_LOGIN_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp, jsonErr.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(fmt.Sprintf("【%s】=>StatusCode:%d", reqUrl, resp.StatusCode))
	if resp.StatusCode != 200 {
		err = errors.New(msg.ERROR_LOGIN_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, platformIp,
			fmt.Sprintf("登录返回异常：%s", util.SerializeToJson(resMap))))
		logrus.Error(err)
		return
	}
	logrus.Info(fmt.Sprintf("【%s】=>resBody:%s", reqUrl, util.SerializeToJson(resMap)))
	return resp.Header["X-Auth-Token"][0], nil
}
func GetHosts(loginInfo contants.LoginInfo) (result map[string]interface{}, err error) {
	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, siteErr.Error()))
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		hostUri := fmt.Sprintf("%s/hosts", siteUri)
		hostRes, hostErr := FusionComputeRequest(hostUri, "GET", nil, loginInfo)
		if hostErr != nil {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, hostErr.Error()))
		}
		if hostRes.StatusCode != 200 {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, util.SerializeToJson(hostRes.ResBody)))
		}
		return hostRes.ResBody, err
	} else {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, "site列表为空！"))
	}

	return result, err
}

func GetStorageResource(loginInfo contants.LoginInfo) (response *FcResult, err error) {
	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, siteErr.Error()))
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		//splitList:=strings.Split(cast.ToString(siteUri), "/")
		//siteId := splitList[len(splitList)-1]
		getStorageResourceUri := fmt.Sprintf("%s/storageresources/queryallstorageresource?limit=5&offset=0&name=tangula", siteUri)
		resourceRes, resourceErr := FusionComputeRequest(getStorageResourceUri, "GET", nil, loginInfo)
		if resourceErr != nil {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, resourceErr.Error()))
			return
		}
		if resourceRes.StatusCode != 200 {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, util.SerializeToJson(resourceRes.ResBody)))
			return
		}
		return resourceRes, err

	} else {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, "site列表为空！"))
	}

	return

}

func AddStorageResource(storageName string, remoteIp string, loginInfo contants.LoginInfo) (response *FcResult, err error) {

	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, loginInfo.Ip, siteErr.Error()))
		return
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		addUrl := fmt.Sprintf("%s/storageresources", siteUri)
		data_channel_list := []map[string]interface{}{}
		data_channel_map := map[string]interface{}{"ip": remoteIp}
		data_channel_list = append(data_channel_list, data_channel_map)
		reqParam := map[string]interface{}{
			"name":        storageName,
			"dataChannel": data_channel_list,
			"storageType": "NAS",
			"vender":      "OTHER",
			"deviceType":  "OTHER",
			"autoscan":    0,
		}
		addRes, addErr := FusionComputeRequest(addUrl, "POST", reqParam, loginInfo)
		if addErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, loginInfo.Ip, addErr.Error()))
			return
		}
		if addRes.StatusCode != 200 {
			if addRes.ResBody["errorCode"] == 10410021 || addRes.ResBody["errorCode"] == 10410066 { //存储名称或ip重复，忽略
				logrus.Info(addRes.ResBody["errorDes"])
			} else {
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, loginInfo.Ip, util.SerializeToJson(addRes.ResBody)))
				return
			}
		}
		return addRes, err
	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED, loginInfo.Ip, "site列表为空！"))
	}
	return

}

func QueryStorageresourcehosts(storageId uint, hostIp string, loginInfo contants.LoginInfo) (response *FcResult, err error) { //查询存储资源关联主机
	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, siteErr.Error()))
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		getStorageResourceUri := fmt.Sprintf("%s/storageresources/%d/querystorageresourcehosts?limit=5&offset=0&status=1&ip=%s", siteUri, storageId, hostIp)
		hostsRes, hostseErr := FusionComputeRequest(getStorageResourceUri, "GET", nil, loginInfo)
		if hostseErr != nil {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, hostseErr.Error()))
			return
		}
		if hostsRes.StatusCode != 200 {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, util.SerializeToJson(hostsRes.ResBody)))
			return
		}
		return hostsRes, err

	} else {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, "site列表为空！"))
	}

	return

}

func ConnectHost(storageId uint, hostUrn string, loginInfo contants.LoginInfo) (response *FcResult, err error) { //关联主机

	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, loginInfo.Ip, siteErr.Error()))
		return
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		addUrl := fmt.Sprintf("%s/storageresources/%d/action/connect?resID=%d", siteUri, storageId, storageId)
		reqParam := map[string]interface{}{
			"hostUrn": hostUrn,
		}
		addRes, addErr := FusionComputeRequest(addUrl, "POST", reqParam, loginInfo)
		//已经关联主机可以重复关联，不报错的
		if addErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, loginInfo.Ip, addErr.Error()))
			return
		}
		if addRes.StatusCode != 200 {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, loginInfo.Ip, util.SerializeToJson(addRes.ResBody)))
			return
		}
		return addRes, err
	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, loginInfo.Ip, "site列表为空！"))
	}
	return
}

func Refresh(hostUrn string, loginInfo contants.LoginInfo) (err error) { //刷新存储设备

	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_REFRESH, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_REFRESH, loginInfo.Ip, siteErr.Error()))
		return
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		refreshUrl := fmt.Sprintf("%s/storageunits/action/refresh", siteUri)
		reqParam := map[string]interface{}{
			"hostUrn": hostUrn,
		}
		refreshRes, refreshErr := FusionComputeRequest(refreshUrl, "POST", reqParam, loginInfo)
		if refreshErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_REFRESH, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_REFRESH, loginInfo.Ip, refreshErr.Error()))
			return
		}
		if refreshRes.StatusCode != 200 {
			if refreshRes.ResBody["errorCode"] != "10410023" { //除了正在刷新中返回
				logrus.Error(util.SerializeToJson(refreshRes.ResBody))
				err = errors.New(msg.ERROR_FUSIONCOMPUTE_REFRESH, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_REFRESH, loginInfo.Ip, refreshRes.ResBody))
				return
			}
		}

		return
	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_REFRESH, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_REFRESH, loginInfo.Ip, "site列表为空！"))
	}
	return
}

func Queryallstorageunit(loginInfo contants.LoginInfo) (response *FcResult, err error) { //查询存储设备
	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, siteErr.Error()))
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		storageunitUri := fmt.Sprintf("%s/storageunits/queryallstorageunit?limit=99&offset=0&useState=all&deviceType=0&type=NAS&name=tangula", siteUri)
		unitRes, unitErr := FusionComputeRequest(storageunitUri, "GET", nil, loginInfo)
		if unitErr != nil {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, unitErr.Error()))
			return
		}
		return unitRes, err

	} else {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, "site列表为空！"))
	}

	return

}

func AddDatastore(name string, hostUrn string, storageUnitUrn string, loginInfo contants.LoginInfo) (response *FcResult, err error) { //添加数据存储

	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, loginInfo.Ip, siteErr.Error()))
		return
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		addUrl := fmt.Sprintf("%s/datastores", siteUri)
		reqParam := map[string]interface{}{
			"hostUrn":        hostUrn,
			"name":           name,
			"storageUnitUrn": storageUnitUrn,
			"useType":        1,
		}
		addRes, addErr := FusionComputeRequest(addUrl, "POST", reqParam, loginInfo)
		if addErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, loginInfo.Ip, addErr.Error()))
			return
		}
		if addRes.StatusCode != 200 {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, loginInfo.Ip, util.SerializeToJson(addRes.ResBody)))
			return
		}
		return addRes, err
	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_DATASTORE, loginInfo.Ip, "site列表为空！"))
	}
	return
}

func DelDatastore(datastoreId uint, hostUrn string, loginInfo contants.LoginInfo) (response *FcResult, err error) { //删除数据存储

	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, loginInfo.Ip, siteErr.Error()))
		return
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		delUrl := fmt.Sprintf("%s/datastores/%d?hostUrn=%s&isForce=false&dataStoreID=%d", siteUri, datastoreId, hostUrn, datastoreId)
		delRes, delErr := FusionComputeRequest(delUrl, "DELETE", nil, loginInfo)
		if delErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, loginInfo.Ip, delErr.Error()))
			return
		}
		if delRes.StatusCode != 200 {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, loginInfo.Ip, util.SerializeToJson(delRes.ResBody)))
			return
		}
		return delRes, err
	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_DEL_DATASTORE, loginInfo.Ip, "site列表为空！"))
	}
	return
}

func QueryDatastore(datastoreName string, loginInfo contants.LoginInfo) (response *FcResult, err error) { //查询存储设备
	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, loginInfo.Ip, siteErr.Error()))
		return
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		queryDatastoreUri := fmt.Sprintf("%s/datastores/ex?page=1&offset=0&limit=5&name=%s", siteUri, datastoreName)
		unitRes, unitErr := FusionComputeRequest(queryDatastoreUri, "GET", nil, loginInfo)
		if unitErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, loginInfo.Ip, unitErr.Error()))
			return
		}
		return unitRes, err

	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_DATASTORE, loginInfo.Ip, "site列表为空！"))
	}

	return

}

func GetTask(taskType string, loginInfo contants.LoginInfo) (response *FcResult, err error) { //查询扫描任务
	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, siteErr.Error()))
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		storageunitUri := fmt.Sprintf("%s/tasks?limit=10&offset=0&type=%s", siteUri, taskType)
		taskRes, taskErr := FusionComputeRequest(storageunitUri, "GET", nil, loginInfo)
		if taskErr != nil {
			err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, taskErr.Error()))
			return
		}
		response = taskRes
		if taskType == "RefreshHostStorageUnitTask" {
			if taskRes.ResBody["errorCode"] == "10000001" { //换老接口再查询
				storageunitUriOld := fmt.Sprintf("%s/tasks?limit=10&offset=0&type=RefreshStorageUnitTask", siteUri)
				taskResOldRes, taskResOldErr := FusionComputeRequest(storageunitUriOld, "GET", nil, loginInfo)
				if taskResOldErr != nil {
					err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, taskErr.Error()))
					return
				}
				response = taskResOldRes

			}
		}
		return response, err

	} else {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, "site列表为空！"))
	}

	return

}

func CrtVm(vmName string, os string, hostUrn string, diskList []map[string]interface{}, loginInfo contants.LoginInfo) (response *FcResult, err error) { //关联主机

	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED, loginInfo.Ip, siteErr.Error()))
		return
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	var osVersion uint
	if os == "Linux" {
		osVersion = 302
	} else {
		osVersion = 202
	}

	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		splitList := strings.Split(cast.ToString(siteUri), "/")
		siteId := splitList[len(splitList)-1]
		crtUrl := fmt.Sprintf("%s/vms?siteID=%s", siteUri, siteId)
		reqParam := map[string]interface{}{
			"name":          vmName,
			"description":   "",
			"location":      hostUrn,
			"parentObjUrn":  fmt.Sprintf("urn:sites:%s", siteId),
			"enableImc":     false,
			"isBindingHost": false,
			"osOptions": map[string]interface{}{
				"osType":    os,
				"osVersion": osVersion,
			},
			"vmConfig": map[string]interface{}{
				"cpu": map[string]interface{}{
					"cpuHotPlug":      0,
					"cpuPolicy":       "shared",
					"cpuThreadPolicy": "prefer",
					"weight":          500,
					"quantity":        1,
					"limit":           0,
					"reservation":     0,
					"coresPerSocket":  1,
				},
				"memory": map[string]interface{}{
					"quantityMB":  2048,
					"hugePage":    "4K",
					"limit":       0,
					"reservation": 0,
					"weight":      20480,
				},
				"disks": diskList,
				"graphicsCard": map[string]interface{}{
					"type": "cirrus",
					"size": 4,
				},
				"properties": map[string]interface{}{
					"antivirusMode":      "",
					"isAutoAdjustNuma":   false,
					"bootFirmware":       "BIOS",
					"bootFirmwareTime":   0,
					"bootOption":         "disk",
					"clockMode":          "freeClock",
					"vmVncKeymapSetting": 7,
					"isHpet":             false,
					"isEnableHa":         false,
					"evsAffinity":        false,
					"secureVmType":       "",
					"realtime":           false,
					"isEnableMemVol":     false,
					"isEnableFt":         false,
					"isAutoUpgrade":      true,
					"emulatorResType":    nil,
					"dpiVmType":          "",
					"cdRomBootOrder":     -1,
					"attachType":         false,
					"enableWatchDog":     false,
				},
				"gpuGroups": []map[string]interface{}{},
			},
			"cpuVendor":         "Intel",
			"autoBoot":          false,
			"isSrcTemplate":     false,
			"isEnableIntegrity": false,
		}
		crtRes, crtErr := FusionComputeRequest(crtUrl, "POST", reqParam, loginInfo)
		if crtErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_CRT_VM, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_CRT_VM, loginInfo.Ip, crtErr.Error()))
			return
		}
		return crtRes, err
	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_CRT_VM, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_CRT_VM, loginInfo.Ip, "site列表为空！"))
	}
	return
}

func GetVmInfo(vmId string, loginInfo contants.LoginInfo) (response *FcResult, err error) { //查询虚拟机信息
	//查询sites
	sitesReqUrl := "/service/sites"
	siteRes, siteErr := FusionComputeRequest(sitesReqUrl, "GET", nil, loginInfo)
	if siteErr != nil {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED, loginInfo.Ip, siteErr.Error()))
	}
	siteList := cast.ToSlice(siteRes.ResBody["sites"])
	if len(siteList) > 0 {
		siteMap := cast.ToStringMap(siteList[0])
		siteUri := siteMap["uri"]
		splitList := strings.Split(cast.ToString(siteUri), "/")
		siteId := splitList[len(splitList)-1]
		vmInfoUri := fmt.Sprintf("%s/vms/%s?siteID=%s&vmID=%s", siteUri, vmId, siteId, vmId)
		vmInfoRes, vmInfoErr := FusionComputeRequest(vmInfoUri, "GET", nil, loginInfo)
		if vmInfoErr != nil {
			err = errors.New(msg.ERROR_FUSIONCOMPUTE_GET_VM_INFO, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_VM_INFO, loginInfo.Ip, vmInfoErr.Error()))
			return
		}
		response = vmInfoRes
		return response, err

	} else {
		err = errors.New(msg.ERROR_FUSIONCOMPUTE_GET_VM_INFO, msg.GetMsg(msg.ERROR_FUSIONCOMPUTE_GET_VM_INFO, loginInfo.Ip, "site列表为空！"))
	}

	return

}

func FusionComputeRequest(url string, method string, reqParam map[string]interface{}, loginInfo contants.LoginInfo) (response *FcResult, err error) {
	//请求fusioncompute平台接口，登录除外
	var fcRes FcResult
	token, loginErr := login_platform(loginInfo.Ip, loginInfo.UserName, loginInfo.PassWord)
	if loginErr != nil {
		err = errors.New(msg.ERROR_LOGIN_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_LOGIN_FUSIONCOMPUTE, loginInfo.Ip, loginErr.Error()))
		logrus.Error(fmt.Sprintf("登录fusioncompute平台错误:%s", err))
		return
	}
	reqUrl := fmt.Sprintf("https://%s:7443%s", loginInfo.Ip, url)
	var req *http.Request

	if reqParam != nil {
		logrus.Info(fmt.Sprintf("reqParam=>:%s", util.SerializeToJson(reqParam)))
		bytesData, bytesErr := json.Marshal(reqParam)
		if bytesErr != nil {
			err = errors.New(msg.ERROR_REQUEST_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_REQUEST_FUSIONCOMPUTE, loginInfo.Ip, bytesErr.Error()))
			logrus.Error(err)
			return
		}
		req, _ = http.NewRequest(method, reqUrl, bytes.NewReader(bytesData))
	} else {
		req, _ = http.NewRequest(method, reqUrl, nil)
	}

	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", loginInfo.Ip))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, resErr := client.Do(req)
	if resErr != nil {
		err = errors.New(msg.ERROR_REQUEST_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_REQUEST_FUSIONCOMPUTE, loginInfo.Ip, resErr.Error()))
		logrus.Error(err)
		return
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		err = errors.New(msg.ERROR_REQUEST_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_REQUEST_FUSIONCOMPUTE, loginInfo.Ip, readErr.Error()))
		logrus.Error(err)
		return
	}
	resMap := make(map[string]interface{})
	jsonErr := json.Unmarshal(body, &resMap)
	if jsonErr != nil {
		err = errors.New(msg.ERROR_REQUEST_FUSIONCOMPUTE, msg.GetMsg(msg.ERROR_REQUEST_FUSIONCOMPUTE, loginInfo.Ip, jsonErr.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(fmt.Sprintf("【%s】=>StatusCode:%d", reqUrl, resp.StatusCode))
	logrus.Info(fmt.Sprintf("【%s】=>resBody:%s", reqUrl, util.SerializeToJson(resMap)))
	fcRes.StatusCode = resp.StatusCode
	fcRes.ResBody = resMap
	return &fcRes, err
}
