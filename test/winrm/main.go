package main

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/service/app/windows"
	"github.com/sirupsen/logrus"
)

func main() {
	var login contants.LoginInfo
	login.Ip = "10.2.18.184"
	login.Port = 5985
	login.UserName = "Administrator"
	login.PassWord = "passwd123."
	client, err := windows.NewWinRmClient(&login)
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}

	var ver contants.VerifyResult
	ver.HostName, ver.Version, ver.Arch, err = client.GetSystemInfo()
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(fmt.Sprintf("Info:%+v", ver))

	status := client.CheackNfsStatus()
	if status {
		logrus.Info(fmt.Sprintf("NFS is RUNNING!"))
	} else {
		logrus.Info(fmt.Sprintf("NFS is STOP!"))
	}
	drives, err := client.GetSpaceDrive()
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(fmt.Sprintf("drives:%+v", drives))
	drive, err := client.GetUnUsedDrive()
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(fmt.Sprintf("drive:%+v", drive))
	path, err := client.GetDesktopPath()
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(fmt.Sprintf("path:%+v", path))

	sharePath := "\\\\192.168.212.32\\mnt\\tangula"
	drivers, err := client.Mount(sharePath)
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(drivers)
	err = client.CreateShortcut("nfs111", sharePath)
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}

	err = client.DeleteShortcut("nfs111")
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	err = client.UnMount(sharePath)
	if err != nil {
		err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, login.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	return
}
