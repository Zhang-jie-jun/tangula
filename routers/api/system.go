package api

import (
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/system"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 获取平台统计信息
// @Description 【登录权限】获取平台统计信息。
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Router /server/dashboard [get]
// @Success 200 {object} Result{response=system.Dashboard}
func GetDashboard(context *gin.Context) {
	// 获取当前登录用户
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	data := system.GetDashboardInfo(userInfo.Account)
	ResponseSuccess(context, data)
}

// @Summary 获取服务器状态信息
// @Description 【登录权限】获取服务器状态信息。
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Router /server/system [get]
// @Success 200 {object} Result{response=system.Server}
func GetServerInfo(context *gin.Context) {
	// 获取当前登录用户
	data, err := system.GetServerInfo()
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 获取ceph集群信息
// @Description 【登录权限】获取ceph集群信息。
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Router /server/ceph [get]
// @Success 200 {object} Result{response=system.CephDetailInfo}
func GetCephInfo(context *gin.Context) {
	// 获取当前登录用户
	data, err := system.GetCephDetailInfo()
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}
