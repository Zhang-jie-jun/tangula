package storage

import (
	"encoding/json"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/host"
	"github.com/Zhang-jie-jun/tangula/internal/dao/instance"
	"github.com/Zhang-jie-jun/tangula/internal/dao/platform"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/internal/dao/script"
	"github.com/Zhang-jie-jun/tangula/internal/system"
	"github.com/Zhang-jie-jun/tangula/internal/vsphere"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/Zhang-jie-jun/tangula/service/app"
	"github.com/Zhang-jie-jun/tangula/service/app/fusioncompute"
	"github.com/Zhang-jie-jun/tangula/service/app/vmware"
	"github.com/Zhang-jie-jun/tangula/service/websockets"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"strconv"
	"time"
)

func MountReplica(param *view.MountParam, user *auth.User) (result map[string]interface{}, err error) {
	var obj replica.Replica
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.MOUNT_REPLICA_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.MOUNT_REPLICA, obj.Name, detail, user.Account, contants.LOG_FAILED)
		}
	}()
	obj, err = replica.ReplicaMgm.FindById(param.MountInfo.ReplicaId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 检查资源权限
	err = service.CheckResource(service.QUERY_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 检查状态
	if obj.Status != contants.NOT_MOUNT {
		err = errors.New(msg.ERROR_REPLICA_STATUS_NOT_MOUNT, msg.GetMsg(msg.ERROR_REPLICA_STATUS_NOT_MOUNT))
		logrus.Error(err)
		return
	}
	// 检查类型
	if !service.CheckType(obj.Type, param.MountInfo.TargetType) {
		err = errors.New(msg.ERROR_REPLICA_TYPE_MISMATCHING, msg.GetMsg(msg.ERROR_REPLICA_TYPE_MISMATCHING))
		logrus.Error(err)
		return
	}
	// 获取目标平台类型、登录信息
	var login contants.LoginInfo
	var targetId uint
	if param.MountInfo.TargetType == contants.LINUX || param.MountInfo.TargetType == contants.WINDOWS {
		hostOrm, e := host.HostMgm.FindById(param.MountInfo.TargetId)
		if e != nil {
			err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, e.Error()))
			logrus.Error(err)
			return
		}
		passWord, e := util.AesDecrypt(hostOrm.PassWord)
		if e != nil {
			err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
			logrus.Error(err)
			return
		}
		login.Ip = hostOrm.Ip
		login.Port = hostOrm.Port
		login.UserName = hostOrm.UserName
		login.PassWord = passWord
		targetId = hostOrm.Id
	} else {
		platformOrm, e := platform.PlatformMgm.FindById(param.MountInfo.TargetId)
		if e != nil {
			err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, e.Error()))
			logrus.Error(err)
			return
		}
		passWord, e := util.AesDecrypt(platformOrm.PassWord)
		if e != nil {
			err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
			logrus.Error(err)
			return
		}
		login.Ip = platformOrm.Ip
		login.Port = platformOrm.Port
		login.UserName = platformOrm.UserName
		login.PassWord = passWord
		targetId = platformOrm.Id

		//检查VMware平台
		if obj.Type == contants.VMWAREVM {
			if cast.ToBool(param.MountInfo.AppConfig["isRegisterVM"]) {
				vms := vmware.NewVmWare(login.Ip, login.UserName, login.PassWord)
				vmList, _, _ := vms.GetAllVmClient()
				for _, vm := range vmList {
					//检查虚拟机名称是否已存在
					if cast.ToString(param.MountInfo.AppConfig["vmName"]) == vm.Name {
						err = errors.New(msg.ERROR_VMNAME_EXSIT, msg.GetMsg(msg.ERROR_VMNAME_EXSIT))
						logrus.Error(err)
						return
					}

					//检查IP地址冲突
					if cast.ToString(param.MountInfo.AppConfig["addr"]) != "" {
						if cast.ToString(param.MountInfo.AppConfig["addr"]) == vm.Ip {
							err = errors.New(msg.ERROR_IP_EXSIT, msg.GetMsg(msg.ERROR_IP_EXSIT))
							logrus.Error(err)
							return
						}
					}

				}

			}
		}
	}
	// 创建执行实例
	var instanceObj instance.Instance
	instanceObj.Type = contants.INSTANCE_MOUNT
	instanceObj.Status = contants.MOUNT_RUNNING
	instanceObj.ReplicaId = obj.Id
	instanceObj.TargetType = param.MountInfo.TargetType
	instanceObj.TargetId = targetId
	instanceObj.MountPoint = ""
	instanceObj, err = instance.InstanceMgm.CreateInstance(instanceObj)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 更新副本状态(防止比异步更新慢)
	tmpResult, _ := replica.ReplicaMgm.FindById(obj.Id)
	if tmpResult.Status == contants.NOT_MOUNT {
		obj.Status = contants.MOUNTING
		obj, err = replica.ReplicaMgm.UpdateReplica(obj)
		if err != nil {
			logrus.Error(err)
			return
		}
		result = obj.TransformMap()
	}
	go execMount(&obj, &instanceObj, &login, &param.MountInfo, user)
	return
}

