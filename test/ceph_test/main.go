package main

import (
	"flag"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
	"github.com/Zhang-jie-jun/tangula/internal/initsvc"
	"github.com/Zhang-jie-jun/tangula/internal/system"
)

var (
	poolName          = "go_pool"
	imageName         = "go_image"
	imageSnapShotName = "go_image_snapshot"
	imageCopyName     = "copy_image"
	imageCloneName    = "clone_image"

	devPath    = ""
	mountPoint = "/mnt/tangula_go_image"
	targetIp   = "192.168.212.51"
)

func CloneOperation() {
	fmt.Printf("Start copy image...\n")
	err := ceph.Client.CopyImage(poolName, imageName, imageCopyName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("copy image success!\n")
	fmt.Printf("Start delete copy image...\n")
	err = ceph.Client.DeleteImage(poolName, imageCopyName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete copy image success!\n")
	fmt.Printf("Start clone image...\n")
	err = ceph.Client.CloneImageBySnapshot(poolName, imageName, imageSnapShotName, imageCloneName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("clone image success!\n")
	fmt.Printf("Start flatten image...\n")
	err = ceph.Client.FlattenImage(poolName, imageCloneName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("flatten image success!\n")
	fmt.Printf("Start delete clone image...\n")
	err = ceph.Client.DeleteImage(poolName, imageCloneName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete clone image success!\n")
}

func CephTest() {
	ClusterInfo, err := ceph.Client.GetClusterInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Cluster Info:%v\n", ClusterInfo)
	fmt.Printf("Start create pool...\n")
	err = ceph.Client.CreatePool(poolName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("create pool success!\n")
	fmt.Printf("Start create image...\n")
	err = ceph.Client.CreateImage(poolName, imageName, 1<<30)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("create image success!\n")
	fmt.Printf("Start create image snapshot...\n")
	err = ceph.Client.CreateSnapshot(poolName, imageName, imageSnapShotName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("create image snapshot success!\n")
	fmt.Printf("Start rollback image snapshot...\n")
	err = ceph.Client.RollbackSnapshot(poolName, imageName, imageSnapShotName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("rollback image snapshot success!\n")
	fmt.Printf("Start protect image snapshot...\n")
	err = ceph.Client.ProtectSnapShot(poolName, imageName, imageSnapShotName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("protect image snapshot success!\n")

	CloneOperation()

	fmt.Printf("Start unprotect image snapshot...\n")
	err = ceph.Client.UnProtectSnapShot(poolName, imageName, imageSnapShotName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("unprotect image snapshot success!\n")
	fmt.Printf("Start delete image snapshot...\n")
	err = ceph.Client.DeleteSnapshot(poolName, imageName, imageSnapShotName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete image snapshot success!\n")
	fmt.Printf("Start delete image...\n")
	err = ceph.Client.DeleteImage(poolName, imageName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete image success!\n")
	fmt.Printf("Start delete pool...\n")
	err = ceph.Client.DeletePool(poolName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete pool success!\n")
}

func MountTest() {
	fmt.Printf("Start create pool...\n")
	err := ceph.Client.CreatePool(poolName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("create pool success!\n")
	fmt.Printf("Start create image...\n")
	err = ceph.Client.CreateImage(poolName, imageName, 1<<30)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("create image success!\n")

	fmt.Printf("Start map image...\n")
	devPath, err = ceph.Client.MapRBDImage(poolName, imageName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("map image success! devPath:%s\n", devPath)
	fmt.Printf("Start format image...\n")
	err = system.SysManage.FormatXFS(devPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("format image success!\n")
	fmt.Printf("Start mount image...\n")
	err = system.SysManage.Mount(devPath, mountPoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("mount image success!\n")
	fmt.Printf("Start add nfs...\n")
	err = system.SysManage.AddNFS(mountPoint, targetIp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("add nfs success!\n")

	// todo 应用额外操作

}

func UnMountTest() {
	fmt.Printf("Start show map image...\n")
	mapInfos, err := ceph.Client.ShowMapRBDImage()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("show map image success!\n")
	fmt.Println(mapInfos)
	for _, iter := range mapInfos {
		if iter.PoolName == poolName && iter.ImageName == imageName {
			devPath = iter.DevPath
		}
	}
	fmt.Printf("Start remove nfs...\n")
	err = system.SysManage.RemoveNFS(mountPoint, targetIp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("remove nfs success!\n")
	fmt.Printf("Start unmount image...\n")
	err = system.SysManage.UnMount(devPath, mountPoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("unmount image success!\n")

	fmt.Printf("Start unmap image...\n")
	err = ceph.Client.UnMapRBDImage(poolName, imageName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("unmap image success!\n")

	fmt.Printf("Start delete image...\n")
	err = ceph.Client.DeleteImage(poolName, imageName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete image success!\n")
	fmt.Printf("Start delete pool...\n")
	err = ceph.Client.DeletePool(poolName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete pool success!\n")
}

func main() {
	contants.ConfPath = flag.String("conf", "../../configs/app.ini", "service config file path")
	flag.Parse()
	// 初始化服务配置
	err := initsvc.LoadResource()
	if err != nil {
		fmt.Printf("Satrt Service Failed! Error:%v\n", err)
		return
	}
	defer initsvc.UnLoadResource()

	err = ceph.InitCeph()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ceph.CloseCeph()
	// ceph测试
	CephTest()
	// 挂载测试
	MountTest()
	// 卸载测试
	UnMountTest()
}
