package api

import (
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 获取指定快照信息
// @Description 【登录权限】获取指定快照信息
// @Tags 快照管理
// @Accept json
// @Produce  json
// @Param id path int true "目标快照ID"
// @Router /store_pool/snapshot/{id} [get]
// @Success 200 {object} Result{response=snapshot.Snapshot}
func GetSnapshotById(context *gin.Context) {
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
	data, err := storage.GetSnapshotById(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 回滚快照
// @Description 【登录权限】回滚快照副本到指定快照，非空闲状态【挂载中、卸载中、卸载失败、已挂载等】副本不允许回滚。
// @Tags 快照管理
// @Accept json
// @Produce  json
// @Param id path int true "目标快照ID"
// @Router /store_pool/snapshot/{id}/rollback [post]
// @Success 200 {object} Result
func RollbackSnapshot(context *gin.Context) {
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
	err = storage.RollbackSnapshot(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 删除快照
// @Description 【登录权限】删除指定快照。
// @Tags 快照管理
// @Accept json
// @Produce  json
// @Param id path int true "目标快照ID"
// @Router /store_pool/snapshot/{id} [delete]
// @Success 200 {object} Result
func DeleteSnapshot(context *gin.Context) {
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
	err = storage.DeleteSnapshot(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 生成镜像
// @Description 【登录权限】根据所选快照生成一个新的镜像。
// @Tags 快照管理
// @Accept json
// @Produce json
// @Param id path int true "目标快照ID"
// @Param name query string true "镜像名称"
// @Param desc query string false "镜像描述" default("")
// @Router /store_pool/snapshot/{id}/image [post]
// @Success 200 {object} Result{response=image.Image}
func CreateImageBySnapshot(context *gin.Context) {
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
	err = storage.CreateImageBySnapshot(idparam.Id, &param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 生成副本
// @Description 【登录权限】通过指定快照生成一个新的副本。
// @Tags 快照管理
// @Accept json
// @Produce  json
// @Param id path int true "目标快照ID"
// @Param name query string true "副本名称"
// @Param desc query string false "副本描述" default("")
// @Router /store_pool/snapshot/{id}/replica [post]
// @Success 200 {object} Result{response=replica.Replica}
func CreateReplicaBySnapshot(context *gin.Context) {
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
	err = storage.CreateReplicaBySnapshot(idparam.Id, &param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}
