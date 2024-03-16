package snapshot

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type SnapshotOrm struct{}

var SnapshotMgm = SnapshotOrm{}

type SnapShotId struct {
	Id uint
}

func (s *SnapshotOrm) FindById(id uint) (snapshot Snapshot, err error) {
	err = dao.ClientDB.Where("id = ?", id).Preload("Replica").Find(&snapshot).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (s *SnapshotOrm) FindByUuid(uuid string) (snapshot Snapshot, err error) {
	err = dao.ClientDB.Where("uuid = ?", uuid).Preload("Replica").Find(&snapshot).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (s *SnapshotOrm) FindByReplicaId(replicaId uint) (res []SnapShotId, err error) {
	err = dao.ClientDB.Raw("SELECT id FROM tangula.tangula_snapshot where replica_id=?", replicaId).Scan(&res).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 根据名称查询，名称可重复，取第一条数据
func (s *SnapshotOrm) FindByName(name string) (snapshot Snapshot, err error) {
	err = dao.ClientDB.Where("name = ?", name).Preload("Replica").First(&snapshot).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 检查同一镜像下是否存在相同名称的快照
func (s *SnapshotOrm) CheckIsExistByNameAndReplicaId(replicaName string, replicaId uint) bool {
	var snapshot Snapshot
	err := dao.ClientDB.Where("name = ? AND replica_id = ?", replicaName, replicaId).First(&snapshot).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (s *SnapshotOrm) CreateSnapshot(snapshot Snapshot) (result Snapshot, err error) {
	err = dao.ClientDB.Create(&snapshot).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = snapshot
	return
}

func (s *SnapshotOrm) UpdateSnapshot(snapshot Snapshot) (result Snapshot, err error) {
	if err = dao.ClientDB.Model(&snapshot).Update(&snapshot).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = snapshot
	return
}

func (s *SnapshotOrm) DeleteSnapshot(id uint) error {
	var snapshot Snapshot
	err := dao.ClientDB.Where("id = ?", id).Find(&snapshot).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&snapshot).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取副本快照列表，支持分页，按名称过滤
func (s *SnapshotOrm) GetSnapshots(index, count int, id uint, filter string) (totalNum int64, snapshots []Snapshot, err error) {
	dao.ClientDB.Model(&snapshots).Where(fmt.Sprintf("replica_id = %d AND name like '%%%s%%'", id, filter)).Count(&totalNum)
	if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
		Where(fmt.Sprintf("replica_id = %d AND name like '%%%s%%'", id, filter)).
		Preload("Replica").Find(&snapshots).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}
