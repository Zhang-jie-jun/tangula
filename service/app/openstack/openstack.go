package openstack

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
)

type Adapter struct {
	Login *contants.LoginInfo
}

func NewAdapter(login *contants.LoginInfo) (adapter *Adapter) {
	return &Adapter{Login: login}
}

func (a *Adapter) VerifyConfigInfo() (*contants.VerifyResult, error) {
	err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	return nil, err
}

func (a *Adapter) Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error) {
	err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	return "", err
}

func (a *Adapter) UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) error {
	err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	return err
}