func execMount(replicaOrm *replica.Replica, instanceObj *instance.Instance, login *contants.LoginInfo,
	mountParam *view.MountInfo, user *auth.User) {
	time.Sleep(2 * time.Second)
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_MOUNT), "")
	// uuid对应为ceph上的名称
	poolName := replicaOrm.Pool.Uuid
	imageName := replicaOrm.Uuid
	// 映射设备
	var devPath string
	// 本机IP
	localIp := system.SysManage.GetLocalIp()
	localMountPoint := fmt.Sprintf("%s/%s/%s", contants.AppCfg.System.MountPath, poolName, imageName)
	sharePath := fmt.Sprintf("%s:%s", localIp, localMountPoint)
	shareIp := getShareIp(login, mountParam)
	logrus.Info("===============================挂载点===============================", localMountPoint)
	logrus.Info("===============================共享路径===============================", sharePath)
	// 此处声明err变量，后续不要再次声明，否则发生错误无法更新副本状态
	var err error
	// 闭包的方式清理资源，更新状态
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New(msg.ERROR_PROGRAM_CRASHES, msg.GetMsg(msg.ERROR_PROGRAM_CRASHES))
			logrus.Error(err)
		}
		if err != nil || e != nil {
			// 回滚
			logrus.Info("回滚副本==============>", replicaOrm.Name)
			_ = system.SysManage.RemoveNFS(localMountPoint, shareIp)
			_ = system.SysManage.UnMount(devPath, localMountPoint)
			_ = ceph.Client.UnMapRBDImage(poolName, imageName)
			detail := msg.GetOperation(msg.MOUNT_REPLICA_FAILED, replicaOrm.Name, err.Error())
			service.CreateLogRecord(msg.MOUNT_REPLICA, replicaOrm.Name, detail, user.Account, contants.LOG_FAILED)
			_ = replica.ReplicaMgm.UpdateReplicaStatus(replicaOrm.Id, contants.NOT_MOUNT)
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR,
				msg.GetOperation(msg.MOUNT_REPLICA_FAILED, replicaOrm.Name, "-"), err.Error())
			instance.InstanceMgm.UpdateInstanceStatus(instanceObj.Id, contants.MOUNT_FAILED)
			websockets.UpdateMsg(replicaOrm.CreateUser)
		}
	}()
	// 检查rbd映射, 如果已经映射则无需再次映射
	mapInfos, _ := ceph.Client.ShowMapRBDImage()
	logrus.Info(mapInfos)
	isMap := true
	for _, iter := range mapInfos {
		if iter.PoolName == poolName && iter.ImageName == imageName {
			isMap = false
			devPath = iter.DevPath
		}
	}
	if isMap {
		// 开始创建rbd映射
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
			msg.GetMsg(msg.INFO_START_RBD_MAP), "")
		devPath, err = ceph.Client.MapRBDImage(poolName, imageName)
		if err != nil {
			logrus.Error(err)
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR,
				msg.GetMsg(msg.INFO_RBD_MAP_FAILED), err.Error())
			return
		}
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
			msg.GetMsg(msg.INFO_RBD_MAP_SUCEESS), "")
	}
	// 格式化文件系统
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_START_FORMAT_REPLICA), "")
	err = system.SysManage.FormatXFS(devPath)
	if err != nil {
		logrus.Error(err)
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR,
			msg.GetMsg(msg.INFO_FORMAT_REPLICA_FAILED), err.Error())
		return
	}
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_FORMAT_REPLICA_SUCCESS), "")

	// 挂载本地文件系统
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_START_MOUNT_FILESYSTEM), "")
	err = system.SysManage.Mount(devPath, localMountPoint)
	if err != nil {
		logrus.Error(err)
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR,
			msg.GetMsg(msg.INFO_MOUNT_FILESYSTEM_FAILED), err.Error())
		return
	}
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_MOUNT_FILESYSTEM_SUCCESS), "")

	//（扩容后）同步文件系统
	err = system.SysManage.GrowfsXFS(devPath)
	if err != nil {
		logrus.Error(err)
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR,
			msg.GetMsg(msg.INFO_GROWFS_REPLICA_FAILED), err.Error())
		return
	}
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_GROWFS_REPLICA_SUCCESS), "")

	time.Sleep(time.Second * 3)

	//添加NFS共享
	err = system.SysManage.AddNFS(localMountPoint, shareIp)
	if err != nil {
		logrus.Error(err)
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR,
			msg.GetMsg(msg.INFO_ADD_NFS_SHARE_FAILED), err.Error())
		return
	}
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_ADD_NFS_SHARE_SUCCESS), "")
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_NFS_SHARE_PATH, localMountPoint), "")

	var remoteMountPoint string

	if mountParam.MountType == contants.APP_DEFAULT {
		if mountParam.IsExecScript {
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
				msg.GetMsg(msg.INFO_MOUNT_THE_WAY, "用户自定义脚本挂载!"), "")
			var scriptObj script.Script
			scriptObj, err = script.ScriptMgm.FindById(mountParam.ScriptId)
			if err != nil {
				instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
					msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_INFO, err.Error()), "")
				logrus.Error(err)
				return
			}
			scriptName := fmt.Sprintf("%s.sh", scriptObj.Uuid)
			scriptPath := fmt.Sprintf("%s/%s", contants.AppCfg.System.ScriptPath, scriptName)

			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
				msg.GetMsg(msg.INFO_START_RUNNING_SHELL_SCRIPT, scriptObj.Name), "")
			var output string
			output, err = system.SysManage.RunShellScript(scriptPath, "mount", sharePath)
			if err != nil {
				instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
					msg.GetMsg(msg.INFO_RUNNING_SHELL_SCRIPT_FAILED, scriptName), "")
				logrus.Error(err)
				return
			}
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
				msg.GetMsg(msg.INFO_RUNNING_SHELL_SCRIPT_SUCCESS, scriptName), "")
			logrus.Info(output)
			remoteMountPoint = "script"
		} else {
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
				msg.GetMsg(msg.INFO_MOUNT_THE_WAY, contants.MountTypeFlags[contants.APP_DEFAULT]), "")
			// 应用处理
			adapter := app.NewAdapter(mountParam.TargetType, login)
			remoteMountPoint, err = adapter.Mount(&sharePath, mountParam, user, instanceObj.Id)
			if err != nil {
				err = errors.New(msg.ERROR_APP_PROCESS_ERROR, msg.GetMsg(msg.ERROR_APP_PROCESS_ERROR, err.Error()))
				logrus.Error(err)
				return
			}
		}
	} else {
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
			msg.GetMsg(msg.INFO_MOUNT_THE_WAY, contants.MountTypeFlags[contants.APP_DEFAULT]), "")
		remoteMountPoint = contants.MountTypeFlags[contants.EXPORT_PATH]
	}

	// todo: 挂载完成[后面不允许有异常]
	// 创建挂载信息
	var mountInfo replica.MountInfo
	mountInfo.ReplicaId = replicaOrm.Id
	mountInfo.TargetType = instanceObj.TargetType
	mountInfo.TargetId = instanceObj.TargetId
	mountParamStr, _ := json.Marshal(mountParam)
	mountInfo.MountParam = string(mountParamStr)
	_, _ = replica.MountInfoMgm.CreateMountInfo(mountInfo)
	// 创建执行记录
	detail := msg.GetOperation(msg.MOUNT_REPLICA_SUCCESS, replicaOrm.Name)
	service.CreateLogRecord(msg.MOUNT_REPLICA, replicaOrm.Name, detail, user.Account, contants.LOG_SUCCESS)
	// 更新实例
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_MOUNT_POINT_NAME, remoteMountPoint), "")
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetOperation(msg.MOUNT_REPLICA_SUCCESS, replicaOrm.Name), "")
	instanceObj.MountPoint = remoteMountPoint
	instanceObj.Status = contants.MOUNT_SUCCESS
	_, _ = instance.InstanceMgm.UpdateInstance(*instanceObj)
	// 更新副本
	replicaOrm.Export = sharePath
	replicaOrm.Status = contants.MOUNTED
	_, _ = replica.ReplicaMgm.UpdateReplica(*replicaOrm)
	websockets.UpdateMsg(replicaOrm.CreateUser)

	return
}

