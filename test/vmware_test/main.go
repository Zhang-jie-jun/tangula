package main

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service/app/vmware"
)

func main() {
	var login contants.LoginInfo
	login.Ip = "192.168.212.52"
	login.Port = 902
	login.UserName = "administrator@vsphere.local"
	login.PassWord = "passwd.com123"
	adapter := vmware.NewAdapter(&login)
	result, err := adapter.VerifyConfigInfo()
	if err != nil {
		fmt.Printf("VerifyConfigInfo error:%v\n", err)
	}
	fmt.Printf("VMware version:%s\n", result.Version)

	temp := map[string]interface{}{"LOCATION_PATH": "VMware开发测试环境/vm/122",
		"COMPUTE_RESOURCE": "VMware开发测试环境/host/cluster/Resources/Jack(开发环境)", "IS_REGISTER_VM": "false"}
	mountParam := &view.MountInfo{AppConfig: temp}
	poolName := "fb84e623c578384dac61330c1450c8b1"
	imageName := "4e6ed01351dd5a4a4fe509c286ff8ec5"
	localMountPoint := fmt.Sprintf("/mnt/tangula/%s/%s", poolName, imageName)
	sharePath := fmt.Sprintf("%s:%s", "192.168.212.32", localMountPoint)
	mountPoint, err := adapter.Mount(&sharePath, mountParam, nil, 0)
	if err != nil {
		fmt.Printf("Mount error:%v\n", err)
	}
	fmt.Printf("mount point:%s\n", mountPoint)
}
