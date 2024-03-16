package vmware

import (
	"bufio"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/instance"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/internal/vsphere"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/vmware/govmomi/object"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Adapter struct {
	Login *contants.LoginInfo
}

func NewAdapter(login *contants.LoginInfo) (adapter *Adapter) {
	return &Adapter{Login: login}
}

func (a *Adapter) VerifyConfigInfo() (result *contants.VerifyResult, err error) {
	if a.Login == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return
	}
	client, err := vsphere.NewClient(a.Login)
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_PLATFORM, msg.GetMsg(msg.ERROR_VERIFY_PLATFORM, a.Login.Ip, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	var res contants.VerifyResult
	res.Version = client.GetVersion()
	return &res, nil
}

func (a *Adapter) Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error) {
	// 设置挂载详细参数
	param, err := setParam(sharePath, mountParam, true)
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
	// 登录VMware平台
	if a.Login == nil {
		err := errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return "", err
	}
	client, err := vsphere.NewClient(a.Login)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	// 检查NFS存储是否存在, 不存在则创建
	isExist := false
	storeNameArr, err := client.GetNasDatastore()
	if err != nil {
		err = errors.New(msg.ERROR_CHEACK_NFS_STORE, msg.GetMsg(msg.ERROR_CHEACK_NFS_STORE, err.Error()))
		logrus.Error(err)
		return "", err
	}
	for _, iter := range storeNameArr {
		if iter == param.StoreName {
			isExist = true
		}
	}
	if !isExist {
		// 创建NFS存储
		var storeInfo vsphere.NasDataStoreParam
		storeInfo.StoreType = param.StoreType
		storeInfo.AccessMode = param.StoreMode  // [readWrite, readOnly]
		storeInfo.LocalPath = param.StoreName   // 存储名称
		storeInfo.RemoteHost = param.RemoteHost // 服务器IP
		storeInfo.RemotePath = param.RemotePath // 共享路径
		//storeInfo.HostPath = param.HostPath     // 主机路径
		storeInfo.HostName = param.HostName // 主机名称
		err = client.CreateNasDatastore(&storeInfo)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
	}

	if param.IsRegisterVm {

		//先检查是否有同ip、登录名称密码的虚拟机（不一定在同一个平台）
		loginInfo := contants.LoginInfo{
			Ip:       param.VmAddr,
			Port:     22,
			UserName: param.VmUsername,
			PassWord: param.VmPassword,
		}
		_, verifyErr := VerifyHostInfo(&loginInfo, param.Os)
		if verifyErr == nil { //没报错,说明登录成功，存在ip冲突的主机
			err = errors.New(msg.ERROR_REGISTER_VM_FAILED, msg.GetMsg(msg.ERROR_REGISTER_VM_FAILED, fmt.Sprintf("检测到存在ip冲突:【%s】，请确认", param.VmAddr)))
			logrus.Error(err)
			return "", err
		}

		// 注册虚拟机
		var registerParam vsphere.RegisterVmParam
		registerParam.HostPath = param.HostPath
		registerParam.ResourcePoolPath = param.ResourcePoolPath
		registerParam.FolderPath = param.LocationPath
		registerParam.VmName = param.VmName
		registerParam.VmxPath = param.VmxPath

		err = client.RegisterVm(&registerParam)
		if err != nil {
			logrus.Error(err)
			_ = client.RemoveNasDatastore(param.StoreName)
			return "", err
		}
		// 配置虚拟机
		var ipAddr vsphere.IpAddr
		ipAddr.Ip = param.VmAddr
		ipAddr.Netmask = param.VmNetmask
		ipAddr.Gateway = param.VmGateway
		ipAddr.Hostname = "tangula"

		logrus.Info(fmt.Sprintf("创建的虚拟机名称=============>%s", param.VmName))
		logrus.Info(fmt.Sprintf("虚拟机配置信息=============>%s", ipAddr))

		vms := NewVmWare(a.Login.Ip, a.Login.UserName, a.Login.PassWord) //govmomi登录vmware平台
		vmList, _, _ := vms.GetAllVmClient()
		for _, vm := range vmList {
			if vm.Name == param.VmName {
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_VM), "")
				v := object.NewVirtualMachine(vms.Client.Client, vm.VM)

				//配置ip
				if param.IsSetIp {
					setIpErr := vms.SetIP(v, ipAddr, param.Os)
					if setIpErr != nil {
						logrus.Error("设置ip报错：", setIpErr)
						instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_SET_IPADDR), setIpErr.Error())
					}
				}

				//打开电源
				if param.IsAutoPowerOn {
					instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_POWERON), "")
					poweronErr := vms.PowerOn(v)
					if poweronErr != nil {
						logrus.Error("打开电源报错：", poweronErr)
						instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_POWER_ON_VM_FAILED), poweronErr.Error())
					}
				}

				//保存虚拟机信息
				/**
				obj, findErr := replica.ReplicaMgm.FindById(mountParam.ReplicaId)
				if findErr != nil {
					err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, findErr.Error()))
					logrus.Error("查找副本出错:", err)
					instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR,
						msg.GetMsg(9999, "查找副本出错:"), findErr.Error())
				}

				initError := SaveHost(param.VmUsername, param.VmPassword, param.VmName, param.Os, ipAddr.Ip, 22, obj.Name, user.Account, vms)
				if initError != nil {
					logrus.Error("保存虚拟机信息出错:", initError)
					instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR,
						msg.GetMsg(msg.ERROR_SAVEHOST), initError.Error())
				}
				instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_SAVEHOST), "")
				**/
				break
			}
		}
	}
	return param.ResourcePoolPath, nil
}

