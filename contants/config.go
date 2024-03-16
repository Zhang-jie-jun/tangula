package contants

// 对应app.ini文件的结构体
type AppConfig struct {
	// 系统配置
	System struct {
		AppName          string `ini:"APP_NAME"`            // 应用程序名称
		RunMode          string `ini:"RUN_MODE"`            // 运行模式
		MountPath        string `ini:"MOUNT_PATH"`          // 镜像挂载路径
		FtpPath          string `ini:"FTP_PATH"`            // FTP服务根路径
		ScriptPath       string `ini:"SCRIPT_PATH"`         // 脚本存放路径
		LogFile          string `ini:"LOG_FILE"`            // 系统log文件
		LogLevel         int    `ini:"LOG_LEVEL"`           // [Panic],[Fatal],[Error],[Warn],[Info],[Debug],[Trace](0-6)
		LogReserveCount  int    `ini:"LOG_RESERVE_COUNT"`   // 保留几个日志文件备份
		LogReserveMaxDay int    `ini:"LOG_RESERVE_MAX_DAY"` // 保留最长时间30天
		SuperAdminName   string `ini:"SUPER_ADMIN_NAME"`    // 系统超级管理员名称
		SuperAdminPswd   string `ini:"SUPER_ADMIN_PSWD"`    // 系统超级管理员密码
		ATName           string `ini:"AT_NAME"`             // 自动化测试账号
		ATPswd           string `ini:"AT_PASSWORD"`         // 自动化测试密码

	}
	// 服务配置
	Server struct {
		Host         string `ini:"HTTP_HOST"`    // 服务连接IP
		Port         int    `ini:"HTTP_PORT"`    // 服务连接端口
		ReadTimeOut  int    `ini:"READ_TIMEOUT"` // 超时时间
		WriteTimeOut int    `ini:"WRITE_TIMEOUT"`
	}
	// 应用配置
	App struct {
		PageSize    int    `ini:"PAGE_SIZE"`
		IdentityKey string `ini:"IDENTITY_KEY"`
	}
	// 数据库配置
	Database struct {
		Type     string `ini:"TYPE"`     // 数据库类型(mysql,sql3)
		User     string `ini:"USER"`     // 数据库用户名
		Password string `ini:"PASSWORD"` // 数据库密码(AES加密)
		Host     string `ini:"HOST"`     // 数据库连接IP
		DBName   string `ini:"NAME"`     // 数据库名称
	}
	// LDAP配置
	LDAP struct {
		Addr         string `ini:"ADDR"`         // 服务地址
		BindUserName string `ini:"BINDUSERNAME"` // 管理用户名称
		BindPassword string `ini:"BINDPASSWORD"` // 管理用户密码
		SearchDN     string `ini:"SEARCHDN"`     // 账号路径
	}
	//CAS配置
	CAS struct {
		Username string `ini:"USERNAME"`
		Password string `ini:"PASSWORD"`
		Url      string `ini:"URL"`
	}

	//jenkins配置
	JENKINS struct {
		Username            string `ini:"USERNAME"`
		Password            string `ini:"PASSWORD"`
		Url                 string `ini:"URL"`
		JobNameCas          string `ini:"JOBNAMECAS"`
		INSTALL_CLIENT_NAME string `ini:"INSTALL_CLIENT_NAME"`
	}
}
