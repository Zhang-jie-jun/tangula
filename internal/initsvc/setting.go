package initsvc

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/ldap"
	"github.com/Zhang-jie-jun/tangula/internal/logger"
	"github.com/go-ini/ini"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func LoadResource() error {
	contants.AppCfg = new(contants.AppConfig)
	cfg, err := ini.Load(*contants.ConfPath)
	if err != nil {
		fmt.Printf("InitCfg Load Failed! Error:%v\n", err)
		return err
	}
	err = cfg.MapTo(contants.AppCfg)
	if err != nil {
		fmt.Printf("InitCfg MapTo Failed! Error:%v\n", err)
		return err
	}
	// 创建服务依赖文件
	err = CheckFileExists(contants.AppCfg.System.FtpPath)
	if err != nil {
		fmt.Printf("Init Path ERROR:%v", err)
		return err
	}
	_ = CheckFileExists(contants.AppCfg.System.MountPath)
	_ = CheckFileExists(contants.AppCfg.System.ScriptPath)
	// 初始化日志
	InitLog()
	// 初始化数据库连接
	err = dao.InitDB(contants.AppCfg.Database.Type, contants.AppCfg.Database.User, contants.AppCfg.Database.Password,
		contants.AppCfg.Database.Host, contants.AppCfg.Database.DBName)
	if err != nil {
		logrus.Infof("InitDB ERROR:%v", err)
		return err
	}
	// 初始化基础角色类型
	var roles []map[string]string
	roles = append(roles, map[string]string{"Name": string(contants.USER), "Desc": string(contants.USER)})
	roles = append(roles, map[string]string{"Name": string(contants.SYSADMIN), "Desc": string(contants.SYSADMIN)})
	roles = append(roles, map[string]string{"Name": string(contants.SUPERADMIN), "Desc": string(contants.SUPERADMIN)})
	err = auth.AuthMgm.BatchCreateRole(roles)
	if err != nil {
		logrus.Infof("BatchCreateRole ERROR:%v", err)
		return err
	}
	// 初始化LDAP服务连接
	err = ldap.InitLDAP(contants.AppCfg.LDAP.Addr, contants.AppCfg.LDAP.BindUserName, contants.AppCfg.LDAP.BindPassword)
	if err != nil {
		logrus.Infof("InitLDAP ERROR:%v", err)
		return err
	}
	// 初始化Ceph连接
	err = ceph.InitCeph()
	if err != nil {
		logrus.Infof("InitCeph ERROR:%v", err)
		return err
	}
	return nil
}

func UnLoadResource() {
	dao.CloseDB()
	ceph.CloseCeph()
}

func CheckFileExists(fileName string) error {
	_, err := os.Stat(fileName)
	if err != nil {
		err = os.MkdirAll(fileName, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

// 日志文件初始化
func InitLog() {
	//file.MakeLog([]string{AppCfg.Sys.LogFile})
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&contants.MyFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.999",
	})

	logrus.SetLevel(logrus.Level(contants.AppCfg.System.LogLevel))
	writer := &logger.Logger{
		Filename:   contants.AppCfg.System.LogFile,          // 日志文件路径
		MaxSize:    50,                                      // megabytes
		MaxBackups: contants.AppCfg.System.LogReserveCount,  //保留7个备份
		MaxAge:     contants.AppCfg.System.LogReserveMaxDay, //最大保留天数30天
		Compress:   false,                                   // disabled by default
		LocalTime:  true,
	}
	mw := io.MultiWriter(os.Stdout, writer)
	logrus.SetOutput(mw)
}
