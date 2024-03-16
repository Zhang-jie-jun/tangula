package platform

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type PlatformOrm struct{}
type TenantOrm struct{}

var PlatformMgm = PlatformOrm{}
var TenantMgm = TenantOrm{}

func (p *PlatformOrm) FindById(id uint) (platform Platform, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&platform).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (p *PlatformOrm) FindByIp(ip string) (platform Platform, err error) {
	err = dao.ClientDB.Where("ip = ?", ip).Find(&platform).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 根据平台名称查询平台信息，名称可重复，取第一条数据
func (p *PlatformOrm) FindByName(name string) (platform Platform, err error) {
	err = dao.ClientDB.Where("name = ?", name).First(&platform).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (p *PlatformOrm) CheckIsExistByIp(ip string) bool {
	var platform Platform
	err := dao.ClientDB.Where("ip = ?", ip).First(&platform).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (p *PlatformOrm) CheckIsExistByName(name string) bool {
	var platform Platform
	err := dao.ClientDB.Where("name = ?", name).First(&platform).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

// 检查同一用户下是否存在相同名称的平台
func (p *PlatformOrm) CheckIsExistByNameAndCreateUser(platformName, createUser string) bool {
	var platform Platform
	err := dao.ClientDB.Where("name = ? AND create_user = ?", platformName, createUser).First(&platform).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (p *PlatformOrm) CreatePlatform(request Platform) (result Platform, err error) {
	err = dao.ClientDB.Create(&request).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = request
	return
}

func (p *PlatformOrm) UpdatePlatform(request Platform) (result Platform, err error) {
	if err = dao.ClientDB.Model(&request).Update(&request).Error; err != nil {
		logrus.Error(err)
		return
	}
	result = request
	return
}

// 删除指定平台
func (p *PlatformOrm) DeletePlatform(id uint) error {
	var platform Platform
	err := dao.ClientDB.Where("id = ?", id).Find(&platform).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&platform).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取平台列表，支持分页，按资源所属、类型、名称、IP过滤
func (p *PlatformOrm) GetPlatformList(index, count int, platformType contants.AppType, filter string) (
	totalNum int64, platforms []Platform, err error) {
	if platformType != 0 {
		dao.ClientDB.Model(&platforms).
			Where(fmt.Sprintf("type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
				platformType, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
				platformType, filter, filter)).Find(&platforms).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	} else {
		dao.ClientDB.Model(&platforms).Where(fmt.Sprintf("name like '%%%s%%' OR ip like '%%%s%%' ",
			filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("name like '%%%s%%' OR ip like '%%%s%%' ", filter, filter)).
			Find(&platforms).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	}
}

// 获取公共平台列表，支持分页，按类型、名称、IP过滤
func (p *PlatformOrm) GetPublicPlatformList(index, count int, platformType contants.AppType, filter string) (
	totalNum int64, platforms []Platform, err error) {
	if platformType != 0 {
		dao.ClientDB.Model(&platforms).
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
				platformType, contants.PUBLIC, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
				platformType, contants.PUBLIC, filter, filter)).Find(&platforms).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	} else {
		dao.ClientDB.Model(&platforms).
			Where(fmt.Sprintf("auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
				contants.PUBLIC, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
				contants.PUBLIC, filter, filter)).Find(&platforms).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	}
}

// 获取私有平台列表，支持分页，按类型、名称、IP过滤
func (p *PlatformOrm) GetPrivatePlatformList(index, count int, platformType contants.AppType, userName,
	filter string) (totalNum int64, platforms []Platform, err error) {
	if platformType != 0 {
		dao.ClientDB.Model(&platforms).
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				platformType, contants.PRIVATE, userName, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				platformType, contants.PRIVATE, userName, filter, filter)).Find(&platforms).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	} else {
		dao.ClientDB.Model(&platforms).
			Where(fmt.Sprintf("auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				contants.PRIVATE, userName, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				contants.PRIVATE, userName, filter, filter)).Find(&platforms).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	}
}

func (p *PlatformOrm) GetPlatformStat(userName string) (PublicNum, PrivateNum int) {
	var platforms []Platform
	dao.ClientDB.Model(&platforms).Where("auth_type = ?", contants.PUBLIC).Count(&PublicNum)
	dao.ClientDB.Model(&platforms).
		Where(fmt.Sprintf("auth_type = %d AND create_user = '%s'", contants.PRIVATE, userName)).Count(&PrivateNum)
	return
}

// 根据请求参数获取平台数量
func (p *PlatformOrm) GetPlatformCountByRequest(request string) (totalNum int64) {
	var platforms []Platform
	dao.ClientDB.Model(&platforms).Where(request).Count(&totalNum)
	return
}

// 根据请求参数获取平台列表, 支持分页
func (p *PlatformOrm) GetPlatformListByRequest(index, count int, request string) (result []map[string]interface{}, err error) {
	var platforms []Platform
	if err := dao.ClientDB.Limit(count).Offset(index).
		Where(request).Find(&platforms).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	for _, iter := range platforms {
		result = append(result, iter.TransformMap())
	}
	return result, err
}
