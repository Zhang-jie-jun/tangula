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

// @Summary 获取已支持主机资源类型
// @Description 【登录权限】获取已支持虚主机资源类型[51.Linux, 52.Windows]。
// @Tags 主机管理
// @Accept json
// @Produce  json
// @Router /resource/host/support [get]
// @Success 200 {object} Result{response=map[int]string}
func GetHostSupportType(context *gin.Context) {
	support := map[contants.AppType]string{
		contants.LINUX:   "Linux",
		contants.WINDOWS: "Windows",
	}
	ResponseSuccess(context, support)
}

// @Summary 添加主机
// @Description 【登录权限】添加一个主机。
// @Tags 主机管理
// @Accept json
// @Produce  json
// @Param name query string true "主机名称"
// @Param desc query string false "描述" default("")
// @Param type query int true "主机类型[51.Linux, 52.Windows]"
// @Param ip query string true "主机访问ip"
// @Param port query int true "主机访问port"
// @Param username query string true "主机登录账户"
// @Param password query string true "主机登录密码"
// @Router /resource/host [post]
// @Success 200 {object} Result{response=host.Host}
func CreateHost(context *gin.Context) {
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	var createParam view.CreateResource
	if err := context.ShouldBind(&createParam); err != nil {
		logrus.Errorf("Error:%v\n", err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	data, err := resource.CreateHost(&createParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 获取所有主机资源
// @Description 【登录权限】获取所有主机资源，用户只能获取属于自身的主机资源(含公共资源)。
// @Tags 主机管理
// @Accept json
// @Produce  json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param auth query int false "权限过滤类型[0.不过滤, 111.私有，222.公共]" default(0)
// @Param type query int false "主机过滤类型[0.不过滤, 51.Linux主机，52.Windows主机]" default(0)
// @Param filter query string false "名称、IP过滤参数[模糊匹配]" default("")
// @Router /resource/host [get]
// @Success 200 {object} Result{response=Paging{data=[]host.Host}}
func GetHost(context *gin.Context) {
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
	totalNum, data, err := resource.GetHostList(&queryParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 获取指定主机资源信息
// @Description 【登录权限】获取指定主机资源信息。
// @Tags 主机管理
// @Accept json
// @Produce  json
// @Param id path int true "目标主机ID"
// @Router /resource/host/{id} [get]
// @Success 200 {object} Result{response=host.Host}
func GetHostById(context *gin.Context) {
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
	data, err := resource.GetHostById(idparam.Id, userInfo)
	if err != nil {
		logrus.Errorf("GetHostById Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 编辑指定主机
// @Description 【登录权限】编辑指定主机，用户只允许编辑属于自身的主机资源，管理员可以编辑公共主机资源。
// @Tags 主机管理
// @Accept json
// @Produce  json
// @Param id path int true "目标主机ID"
// @Param desc query string false "描述" default("")
// @Param ip query string true "主机访问ip"
// @Param port query int true "主机访问port"
// @Param username query string true "主机登录账户"
// @Param password query string true "主机登录密码"
// @Router /resource/host/{id} [put]
// @Success 200 {object} Result{response=host.Host}
func EditHost(context *gin.Context) {
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
	data, err := resource.EditHost(idparam.Id, &updateParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 删除指定主机
// @Description 【登录权限】删除指定主机，用户只允许删除属于自身的主机资源，管理员可以删除公共主机资源。
// @Tags 主机管理
// @Accept json
// @Produce  json
// @Param id path int true "目标主机ID"
// @Router /resource/host/{id} [delete]
// @Success 200 {object} Result
func DeleteHost(context *gin.Context) {
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
	err = resource.DeleteHost(idparam.Id, userInfo)
	if err != nil {
		logrus.Errorf("DeleteHost Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 发布指定主机资源
// @Description 【登录权限】用户发布属于自身的主机资源，发布后主机资源不在属于用户。
// @Tags 主机管理
// @Accept json
// @Produce  json
// @Param id path int true "目标主机ID"
// @Router /resource/host/{id}/publish [post]
// @Success 200 {object} Result
func PublishHost(context *gin.Context) {
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
	data, err := resource.PublishHost(idparam.Id, userInfo)
	if err != nil {
		logrus.Errorf("PublishHost Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

func UpdateHost(context *gin.Context) {
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
	updErr := resource.UpdHost(idparam.Id, userInfo)
	if updErr != nil {
		logrus.Errorf("PublishHost Error:%v", updErr)
		ResponseFailure(context, updErr)
		return
	}
	ResponseSuccess(context, nil)
}

func DeployApp(context *gin.Context) {
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	var deployParam view.DeployParam
	if paramErr := context.ShouldBind(&deployParam); paramErr != nil {
		logrus.Errorf("Error:%v\n", paramErr)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, paramErr.Error()), nil)
		return
	}
	data, err := resource.DeployApps(deployParam.HostId, &deployParam, userInfo, context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

func UpdDeploy(context *gin.Context) {
	var updparam view.UpdDeployParam
	if err := context.ShouldBind(&updparam); err != nil {
		logrus.Errorf("Error:%v\n", err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}

	updErr := resource.UpdDeploy(updparam.Id, updparam.Status, updparam.Log)
	if updErr != nil {
		logrus.Errorf("UpdDeploy Error:%v", updErr)
		ResponseFailure(context, updErr)
		return
	}
	ResponseSuccess(context, nil)
}

func GetDeployRecord(context *gin.Context) {
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
	totalNum, data, err := resource.GetDeployRecord(idparam.Id, &queryParam)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
	return
}
