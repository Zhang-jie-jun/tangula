package jwt

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/ldap"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var globalUser string

// 定义JWT中间件方法
func JWTMiddlewareInit(jwtAuthorizator Authorizator) (authMiddleware *jwt.GinJWTMiddleware) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 60,
		MaxRefresh:  time.Hour,
		IdentityKey: contants.AppCfg.App.IdentityKey,
		PayloadFunc: func(data interface{}) (result jwt.MapClaims) {
			if user, ok := data.(*auth.User); ok {
				// 通过用户名获取用户信息
				var userInfo auth.User
				var err error
				if user.Account == contants.AppCfg.System.SuperAdminName {
					role, err := auth.AuthMgm.GetSuperRole()
					if err != nil {
						return result
					}
					userInfo = auth.User{
						Name:    fmt.Sprintf("超级管理员(%s)", user.Account),
						Account: user.Account,
						Mail:    user.Account,
						Phone:   "",
						RoleId:  role.Id,
					}
				} else if user.Account == contants.AppCfg.System.ATName {
					role, err := auth.AuthMgm.GetSuperRole()
					if err != nil {
						return result
					}
					userInfo = auth.User{
						Name:    fmt.Sprintf("自动化测试(%s)", user.Account),
						Account: user.Account,
						Mail:    user.Account,
						Phone:   "",
						RoleId:  role.Id,
					}
				} else {
					userInfo, err = auth.AuthMgm.FindByAccount(user.Account)
					if err != nil {
						return result
					}
				}
				globalUser = userInfo.Account
				// 设置jwt身份校验信息
				result = jwt.MapClaims{
					"name":    userInfo.Name,
					"account": userInfo.Account,
					"mail":    userInfo.Mail,
					"phone":   userInfo.Phone,
					"role_id": userInfo.RoleId,
				}
			}
			return result
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			Name, ok1 := claims["name"].(string)
			Account, ok2 := claims["account"].(string)
			Mail, ok3 := claims["mail"].(string)
			Phone, ok4 := claims["phone"].(string)
			RoleId, ok5 := claims["role_id"].(float64)
			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
				return nil
			}
			//Set the identity
			return &auth.User{
				Name:    Name,
				Account: Account,
				Mail:    Mail,
				Phone:   Phone,
				RoleId:  uint(RoleId),
			}
		},
		Authenticator: func(c *gin.Context) (user interface{}, err error) {
			var loginInfo view.UserParam
			// 闭包的方式记录操作日志
			defer func() {
				if err != nil {
					detail := msg.GetOperation(msg.LOGIN_FAILED, loginInfo.UserName, err.Error())
					service.CreateLogRecord(msg.LOGIN, loginInfo.UserName, detail, loginInfo.UserName, contants.LOG_FAILED)
				} else {
					detail := msg.GetOperation(msg.LOGIN_SUCCESS, loginInfo.UserName)
					service.CreateLogRecord(msg.LOGIN, loginInfo.UserName, detail, loginInfo.UserName, contants.LOG_SUCCESS)
				}
			}()
			if err = c.ShouldBind(&loginInfo); err != nil {
				logrus.Errorf("Error:%v\n", err)
				err = jwt.ErrMissingLoginValues
				return
			}
			// 超级管理员认证
			if loginInfo.UserName == contants.AppCfg.System.SuperAdminName {
				passWord, e := util.AesDecrypt(contants.AppCfg.System.SuperAdminPswd)
				if e != nil {
					logrus.Error(e)
					err = jwt.ErrMissingLoginValues
					return
				}
				if loginInfo.PassWord == passWord {
					user = &auth.User{Account: loginInfo.UserName}
					return
				}
			} else if loginInfo.UserName == contants.AppCfg.System.ATName {
				passWord, e := util.AesDecrypt(contants.AppCfg.System.ATPswd)
				if e != nil {
					logrus.Error(e)
					err = jwt.ErrMissingLoginValues
					return
				}
				if loginInfo.PassWord == passWord {
					user = &auth.User{Account: loginInfo.UserName}
					return
				}
			} else {
				if ldap.Login(loginInfo.UserName, loginInfo.PassWord) {
					user = &auth.User{Account: loginInfo.UserName}
					return
				}
			}
			err = jwt.ErrMissingLoginValues
			return
		},
		//receives identity and handles authorization logic
		Authorizator: jwtAuthorizator.HandleAuthorizator,
		//handles unauthorized logic
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":     msg.ERROR,
				"message":  msg.GetMsg(msg.ERROR, message),
				"response": nil})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			data := map[string]interface{}{"token": token, "expire": expire}
			c.JSON(http.StatusOK, gin.H{"code": msg.SUCCESS, "message": msg.GetMsg(msg.SUCCESS), "response": data})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			data := map[string]interface{}{"token": token, "expire": expire}
			c.JSON(http.StatusOK, gin.H{"code": msg.SUCCESS, "message": msg.GetMsg(msg.SUCCESS), "response": data})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			var account string
			if globalUser != "" {
				account = globalUser
			} else {
				account = "未知用户"
			}
			//claims := jwt.ExtractClaims(c)
			//account, _ := claims["account"].(string)
			detail := msg.GetOperation(msg.LOGOUT_SUCCESS, account)
			//service.CreateLogRecord(msg.LOGOUT, account, detail, account, contants.LOG_SUCCESS)
			service.CreateLogRecord(msg.LOGOUT, account, detail, globalUser, contants.LOG_SUCCESS)
			c.JSON(http.StatusOK, gin.H{"code": msg.SUCCESS, "message": msg.GetMsg(msg.SUCCESS), "response": nil})
		},
		// TokenLookup is a string in the system of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Tangula"
		TokenHeadName: "Tangula",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		logrus.Fatal("JWT Error:" + err.Error())
	}
	return
}
