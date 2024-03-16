package host

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type HostOrm struct{}

var HostMgm = HostOrm{}

func (h *HostOrm) FindById(id uint) (host Host, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&host).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (h *HostOrm) FindByIp(ip string) (host Host, err error) {
	err = dao.ClientDB.Where("ip = ?", ip).Find(&host).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 根据主机名称查询主机信息，名称可重复，取第一条数据
func (h *HostOrm) FindByName(name string) (host Host, err error) {
	err = dao.ClientDB.Where("name = ?", name).First(&host).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (h *HostOrm) CheckIsExistByIp(ip string) bool {
	var host Host
	err := dao.ClientDB.Where("ip = ?", ip).First(&host).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (h *HostOrm) CheckIsExistByName(name string) bool {
	var host Host
	err := dao.ClientDB.Where("name = ?", name).First(&host).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

// 检查同一用户下是否存在相同名称的主机
func (h *HostOrm) CheckIsExistByNameAndCreateUser(hostName, createUser string) bool {
	var host Host
	err := dao.ClientDB.Where("name = ? AND create_user = ?", hostName, createUser).First(&host).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (h *HostOrm) CreateHost(host Host) (result Host, err error) {
	err = dao.ClientDB.Create(&host).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = host
	return
}

func (h *HostOrm) UpdateHost(host Host) (result Host, err error) {
	if err = dao.ClientDB.Model(&host).Update(&host).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = host
	return
}

// 删除指定主机
func (h *HostOrm) DeleteHost(id uint) error {
	var host Host
	err := dao.ClientDB.Where("id = ?", id).Find(&host).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&host).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取主机列表，支持分页，按资源所属、名称、IP过滤
func (h *HostOrm) GetHostList(index, count int, hostType contants.AppType, filter string) (totalNum int64, hosts []Host, err error) {
	if hostType != 0 {
		dao.ClientDB.Model(&hosts).Where(fmt.Sprintf("type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
			hostType, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND (name like '%%%s%%' OR ip like '%%%s%%') AND status=1 ",
				hostType, filter, filter)).Find(&hosts).Error; err != nil {
			logrus.Errorf("Error:%v", err)
			return
		}
	} else {
		dao.ClientDB.Model(&hosts).Where(fmt.Sprintf("name like '%%%s%%' OR ip like '%%%s%%' ", filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("name like '%%%s%%' OR ip like '%%%s%%' ", filter, filter)).Find(&hosts).Error; err != nil {
			logrus.Errorf("Error:%v", err)
			return
		}
	}
	return
}

// 获取公共主机列表，支持分页，按名称、IP过滤
func (h *HostOrm) GetPublicHostList(index, count int, hostType contants.AppType, filter string) (totalNum int64, hosts []Host, err error) {
	if hostType != 0 {
		dao.ClientDB.Model(&hosts).Where(fmt.Sprintf("type = %d AND auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
			hostType, contants.PUBLIC, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')",
				hostType, contants.PUBLIC, filter, filter)).Find(&hosts).Error; err != nil {
			logrus.Error(err)
			return
		}
	} else {
		dao.ClientDB.Model(&hosts).Where(fmt.Sprintf("auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')", contants.PUBLIC, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("auth_type = %d AND (name like '%%%s%%' OR ip like '%%%s%%')", contants.PUBLIC, filter, filter)).Find(&hosts).Error; err != nil {
			logrus.Error(err)
			return
		}
	}
	return
}

// 获取私有主机列表，支持分页，按名称、IP过滤
func (h *HostOrm) GetPrivateHostList(index, count int, hostType contants.AppType, userName,
	filter string) (totalNum int64, hosts []Host, err error) {
	if hostType != 0 {
		dao.ClientDB.Model(&hosts).
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				hostType, contants.PRIVATE, userName, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				hostType, contants.PRIVATE, userName, filter, filter)).Find(&hosts).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	} else {
		dao.ClientDB.Model(&hosts).
			Where(fmt.Sprintf("auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				contants.PRIVATE, userName, filter, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("auth_type = %d AND create_user = '%s' AND (name like '%%%s%%' OR ip like '%%%s%%')",
				contants.PRIVATE, userName, filter, filter)).Find(&hosts).Error; err != nil {
			logrus.Error(err)
			return
		}
	}

	return
}

func (h *HostOrm) GetHostStat(userName string) (PublicNum, PrivateNum int) {
	var hosts []Host
	dao.ClientDB.Model(&hosts).Where("auth_type = ?", contants.PUBLIC).Count(&PublicNum)
	dao.ClientDB.Model(&hosts).
		Where(fmt.Sprintf("auth_type = %d AND create_user = '%s'", contants.PRIVATE, userName)).Count(&PrivateNum)
	return
}

// 根据请求参数获取主机数量
func (h *HostOrm) GetHostCountByRequest(request string) (totalNum int64) {
	var hosts []Host
	dao.ClientDB.Model(&hosts).Where(request).Count(&totalNum)
	return
}

// 根据请求参数获取主机列表, 支持分页
func (h *HostOrm) GetHostListByRequest(index, count int, request string) (hosts []Host, err error) {
	if err = dao.ClientDB.Limit(count).Offset(index).
		Where(request).Find(&hosts).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (h *HostOrm) CreateDeployApp(deployApp DeployApp) (result DeployApp, err error) {
	err = dao.ClientDB.Create(&deployApp).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = deployApp
	return
}

func (h *HostOrm) UpdateDeployApp(deployApp DeployApp) (result DeployApp, err error) {
	if err = dao.ClientDB.Model(&deployApp).Update(&deployApp).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = deployApp
	return
}

func (h *HostOrm) FindLastDeployById(id uint) (deployApp DeployApp, err error) {
	err = dao.ClientDB.Where("id = ?", id).Last(&deployApp).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (h *HostOrm) FindLastDeployByHost(id uint) (deployApp DeployApp, err error) {
	err = dao.ClientDB.Where("hostId = ?", id).Last(&deployApp).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (h *HostOrm) GetDeployAppByHostId(hostId uint, index, count int) (totalNum int64, deployApp []DeployApp, err error) {
	dao.ClientDB.Model(&deployApp).Where("hostId = ?", hostId).Count(&totalNum)
	if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
		Where("hostId = ?", hostId).Find(&deployApp).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}
