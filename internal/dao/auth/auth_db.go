package auth

import (
	"bytes"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/sirupsen/logrus"
	"time"
)

type Auth struct{}

var AuthMgm = Auth{}

func (auth *Auth) FindById(id uint) (user User, err error) {
	err = dao.ClientDB.Where("id = ?", id).Preload("Role").Find(&user).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) FindByName(userName string) (user User, err error) {
	err = dao.ClientDB.Where("name = ?", userName).Preload("Role").Find(&user).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) FindByAccount(account string) (user User, err error) {
	err = dao.ClientDB.Where("account = ?", account).Preload("Role").Find(&user).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) FindByMail(mail string) (user User, err error) {
	err = dao.ClientDB.Where("mail = ?", mail).Preload("Role").Find(&user).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 根据用户账号检查用户是否存在
func (auth *Auth) CheckUserExistByAccount(account string) bool {
	var user User
	err := dao.ClientDB.Select("account").Where("account = ?", account).First(&user).Error
	if err != nil {
		logrus.Error(err)
		return false
	}
	if len(user.Account) > 0 {
		return true
	}
	return false
}

// 根据用户邮箱检查用户是否存在
func (auth *Auth) CheckUserExistByMail(mail string) bool {
	var user User
	dao.ClientDB.Select("mail").Where("mail = ?", mail).First(&user)
	if len(user.Mail) > 0 {
		return true
	}
	return false
}

func (auth *Auth) CreateUser(name, account, mail, phone string) (user User, err error) {
	if ok := auth.CheckUserExistByMail(mail); ok {
		err = errors.New(msg.ERROR_USER_IS_EXIST, msg.GetMsg(msg.ERROR_USER_IS_EXIST))
		logrus.Error(err)
		return
	}
	var role Role
	err = dao.ClientDB.Where("name = ?", contants.USER).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	user = User{Name: name, Account: account, Mail: mail, Phone: phone, Role: role, RoleId: role.Id}
	err = dao.ClientDB.Create(&user).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) CreateOrUpdateUser(name, account, mail, phone string) (user User, err error) {
	if ok := auth.CheckUserExistByMail(mail); ok {
		err = dao.ClientDB.Where("mail = ?", mail).Preload("Role").Find(&user).Error
		if err != nil {
			logrus.Error(err)
			return
		}
		user.Name = name
		user.Account = account
		user.Mail = mail
		user.Phone = phone
		user.UsageCount = user.UsageCount + 1
		if err = dao.ClientDB.Model(&user).Update(&user).Error; err != nil {
			logrus.Error(err)
			return
		}
	} else {
		var role Role
		err = dao.ClientDB.Where("name = ?", contants.USER).Find(&role).Error
		if err != nil {
			logrus.Error(err)
			return
		}
		user = User{Name: name, Account: account, Mail: mail, Phone: phone, Role: role, RoleId: role.Id}
		err = dao.ClientDB.Create(&user).Error
		if err != nil {
			logrus.Error(err)
			return
		}
	}
	return
}

func (auth *Auth) BatchCreateUser(users []map[string]string) (err error) {
	if len(users) == 0 {
		return
	}
	var role Role
	err = dao.ClientDB.Where("name = ?", contants.USER).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	var buffer bytes.Buffer
	sql := "insert into `tangula_user` (`name`,`account`,`mail`,`phone`,`created_time`,`updated_time`,`role_id`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	createTime := time.Now().Format("2006-01-02 15:04:05")
	updateTime := createTime
	for i, user := range users {
		if i == len(users)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%d');",
				user["Name"], user["Account"], user["Mail"], user["Phone"], createTime, updateTime, role.Id))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%d'),",
				user["Name"], user["Account"], user["Mail"], user["Phone"], createTime, updateTime, role.Id))
		}
	}
	return dao.ClientDB.Exec(buffer.String()).Error
}

// 获取用户列表，根据角色、名称，用户、邮箱过滤(like模糊查询，后续需要考虑性能问题)
func (auth *Auth) GetUserList(index, count int, roleName, filter string) (totalNum int64, users []User, err error) {
	if roleName != "" {
		var role Role
		err = dao.ClientDB.Where("name = ?", roleName).Find(&role).Error
		if err != nil {
			logrus.Error(err)
			return
		}
		dao.ClientDB.Model(&users).
			Where(fmt.Sprintf("role_id = %d AND (name like '%%%s%%' OR account like '%%%s%%' OR mail like '%%%s%%')",
				role.Id, filter, filter, filter)).Count(&totalNum)
		if err := dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
			Where(fmt.Sprintf("role_id = %d AND (name like '%%%s%%' OR account like '%%%s%%' OR mail like '%%%s%%')",
				role.Id, filter, filter, filter)).Preload("Role").Find(&users).Error; err != nil {
			logrus.Error(err)
			return totalNum, users, err
		}
	} else {
		dao.ClientDB.Model(&users).Where(fmt.Sprintf("name like '%%%s%%' OR account like '%%%s%%' OR mail like '%%%s%%'",
			filter, filter, filter)).Count(&totalNum)
		if err := dao.ClientDB.Limit(count).Offset(index).Order("id DESC", true).
			Where(fmt.Sprintf("name like '%%%s%%' OR account like '%%%s%%' OR mail like '%%%s%%'",
				filter, filter, filter)).Preload("Role").Find(&users).Error; err != nil {
			logrus.Error(err)
			return totalNum, users, err
		}
	}
	return totalNum, users, err
}