func UnMountReplica(id uint, user *auth.User) (result map[string]interface{}, err error) {
	var obj replica.Replica
	// 闭包的方式记录操作日志
	defer func() {
		e := recover()
		if err != nil || e != nil {
			detail := msg.GetOperation(msg.UNMOUNT_REPLICA_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.UNMOUNT_REPLICA, obj.Name, detail, user.Account, contants.LOG_FAILED)
		}
	}()
	obj, err = replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 检查资源权限
	logrus.Info("-------------------检查资源权限---------------------", obj.Name)
	err = service.CheckResource(service.QUERY_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 获取挂载信息
	logrus.Info("-------------------获取挂载信息---------------------", obj.Name)
	mountInfo, getMountInfoErr := replica.MountInfoMgm.FindMountInfoByReplicaId(id)
	if getMountInfoErr != nil {
		logrus.Error(getMountInfoErr)
		//直接更新为未挂载状态
		logrus.Info("-------------------直接更新为未挂载状态---------------------", obj.Name)
		obj.Export = "-"
		obj.Status = contants.NOT_MOUNT
		_, _ = replica.ReplicaMgm.UpdateReplica(obj)
		websockets.UpdateMsg(obj.CreateUser)
		return
	}
	// 检查状态
	if obj.Status != contants.MOUNTED {
		err = errors.New(msg.ERROR_REPLICA_STATUS_IS_MOUNT, msg.GetMsg(msg.ERROR_REPLICA_STATUS_IS_MOUNT))
		logrus.Error(err)
		return
	}
	// 获取目标平台登录信息
	logrus.Info("-------------------获取目标平台登录信息---------------------", obj.Name)
	var login contants.LoginInfo
	if mountInfo.TargetType == contants.LINUX || mountInfo.TargetType == contants.WINDOWS {
		var hostOrm host.Host
		hostOrm, err = host.HostMgm.FindById(mountInfo.TargetId)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_INFO, msg.GetMsg(msg.ERROR_GET_HOST_INFO, err.Error()))
			logrus.Error(err)
			return
		}
		passWord, e := util.AesDecrypt(hostOrm.PassWord)
		if e != nil {
			err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
			logrus.Error(err)
			return
		}
		login.Ip = hostOrm.Ip
		login.Port = hostOrm.Port
		login.UserName = hostOrm.UserName
		login.PassWord = passWord
	} else {
		var platformOrm platform.Platform
		platformOrm, err = platform.PlatformMgm.FindById(mountInfo.TargetId)
		if err != nil {
			err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error()))
			logrus.Error(err)
			return
		}
		passWord, e := util.AesDecrypt(platformOrm.PassWord)
		if e != nil {
			err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
			logrus.Error(err)
			return
		}
		login.Ip = platformOrm.Ip
		login.Port = platformOrm.Port
		login.UserName = platformOrm.UserName
		login.PassWord = passWord
	}
	// 创建执行实例
	var instanceObj instance.Instance
	instanceObj.Type = contants.INSTANCE_UNMOUNT
	instanceObj.Status = contants.UMOUNT_RUNNING
	instanceObj.ReplicaId = obj.Id
	instanceObj.TargetType = mountInfo.TargetType
	instanceObj.TargetId = mountInfo.TargetId
	instanceObj.MountPoint = ""
	instanceObj, err = instance.InstanceMgm.CreateInstance(instanceObj)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 更新副本状态
	logrus.Info("-------------------更新副本状态---------------------", obj.Name)
	tmpResult, _ := replica.ReplicaMgm.FindById(obj.Id)
	if tmpResult.Status == contants.MOUNTED {
		obj.Status = contants.UNMOUNTING
		obj, err = replica.ReplicaMgm.UpdateReplica(obj)
		if err != nil {
			logrus.Error(err)
			return
		}
		result = obj.TransformMap()
	}

	go execUnMount(&obj, &instanceObj, &login, &mountInfo, user)
	return
}

