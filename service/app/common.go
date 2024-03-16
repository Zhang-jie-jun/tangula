package app

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/app/cas"
	"github.com/Zhang-jie-jun/tangula/service/app/fusioncompute"
	"github.com/Zhang-jie-jun/tangula/service/app/linux"
	"github.com/Zhang-jie-jun/tangula/service/app/vmware"
	"github.com/Zhang-jie-jun/tangula/service/app/windows"
)

// 应用适配器(应用需要实现认证，挂载，卸载逻辑)
type Adapter interface {
	VerifyConfigInfo() (result *contants.VerifyResult, err error)
	Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error)
	UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) error
}

func NewAdapter(appType contants.AppType, login *contants.LoginInfo) (adapter Adapter) {
	switch appType {
	case contants.LINUX:
		adapter = linux.NewAdapter(login)
	case contants.WINDOWS:
		adapter = windows.NewAdapter(login)
	case contants.VMWARE:
		adapter = vmware.NewAdapter(login)
	case contants.CAS:
		adapter = cas.NewAdapter(login)
	case contants.FUSIONCOMPUTE:
		adapter = fusioncompute.NewAdapter(login)
	default:
		adapter = NewDemo(login)
	}
	return adapter
}

type Demo struct {
	Login *contants.LoginInfo
}

func NewDemo(login *contants.LoginInfo) (demo *Demo) {
	return &Demo{Login: login}
}

func (d *Demo) VerifyConfigInfo() (*contants.VerifyResult, error) {
	err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	return nil, err
}

func (d *Demo) Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error) {
	err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	return "", err
}

func (d *Demo) UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) error {
	err := errors.New(msg.ERROR_NONSUPPORT_APP_TYPE, msg.GetMsg(msg.ERROR_NONSUPPORT_APP_TYPE))
	return err
}
