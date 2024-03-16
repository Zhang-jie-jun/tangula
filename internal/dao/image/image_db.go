package image

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type ImageOrm struct{}

var ImageMgm = ImageOrm{}

func (i *ImageOrm) FindById(id uint) (image Image, err error) {
	err = dao.ClientDB.Where("id = ?", id).Preload("Pool").Find(&image).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (i *ImageOrm) FindByUuid(uuid string) (image Image, err error) {
	err = dao.ClientDB.Where("uuid = ?", uuid).Preload("Pool").Find(&image).Error
	if err != nil {
		logrus.Errorf("FindByUuid Error:%v", err)
		return
	}
	return
}

// 根据名称查询，名称可重复，取第一条数据
func (i *ImageOrm) FindByName(name string) (image Image, err error) {
	err = dao.ClientDB.Where("name = ?", name).Preload("Pool").First(&image).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 检查是否存在相同名称的镜像
func (i *ImageOrm) CheckIsExistByName(storePoolName string) bool {
	var image Image
	err := dao.ClientDB.Where("name = ?", storePoolName).First(&image).Error
	if err != nil {
		logrus.Info(err)
		return false
	}
	return true
}

func (i *ImageOrm) CreateImage(image Image) (result Image, err error) {
	err = dao.ClientDB.Create(&image).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = image
	return
}

func (i *ImageOrm) UpdateImage(image Image) (result Image, err error) {
	if err = dao.ClientDB.Model(&image).Update(&image).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = image
	return
}

func (i *ImageOrm) UpdateImageStatus(id uint, status contants.ImageStatus) error {
	var image Image
	err := dao.ClientDB.Where("id = ?", id).Find(&image).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	image.Status = status
	if err = dao.ClientDB.Model(&image).Update(&image).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (i *ImageOrm) UpdateImageStatusAndUuid(id uint, uuid string, status contants.ImageStatus) error {
	var image Image
	err := dao.ClientDB.Where("id = ?", id).Find(&image).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	image.Status = status
	image.Uuid = uuid
	if err = dao.ClientDB.Model(&image).Update(&image).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (i *ImageOrm) DeleteImage(id uint) error {
	var image Image
	err := dao.ClientDB.Where("id = ?", id).Find(&image).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&image).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取平台列表，支持分页，按资源所属、类型、名称、IP过滤
func (i *ImageOrm) GetImageList(index, count int, imageType contants.ImageType, filter string) (
	totalNum int64, images []Image, err error) {
	if imageType != 0 {
		dao.ClientDB.Model(&images).
			Where(fmt.Sprintf("type = %d AND name like '%%%s%%'", imageType, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND name like '%%%s%%'", imageType, filter)).
			Preload("Pool").Find(&images).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	} else {
		dao.ClientDB.Model(&images).Where(fmt.Sprintf("name like '%%%s%%'", filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("name like '%%%s%%'", filter)).
			Preload("Pool").Find(&images).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	}
}

// 获取公共平台列表，支持分页，按类型、名称、IP过滤
func (i *ImageOrm) GetPublicImageList(index, count int, imageType contants.ImageType, filter string) (
	totalNum int64, images []Image, err error) {
	if imageType != 0 {
		dao.ClientDB.Model(&images).
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND name like '%%%s%%'",
				imageType, contants.PUBLIC, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND name like '%%%s%%'",
				imageType, contants.PUBLIC, filter)).Preload("Pool").Find(&images).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	} else {
		dao.ClientDB.Model(&images).
			Where(fmt.Sprintf("auth_type = %d AND name like '%%%s%%'",
				contants.PUBLIC, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("auth_type = %d AND name like '%%%s%%'",
				contants.PUBLIC, filter)).Preload("Pool").Find(&images).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	}
}

// 获取私有平台列表，支持分页，按类型、名称、IP过滤
func (i *ImageOrm) GetPrivateImageList(index, count int, userName string, imageType contants.ImageType,
	filter string) (totalNum int64, images []Image, err error) {
	if imageType != 0 {
		dao.ClientDB.Model(&images).
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND create_user = '%s' AND name like '%%%s%%'",
				imageType, contants.PRIVATE, userName, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("type = %d AND auth_type = %d AND create_user = '%s' AND name like '%%%s%%'",
				imageType, contants.PRIVATE, userName, filter)).Preload("Pool").Find(&images).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	} else {
		dao.ClientDB.Model(&images).
			Where(fmt.Sprintf("auth_type = %d AND create_user = '%s' AND name like '%%%s%%'",
				contants.PRIVATE, userName, filter)).Count(&totalNum)
		if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("auth_type = %d AND create_user = '%s' AND name like '%%%s%%'",
				contants.PRIVATE, userName, filter)).Preload("Pool").Find(&images).Error; err != nil {
			logrus.Error(err)
			return
		}
		return
	}
}

func (i *ImageOrm) GetImageStat(userName string) (PublicNum, PrivateNum int) {
	var images []Image
	dao.ClientDB.Model(&images).Where("auth_type = ?", contants.PUBLIC).Count(&PublicNum)
	dao.ClientDB.Model(&images).
		Where(fmt.Sprintf("auth_type = %d AND create_user = '%s'", contants.PRIVATE, userName)).Count(&PrivateNum)
	return
}

func (i *ImageOrm) GetImageByPoolId(poolId uint) ([]Image, error) {
	var images []Image
	if err := dao.ClientDB.Where("pool_id = ?", poolId).Find(&images).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return images, nil
}

func (i *ImageOrm) DeleteImageByPoolId(poolId uint) error {
	var images Image
	err := dao.ClientDB.Where("pool_id = ?", poolId).Find(&images).Error
	if err != nil {
		logrus.Info("存储池不包含镜像:", poolId)
	} else {
		err = dao.ClientDB.Delete(&images).Error
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}