func BatchUnMountReplica(param *view.BatchUnmountParam, user *auth.User) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	var res replica.Replica
	var err error
	for _, replicaId := range param.UnMountInfo {
		//先检查副本状态
		res, err = replica.ReplicaMgm.FindById(replicaId)
		if err != nil {
			logrus.Error(err)
			tmpMap := make(map[string]interface{})
			tmpMap[strconv.Itoa(int(replicaId))] = err.Error()
			result = append(result, tmpMap)
			continue
		}
		if res.Status != contants.MOUNTED {
			logrus.Error(err)
			tmpMap := make(map[string]interface{})
			tmpMap[strconv.Itoa(int(replicaId))] = "副本:" + strconv.Itoa(int(replicaId)) + "不是已挂载状态"
			result = append(result, tmpMap)
			continue
		}
		data, err := UnMountReplica(replicaId, user)
		if err != nil {
			if len(data) == 0 {
				tmpMap := make(map[string]interface{})
				tmpMap[strconv.Itoa(int(replicaId))] = err.Error()
				result = append(result, tmpMap)
			} else {
				data[strconv.Itoa(int(replicaId))] = err.Error()
				result = append(result, data)
			}

		}
	}
	return result, nil
}

func execUnMount(replicaOrm *replica.Replica, instanceObj *instance.Instance, login *contants.LoginInfo,
	mountInfo *replica.MountInfo, user *auth.User) {
	time.Sleep(2 * time.Second)
	// 此处声明err变量，后续不要再次声明，否则发生错误无法更新副本状态
	var err error
	// 闭包的方式更新失败状态
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New(msg.ERROR_PROGRAM_CRASHES, msg.GetMsg(msg.ERROR_PROGRAM_CRASHES))
			logrus.Error(err)
		}
		if err != nil || e != nil {
			detail := msg.GetOperation(msg.UNMOUNT_REPLICA_FAILED, replicaOrm.Name, err.Error())
			service.CreateLogRecord(msg.UNMOUNT_REPLICA, replicaOrm.Name, detail, user.Account, contants.LOG_FAILED)
			_ = replica.ReplicaMgm.UpdateReplicaStatus(replicaOrm.Id, contants.MOUNTED)
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR, msg.GetOperation(msg.UNMOUNT_REPLICA_FAILED, replicaOrm.Name, "-"), err.Error())
			instance.InstanceMgm.UpdateInstanceStatus(instanceObj.Id, contants.UMOUNT_FAILED)
			websockets.UpdateMsg(replicaOrm.CreateUser)
		}
	}()
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO,
		msg.GetMsg(msg.INFO_START_UNMOUNT), "")
	// uuid对应为ceph上的名称
	poolName := replicaOrm.Pool.Uuid
	imageName := replicaOrm.Uuid
	// 本地挂载路径
	localMountPoint := fmt.Sprintf("%s/%s/%s", contants.AppCfg.System.MountPath, poolName, imageName)
	// 本机IP
	localIp := system.SysManage.GetLocalIp()
	sharePath := fmt.Sprintf("%s:%s", localIp, localMountPoint)

	var mountParam view.MountInfo
	paramErr := json.Unmarshal([]byte(mountInfo.MountParam), &mountParam)
	if paramErr != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, paramErr.Error()))
		logrus.Error(err)
		return
	}
	if mountParam.MountType == contants.APP_DEFAULT {
		if mountParam.IsExecScript {
			scriptObj, scriptErr := script.ScriptMgm.FindById(mountParam.ScriptId)
			if scriptErr != nil {
				err = errors.New(msg.ERROR_GET_SCRIPT_FILE_INFO, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_INFO, scriptErr.Error()))
				instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_INFO, err.Error()), "")
				logrus.Error(err)
				return
			}
			scriptName := fmt.Sprintf("%s.sh", scriptObj.Uuid)
			scriptPath := fmt.Sprintf("%s/%s", contants.AppCfg.System.ScriptPath, scriptName)

			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_RUNNING_SHELL_SCRIPT, scriptName), "")
			output, runShellErr := system.SysManage.RunShellScript(scriptPath, "unmount", sharePath)
			if runShellErr != nil {
				err = errors.New(msg.INFO_RUNNING_SHELL_SCRIPT_FAILED, msg.GetMsg(msg.INFO_RUNNING_SHELL_SCRIPT_FAILED, runShellErr.Error()))
				instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_RUNNING_SHELL_SCRIPT_FAILED, scriptName), "")
				logrus.Error(err)
				return
			}
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_RUNNING_SHELL_SCRIPT_SUCCESS, scriptName), "")
			logrus.Info(output)
		} else {
			// todo: 应用处理
			adapter := app.NewAdapter(mountInfo.TargetType, login)
			unMountErr := adapter.UnMount(&sharePath, &mountParam, instanceObj.Id)
			if unMountErr != nil {
				err = errors.New(msg.ERROR_APP_PROCESS_ERROR, msg.GetMsg(msg.ERROR_APP_PROCESS_ERROR, unMountErr.Error()))
				logrus.Error(err)
				return
			}
		}
	}
	time.Sleep(time.Second * 3)
	// 映射设备
	var devPath string
	// 获取映射设备
	mapInfos, _ := ceph.Client.ShowMapRBDImage()
	logrus.Info(mapInfos)
	for _, iter := range mapInfos {
		if iter.PoolName == poolName && iter.ImageName == imageName {
			devPath = iter.DevPath
		}
	}
	if devPath == "" { //获取不到映射，跳过umount和unmap过程
		logrus.Error("获取不到映射，跳过umount和unmap过程")
		//return
	} else {
		// 取消NFS共享
		logrus.Info("-------------------取消NFS共享---------------------", replicaOrm.Name)
		shareIp := getShareIp(login, &mountParam)
		shareErr := system.SysManage.RemoveNFS(localMountPoint, shareIp)
		if shareErr != nil {
			err = errors.New(msg.INFO_REMOVE_NFS_SHARE_FAILED, msg.GetMsg(msg.INFO_REMOVE_NFS_SHARE_FAILED, shareErr.Error()))
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR, msg.GetMsg(msg.INFO_REMOVE_NFS_SHARE_FAILED), err.Error())
			logrus.Error(err)
			return
		}
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_REMOVE_NFS_SHARE_SUCCESS), "")
		time.Sleep(time.Second * 3)
		// 卸载本地文件系统
		logrus.Info("-------------------卸载本地文件系统---------------------", replicaOrm.Name)
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_UNMOUNT_FILESYSTEM), "")
		unMountErr := system.SysManage.UnMount(devPath, localMountPoint)
		if unMountErr != nil {
			// 回滚
			// _ = system.SysManage.AddNFS(localMountPoint, login.Ip)
			instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_ERROR, msg.GetMsg(msg.INFO_UNMOUNT_FILESYSTEM_FAILED), unMountErr.Error())
			logrus.Error(unMountErr)
			return
		}
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_UNMOUNT_FILESYSTEM_SUCCESS), "")
		time.Sleep(time.Second * 3)

		// 取消rbd映射
		logrus.Info("-------------------取消rbd映射---------------------", replicaOrm.Name)
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_RBD_UN_MAP), "")
		unmapErr := ceph.Client.UnMapRBDImage(poolName, imageName)
		if unmapErr != nil {
			logrus.Error(unmapErr)
			//return #取消rbd映射失败不影响解除挂载
		}
		instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, msg.GetMsg(msg.INFO_RBD_UN_MAP_SUCEESS), "")

		//FusionCompute 重新扫描一次
		if mountInfo.TargetType == contants.FUSIONCOMPUTE {
			refreshErr := fusioncompute.Refresh(cast.ToString(mountParam.AppConfig["fcHostUrn"]), *login)
			if refreshErr != nil {
				logrus.Error(refreshErr)
			}
		}

	}
	// todo: 卸载完成[后面不允许有异常]
	// 删除挂载信息
	logrus.Info("-------------------删除挂载信息---------------------", replicaOrm.Name)
	_ = replica.MountInfoMgm.DeleteMountInfoByReplicaId(replicaOrm.Id)
	// 创建执行记录
	logrus.Info("-------------------创建执行记录---------------------", replicaOrm.Name)
	detail := msg.GetOperation(msg.UNMOUNT_REPLICA_SUCCESS, replicaOrm.Name)
	service.CreateLogRecord(msg.UNMOUNT_REPLICA, replicaOrm.Name, detail, user.Account, contants.LOG_SUCCESS)
	// 更新实例
	logrus.Info("-------------------更新实例---------------------", replicaOrm.Name)
	info := msg.GetOperation(msg.UNMOUNT_REPLICA_SUCCESS, replicaOrm.Name)
	instance.InstanceLogMgm.PushInstanceLog(instanceObj.Id, contants.LOG_INFO, info, "")
	instanceObj.Status = contants.UMOUNT_SUCCESS
	_, _ = instance.InstanceMgm.UpdateInstance(*instanceObj)
	// 更新副本
	logrus.Info("-------------------更新副本状态---------------------", replicaOrm.Name)
	replicaOrm.Export = "-"
	replicaOrm.Status = contants.NOT_MOUNT
	_, _ = replica.ReplicaMgm.UpdateReplica(*replicaOrm)
	websockets.UpdateMsg(replicaOrm.CreateUser)
	return
}

