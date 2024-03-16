package cas

import (
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"strings"
	"time"
)

// 挂载详细信息
type MountDetailInfo struct {
	HostId       string
	HostName     string
	VsId         uint
	ClusterId    uint
	HostpoolId   uint
	StorePath    string
	StoreName    string
	RemoteHost   string
	RemotePath   string
	IsRegisterVM bool
	IsCrtByJson  bool
	Title        string
	PoolName     string
}

var crtCasReqData = map[string]interface{}{
	"name":          "",       //虚拟机名称。类型：String。要求非空。只允许包含字母、数字，“-_.” 。名称的首字符不能为减号及点。长度不超过64。
	"title":         "",       //虚拟机显示名称。类型：String。若为空，将用name的值赋值给title。只允许输入汉字，字母，数字，减号，下划线，空格及点，并且不允许全部为空格长度不超过64。
	"hostId":        "",       // 物理机ID。类型：Long。与clusterId至少填其一，若都填写请保证hostId在clusterId管辖下。若填写，须填写CVM中已有已有的物理机ID。
	"clusterId":     1,        //集群ID。类型：Long。与hostId至少填其一。若填写，须填写CVM中已有已有的clusterId。若填写了hostId，将根据hostId所在集群Id进行值覆盖。
	"system":        1,        //虚拟机安装的操作系统
	"osVersion":     "",       // 虚拟机安装的操作系统版本
	"autoMigrate":   0,        //是否允许虚拟机自动迁移
	"memory":        2048,     // 内存大小
	"memoryInit":    2,        // 内存值。类型：Double。必须。与memoryUnit联合使用。单位为由memoryUnit决定,换算成MB后，必须与memory指定的值相同。取值512MB-主机允许的最大内存。
	"memoryUnit":    "GB",     // 内存单位。类型：String。必须。与memoryInit联合使用。取值：MB、GB。
	"memoryBacking": 0,        //内存预留百分比, 范围：0--100即(0%--100%),0表示不预留，100表示全部预留。
	"autoMem":       0,        // 写入domainxml文件中。是否开启内存气球功能。 0-关闭， 1-开启。
	"hugepage":      false,    // 是否开启大页。
	"cpuSockets":    1,        // 虚拟机CPU socket个数
	"cpuCores":      1,        //虚拟机CPU核数
	"formatEnable":  1,        // 是否可用。设置为1
	"maxCpuSocket":  2,        // cpu数量最大值。
	"cpuMode":       "custom", // CPU工作模式。类型：String。必须。可能的取值：custom（兼容模式）、host-model（主机匹配模式）、host-passthrough（主机直通模式）
	"cpuShares":     512,      // 虚拟机在物理机上CPU调度优先级。 枚举值：1024：高，512：中，256：低
	"blkiotune":     300,      // 虚拟机I/O优先级，枚举值：500：高，300：中，200：低。
	"osBit":         "x86_64", //虚拟机BIT。体系结构：64位和32位。取值： x86_64 x86。
	"storage":       "",       //存储信息
	"network":       "",       //网络信息
}

func CrtStoragePoolParams(casReq *MountDetailInfo) map[string]interface{} {
	reqData := map[string]interface{}{
		"hostId":      casReq.HostId,
		"name":        casReq.StoreName,
		"title":       casReq.StoreName,
		"type":        "netfs",
		"path":        "/vms/" + casReq.StoreName,
		"hostIp":      casReq.RemoteHost,
		"remoteDir":   casReq.RemotePath,
		"rsFsLunInfo": []map[string]interface{}{},
	}
	logrus.Info(reqData)
	return reqData
}

