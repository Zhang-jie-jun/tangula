package main

import (
	"flag"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/Zhang-jie-jun/tangula/internal/initsvc"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Host_Test() {
	// 创建
	var host Host
	host.Name = "test"
	host.Desc = "this is test host!"
	host.HostName = "local.host"
	host.Type = 1
	host.Os = "Centos 7.5"
	host.Arch = "x86_64"
	host.Ip = "127.0.0.1"
	host.Port = 22
	host.UserName = "root"
	host.PassWord = "passwd.com123"
	host.CreateUser = "jack"
	host.AuthType = 111
	host.CreatedAt = time.Now()
	err := dao.ClientDB.Create(&host).Error
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(host)
	// 查询
	err = dao.ClientDB.Where("id = ?", host.Id).Find(&host).Error
	if err != nil {
		fmt.Printf("FindById Error:%v", err)
		return
	}
	fmt.Println(host)
	// 删除
	err = dao.ClientDB.Where("id = ?", host.Id).Find(&host).Error
	if err != nil {
		fmt.Printf("Delete host error:%v", err)
		return
	}
	err = dao.ClientDB.Delete(&host).Error
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(host)
	res := host.TransformMap()
	fmt.Println(res)
}

func Platform_Test() {
	// 创建
	var platform Platform
	platform.Name = "test"
	platform.Desc = "this is test host!"
	platform.Type = 11
	platform.Ip = "127.0.0.1"
	platform.Port = 22
	platform.UserName = "root"
	platform.PassWord = "passwd.com123"
	platform.Version = "6.7.3"
	platform.CreateUser = "jack"
	platform.AuthType = 111
	err := dao.ClientDB.Create(&platform).Error
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(platform)
	// 查询
	err = dao.ClientDB.Where("id = ?", platform.Id).Find(&platform).Error
	if err != nil {
		fmt.Printf("FindById Error:%v", err)
		return
	}
	fmt.Println(platform)
	// 删除
	//err = dao.ClientDB.Where("id = ?", platform.Id).Find(&platform).Error
	//if err != nil {
	//	fmt.Printf("Delete platform error:%v", err)
	//	return
	//}
	//err = dao.ClientDB.Delete(&platform).Error
	//if err != nil {
	//	fmt.Printf("Error:%v", err)
	//	return
	//}
	fmt.Println(platform)
	res := platform.TransformMap()
	fmt.Println(res)
}

func Tenant_Test() {
	// 创建
	var tenant Tenant
	tenant.Name = "test"
	tenant.Desc = "this is test host!"
	tenant.DomainId = "214c9e1e-a2e2-4e89-859d-d279eef520d9"
	tenant.UserName = "root"
	tenant.PassWord = "passwd.com123"
	tenant.PlatformId = 1
	err := dao.ClientDB.Create(&tenant).Error
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(tenant)
	// 查询
	err = dao.ClientDB.Where("id = ?", tenant.Id).Preload("Platform").Find(&tenant).Error
	if err != nil {
		fmt.Printf("FindById Error:%v", err)
		return
	}
	fmt.Println(tenant)
	res := tenant.TransformMap()
	fmt.Println(res)
	// 删除
	err = dao.ClientDB.Where("id = ?", tenant.Id).Find(&tenant).Error
	if err != nil {
		fmt.Printf("Delete tenant error:%v", err)
		return
	}
	err = dao.ClientDB.Delete(&tenant).Error
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(tenant)
}

func StorePool_Test() {
	var storePool StorePool
	storePool.Name = "pool1"
	storePool.Desc = "this is pool1!"
	storePool.Uuid = util.GenerateGuid()
	storePool.CreateUser = "jack.zhang"
	err := dao.ClientDB.Create(&storePool).Error
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(storePool)

	var storePools []StorePool
	totalNum := 0
	dao.ClientDB.Model(&storePools).Where(fmt.Sprintf("name like '%%%s%%'", "")).Count(&totalNum)
	if err := dao.ClientDB.Limit(15).Offset(0).
		Where(fmt.Sprintf("name like '%%%s%%'", "")).Find(&storePools).Error; err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(storePools)

	err = dao.ClientDB.Where("id = ?", 1).Find(&storePool).Error
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	fmt.Println(storePool)
}

type replicaInfo struct {
	Id   uint
	Name string
	Desc string
}

func GetBySql() (res []view.VmFileRecoveryResult, err error) {
	fmt.Printf("开始查询")
	err = dao.ClientDB.Raw("SELECT sequence,pass,failed,duration,startTime from case_manager.vmfile_recovery_test where id = (SELECT max(id)  FROM case_manager.vmfile_recovery_test where title=?) ", "总计").Scan(&res).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func main() {
	contants.ConfPath = flag.String("conf", "D:\\mygit\\tangula\\configs\\app.ini", "config file path")
	flag.Parse()
	// 初始化服务配置
	err := initsvc.LoadResource()
	if err != nil {
		fmt.Printf("InitCfg Failed! Error:%v\n", err)
		os.Exit(1)
	}
	defer initsvc.UnLoadResource()
	// 主机测试
	//Host_Test()
	// 平台测试
	//Platform_Test()
	// 租户测试
	//Tenant_Test()
	// 存储池测试
	// StorePool_Test()
	res, sqlErr := GetBySql()
	if sqlErr != nil {
		logrus.Error(sqlErr)
	}
	logrus.Info(res)
	for _, i := range res {
		fmt.Println(i.Pass, i.StartTime)
	}

}