func BatchMountReplica(param *view.BatchMountParam, user *auth.User) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for _, mountInfo := range param.MountInfos {
		mountParam := &view.MountParam{MountInfo: mountInfo}
		data, err := MountReplica(mountParam, user)
		if err != nil {
			data[strconv.Itoa(int(mountInfo.ReplicaId))] = err.Error()
		}
		result = append(result, data)
	}
	return result, nil
}

func GetInstances(id uint, queryParam *view.PageQueryParam) (int64, []map[string]interface{}, error) {
	totalNum, instances, err := instance.InstanceMgm.GetInstancesByReplicaId(id, queryParam.Index, queryParam.Count)
	if err != nil {
		err = errors.New(msg.ERROR_GET_LOG_RECORD_INFO, msg.GetMsg(msg.ERROR_GET_LOG_RECORD_INFO, err.Error()))
		logrus.Error(err)
		return 0, nil, err
	}
	var result []map[string]interface{}
	for _, iter := range instances {
		instanceMap := iter.TransformMap()
		result = append(result, instanceMap)
		if iter.TargetType == contants.LINUX || iter.TargetType == contants.WINDOWS {
			targetInfo, e := host.HostMgm.FindById(iter.TargetId)
			if e != nil {
				logrus.Error(e)
				continue
			}
			instanceMap["targetInfo"] = targetInfo.TransformMap()
		} else {
			targetInfo, e := platform.PlatformMgm.FindById(iter.TargetId)
			if e != nil {
				logrus.Error(e)
				continue
			}
			instanceMap["targetInfo"] = targetInfo.TransformMap()
		}
	}
	return totalNum, result, nil
}

