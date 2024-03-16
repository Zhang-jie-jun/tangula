package api

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/resource"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 获取已支持虚拟化(云)平台类型
// @Description 【登录权限】获取已支持虚拟化(云)平台类型[11.VMware，12.Cas, 13.HCS，14.OpenStack]。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Router /resource/platform/support [get]
// @Success 200 {object} Result{response=map[int]string}
func GetPlatformSupportType(context *gin.Context) {
	support := map[contants.AppType]string{
		contants.VMWARE:        "VMware",
		contants.CAS:           "CAS",
		contants.FUSIONCOMPUTE: "FusionCompute",
	}
	ResponseSuccess(context, support)
}

// @Summary 添加虚拟化(云)平台
// @Description 【登录权限】添加一个虚拟化(云)平台。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param name query string true "平台名称"
// @Param desc query string false "描述" default("")
// @Param type query int true "平台类型[11.VMware，12.Cas, 13.HCS，14.OpenStack]"
// @Param ip query string true "平台访问ip"
// @Param port query int true "平台访问port"
// @Param username query string true "平台登录账户"
// @Param password query string true "平台登录密码"
// @Router /resource/platform [post]
// @Success 200 {object} Result{response=platform.Platform}
func CreatePlatform(context *gin.Context) {
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	var createParam view.CreateResource
	if err := context.ShouldBind(&createParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	data, err := resource.CreatePlatform(&createParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 获取所有虚拟化(云)平台资源
// @Description 【登录权限】获取所有虚拟化(云)平台资源，用户只能获取属于自身的平台资源(含公共资源)。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param auth query int false "平台权限类型[0.不过滤类型，111.私有平台，222.公共平台]" default(0)
// @Param type query int false "平台类型过滤参数[0.不过滤类型，11.VMWARE，12.CAS，13.HCS，14.OPENSTACK]" default(0)
// @Param filter query string false "名称、IP搜索过滤参数" default("")
// @Router /resource/platform [get]
// @Success 200 {object} Result{response=Paging{data=[]platform.Platform}}
func GetPlatform(context *gin.Context) {
	var queryParam view.QueryResourceParam
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
	totalNum, data, err := resource.GetPlatform(&queryParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 获取指定虚拟化(云)平台资源信息
// @Description 【登录权限】获取指定虚拟化(云)平台资源信息。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标平台ID"
// @Router /resource/platform/{id} [get]
// @Success 200 {object} Result{response=platform.Platform}
func GetPlatformById(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Errorf("Error:%v\n", err)
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
	data, err := resource.GetPlatformById(idparam.Id, userInfo)
	if err != nil {
		logrus.Errorf("GetPlatformById Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 编辑指定虚拟化(云)平台
// @Description 【登录权限】编辑指定虚拟化(云)平台，用户只允许编辑属于自身的平台资源，管理员可以编辑公共平台资源， 超级管理员可用编辑所有平台资源。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标平台ID"
// @Param desc query string false "描述" default("")
// @Param ip query string true "平台访问ip"
// @Param port query int true "平台访问port"
// @Param username query string true "平台登录账户"
// @Param password query string true "平台登录密码"
// @Router /resource/platform/{id} [put]
// @Success 200 {object} Result{response=platform.Platform}
func EditPlatform(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Errorf("Error:%v\n", err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var updateParam view.UpdateResource
	if err := context.ShouldBind(&updateParam); err != nil {
		logrus.Errorf("Error:%v\n", err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	logrus.Infof("Param:%v\n", updateParam)
	data, err := resource.EditPlatform(idparam.Id, &updateParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 删除指定虚拟化(云)平台
// @Description 【登录权限】删除指定虚拟化(云)平台，用户只允许删除属于自身的平台资源，管理员可以删除公共平台资源， 超级管理员可用删除所有平台资源。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标平台ID"
// @Router /resource/platform/{id} [delete]
// @Success 200 {object} Result
func DeletePlatform(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Errorf("Error:%v\n", err)
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
	err = resource.DeletePlatform(idparam.Id, userInfo)
	if err != nil {
		logrus.Errorf("DeletePlatform Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 发布指定虚拟化(云)平台资源
// @Description 【登录权限】用户发布属于自身的虚拟化(云)平台资源，发布后平台资源不在属于用户。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标平台ID"
// @Router /resource/platform/{id}/publish [post]
// @Success 200 {object} Result
func PublishPlatform(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Errorf("Error:%v\n", err)
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
	data, err := resource.PublishPlatform(idparam.Id, userInfo)
	if err != nil {
		logrus.Errorf("PublishPlatform Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 获取VMware平台数据源
// @Description 【登录权限】获取VMware平台数据源信息，用户只能获取属于自身的平台资源(含公共资源)。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标平台ID"
// @Param fullPath query string true "数据源路径[初始值传空]"
// @Param showType query int true "数据源浏览方式[1.按计算资源浏览, 2.按虚拟机和模板浏览, 3.按存储浏览， 4.按网络浏览]"
// @Param isGetAll query bool false "是否获取所有" default(false)
// @Param isGetVm query bool false "是否获取虚拟机" default(false)
// @Router /resource/platform/{id}/vmware/datasources [get]
// @Success 200 {object} Result{response=Paging{data=[]vsphere.DataSource}}
func GetVMwareDataSources(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var dataSourcesParam view.VMwareDataSources
	if err := context.ShouldBind(&dataSourcesParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	totalNum, data, err := resource.GetVMwareDataSources(idparam.Id, &dataSourcesParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary VMware根据位置路径获取主机列表
// @Param id path int true "目标平台ID"
// @Param path query string true "位置路径"
// @Router /resource/platform/{id}/vmware/hostsByPath [get]
func GetVMwareHostsByPath(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	var vmwareHostParam view.VMwareHostParam
	if err := context.ShouldBind(&vmwareHostParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	response, err := resource.GetVMwareHostsByPath(idparam.Id, vmwareHostParam.Path, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, response)

}

func GetCasHostList(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	_, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response, err := resource.GetCasHosts(idparam.Id)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, response)
}

func GetCasHostInfo(context *gin.Context) {
	var idparam view.CasHostParam
	if err := context.ShouldBind(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}

	response, err := resource.GetCasHostDetails(idparam.Id, idparam.HostId)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, response)
}

func GetCasVm(context *gin.Context) {
	var idparam view.CasHostParam
	if err := context.ShouldBind(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response, err := resource.GetCasVmList(idparam.Id, idparam.HostId, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, response)
}

func GetCasStorage(context *gin.Context) {
	var idparam view.CasHostParam
	if err := context.ShouldBind(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response, err := resource.GetCasStorageList(idparam.Id, idparam.HostId, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, response)
}

func GetFcHostList(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	_, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response, err := resource.GetFcHosts(idparam.Id)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, response)
}
