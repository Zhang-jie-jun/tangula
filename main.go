package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/initsvc"
	"github.com/Zhang-jie-jun/tangula/routers"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Tangula
// @version 1.0
// @description Tangula API接口文档
// @description 唐古拉山脉，长江正源，滋润万物。
// @description Tangula 寓意数据的源头，为研发团队提供即取即用的数据服务，帮助团队高效产出。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email zhang.jiejun@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8000
// @BasePath
func main() {
	contants.ConfPath = flag.String("conf", "./configs/app.ini", "service config file path")
	flag.Parse()
	// 初始化服务配置
	err := initsvc.LoadResource()
	if err != nil {
		fmt.Printf("Satrt Service Failed! Error:%v\n", err)
		os.Exit(1)
	}
	defer initsvc.UnLoadResource()
	logrus.Info("==========================Server starting==============================")
	// 加载路由
	router := routers.InitRouter()
	// 设置服务启动参数
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", contants.AppCfg.Server.Host, contants.AppCfg.Server.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(contants.AppCfg.Server.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(contants.AppCfg.Server.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 启动服务
	go func() {

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}

	}()

	// 优雅的关闭(或重启)服务
	// 5秒后优雅Shutdown服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Block until a signal is received.
	s := <-quit
	logrus.Info("Shutdown Server ...", s)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Info("Server exiting")
}
