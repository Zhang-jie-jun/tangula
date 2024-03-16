//go:build !linux
// +build !linux

package ceph

import (
	"github.com/sirupsen/logrus"
)

func InitCeph() (err error) {
	logrus.Warnf("Connect: connect error:%v\n", "windows nonsupport ceph!")
	return nil
}

func CloseCeph() {

}
