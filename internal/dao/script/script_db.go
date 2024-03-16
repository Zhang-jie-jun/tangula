package script

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/sirupsen/logrus"
)

type ScriptOrm struct{}

var ScriptMgm = ScriptOrm{}

func (s *ScriptOrm) FindById(id uint) (script Script, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&script).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (s *ScriptOrm) FindByUuid(uuid string) (script Script, err error) {
	logrus.Infof("uuid:%s", uuid)
	err = dao.ClientDB.Where("uuid = ?", uuid).Find(&script).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (s *ScriptOrm) CreateScript(script Script) (result Script, err error) {
	err = dao.ClientDB.Create(&script).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	result = script
	return
}

func (s *ScriptOrm) UpdateScript(script Script) (result Script, err error) {
	if err = dao.ClientDB.Model(&script).Update(&script).Error; err != nil {
		logrus.Error(err)
		return result, err
	}
	result = script
	return
}

func (s *ScriptOrm) DeleteScript(id uint) error {
	var script Script
	err := dao.ClientDB.Where("id = ?", id).Find(&script).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = dao.ClientDB.Delete(&script).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 获取脚本列表，支持分页，按名称过滤
func (s *ScriptOrm) GetScripts(index, count int, user, filter string) (totalNum int64, script []Script, err error) {
	dao.ClientDB.Model(&script).Where(fmt.Sprintf("create_user = '%s' AND name like '%%%s%%'",
		user, filter)).Count(&totalNum)
	if err = dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
		Where(fmt.Sprintf("create_user = '%s' AND name like '%%%s%%'", user, filter)).
		Find(&script).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (s *ScriptOrm) GetScriptStat(userName string) (TotalNum int) {
	var scripts []Script
	dao.ClientDB.Model(&scripts).Where("create_user = ?", userName).Count(&TotalNum)
	return
}
