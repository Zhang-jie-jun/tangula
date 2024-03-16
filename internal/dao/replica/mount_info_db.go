package replica

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type MountInfoOrm struct{}

var MountInfoMgm = MountInfoOrm{}

func (m *MountInfoOrm) FindMountInfoByReplicaId(replicaId uint) (mountInfo MountInfo, err error) {
	err = dao.ClientDB.Where("replica_id = ?", replicaId).Find(&mountInfo).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (m *MountInfoOrm) CreateMountInfo(mountInfo MountInfo) (result MountInfo, err error) {
	err = dao.ClientDB.Create(&mountInfo).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = mountInfo
	return
}

func (m *MountInfoOrm) DeleteMountInfoByReplicaId(replicaId uint) error {
	var mountInfo MountInfo
	err := dao.ClientDB.Where("replica_id = ?", replicaId).Find(&mountInfo).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&mountInfo).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (m *MountInfoOrm) GetMountInfoByTarget(targetType contants.AppType, targetId uint) ([]MountInfo, error) {
	var mountInfos []MountInfo
	if err := dao.ClientDB.Where("target_type = ? AND target_id = ?", targetType, targetId).
		Find(&mountInfos).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return mountInfos, nil
}
