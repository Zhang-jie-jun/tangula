package api

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/auth"
	"github.com/Zhang-jie-jun/tangula/service/logmgm"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary 获取所有日志记录
// @Description 【登录权限】获取所有操作日志记录，普通用户只能获取自身的操作日志记录，系统管理员平台上所有用户的操作日志记录。
// @Tags 日志管理
// @Accept json
// @Produce json
// @Param index query int false "分页索引" default(0)
// @Param count query int false "分页数量" default(15)
// @Param user query string false "用户名称" default("")
// @Router /log/record [get]
// @Success 200 {object} Result{response=Paging{data=[]log.LogRecord}}
func GetLogRecord(context *gin.Context) {
	var queryParam view.LogRecordQueryParam
	if err := context.ShouldBind(&queryParam); err != nil {
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
	// 赋默认值
	if queryParam.Count == 0 {
		queryParam.Count = 15
	}
	// 普通用户只能获取属于自身的记录
	if userInfo.Role.Name == contants.USER {
		queryParam.User = userInfo.Account
	}
	totalNum, data, err := logmgm.GetLogRecord(&queryParam)
	if err != nil {
		logrus.Errorf("Error:%v", err)
		ResponseFailure(context, err)
		return
	}
	response := Paging{Data: data, TotalNum: totalNum}
	ResponseSuccess(context, response)
}
