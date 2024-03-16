package instance

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type InstanceOrm struct{}

var InstanceMgm = InstanceOrm{}

func (i *InstanceOrm) FindById(id uint) (instance Instance, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&instance).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (i *InstanceOrm) CreateInstance(instance Instance) (result Instance, err error) {
	err = dao.ClientDB.Create(&instance).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = instance
	return
}

func (i *InstanceOrm) UpdateInstance(instance Instance) (result Instance, err error) {
	if err = dao.ClientDB.Model(&instance).Update(&instance).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = instance
	return
}

func (i *InstanceOrm) UpdateInstanceStatus(id uint, status contants.InstanceStatus) {
	var instance Instance
	err := dao.ClientDB.Where("id = ?", id).Find(&instance).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	instance.Status = status
	if err = dao.ClientDB.Model(&instance).Update(&instance).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (i *InstanceOrm) DeleteInstance(id uint) error {
	var instance Instance
	err := dao.ClientDB.Where("id = ?", id).Find(&instance).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&instance).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (i *InstanceOrm) DeleteInstancesByReplicaId(replicaId uint) error {
	var instances []Instance
	err := dao.ClientDB.Order("id DESC").Where("replica_id = ?", replicaId).Find(&instances).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&instances).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (i *InstanceOrm) GetInstancesByReplicaId(replicaId uint, index, count int) (totalNum int64, instances []Instance, err error) {
	dao.ClientDB.Model(&instances).Where("replica_id = ?", replicaId).Count(&totalNum)
	if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
		Where("replica_id = ?", replicaId).Find(&instances).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}
