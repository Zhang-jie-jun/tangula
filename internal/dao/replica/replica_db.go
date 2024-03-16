package replica

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type ReplicaOrm struct{}

var ReplicaMgm = ReplicaOrm{}

func (r *ReplicaOrm) FindById(id uint) (replica Replica, err error) {
	err = dao.ClientDB.Where("id = ?", id).Preload("Pool").Find(&replica).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (r *ReplicaOrm) FindByUuid(uuid string) (replica Replica, err error) {
	err = dao.ClientDB.Where("uuid = ?", uuid).Preload("Pool").Find(&replica).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 根据名称查询，名称可重复，取第一条数据
func (r *ReplicaOrm) FindByName(name string) (replica Replica, err error) {
	err = dao.ClientDB.Where("name = ?", name).Preload("Pool").First(&replica).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 检查同一用户下是否存在相同名称的镜像
func (r *ReplicaOrm) CheckIsExistByNameAndCreateUser(replicaName, createUser string) bool {
	var replica Replica
	err := dao.ClientDB.Where("name = ? AND create_user = ?", replicaName, createUser).First(&replica).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (r *ReplicaOrm) CreateReplica(replica Replica) (result Replica, err error) {
	err = dao.ClientDB.Create(&replica).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = replica
	return
}

func (r *ReplicaOrm) UpdateReplica(replica Replica) (result Replica, err error) {
	if err = dao.ClientDB.Model(&replica).Update(&replica).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = replica
	return
}

func (r *ReplicaOrm) UpdateReplicaStatus(id uint, status contants.ImageStatus) error {
	var replica Replica
	err := dao.ClientDB.Where("id = ?", id).Find(&replica).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	replica.Status = status
	if err = dao.ClientDB.Model(&replica).Update(&replica).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r *ReplicaOrm) UpdateReplicaStatusAndUuid(id uint, uuid string, status contants.ImageStatus) error {
	var replica Replica
	err := dao.ClientDB.Where("id = ?", id).Find(&replica).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	replica.Status = status
	replica.Uuid = uuid
	if err = dao.ClientDB.Model(&replica).Update(&replica).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r *ReplicaOrm) DeleteReplica(id uint) error {
	var replica Replica
	err := dao.ClientDB.Where("id = ?", id).Find(&replica).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&replica).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取副本列表，支持分页，按名称过滤
func (r *ReplicaOrm) GetReplica(index, count int, user string, status contants.ImageStatus,
	imageType contants.ImageType, filter string, uuid string) (totalNum int64, replica []Replica, err error) {
	if status != 0 {
		dao.ClientDB.Model(&replica).Where(fmt.Sprintf("create_user = '%s' AND status = %d AND name like '%%%s%%' AND uuid like '%%%s%%'",
			user, status, filter, uuid)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("create_user = '%s' AND status = %d AND name like '%%%s%%' AND uuid like '%%%s%%'", user, status, filter, uuid)).
			Preload("Pool").Find(&replica).Error; err != nil {
			logrus.Error(err)
			return
		}
	} else if imageType != 0 {
		dao.ClientDB.Model(&replica).Where(fmt.Sprintf("create_user = '%s' AND type = %d AND name like '%%%s%%' AND uuid like '%%%s%%'",
			user, imageType, filter, uuid)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("create_user = '%s' AND type = %d AND name like '%%%s%%' AND uuid like '%%%s%%'", user, imageType, filter, uuid)).
			Preload("Pool").Find(&replica).Error; err != nil {
			logrus.Error(err)
			return
		}
	} else {
		dao.ClientDB.Model(&replica).Where(fmt.Sprintf("create_user = '%s' AND name like '%%%s%%' AND uuid like '%%%s%%'",
			user, filter, uuid)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("create_user = '%s' AND name like '%%%s%%' AND uuid like '%%%s%%'", user, filter, uuid)).
			Preload("Pool").Find(&replica).Error; err != nil {
			logrus.Error(err)
			return
		}
	}
	return
}

func (r *ReplicaOrm) GetReplicaByPoolId(poolId uint) ([]Replica, error) {
	var replicas []Replica
	if err := dao.ClientDB.Where("pool_id = ?", poolId).Find(&replicas).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return replicas, nil
}

func (r *ReplicaOrm) DeleteReplicaByPoolId(poolId uint) error {
	var replicas Replica
	err := dao.ClientDB.Where("pool_id = ?", poolId).Find(&replicas).Error
	if err != nil {
		logrus.Info("存储池不包含副本:", poolId)
	} else {
		err = dao.ClientDB.Delete(&replicas).Error
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func (r *ReplicaOrm) GetReplicaStat(userName string) (TotalNum int) {
	var replicas []Replica
	dao.ClientDB.Model(&replicas).Where("create_user = ?", userName).Count(&TotalNum)
	return
}
