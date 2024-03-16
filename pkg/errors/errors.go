package errors

import (
	"encoding/json"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
)

type errors struct {
	Code msg.Massage `json:"code"`
	Msg  string      `json:"message"`
}

func (e *errors) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

func (e *errors) GetCode() msg.Massage {
	return e.Code
}

func New(code msg.Massage, msg string) *errors {
	return &errors{
		Code: code,
		Msg:  msg,
	}
}
