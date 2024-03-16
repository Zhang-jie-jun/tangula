package dao

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	ClientDB *gorm.DB // 数据库引擎
)

func InitDB(DBType, UserName, PassWord, Host, DBName string) error {
	passWord, err := util.AesDecrypt(PassWord)
	if err != nil {
		return err
	}
	linkInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		UserName,
		passWord,
		Host,
		DBName)
	ClientDB, err = gorm.Open(DBType, linkInfo)
	if err != nil {
		logrus.Errorf("Error:%v\n", err)
		return err
	}
	if ClientDB == nil {
		logrus.Errorf("gorm open failed(db is nil)!")
		return err
	}

	ClientDB.SingularTable(true)
	ClientDB.DB().SetMaxIdleConns(10)
	ClientDB.DB().SetMaxOpenConns(100)
	return nil
}

func CloseDB() {
	err := ClientDB.Close()
	if err != nil {
		logrus.Errorf("Close DB Error:%v\n", err)
	}
}
