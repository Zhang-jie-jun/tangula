package api

import (
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/script"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 浏览脚本列表
// @Description 【登录权限】浏览脚本列表。
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param filter query string false "名称搜索过滤参数" default("")
// @Router /script/browse [get]
// @Success 200 {object} Result{response=[]script.Script}
func GetScripts(context *gin.Context) {
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
	totalNum, data, err := script.GetScripts(&queryParam, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}

// @Summary 上传脚本
// @Description 【登录权限】上传脚本。
// @Tags 脚本管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData string true "目标脚本文件"
// @Param desc query string false "脚本描述" default("")
// @Router /script/upload [post]
// @Success 200 {object} Result{response=script.Script}
func UploadScript(context *gin.Context) {
	data, err := script.UploadScript(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, data)
}

// @Summary 下载脚本
// @Description 下载脚本。
// @Tags 脚本管理
// @Accept multipart/form-data
// @Produce  json
// @Param id path int true "指定脚本ID"
// @Router /script/{id}/download [get]
// @Success 200 {object} Result{}
func DownloadScript(context *gin.Context) {
	err := script.DownloadScript(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
}

// @Summary 删除脚本
// @Description 删除脚本。
// @Tags 脚本管理
// @Accept json
// @Produce  json
// @Param id path int true "指定脚本ID"
// @Router /script/{id}/delete [delete]
// @Success 200 {object} Result{}
func DeleteScript(context *gin.Context) {
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
	err = script.DeleteScript(idparam.Id, userInfo)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, nil)
}

// @Summary 查看脚本内容
// @Description 【登录权限】查看脚本内容。
// @Tags 脚本管理
// @Accept json
// @Produce  json
// @Param id path int true "指定脚本ID"
// @Router /script/{id}/content [get]
// @Success 200 {object} Result{response=[]byte}
func GetScriptContent(context *gin.Context) {
	err := script.GetScriptContent(context)
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
}

func MetaCompare(context *gin.Context) {
	var metaParam view.MetaParam
	if err := context.ShouldBind(&metaParam); err != nil {
		logrus.Error(err)
		ResponseCustom(context, msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error()), nil)
		return
	}
	result, buildErr := script.MetaCompare(&metaParam)
	if buildErr != nil {
		logrus.Error(buildErr)
		ResponseFailure(context, buildErr)
		return
	}
	ResponseSuccess(context, result)
}

func GetVmFileRecoveryReport(context *gin.Context) {
	result, err := script.GetVmFileRecoveryReportImpl()
	if err != nil {
		logrus.Error(err)
		ResponseFailure(context, err)
		return
	}
	ResponseSuccess(context, result)

}
