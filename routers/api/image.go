package api

import (
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 获取所有镜像信息
// @Description 【登录权限】获取所有镜像信息，普通用户与系统管理员只能获取公共镜像或自身创建的镜像信息，超级管理员可以获取平台上所有的镜像信息。
// @Tags 镜像管理
// @Accept json
// @Produce  json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param auth query int false "镜像权限过滤参数[0.不过滤类型，111.私有镜像，222.公有镜像]" default(0)
// @Param type query int false "镜像类型过滤参数[0.不过滤类型，1001.FILE，1002.ISO，1003.CUSTOME，1004.VMWAREVM，1005.CASVM，1006.HCSVM，1007.OPENSTACKVM]" default(0)
// @Param filter query string false "名称搜索过滤参数" default("")
// @Router /store_pool/image [get]
// @Success 200 {object} Result{response=Paging{data=[]image.Image}}
func GetImage(context *gin.Context) {
	var queryParam view.ImageQueryParam
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
	totalNum, data, err := storage.GetImage(&queryParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 获取指定镜像
// @Description 【登录权限】获取指定镜像信息，普通用户与系统管理员只能获取公共镜像或自身创建的镜像信息，超级管理员可以获取平台上所有的镜像信息。
// @Tags 镜像管理
// @Accept json
// @Produce  json
// @Param id path int true "目标镜像ID"
// @Router /store_pool/image/{id} [get]
// @Success 200 {object} Result{response=image.Image}
func GetImageById(context *gin.Context) {
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
	data, err := storage.GetImageById(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 删除镜像
// @Description 【登录权限】删除指定镜像，普通用户自能删除属于自身的镜像，管理员可删除自身及公有镜像，超级管理员可删除所有镜像。
// @Tags 镜像管理
// @Accept json
// @Produce  json
// @Param id path int true "目标镜像ID"
// @Router /store_pool/image/{id} [delete]
// @Success 200 {object} Result
func DeleteImage(context *gin.Context) {
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
	err = storage.DeleteImage(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 发布镜像
// @Description 【登录权限】发布指定镜像，只能发布属于用户自身的镜像。
// @Tags 镜像管理
// @Accept json
// @Produce  json
// @Param id path int true "目标镜像ID"
// @Router /store_pool/image/{id}/publish [post]
// @Success 200 {object} Result{response=image.Image}
func PublishImage(context *gin.Context) {
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
	data, err := storage.PublishImage(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 生成副本
// @Description 【登录权限】通过指定镜像生成副本。
// @Tags 镜像管理
// @Accept json
// @Produce  json
// @Param id path int true "目标镜像ID"
// @Param name query string true "副本名称"
// @Param desc query string false "副本描述" default("")
// @Router /store_pool/image/{id}/replica [post]
// @Success 200 {object} Result{response=replica.Replica}
func CreateReplicaByImage(context *gin.Context) {
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
	data, crtErr := storage.CreateReplicaByImage(idparam.Id, &param, userInfo)
	if crtErr != nil {
		logrus.Error(crtErr)
		ResponseFailure(context, crtErr)
		return
	}
	ResponseSuccess(context, data)
}

func EditImage(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var param view.ImageEditParam
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
	data, err := storage.EditImage(idparam.Id, &param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}
