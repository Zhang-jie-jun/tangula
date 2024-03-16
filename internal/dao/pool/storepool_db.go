package pool

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type StorePoolOrm struct{}

var StorePoolMgm = StorePoolOrm{}

func (s *StorePoolOrm) FindById(id uint) (storePool StorePool, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&storePool).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (s *StorePoolOrm) FindByUuid(uuid string) (storePool StorePool, err error) {
	logrus.Infof("uuid:%s", uuid)
	err = dao.ClientDB.Where("uuid = ?", uuid).Find(&storePool).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 根据名称查询，名称可重复，取第一条数据
func (s *StorePoolOrm) FindByName(name string) (storePool StorePool, err error) {
	err = dao.ClientDB.Where("name = ?", name).First(&storePool).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 检查同一用户下是否存在相同名称的存储池
func (s *StorePoolOrm) CheckIsExistByNameAndCreateUser(storePoolName, createUser string) bool {
	var storePool StorePool
	err := dao.ClientDB.Where("name = ? AND create_user = ?", storePoolName, createUser).First(&storePool).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (s *StorePoolOrm) CreateStorePool(storePool StorePool) (result StorePool, err error) {
	err = dao.ClientDB.Create(&storePool).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = storePool
	return
}

func (s *StorePoolOrm) UpdateStorePool(storePool StorePool) (result StorePool, err error) {
	if err = dao.ClientDB.Model(&storePool).Update(&storePool).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = storePool
	return
}

func (s *StorePoolOrm) DeleteStorePool(id uint) error {
	var storePool StorePool
	err := dao.ClientDB.Where("id = ?", id).Find(&storePool).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&storePool).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取存储池列表，支持分页，按名称过滤
func (s *StorePoolOrm) GetStorePools(index, count int, filter string) (totalNum int64, storePools []StorePool, err error) {
	dao.ClientDB.Model(&storePools).Where(fmt.Sprintf("name like '%%%s%%'", filter)).Count(&totalNum)
	if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
		Where(fmt.Sprintf("name like '%%%s%%'", filter)).Find(&storePools).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (p *StorePoolOrm) GetPoolStat() (TotalNum int) {
	var storePools []StorePool
	dao.ClientDB.Model(&storePools).Count(&TotalNum)
	return
}
