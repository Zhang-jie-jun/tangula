package log

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type LogOrm struct{}

var LogMgm = LogOrm{}

func (l *LogOrm) FindRecordById(id uint) (record LogRecord, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&record).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (l *LogOrm) FindRecordByName(name string) (record LogRecord, err error) {
	err = dao.ClientDB.Where("name = ?", name).First(&record).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (l *LogOrm) CreateRecord(record LogRecord) (result LogRecord, err error) {
	err = dao.ClientDB.Create(&record).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = record
	return
}

func (l *LogOrm) UpdateRecord(record LogRecord) (result LogRecord, err error) {
	if err = dao.ClientDB.Model(&record).Update(&record).Error; err != nil {
		logrus.Error(err)
		return
	}
	result = record
	return
}

func (l *LogOrm) DeleteRecord(id uint) error {
	var record LogRecord
	err := dao.ClientDB.Where("id = ?", id).Find(&record).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&record).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (l *LogOrm) BatchDeleteByUser(user string) error {
	var records []LogRecord
	err := dao.ClientDB.Where("user = ?", user).Find(&records).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&records).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取操作记录列表，支持分页，操作人过滤
func (l *LogOrm) GetHostListByUser(index, count int, user string) (totalNum int64, records []LogRecord, err error) {
	dao.ClientDB.Model(&records).Where(fmt.Sprintf(" user like '%%%s%%' ", user)).Count(&totalNum)
	if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
		Where(fmt.Sprintf(" user like '%%%s%%' ", user)).Find(&records).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}
