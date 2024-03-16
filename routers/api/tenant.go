package api

import (
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 添加云平台租户
// @Description 【登录权限】添加云平台租户。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param platformId query int true "云平台Id"
// @Param name query string true "租户名称"
// @Param domainId query string true "访问域Id"
// @Param username query string true "租户登录账户"
// @Param password query string true "租户登录密码"
// @Router /resource/platform/tenant [post]
// @Success 200 {object} Result{response=platform.Tenant}
func CreatePlatformTenant(context *gin.Context) {
	logrus.Info(context)
	ResponseCustom(context, msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE), nil)
}

// @Summary 获取云平台所有租户
// @Description 【登录权限】获取指定云平台的所有租户信息。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标云平台ID"
// @Router /resource/platform/{id}/tenant [get]
// @Success 200 {object} Result{response=Paging{data=[]platform.Tenant}}
func GetPlatformTenant(context *gin.Context) {
	logrus.Info(context)
	ResponseCustom(context, msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE), nil)
}

// @Summary 获取指定云平台租户信息
// @Description 【登录权限】获取指定租户信息。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标云平台租户ID"
// @Router /resource/platform/tenant/:id [get]
// @Success 200 {object} Result{response=platform.Tenant}
func GetPlatformTenantById(context *gin.Context) {
	logrus.Info(context)
	ResponseCustom(context, msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE), nil)
}

// @Summary 编辑指定云平台租户
// @Description 【登录权限】编辑指定云平台租户信息。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标云平台租户ID"
// @Param domainId query string true "访问域Id"
// @Param username query string true "租户登录账户"
// @Param password query string true "租户登录密码"
// @Router /resource/platform/tenant/:id [put]
// @Success 200 {object} Result{response=platform.Tenant}
func EditPlatformTenant(context *gin.Context) {
	logrus.Info(context)
	ResponseCustom(context, msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE), nil)
}

// @Summary 删除指定云平台租户
// @Description 【登录权限】删除指定云平台租户。
// @Tags 平台管理
// @Accept json
// @Produce  json
// @Param id path int true "目标云平台租户ID"
// @Router /resource/platform/tenant/:id [delete]
// @Success 200 {object} Result
func DeletePlatformTenant(context *gin.Context) {
	logrus.Info(context)
	ResponseCustom(context, msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE), nil)
}