func SetCasVmParamsByJson(casReq *MountDetailInfo, jsonMap util.Casjson) (reqData map[string]interface{}, err error) {

	dist_data_list := []map[string]interface{}{}
	for _, diskInfo := range jsonMap.DiskList {
		logrus.Info("diskInfo===>:", diskInfo)
		splitList := strings.Split(cast.ToString(diskInfo["path"]), "/")
		tmpPath := splitList[len(splitList)-1]

		var tmpDiskmap map[string]interface{}
		tmpDiskmap = make(map[string]interface{})
		tmpDiskmap["assignType"] = 0
		tmpDiskmap["device"] = "disk"
		tmpDiskmap["targetBus"] = "virtio"
		tmpDiskmap["type"] = "file"
		tmpDiskmap["driveType"] = "qcow2"
		tmpDiskmap["mode"] = cast.ToInt(diskInfo["mode"])
		tmpDiskmap["capacity"] = cast.ToInt(diskInfo["size"])
		tmpDiskmap["clusterSize"] = cast.ToInt(diskInfo["clusterSize"])
		tmpDiskmap["poolName"] = casReq.StoreName
		tmpDiskmap["storeFile"] = "/vms/" + casReq.StoreName + "/" + tmpPath

		dist_data_list = append(dist_data_list, tmpDiskmap)
		logrus.Info("存储信息--------------------:", tmpDiskmap)
	}

	networkList := []map[string]interface{}{}
	for _, networkInfo := range jsonMap.NetworkList {
		logrus.Info("networkInfo===>:", networkInfo)

		var tmpNetworkMap map[string]interface{}
		tmpNetworkMap = make(map[string]interface{})
		tmpNetworkMap["vsId"] = cast.ToInt(networkInfo["vsId"])
		tmpNetworkMap["vsName"] = "vswitch0"
		tmpNetworkMap["profileId"] = 1
		tmpNetworkMap["isKernelAccelerated"] = 1
		tmpNetworkMap["deviceModel"] = "virtio"
		networkList = append(networkList, tmpNetworkMap)
		logrus.Info("网卡信息-------------------:", tmpNetworkMap)
	}

	crtCasReqData["name"] = jsonMap.Name
	crtCasReqData["title"] = jsonMap.Title
	crtCasReqData["hostId"] = casReq.HostId
	crtCasReqData["storage"] = dist_data_list
	crtCasReqData["network"] = networkList
	logrus.Info(crtCasReqData)

	return crtCasReqData, err
}

func SetCasVmParams(title string, hostId string, vsId uint, diskInfoList []view.DiskInfo) (reqData map[string]interface{}, err error) {

	if len(diskInfoList) == 0 {
		err = errors.New(msg.ERRIR_CREATE_CAS, msg.GetMsg(msg.ERRIR_CREATE_CAS, title, "磁盘列表为空!"))
		return nil, err
	}
	dist_data_list := []map[string]interface{}{}
	vmName := diskInfoList[0].Name
	currentTime := time.Now()
	currentTime_str := currentTime.Format("20060102150405")
	for _, disk := range diskInfoList {
		var diskInfo map[string]interface{}
		diskInfo = make(map[string]interface{})
		diskInfo["assignType"] = 0                                        //分配类型:0指定 1动态分配。
		diskInfo["device"] = "disk"                                       // 存储类型，枚举值：disk：硬盘， cdrom：光驱， floppy：软驱
		diskInfo["targetBus"] = "virtio"                                  //存储设备总线类型，枚举值：ide：IDE硬盘， scsi：SCSI硬盘， usb：USB硬盘;，virtio：Virtio
		diskInfo["type"] = "file"                                         // 文件类型。枚举值：file 文件， block 块设备
		diskInfo["driveType"] = "qcow2"                                   // 存储卷类型。必须和使用的存储卷的类型保持一致，枚举值：qcow2，raw
		diskInfo["cacheType"] = "directsync"                              // 磁盘缓存方式。可取值：directsync表示直接读写、writeback表示二级虚拟机缓存、writethroug表示一级物理缓存、none表示一级虚拟缓存。
		diskInfo["storeFile"] = "/vms/" + disk.PoolName + "/" + disk.Name //存储文件路径

		dist_data_list = append(dist_data_list, diskInfo)
	}

	vswitchList := []map[string]interface{}{
		{
			"vsId":                vsId,
			"vsName":              "vswitch0",
			"profileId":           1,
			"isKernelAccelerated": 1,
			"deviceModel":         "virtio",
		},
	}
	crtCasReqData["name"] = vmName + "-" + currentTime_str
	crtCasReqData["title"] = title
	crtCasReqData["hostId"] = hostId
	crtCasReqData["storage"] = dist_data_list
	crtCasReqData["network"] = vswitchList
	logrus.Info("创建CAS虚拟机请求参数:", util.SerializeToJson(crtCasReqData))

	return crtCasReqData, err
}
