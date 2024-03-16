package api

import (
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 用户登录
// @Description 用户登录获取token。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Param username query string true "用户名称"
// @Param password query string true "用户密码"
// @Router /login [post]
// @Success 200 {object} Result{response=Token}
func Login(context *gin.Context) {
	logrus.Info(context)
}

// @Summary 用户注销
// @Description 【登录权限】用户注销登录。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Router /logout [post]
// @Success 200 {object} Result
func Logout(context *gin.Context) {
	logrus.Info(context)
}

// @Summary 刷新token
// @Description 【登录权限】刷新token有效期。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Router /auth/token [get]
// @Success 200 {object} Result{response=Token}
func Refresh(context *gin.Context) {
	logrus.Info(context)
}

// @Summary 获取用户信息
// @Description 【登录权限】获取指定用户信息(用户ID，用户名称，用户账号, 用户邮箱，联系方式, 用户角色，用户状态等)。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Router /auth/user [get]
// @Success 200 {object} Result{response=auth.User}
func GetUserInfo(context *gin.Context) {
	userInfo, err := auth.GetUserInfo(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, userInfo)
}

// @Summary 获取指定用户信息
// @Description 【管理员权限】获取指定用户信息(用户ID，用户名称，用户账号, 用户邮箱，联系方式, 用户角色，用户状态等)。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Param id path int true "目标用户ID"
// @Router /auth/user/{id} [get]
// @Success 200 {object} Result{response=auth.User}
func GetAssignUserInfo(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	user, err := auth.GetAssignUserInfo(idparam.Id)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, user)
}

// @Summary 获取用户列表
// @Description 【超级管理员权限】获取平台上所有用户信息，支持分页、按角色类型、用户名称、用户账户、用户邮箱过滤。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param roleName query string false "角色过滤参数[普通用户、系统管理员]" default("")
// @Param filter query string false "用户名称、用户账号、用户邮箱过滤参数[模糊匹配]" default("")
// @Router /auth/users [get]
// @Success 200 {object} Result{response=Paging{data=[]auth.User}}
func GetUserList(context *gin.Context) {
	var queryParam view.UserQueryParam
	if err := context.ShouldBind(&queryParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 赋默认值
	if queryParam.Count == 0 {
		queryParam.Count = 15
	}
	totalNum, data, err := auth.GetUserList(&queryParam)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 编辑用户角色
// @Description 【管理员权限】编辑用户角色。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Param id path int true "目标用户ID"
// @Param roleName query string true "角色名称[普通用户、系统管理员]"
// @Router /auth/user/{id}/set_role [put]
// @Success 200 {object} Result
func SetUserRole(context *gin.Context) {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	var queryParam view.RoleParam
	if err := context.ShouldBind(&queryParam); err != nil {
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
	user, err := auth.SetUserRole(idparam.Id, queryParam.RoleType, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, user)
}

// @Summary 获取角色列表
// @Description 【管理员权限】获取角色列表。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param filter query string false "名称过滤参数" default("")
// @Router /auth/role [get]
// @Success 200 {object} Result{response=Paging{data=[]auth.Role}}
func GetRoleList(context *gin.Context) {
	var queryParam view.UserQueryParam
	if err := context.ShouldBind(&queryParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	// 赋默认值
	if queryParam.Count == 0 {
		queryParam.Count = 15
	}
	totalNum, data, err := auth.GetRoleList(&queryParam)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 启用用户
// @Description 【系统管理员权限】启用用户权限。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Param id path int true "目标用户ID"
// @Router /auth/user/{id}/enable [put]
// @Success 200 {object} Result
func EnableUser(context *gin.Context) {
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
	user, err := auth.EnableUser(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, user)
}

// @Summary 禁用用户
// @Description 【系统管理员权限】禁用用户权限，禁用后用户不可登录平台、操作平台资源。
// @Tags 用户鉴权
// @Accept json
// @Produce json
// @Param id path int true "目标用户ID"
// @Router /auth/user/{id}/disable [put]
// @Success 200 {object} Result
func DisableUser(context *gin.Context) {
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
	user, err := auth.DisableUser(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, user)
}