func GetInstanceLogs(id uint) ([]map[string]interface{}, error) {
	details, err := instance.InstanceLogMgm.GetLogsByInstanceId(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_LOG_DETAIL_INFO, msg.GetMsg(msg.ERROR_GET_LOG_DETAIL_INFO, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	var result []map[string]interface{}
	for _, iter := range details {
		result = append(result, iter.TransformMap())
	}
	return result, err
}

func getShareIp(login *contants.LoginInfo, mountParam *view.MountInfo) string {
	if mountParam.MountType == contants.APP_DEFAULT {
		switch mountParam.TargetType {
		case contants.VMWARE:
			client, err := vsphere.NewClient(login)
			if err != nil {
				logrus.Error(err)
				return login.Ip
			}
			// 如果是vcenter需要共享给vcenter管理的esxi主机，否则无法创建Nas存储
			if client.IsvCenter() {
				/**
				var hostPath string
				if !strings.Contains(cast.ToString(mountParam.AppConfig["computeResource"]), "/Resources/") {
					// 添加隐藏路径
					hostPath = cast.ToString(mountParam.AppConfig["computeResource"])
				} else {
					hostPath = strings.Split(cast.ToString(mountParam.AppConfig["computeResource"]), "/Resources/")[0]
				}
				mountIp, err := client.GetHostNameByPath(hostPath)
				if err != nil {
					logrus.Error(err)
					return login.Ip
				}
				**/

				return cast.ToString(mountParam.AppConfig["vmwareHost"])
			}
			return login.Ip
		case contants.CAS:
			return login.Ip
		case contants.FUSIONCOMPUTE:
			return cast.ToString(mountParam.AppConfig["fcHostIp"])
		case contants.LINUX:
			return login.Ip
		case contants.WINDOWS:
			return login.Ip
		}
	}
	return login.Ip
}
