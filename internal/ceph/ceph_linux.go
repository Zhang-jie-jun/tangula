//go:build linux
// +build linux

package ceph

import "C"
import (
	"github.com/ceph/go-ceph/rados"
	"github.com/sirupsen/logrus"
)

var (
	Connect *rados.Conn //连接引擎
)

func InitCeph() error {
	err := Client.loginCeph()
	if err != nil {
		logrus.Errorf("Connect: connect error:%v\n", err)
		return err
	}
	return nil
}

func CloseCeph() {
	if Connect != nil {
		Connect.Shutdown()
	}
}