// 卸载时谨慎处理逻辑，防止副本一直卸载不掉的情况
func (a *Adapter) UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) error {
	// 设置参数
	param, err := setParam(sharePath, mountParam, false)
	if err != nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, err.Error()))
		logrus.Error(err)
		return err
	}
	if param == nil {
		err = errors.New(msg.ERROR_ANALYTIC_MOUNT_PARAMETER,
			msg.GetMsg(msg.ERROR_ANALYTIC_MOUNT_PARAMETER, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return err
	}
	// 登录VMware平台
	if a.Login == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return err
	}
	client, err := vsphere.NewClient(a.Login)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if param.IsRegisterVm {
		if param.VmName != "" {
			//关闭电源
			vms := NewVmWare(a.Login.Ip, a.Login.UserName, a.Login.PassWord) //govmomi登录vmware平台
			vmList, _, _ := vms.GetAllVmClient()
			for _, vm := range vmList {
				if vm.Name == param.VmName {
					logrus.Info(fmt.Sprintf("虚拟机电源:%s", vm.PowerState))
					if vm.PowerState == "poweredOn" {
						logrus.Info("开始关闭虚拟机电源:", vm.Name)
						instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_POWEROFF, param.VmName), "")
						v := object.NewVirtualMachine(vms.Client.Client, vm.VM)
						powerOffErr := vms.PowerOff(v)
						if powerOffErr != nil {
							logrus.Error("关闭虚拟机电源报错：", powerOffErr)
							instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_POWER_OFF_VM_FAILED), powerOffErr.Error())
							return powerOffErr
						}
					}
					time.Sleep(time.Second * 2)
					instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_UN_REGISTER_VM_BY_VMWARE, param.VmName), "")
					err = client.UnRegisterVm(vm.Uuid)
					if err != nil {
						logrus.Error(err)
						if strings.Contains(cast.ToString(err), "查找对象发生错误") {
							logrus.Info("虚拟化平台上查询不到虚拟机信息，继续卸载动作")
						} else {
							instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_ERROR, msg.GetMsg(msg.ERROR_UN_REGISTER_VM_FAILED), err.Error())
							return err
						}
					}
					time.Sleep(time.Second * 2)
					break
				}
			}
		} else {
			logrus.Info("无法获取虚拟机信息，继续卸载动作")
		}
	}
	// 检查NFS存储是否存在, 存在则卸载
	isExist := false
	storeNameArr, err := client.GetNasDatastore()
	if err != nil {
		err = errors.New(msg.ERROR_CHEACK_NFS_STORE, msg.GetMsg(msg.ERROR_CHEACK_NFS_STORE, err.Error()))
		logrus.Error(err)
		return err
	}
	for _, iter := range storeNameArr {
		if iter == param.StoreName {
			isExist = true
		}
	}
	if isExist {
		// 卸载NFS存储
		logrus.Info(fmt.Sprintf("开始卸载存储【%s】", param.StoreName))
		instance.InstanceLogMgm.PushInstanceLog(instanceId, contants.LOG_INFO, msg.GetMsg(msg.INFO_START_DELETE_NFS_STORE, param.StoreName), "")
		err = client.RemoveNasDatastore(param.StoreName)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func setParam(sharePath *string, mountParam *view.MountInfo, isMount bool) (*MountDetailInfo, error) {
	var detailInfo MountDetailInfo
	// 设置位置路径，主机路径，资源池路径
	if !strings.Contains(cast.ToString(mountParam.AppConfig["locationPath"]), "/vm/") {
		// 添加隐藏路径
		detailInfo.LocationPath = cast.ToString(mountParam.AppConfig["locationPath"]) + "/vm"
	} else {
		detailInfo.LocationPath = cast.ToString(mountParam.AppConfig["locationPath"])
	}

	if !strings.Contains(cast.ToString(mountParam.AppConfig["computeResource"]), "/Resources/") {
		// 添加隐藏路径
		detailInfo.HostPath = cast.ToString(mountParam.AppConfig["computeResource"])
		detailInfo.ResourcePoolPath = cast.ToString(mountParam.AppConfig["computeResource"]) + "/Resources"
	} else {
		detailInfo.HostPath = strings.Split(cast.ToString(mountParam.AppConfig["computeResource"]), "/Resources/")[0]
		detailInfo.ResourcePoolPath = cast.ToString(mountParam.AppConfig["computeResource"])
	}

	detailInfo.HostName = cast.ToString(mountParam.AppConfig["vmwareHost"])

	// 设置控制参数
	detailInfo.IsRegisterVm = cast.ToBool(mountParam.AppConfig["isRegisterVM"])
	// 设置存储信息
	start1 := strings.LastIndex(*sharePath, ":")
	remoteHost := string([]byte(*sharePath)[0:start1])  // 截取IP
	remotePath := string([]byte(*sharePath)[start1+1:]) // 截取路径：mnt/tangula/pool_uuid/image_uuid
	start2 := strings.LastIndex(remotePath, "/")
	detailInfo.StoreName = fmt.Sprintf("tangula_%s", string([]byte(remotePath)[start2+1:])) // 截取image_uuid
	detailInfo.StoreType = "NFS"
	detailInfo.StoreMode = "readWrite"
	detailInfo.RemoteHost = remoteHost
	detailInfo.RemotePath = remotePath

	// 设置虚拟机信息
	if detailInfo.IsRegisterVm {

		if isMount { //挂载时检查磁盘文件
			checkerr := checkFolderNum(detailInfo.RemotePath)
			if checkerr != nil {
				return nil, checkerr
			}
		}
		vmBaseInfo, err := getVmInfo(detailInfo.RemotePath)
		if err != nil {
			if isMount { //卸载场景忽视报错
				err = errors.New(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()))
				logrus.Error(err)
				return nil, err
			}

		} else {
			// vmPath路径转换为平台能够识别的路径
			temp := strings.Split(vmBaseInfo.VmxPath, detailInfo.RemotePath+"/")[1]
			datastorePath := fmt.Sprintf("[%s] %s", detailInfo.StoreName, temp)

			detailInfo.VmxPath = datastorePath
			detailInfo.VmUuid = vmBaseInfo.VmUuid

			//设置虚拟机操作系统
			detailInfo.Os = cast.ToString(mountParam.AppConfig["os"])

			//设置虚拟机名称
			if cast.ToString(mountParam.AppConfig["vmName"]) != "" {
				detailInfo.VmName = cast.ToString(mountParam.AppConfig["vmName"])
			} else { //取副本名称
				replicaObj, findErr := replica.ReplicaMgm.FindById(mountParam.ReplicaId)
				if findErr != nil {
					err = errors.New(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, "查询副本失败："+findErr.Error()))
					logrus.Error(err)
					return nil, err
				}
				detailInfo.VmName = replicaObj.Name

			}

			// 设置虚拟机登录用户名密码
			detailInfo.VmUsername = cast.ToString(mountParam.AppConfig["username"])
			detailInfo.VmPassword = cast.ToString(mountParam.AppConfig["password"])

			// 设置虚拟机IP
			detailInfo.IsSetIp = cast.ToBool(mountParam.AppConfig["isSetIp"])
			detailInfo.VmAddr = cast.ToString(mountParam.AppConfig["addr"])
			detailInfo.VmNetmask = cast.ToString(mountParam.AppConfig["netMask"])
			detailInfo.VmGateway = cast.ToString(mountParam.AppConfig["gateWay"])

			//是否开机
			detailInfo.IsAutoPowerOn = cast.ToBool(mountParam.AppConfig["powerOn"])

		}
	}

	logrus.Info("setParam:", util.SerializeToJson(detailInfo))
	util.FormatJson(detailInfo)
	return &detailInfo, nil
}

