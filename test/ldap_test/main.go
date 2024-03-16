package main

import (
	"flag"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/initsvc"
	"github.com/Zhang-jie-jun/tangula/internal/ldap"
	"os"
)

var (
	UserName = "jack.zhang"
	PassWord = "password"
)

func main() {
	contants.ConfPath = flag.String("conf", "E:\\gocode\\tangula\\configs\\app.ini", "config file path")
	flag.Parse()
	// 初始化服务配置
	err := initsvc.LoadResource()
	if err != nil {
		fmt.Printf("InitCfg Failed! Error:%v\n", err)
		os.Exit(1)
	}
	defer initsvc.UnLoadResource()

	ok := ldap.Login(UserName, PassWord)
	if ok {
		fmt.Printf("login success!")
	} else {
		fmt.Printf("login failed!")
	}
}
