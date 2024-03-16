package api

import (
	"encoding/json"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Token struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

type Paging struct {
	Data     interface{} `json:"data"`
	TotalNum int64       `json:"totalNum"`
}

type Result struct {
	Code     msg.Massage `json:"code"`
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}

func ResponseSuccess(context *gin.Context, data interface{}) {
	result := Result{Code: msg.SUCCESS, Message: msg.GetMsg(msg.SUCCESS), Response: data}
	context.JSON(http.StatusOK, result)
}

func ResponseFailure(context *gin.Context, err error) {
	e := errors.New(msg.UNDEFINED, msg.GetMsg(msg.UNDEFINED, err.Error()))
	if err != nil {
		_ = json.Unmarshal([]byte(err.Error()), &e)
	}
	result := Result{Code: e.GetCode(), Message: e.Error(), Response: nil}
	context.JSON(http.StatusOK, result)
}

func ResponseCustom(context *gin.Context, code msg.Massage, msg string, data interface{}) {
	result := Result{Code: code, Message: msg, Response: data}
	context.JSON(http.StatusOK, result)
}
