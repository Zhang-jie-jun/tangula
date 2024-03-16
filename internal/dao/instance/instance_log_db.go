package instance

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/sirupsen/logrus"
)

type InstanceLogOrm struct{}

var InstanceLogMgm = InstanceLogOrm{}

// 插入执行日志
func (i *InstanceLogOrm) PushInstanceLog(id uint, level contants.Level, info, detail string) {
	var log InstanceLog
	log.Level = level
	log.Info = info
	log.Detail = detail
	log.InstanceId = id
	log, err := i.CreateLog(log)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_LOG_DETAIL_INFO, msg.GetMsg(msg.ERROR_CREATE_LOG_DETAIL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (i *InstanceLogOrm) FindById(id uint) (instance Instance, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&instance).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (i *InstanceLogOrm) CreateLog(instanceLog InstanceLog) (result InstanceLog, err error) {
	err = dao.ClientDB.Create(&instanceLog).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = instanceLog
	return
}

func (i *InstanceLogOrm) DeleteLog(id uint) error {
	var instanceLog InstanceLog
	err := dao.ClientDB.Where("id = ?", id).Find(&instanceLog).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&instanceLog).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (i *InstanceLogOrm) DeleteLogsByInstanceId(instanceId uint) error {
	var instanceLog []InstanceLog
	err := dao.ClientDB.Where("instance_id = ?", instanceId).Find(&instanceLog).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&instanceLog).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (i *InstanceLogOrm) GetLogsByInstanceId(instanceId uint) (logs []InstanceLog, err error) {
	if err = dao.ClientDB.Where("instance_id = ?", instanceId).Find(&logs).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}