func checkFolderNum(path string) (err error) {
	dirSub, readErr := ioutil.ReadDir(path)
	if readErr != nil {
		return readErr
	}
	folderNums := 0
	for _, iter := range dirSub {
		if iter.IsDir() {
			folderNums += 1
		}
	}
	if folderNums != 1 {
		err = errors.New(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, "副本中磁盘文件目录为空或多于1个"))
	}
	logrus.Info(fmt.Sprintf("副本中磁盘文件目录数量:%d", folderNums))
	return err
}
func getVmxPath(path string) (string, error) {
	dirSub, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	//此处不会抛错
	for _, iter := range dirSub {
		if iter.IsDir() {
			vmxPath, err := getVmxPath(path + "/" + iter.Name())
			if err != nil {
				return "", err
			}
			if vmxPath != "" {
				return vmxPath, nil
			}
		} else {
			if strings.Contains(iter.Name(), ".vmx") {
				return path + "/" + iter.Name(), nil
			}
		}
	}
	return "", nil
}

func getVmInfo(sharePath string) (*vmBaseInfo, error) {
	var vmBaseInfo vmBaseInfo
	vmxPath, err := getVmxPath(sharePath)

	if err != nil {
		err = errors.New(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	if vmxPath == "" {
		err = errors.New(msg.ERROR_GET_VM_BASE_INFO,
			msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_VMX_FILE_NOT_EXIST)))
		logrus.Error(err)
		return nil, err
	}
	vmBaseInfo.VmxPath = vmxPath
	fp, err := os.Open(vmBaseInfo.VmxPath)
	if err != nil {
		err = errors.New(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		if e := fp.Close(); e != nil {
			logrus.Error(err)
		}
	}()
	rd := bufio.NewReader(fp)
	for {
		lineByte, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			err = errors.New(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		// 解析虚拟机uuid
		// "42 28 8f 5f e2 59 36 c3-b9 ac 2b b3 8e ca 4c 8b" ===> "42288f5f-e259-36c3-b9ac-2bb38eca4c8b"
		if strings.Contains(string(lineByte), "vc.uuid") {
			tempSlice := strings.Split(string(lineByte), "\"")
			if len(tempSlice) != 3 {
				err = errors.New(msg.ERROR_GET_VM_BASE_INFO,
					msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_UUID_BY_VMX)))
				logrus.Error(err)
				return nil, err
			}
			tempSlice = strings.Split(tempSlice[1], "-")
			if len(tempSlice) != 2 {
				err = errors.New(msg.ERROR_GET_VM_BASE_INFO,
					msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_UUID_BY_VMX)))
				logrus.Error(err)
				return nil, err
			}
			tempSlice = strings.Fields(tempSlice[0] + tempSlice[1])
			if len(tempSlice) != 15 {
				err = errors.New(msg.ERROR_GET_VM_BASE_INFO,
					msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_UUID_BY_VMX)))
				logrus.Error(err)
				return nil, err
			}
			var uuidStr string
			for _, iter := range tempSlice {
				uuidStr += iter
			}
			if len(uuidStr) != 32 {
				err = errors.New(msg.ERROR_GET_VM_BASE_INFO,
					msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_UUID_BY_VMX)))
				logrus.Error(err)
				return nil, err
			}
			vmBaseInfo.VmUuid = fmt.Sprintf("%s-%s-%s-%s-%s", uuidStr[0:8], uuidStr[8:12],
				uuidStr[12:16], uuidStr[16:20], uuidStr[20:])
		}
		// 解析虚拟机名称
		/**
		if strings.Contains(string(lineByte), "displayName") {
			temp := strings.Split(string(lineByte), "\"")
			if len(temp) != 3 {
				err = errors.New(msg.ERROR_GET_VM_BASE_INFO,
					msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_NAME_BY_VMX)))
				logrus.Error(err)
				return nil, err
			}
			vmBaseInfo.VmName = temp[1]
		}
		**/
	}
	if vmBaseInfo.VmUuid == "" {
		err = errors.New(msg.ERROR_GET_VM_BASE_INFO,
			msg.GetMsg(msg.ERROR_GET_VM_BASE_INFO, msg.GetMsg(msg.ERROR_GET_VM_UUID_BY_VMX)))
		logrus.Error(err)
		return nil, err
	}
	logrus.Info(fmt.Sprintf("vmBaseInfo:%+v\n", vmBaseInfo))
	return &vmBaseInfo, nil
}
