package windows

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/masterzen/winrm"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

type winRmClient struct {
	loginInfo *contants.LoginInfo
	client    *winrm.Client
}

func NewWinRmClient(login *contants.LoginInfo) (client *winRmClient, err error) {
	c := &winRmClient{loginInfo: login}
	err = c.Login()
	if err == nil {
		client = c
	}
	return
}

func (c *winRmClient) Login() (err error) {
	if c.client != nil {
		return nil
	}
	endpoint := winrm.NewEndpoint(
		c.loginInfo.Ip,
		int(c.loginInfo.Port),
		false,
		false,
		nil,
		nil,
		nil,
		60*time.Second)
	c.client, err = winrm.NewClient(endpoint, c.loginInfo.UserName, c.loginInfo.PassWord)
	if err != nil {
		err = errors.New(msg.ERROR_LOGIN_REMOTE_HOST, msg.GetMsg(msg.ERROR_LOGIN_REMOTE_HOST, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	if c.client == nil {
		err = errors.New(msg.ERROR_LOGIN_REMOTE_HOST,
			msg.GetMsg(msg.ERROR_LOGIN_REMOTE_HOST, c.loginInfo.Ip, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return
	}
	return
}

func (c *winRmClient) GetSystemInfo() (hostName, hostOs, arch string, err error) {
	cmd := "systeminfo"
	result, err := c.RunCommond(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REMOTE_HOST_INFO,
			msg.GetMsg(msg.ERROR_GET_REMOTE_HOST_INFO, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	for _, line := range result {
		if strings.Contains(line, "主机名:") || strings.Contains(line, "Host Name:") {
			hostName = strings.Split(line, ":")[1]
		} else if strings.Contains(line, "OS 名称:") || strings.Contains(line, "OS Name:") {
			hostOs = strings.Split(line, ":")[1]
		} else if strings.Contains(line, "系统类型:") || strings.Contains(line, "System Type:") {
			arch = strings.Split(line, ":")[1]
		}
	}
	return
}

func (c *winRmClient) Mount(sharePath string) (string, error) {
	// 1.检查服务状态
	if !c.CheackNfsStatus() {
		err := errors.New(msg.ERROR_REMOTE_HOST_NFS_NOT_START,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_NFS_NOT_START, c.loginInfo.Ip))
		logrus.Error(err)
		return "", err
	}
	// 2.获取可用盘符
	drive, err := c.GetUnUsedDrive()
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	// 3.挂载nfs共享到指定盘符
	cmd := fmt.Sprintf("net use %s %s", drive, sharePath)
	result, err := c.RunCommond(cmd)
	if err != nil {
		err := errors.New(msg.ERROR_REMOTE_HOST_MOUNT_FAILED,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_MOUNT_FAILED, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return "", err
	}
	if len(result) != 1 || !strings.Contains(result[0], "The command completed successfully.") {
		err := errors.New(msg.ERROR_REMOTE_HOST_MOUNT_FAILED, msg.GetMsg(msg.ERROR_REMOTE_HOST_MOUNT_FAILED,
			c.loginInfo.Ip, msg.GetMsg(msg.ERROR_EXPECTED_RESULT_IS_INCORRECT)))
		logrus.Error(err)
		return "", err
	}
	return drive, nil
}

func (c *winRmClient) UnMount(sharePath string) error {
	// 1.根据共享路径获取盘符
	drive, err := c.GetDriveByPath(sharePath)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 2.卸载NFS挂载
	cmd := fmt.Sprintf("net use %s /delete", drive)
	result, err := c.RunCommond(cmd)
	if err != nil {
		err := errors.New(msg.ERROR_REMOTE_HOST_UN_MOUNT_FAILED,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_UN_MOUNT_FAILED, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return err
	}
	if len(result) != 1 || !strings.Contains(result[0], "was deleted successfully.") {
		err := errors.New(msg.ERROR_REMOTE_HOST_UN_MOUNT_FAILED, msg.GetMsg(msg.ERROR_REMOTE_HOST_UN_MOUNT_FAILED,
			c.loginInfo.Ip, msg.GetMsg(msg.ERROR_EXPECTED_RESULT_IS_INCORRECT)))
		logrus.Error(err)
		return err
	}
	return nil
}

func (c *winRmClient) CreateShortcut(shortcutName, sharePath string) error {
	// 获取桌面路径
	path, err := c.GetDesktopPath()
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 创建桌面快捷方式 cmd: mklink /D C:\\Users\\Administrator\\Desktop\\dirName \\\\192.168.212.32\\mnt\\share
	cmd := fmt.Sprintf("mklink /D %s\\%s %s", path, shortcutName, sharePath)
	result, err := c.RunCommond(cmd)
	if err != nil {
		err := errors.New(msg.ERROR_REMOTE_HOST_CREATE_SHORTCUT,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_CREATE_SHORTCUT, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return err
	}
	if len(result) != 1 || !strings.Contains(result[0], "symbolic link created for") {
		logrus.Error(result)
	}
	return nil
}

func (c *winRmClient) DeleteShortcut(shortcutName string) error {
	// 获取桌面路径
	path, err := c.GetDesktopPath()
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 删除桌面快捷方式
	cmd := fmt.Sprintf("rd/s/q %s\\%s", path, shortcutName)
	_, err = c.RunCommond(cmd)
	if err != nil {
		err := errors.New(msg.ERROR_REMOTE_HOST_DELETE_SHORTCUT,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_DELETE_SHORTCUT, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return err
	}
	return nil
}

func (c *winRmClient) CheackNfsStatus() bool {
	cmd := "sc query nfsclnt"
	result, err := c.RunCommond(cmd)
	if err != nil {
		logrus.Error(err)
		return false
	}
	var nfsStatus string
	for _, line := range result {
		if strings.Contains(line, "STATE:") {
			nfsStatus = strings.Split(line, ":")[1]
		}
	}
	if strings.Contains(nfsStatus, "RUNNING") {
		return true
	}
	return false
}

func (c *winRmClient) GetDesktopPath() (string, error) {
	cmd := "chdir"
	result, err := c.RunCommond(cmd)
	if err != nil {
		err := errors.New(msg.ERROR_REMOTE_HOST_GET_DESKTOP_PATH,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_GET_DESKTOP_PATH, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return "", err
	}
	if len(result) != 1 {
		err := errors.New(msg.ERROR_REMOTE_HOST_GET_DESKTOP_PATH, msg.GetMsg(msg.ERROR_REMOTE_HOST_GET_DESKTOP_PATH,
			c.loginInfo.Ip, msg.GetMsg(msg.ERROR_EXPECTED_RESULT_IS_INCORRECT)))
		logrus.Error(err)
		return "", err
	}
	path := result[0] + "\\Desktop"
	return path, nil
}

// 获取一个未使用的盘符
func (c *winRmClient) GetUnUsedDrive() (string, error) {
	// Windows盘符最多只能使用26个，A B为软盘专用，实际只有24个
	devPool := []string{"C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	// 获取已用盘符
	drives, err := c.GetSpaceDrive()
	if err != nil {
		err := errors.New(msg.ERROR_GET_REMOTE_HOST_DRIVE,
			msg.GetMsg(msg.ERROR_GET_REMOTE_HOST_DRIVE, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return "", err
	}
	// 将切片转为map，方便查找
	drivesMap := make(map[string]string)
	for _, iter := range drives {
		drivesMap[iter] = iter
	}
	for _, dev := range devPool {
		_, ok := drivesMap[dev]
		if !ok {
			drive := dev + ":"
			return drive, nil
		}
	}
	err = errors.New(msg.ERROR_REMOTE_HOST_NO_DRIVE,
		msg.GetMsg(msg.ERROR_REMOTE_HOST_NO_DRIVE, c.loginInfo.Ip))
	return "", err
}

// 获取已使用盘符
func (c *winRmClient) GetSpaceDrive() ([]string, error) {
	var drives []string
	cmdFs := "fsutil fsinfo drives"
	result, err := c.RunCommond(cmdFs)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for _, line := range result {
		if strings.Contains(line, "Drives:") {
			temp := strings.Split(line, ": ")[1]
			temp = util.DeleteSpace(temp)
			drives = strings.Split(temp, ":\\")
			drives = drives[:len(drives)-1]
		}
	}
	cmdNetuse := "net use"
	result, err = c.RunCommond(cmdNetuse)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for _, line := range result {
		if strings.Contains(line, "Unavailable") {
			temp := strings.Split(line, ":")[0]
			temp = strings.Split(temp, " ")[1]
			notExist := true
			for _, iter := range drives {
				if iter == temp {
					notExist = false
				}
			}
			if notExist {
				drives = append(drives, temp)
			}
		}
	}
	return drives, nil
}

func (c *winRmClient) GetDriveByPath(path string) (string, error) {
	cmd := "net use"
	result, err := c.RunCommond(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REMOTE_HOST_DRIVE_BY_PATH,
			msg.GetMsg(msg.ERROR_GET_REMOTE_HOST_DRIVE_BY_PATH, c.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return "", err
	}
	for _, line := range result {
		if strings.Contains(line, path) {
			temp := strings.Split(line, ":")[0]
			temp = strings.Split(temp, " ")[1]
			drive := temp + ":"
			return drive, nil
		}
	}
	err = errors.New(msg.ERROR_GET_REMOTE_HOST_DRIVE_BY_PATH, msg.GetMsg(msg.ERROR_GET_REMOTE_HOST_DRIVE_BY_PATH,
		c.loginInfo.Ip, msg.GetMsg(msg.ERROR_NOT_FOUND_WITH_DRIVER)))
	return "", err
}

func (c *winRmClient) RunCommond(cmd string) ([]string, error) {
	var lines []string
	var stdout, stderr bytes.Buffer
	_, err := c.client.Run(cmd, &stdout, &stderr)
	if err != nil {
		errInfo := err.Error() + ":" + stderr.String()
		err = errors.New(msg.ERROR_RUN_REMOTE_COMMAND,
			msg.GetMsg(msg.ERROR_RUN_REMOTE_COMMAND, c.loginInfo.Ip, cmd, errInfo))
		logrus.Error(err)
		return nil, err
	}
	reader := bufio.NewReader(&stdout)
	for {
		lineByte, _, err1 := reader.ReadLine()
		if err1 != nil {
			if err1 == io.EOF {
				break
			}
			err = errors.New(msg.ERROR_RUN_REMOTE_COMMAND, msg.GetMsg(msg.ERROR_RUN_REMOTE_COMMAND,
				c.loginInfo.Ip, cmd, msg.GetMsg(msg.ERROR_READ_STDOUT_FAILED)))
			break
		}
		// 去掉多余的换行符
		line := strings.Replace(string(lineByte), "\n", "", -1)
		// 去掉多余的空格
		line = util.DeleteExtraSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	logrus.Info(lines)
	return lines, nil
}
