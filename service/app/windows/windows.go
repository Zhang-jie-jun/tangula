package windows

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/sirupsen/logrus"
	"strings"
)

type Adapter struct {
	Login *contants.LoginInfo
}

func NewAdapter(login *contants.LoginInfo) (adapter *Adapter) {
	return &Adapter{Login: login}
}

func (a *Adapter) VerifyConfigInfo() (result *contants.VerifyResult, err error) {
	if a.Login == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return
	}
	client, err := NewWinRmClient(a.Login)
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, a.Login.Ip, err.Error()))
		logrus.Error(err)
		return
	}

	var ver contants.VerifyResult
	ver.HostName, ver.OSType, ver.Arch, err = client.GetSystemInfo()
	result = &ver
	logrus.Info(result)
	return
}

func (a *Adapter) Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error) {
	// 路径转换
	start := strings.LastIndex(*sharePath, ":")
	ip := string([]byte(*sharePath)[:start])
	path := string([]byte(*sharePath)[start+1:])
	tempArr := strings.Split(path, "/")
	winSharePath := fmt.Sprintf("\\\\%s", ip)
	for _, iter := range tempArr {
		if iter == "" {
			continue
		}
		winSharePath = fmt.Sprintf("%s\\%s", winSharePath, iter)
	}
	// 获取imageUuid用于创建快捷方式名称
	index := strings.LastIndex(*sharePath, "/")
	shortcutName := fmt.Sprintf("tangula_%s", string([]byte(*sharePath)[index+1:]))
	// 创建winrm客户端
	if a.Login == nil {
		err := errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return "", err
	}
	client, err := NewWinRmClient(a.Login)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	// 挂载NFS存储
	mountPoint, err := client.Mount(winSharePath)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	// 创建桌面快捷方式
	err = client.CreateShortcut(shortcutName, winSharePath)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return mountPoint, nil
}

func (a *Adapter) UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) error {
	// 路径转换
	start := strings.LastIndex(*sharePath, ":")
	ip := string([]byte(*sharePath)[:start])
	path := string([]byte(*sharePath)[start+1:])
	tempArr := strings.Split(path, "/")
	winSharePath := fmt.Sprintf("\\\\%s", ip)
	for _, iter := range tempArr {
		if iter == "" {
			continue
		}
		winSharePath = fmt.Sprintf("%s\\%s", winSharePath, iter)
	}
	// 获取快捷方式名称
	index := strings.LastIndex(*sharePath, "/")
	shortcutName := fmt.Sprintf("tangula_%s", string([]byte(*sharePath)[index+1:]))
	// 创建winrm客户端
	if a.Login == nil {
		err := errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return err
	}
	client, err := NewWinRmClient(a.Login)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 删除桌面快捷方式
	err = client.DeleteShortcut(shortcutName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 卸载NFS存储
	err = client.UnMount(winSharePath)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
