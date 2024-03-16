package linux

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"strings"
)

type Adapter struct {
	Login *contants.LoginInfo
}

func NewAdapter(login *contants.LoginInfo) (adapter *Adapter) {
	return &Adapter{Login: login}
}

// 登录认证并获取主机信息
func (a *Adapter) VerifyConfigInfo() (result *contants.VerifyResult, err error) {
	if a.Login == nil {
		err = errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return
	}
	client, err := NewSshClient(a.Login)
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, a.Login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	defer client.Logout()

	var ver contants.VerifyResult
	ver.HostName = client.GetHostName()
	ver.OSType = client.GetHostOs()
	ver.Arch = client.GetHostArch()
	result = &ver
	logrus.Info(result)
	return
}

func (a *Adapter) Mount(sharePath *string, mountParam *view.MountInfo, user *auth.User, instanceId uint) (string, error) {
	// 解析挂载参数
	var mountPoint string
	if mountParam.AppConfig != nil {
		if _, ok := mountParam.AppConfig["mountPoint"]; ok {
			mountPoint = cast.ToString(mountParam.AppConfig["mountPoint"])
		}
	} else {
		mountParam.AppConfig = make(map[string]interface{})
	}
	if mountPoint == "" {
		start := strings.LastIndex(*sharePath, ":")
		temp := string([]byte(*sharePath)[start+1:])
		mountParam.AppConfig["mountPoint"] = temp
		mountPoint = temp
	}
	logrus.Info(mountPoint)
	// 挂载远程主机文件系统
	if a.Login == nil {
		err := errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return "", err
	}
	client, err := NewSshClient(a.Login)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	defer client.Logout()
	err = client.Mount(*sharePath, mountPoint)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return mountPoint, err
}

func (a *Adapter) UnMount(sharePath *string, mountParam *view.MountInfo, instanceId uint) error {
	start := strings.LastIndex(*sharePath, ":")
	temp := string([]byte(*sharePath)[start+1:])
	// 解析挂载参数
	mountPoint := cast.ToString(mountParam.AppConfig["mountPoint"])
	if mountPoint == "" {
		mountParam.AppConfig["mountPoint"] = temp
		mountPoint = temp
	}
	delMp := false
	// 判断是否是自动生成的挂载点，若是则需要清理挂载目录
	if mountPoint == temp {
		delMp = true
	}
	// 卸载远程主机文件系统
	if a.Login == nil {
		err := errors.New(msg.ERROR_INVALID_PARAMS,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		//return err
		return nil
	}
	client, err := NewSshClient(a.Login)
	if err != nil {
		logrus.Error(err)
		//return err
		return nil
	}
	defer client.Logout()
	err = client.UnMount(mountPoint, delMp)
	if err != nil {
		logrus.Error(err)
		//return err  暂时无视
	}
	return nil
}