func (auth *Auth) GetUserStat() (TotalNum, ActiveNum, TodayNum int) {
	var users []User
	dao.ClientDB.Model(&users).Count(&TotalNum)
	dao.ClientDB.Model(&users).Where("usage_count != 0").Count(&ActiveNum)
	k := time.Now()
	sd, err := time.ParseDuration("-24h")
	if err != nil {
		logrus.Error(err)
		return
	}
	dao.ClientDB.Model(&users).Where("updated_time >= ? ", k.Add(sd)).Count(&TodayNum)
	return
}

// 编辑用户角色
func (auth *Auth) SetUserRole(id uint, roleType contants.RoleType) (user User, err error) {
	var role Role
	if role, err = auth.FindRoleByName(roleType); err != nil {
		err = errors.New(msg.ERROR_USER_NOT_EXIST, msg.GetMsg(msg.ERROR_USER_NOT_EXIST))
		return
	}
	if user, err = auth.FindById(id); err != nil {
		err = errors.New(msg.ERROR_USER_NOT_EXIST, msg.GetMsg(msg.ERROR_USER_NOT_EXIST))
		return
	}
	user.RoleId = role.Id
	user.Role = role
	if err = dao.ClientDB.Model(&user).Update(&user).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 启用用户
func (auth *Auth) EnableUser(id uint) (user User, err error) {
	if user, err = auth.FindById(id); err != nil {
		err = errors.New(msg.ERROR_USER_NOT_EXIST, msg.GetMsg(msg.ERROR_USER_NOT_EXIST))
		return
	}
	if user.Status != contants.DISABLE {
		return
	}
	user.Status = contants.ENABLE
	if err = dao.ClientDB.Model(&user).Update(&user).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 禁用用户
func (auth *Auth) DisableUser(id uint) (user User, err error) {
	if user, err = auth.FindById(id); err != nil {
		err = errors.New(msg.ERROR_USER_NOT_EXIST, msg.GetMsg(msg.ERROR_USER_NOT_EXIST))
		return
	}
	if user.Status != contants.ENABLE {
		return
	}
	user.Status = contants.DISABLE
	if err = dao.ClientDB.Model(&user).Update(&user).Error; err != nil {
		logrus.Error(err)
		return
	}
	return
}

// Role管理
func (auth *Auth) FindRoleById(id uint) (role Role, err error) {
	err = dao.ClientDB.Where("id = ?", id).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) FindRoleByName(roleName contants.RoleType) (role Role, err error) {
	err = dao.ClientDB.Where("name = ?", roleName).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) GetUserRole() (role Role, err error) {
	err = dao.ClientDB.Where("name = ?", contants.USER).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) GetSysRole() (role Role, err error) {
	err = dao.ClientDB.Where("name = ?", contants.SYSADMIN).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (auth *Auth) GetSuperRole() (role Role, err error) {
	err = dao.ClientDB.Where("name = ?", contants.SUPERADMIN).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 获取角色列表，根据名称过滤
func (auth *Auth) GetRoleList(index, count int, filter string) (totalNum int64, roles []Role, err error) {
	dao.ClientDB.Model(&roles).Where(fmt.Sprintf(" name like '%%%s%%' ", filter)).Count(&totalNum)
	if err := dao.ClientDB.Limit(count).Offset(index).Order("id DESC").
		Where(fmt.Sprintf(" name like '%%%s%%' ", filter)).Find(&roles).Error; err != nil {
		logrus.Error(err)
		return totalNum, roles, err
	}
	return totalNum, roles, err
}

func (auth *Auth) BatchCreateRole(users []map[string]string) (err error) {
	if len(users) == 0 {
		return
	}
	var buffer bytes.Buffer
	sql := "insert into `tangula_role` (`name`,`desc`,`created_time`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		logrus.Error(err)
		return err
	}
	createTime := time.Now().Format("2006-01-02 15:04:05")
	for i, user := range users {
		var role Role
		_ = dao.ClientDB.Where("name = ?", user["Name"]).Find(&role).Error
		// 已存在则跳过
		if role.Id != 0 {
			continue
		}
		if i == len(users)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s');",
				user["Name"], user["Desc"], createTime))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s'),",
				user["Name"], user["Desc"], createTime))
		}
	}
	if buffer.String() == sql {
		return
	}
	return dao.ClientDB.Exec(buffer.String()).Error
}

func (auth *Auth) CreateRole(name contants.RoleType, desc string) (role Role, err error) {
	err = dao.ClientDB.Where("name = ?", name).Find(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	if role.Id != 0 {
		err = errors.New(msg.ERROR_ROLE_IS_EXIST, msg.GetMsg(msg.ERROR_ROLE_IS_EXIST))
		logrus.Error(err)
		return
	}
	role.Name = name
	role.Desc = desc
	err = dao.ClientDB.Create(&role).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}
