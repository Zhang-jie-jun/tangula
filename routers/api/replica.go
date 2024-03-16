package api

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 获取已支持副本镜像类型
// @Description 【登录权限】获取已支持虚副本镜像类型[1001.File, 1002.ISO, 1003.CUSTOME, 1004.VMware, 1005.Cas, 1006.HCS, 1007.OpenStack]。
// @Tags 副本管理
// @Accept json
// @Produce  json
// @Router /store_pool/replica/support [get]
// @Success 200 {object} Result{response=map[int]string}
func GetReplicaSupportType(context *gin.Context) {
	ResponseSuccess(context, contants.ImageFlags)
}

// @Summary 获取已支持副本挂载类型
// @Description 【登录权限】获取已支持副本挂载类型[1.应用默认挂载， 2.仅导出共享路径]。
// @Tags 副本管理
// @Accept json
// @Produce  json
// @Router /store_pool/replica/mount_type [get]
// @Success 200 {object} Result{response=map[int]string}
func GetReplicaMountType(context *gin.Context) {
	ResponseSuccess(context, contants.MountTypeFlags)
}

// @Summary 创建副本
// @Description 【登录权限】创建一个新副本。
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param name query string true "副本名称"
// @Param type query int true "副本类型[1001.File, 1002.ISO, 1003.CUSTOME, 1004.VMware, 1005.Cas, 1006.HCS, 1007.OpenStack]"
// @Param desc query string false "描述" default("")
// @Param size query int true "副本大小[单位：GB]"
// @Param storePoolId query string true "所在存储池uuid"
// @Router /store_pool/replica [post]
// @Success 200 {object} Result{response=replica.Replica}
func CreateReplica(context *gin.Context) {
	var param view.ReplicaCreateParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.CreateReplica(&param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 编辑指定副本
// @Description 【登录权限】编辑指定副本。
// @Tags 副本管理
// @Accept json
// @Produce  json
// @Param id path int true "目标副本ID"
// @Param desc query string false "描述" default("")
// @Param size query int false "副本大小[单位：GB]" default(0)
// @Router /store_pool/replica/{id} [put]
// @Success 200 {object} Result{response=replica.Replica}
func EditReplica(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var param view.ReplicaEditParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.EditReplica(idparam.Id, &param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 获取所有副本信息
// @Description 【登录权限】获取所有副本信息，普通用户与系统管理员只能获取自身创建的副本信息，超级管理员可以获取平台上所有的副本信息。
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param status query int false "副本状态过滤参数[0.不过滤，1024.未挂载，2048.挂载中，4096.已挂载"，8192.卸载中] default(0)
// @Param dataType query int false "镜像类型过滤参数[0.不过滤类型，1001.FILE，1002.ISO，1003.CUSTOME，1004.VMWAREVM，1005.CASVM，1006.HCSVM，1007.OPENSTACKVM]" default(0)
// @Param filter query string false "名称搜索过滤参数" default("")
// @Router /store_pool/replica [get]
// @Success 200 {object} Result{response=Paging{data=[]replica.Replica}}
func GetReplica(context *gin.Context) {
	var queryParam view.ReplicaQueryParam
	if err := context.ShouldBind(&queryParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	// 赋默认值
	if queryParam.Count == 0 {
		queryParam.Count = 15
	}
	totalNum, data, err := storage.GetReplica(&queryParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 获取指定副本
// @Description 【登录权限】获取指定副本信息，普通用户与系统管理员只能获取自身创建的副本信息，超级管理员可以获取平台上所有的副本信息。
// @Tags 镜像管理
// @Accept json
// @Produce json
// @Param id path int true "目标副本ID"
// @Router /store_pool/replica/{id} [get]
// @Success 200 {object} Result{response=replica.Replica}
func GetReplicaById(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.GetReplicaById(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 删除指定副本
// @Description 【登录权限】删除指定副本，普通用户与系统管理员只能删除属于自身的副本，超级管理员可删除所有副本。
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param id path int true "目标副本ID"
// @Router /store_pool/replica/{id} [delete]
// @Success 200 {object} Result
func DeleteReplica(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	err = storage.DeleteReplica(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

func DeleteReplicaBatch(context *gin.Context) {
	var idparam view.BatchDeleteReplicaParam
	if err := context.ShouldBind(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.BatchDeleteReplica(&idparam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 生成镜像
// @Description 【登录权限】根据所选副本生成一个新的镜像。
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param id path int true "目标副本ID"
// @Param name query string true "镜像名称"
// @Param desc query string false "镜像描述" default("")
// @Router /store_pool/replica/{id}/image [post]
// @Success 200 {object} Result{response=image.Image}
func CreateImageByReplica(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var param view.ImageCreateParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	err = storage.CreateImageByReplica(idparam.Id, &param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 创建快照
// @Description 【登录权限】转换指定副本为镜像。
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param id path int true "目标副本ID"
// @Param name query string true "快照名称"
// @Param desc query string false "快照描述" default("")
// @Router /store_pool/replica/{id}/snapshot [post]
// @Success 200 {object} Result{response=snapshot.Snapshot}
func CreateSnapshotByReplica(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var param view.ImageCreateParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.CreateSnapshotByReplica(idparam.Id, &param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 获取指定副本所有快照信息
// @Description 【登录权限】获取指定副本所有快照信息
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param id path int true "目标副本ID"
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param filter query string false "名称过滤参数" default("")
// @Router /store_pool/replica/{id}/snapshot [get]
// @Success 200 {object} Result{response=Paging{data=[]snapshot.Snapshot}}
func GetSnapshotByReplica(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var queryParam view.QueryParam
	if err := context.ShouldBind(&queryParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	// 赋默认值
	if queryParam.Count == 0 {
		queryParam.Count = 15
	}
	totalNum, data, err := storage.GetSnapshotByReplica(idparam.Id, &queryParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 删除指定副本所有快照
// @Description 【登录权限】删除指定副本所有快照
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param id path int true "目标副本ID"
// @Router /store_pool/replica/{id}/snapshot [delete]
// @Success 200 {object} Result
func DeleteSnapshotByReplica(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	err = storage.DeleteSnapshotByReplica(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 挂载副本
// @Description 【登录权限】挂载副本到目标平台。
// @Description
// @Description 【emphasis】 参数AppConfig需要根据应用区分，组织不同形式的[key:value]字符串。
// @Description ----Linux平台：
// @Description --------【可选参数:string】mountPoint：目标挂载位置
// @Description ----VMware平台：
// @Description --------【必填参数:string】locationPath：挂载清单位置
// @Description --------【必填参数:string】computeResource：挂载计算资源
// @Description --------【必填参数:bool】isRegisterVM：是否注册虚拟机
// @Description --------【选填参数:string】vmName：指定虚拟机注册名称
// @Description --------【选填参数:bool】powerOn：虚拟机注册后是否打开电源
// @Description --------【选填参数:string】hostName：指定虚拟机操作系统主机名
// @Description --------【选填参数:string】addr：指定虚拟机操作系统IP地址
// @Description --------【选填参数:string】netMask：指定虚拟机操作系统网关
// @Description --------【选填参数:string】gateWay：指定虚拟机操作系统掩码
// @Description ----CAS平台：
// @Description --------【必填参数:string】name：虚拟机名称
// @Description --------【必填参数:int】hostId：主机id
// @Description --------【必填参数:int】system：操作系统编号（可以为1）
// @Description --------【选填参数:string】osVersion：操作系统名称，可以为空字符串
// @Description --------【选填参数:int】vsId：虚拟机交换机id，一般为1
// @Description --------【选填参数:bool】isRegisterVM：是否新建虚拟机
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param  body body view.MountParam true "挂载参数"
// @Router /store_pool/replica/mount [post]
// @Success 200 {object} Result{response=replica.Replica}
func MountReplica(context *gin.Context) {
	var param view.MountParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.MountReplica(&param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

func UnMountReplica(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Info(idparam)
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.UnMountReplica(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 批量挂载副本
// @Description 【登录权限】批量挂载副本到目标平台。
// @Description
// @Description 【emphasis】 参数AppConfig需要根据应用区分，组织不同形式的[key:value]字符串。
// @Description ----Linux平台：
// @Description --------【可选参数:string】mountPoint：目标挂载位置
// @Description ----VMware平台：
// @Description --------【必填参数:string】locationPath：挂载清单位置
// @Description --------【必填参数:string】computeResource：挂载计算资源
// @Description --------【必填参数:bool】isRegisterVM：是否注册虚拟机
// @Description --------【选填参数:string】vmName：指定虚拟机注册名称
// @Description --------【选填参数:bool】powerOn：虚拟机注册后是否打开电源
// @Description --------【选填参数:string】hostName：指定虚拟机操作系统主机名
// @Description --------【选填参数:string】addr：指定虚拟机操作系统IP地址
// @Description --------【选填参数:string】netMask：指定虚拟机操作系统网关
// @Description --------【选填参数:string】gateWay：指定虚拟机操作系统掩码
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param  body body view.BatchMountParam true "批量挂载参数"
// @Router /store_pool/replica/batch/mount [post]
// @Success 200 {object} Result{response=[]replica.Replica}
func BatchMountReplica(context *gin.Context) {
	var param view.BatchMountParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.BatchMountReplica(&param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 批量卸载副本
// @Description 【登录权限】批量卸载副本。

// @Tags 副本管理
// @Accept json
// @Produce json
// @Param  body body view.BatchUnmountParam true "批量挂载参数"
// @Router /store_pool/replica/batch/unmount [post]
// @Success 200 {object} Result{response=[]replica.Replica}
func UnMountReplicaBatch(context *gin.Context) {
	var param view.BatchUnmountParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data, err := storage.BatchUnMountReplica(&param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 获取副本执行实例记录
// @Description 【登录权限】获取副本执行记录。
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param id path int true "目标副本ID"
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Router /store_pool/replica/:id/mount/instance [get]
// @Success 200 {object} Result{response=instance.Instance}
func GetInstances(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var queryParam view.PageQueryParam
	if err := context.ShouldBind(&queryParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 赋默认值
	if queryParam.Count == 0 {
		queryParam.Count = 15
	}
	totalNum, data, err := storage.GetInstances(idparam.Id, &queryParam)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
	return
}

// @Summary 获取指定副本执行实例详情
// @Description 【登录权限】获取执行记录详情。
// @Tags 副本管理
// @Accept json
// @Produce json
// @Param id path int true "目标实例ID"
// @Router /store_pool/replica/mount/instance/:id/logs [get]
// @Success 200 {object} Result{response=instance.InstanceLog}
func GetInstanceLogs(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	data, err := storage.GetInstanceLogs(idparam.Id)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: int64(len(data))}
	ResponseSuccess(context, response)
	return
}

func UploadFile(context *gin.Context) {
	var idparam view.UpdIdParam
	if err := context.ShouldBind(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	updErr := storage.UploadJson(idparam.Id, context, userInfo)
	if updErr != nil {
		logrus.Error(updErr)
		ResponseFailure(context, updErr)
		return
	}
	ResponseSuccess(context, nil)
}

func ViewFile(context *gin.Context) {
	updErr := storage.ViewJson(context)
	if updErr != nil {
		logrus.Error(updErr)
		ResponseFailure(context, updErr)
		return
	}
}

func DoCompability(context *gin.Context) {
	var jobParam view.DoCasJobParam
	if err := context.ShouldBind(&jobParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	result, buildErr := storage.DoCasJob(&jobParam, context)
	if buildErr != nil {
		logrus.Error(buildErr)
		ResponseFailure(context, buildErr)
		return
	}
	ResponseSuccess(context, result)
}
