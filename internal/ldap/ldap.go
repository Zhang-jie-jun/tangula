package ldap

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
)

func InitLDAP(Addr, userName, passWord string) error {
	passWord, err := util.AesDecrypt(passWord)
	if err != nil {
		logrus.Error(err)
		return err
	}
	conn, err := ldap.Dial("tcp", Addr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer conn.Close()
	err = conn.Bind(userName, passWord)
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = GetAllUserInfo(conn)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return err
}

// 服务启动时获取所有用户并入库
func GetAllUserInfo(conn *ldap.Conn) error {
	searchRequest := ldap.NewSearchRequest(
		contants.AppCfg.LDAP.SearchDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 60, false,
		"(&(objectClass=organizationalPerson))",
		[]string{"sAMAccountName", "displayName", "mail", "mobile"},
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		logrus.Error(err)
		return err
	}
	var users []map[string]string
	for _, entry := range sr.Entries {
		name := entry.GetAttributeValue("displayName")
		account := entry.GetAttributeValue("sAMAccountName")
		mail := entry.GetAttributeValue("mail")
		phone := entry.GetAttributeValue("mobile")
		if ok := auth.AuthMgm.CheckUserExistByMail(mail); ok {
			// 用户已存在则跳过
			continue
		}
		user := map[string]string{"Name": name, "Account": account, "Mail": mail, "Phone": phone}
		users = append(users, user)
	}
	err = auth.AuthMgm.BatchCreateUser(users)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func Login(userName, passWord string) bool {
	// 建立连接
	conn, err := ldap.Dial("tcp", contants.AppCfg.LDAP.Addr)
	if err != nil {
		logrus.Error(err)
		return false
	}
	defer conn.Close()
	BindPassWord, err := util.AesDecrypt(contants.AppCfg.LDAP.BindPassword)
	if err != nil {
		logrus.Error(err)
		return false
	}
	err = conn.Bind(contants.AppCfg.LDAP.BindUserName, BindPassWord)
	if err != nil {
		logrus.Error(err)
		return false
	}
	// 获取登录用户信息
	filter := fmt.Sprintf("(&(objectClass=organizationalPerson)(sAMAccountName=%s))", userName)
	searchRequest := ldap.NewSearchRequest(
		contants.AppCfg.LDAP.SearchDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn", "sAMAccountName", "displayName", "mail", "mobile"},
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		logrus.Error(err)
		return false
	}
	if len(sr.Entries) != 1 {
		return false
	}
	// 认证用户密码
	userDN := sr.Entries[0].DN
	err = conn.Bind(userDN, passWord)
	if err != nil {
		logrus.Error(err)
		return false
	}

	// 登录成功后添加或更新用户信息到数据库
	name := sr.Entries[0].GetAttributeValue("displayName")
	account := sr.Entries[0].GetAttributeValue("sAMAccountName")
	mail := sr.Entries[0].GetAttributeValue("mail")
	phone := sr.Entries[0].GetAttributeValue("mobile")
	_, err = auth.AuthMgm.CreateOrUpdateUser(name, account, mail, phone)
	if err != nil {
		logrus.Error(err)
	}

	return true
}
