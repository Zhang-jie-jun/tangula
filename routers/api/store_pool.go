package api

import (
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 创建存储池
// @Description 【管理员权限】创建一个新存储池。
// @Tags 存储池管理
// @Accept json
// @Produce json
// @Param name query string true "存储池名称"
// @Param desc query string false "存储池描述" default("")
// @Router /resource/store_pool [post]
// @Success 200 {object} Result{response=pool.StorePool}
func CreateStorePool(context *gin.Context) {
	var param view.NameParam
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
	data, err := storage.CreateStorePool(&param, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 删除存储池
// @Description 【管理员权限】删除指定存储池。
// @Tags 存储池管理
// @Accept json
// @Produce json
// @Param id path int true "目标存储池ID"
// @Router /resource/store_pool/{id} [delete]
// @Success 200 {object} Result
func DeleteStorePool(context *gin.Context) {
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
	err = storage.DeleteStorePool(idparam.Id, userInfo)
	if err != nil {
		logrus.Errorf("DeleteStorePool Error:%v", err)
		ResponseFailure(context, err)
		return
	}

	ResponseSuccess(context, nil)
}

// @Summary 获取所有存储池
// @Description 【登录权限】获取平台上所有存储池信息。
// @Tags 存储池管理
// @Accept json
// @Produce json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param filter query string false "名称搜索过滤参数" default("")
// @Router /resource/store_pool [get]
// @Success 200 {object} Result{response=Paging{data=[]pool.StorePool}}
func GetStorePool(context *gin.Context) {
	var param view.QueryParam
	if err := context.ShouldBind(&param); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 赋默认值
	if param.Count == 0 {
		param.Count = 15
	}
	totalNum, data, err := storage.GetStorePools(&param)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 获取指定存储池
// @Description 【登录权限】获取指定存储池信息。
// @Tags 存储池管理
// @Accept json
// @Produce  json
// @Param id path int true "目标存储池ID"
// @Router /resource/store_pool/{id} [get]
// @Success 200 {object} Result{response=pool.StorePool}
func GetStorePoolById(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Errorf("Error:%v\n", err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	data, err := storage.GetStorePoolById(idparam.Id)
	if err != nil {
		logrus.Errorf("GetStorePoolById Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}
