package linux

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"time"
)

type sshClient struct {
	loginInfo *contants.LoginInfo
	client    *ssh.Client
}

func NewSshClient(login *contants.LoginInfo) (client *sshClient, err error) {
	c := &sshClient{loginInfo: login}
	err = c.Login()
	if err == nil {
		client = c
	}
	return
}

func (s *sshClient) Login() (err error) {
	config := ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            s.loginInfo.UserName,
		Auth:            []ssh.AuthMethod{ssh.Password(s.loginInfo.PassWord)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", s.loginInfo.Ip, s.loginInfo.Port)
	s.client, err = ssh.Dial("tcp", addr, &config)
	if err != nil {
		err = errors.New(msg.ERROR_LOGIN_REMOTE_HOST, msg.GetMsg(msg.ERROR_LOGIN_REMOTE_HOST, s.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	if s.client == nil {
		err = errors.New(msg.ERROR_LOGIN_REMOTE_HOST,
			msg.GetMsg(msg.ERROR_LOGIN_REMOTE_HOST, s.loginInfo.Ip, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return
	}
	return
}

func (s *sshClient) Logout() {
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			logrus.Error(err)
		}
	}
}

// todo: 禁止执行阻塞式or交互式命令(如:top，ping等)
func (s *sshClient) RemoteRunCmd(cmd string) (result []string, err error) {
	if s.client == nil {
		err = s.Login()
		if err != nil {
			return
		}
	}
	session, err := s.client.NewSession()
	if err != nil {
		return
	}
	defer func() {
		e := session.Close()
		if e != nil && e != io.EOF {
			logrus.Error(err)
		}
	}()
	stdout, err := session.CombinedOutput(cmd)
	if err == nil || err == io.EOF {
		// 按行获取执行结果
		r := bytes.NewReader(stdout)
		reader := bufio.NewReader(r)
		for {
			lineByte, _, err1 := reader.ReadLine()
			if err1 != nil {
				if err1 == io.EOF {
					break
				}
				err = errors.New(msg.ERROR_RUN_REMOTE_COMMAND, msg.GetMsg(msg.ERROR_RUN_REMOTE_COMMAND,
					s.loginInfo.Ip, cmd, msg.GetMsg(msg.ERROR_READ_STDOUT_FAILED)))
				logrus.Error(err)
				return
			}
			// 去掉多余的换行符
			line := strings.Replace(string(lineByte), "\n", "", -1)
			result = append(result, line)
		}
		return result, nil
	}
	err = errors.New(msg.ERROR_RUN_REMOTE_COMMAND,
		msg.GetMsg(msg.ERROR_RUN_REMOTE_COMMAND, s.loginInfo.Ip, cmd, err.Error()))
	logrus.Error(err)
	return
}

func (s *sshClient) GetHostName() (hostName string) {
	cmd := "hostname"
	stdout, err := s.RemoteRunCmd(cmd)
	if err != nil {
		logrus.Error(err)
		return "unknown"
	}
	if len(stdout) == 0 {
		return "unknown"
	}
	if stdout[0] == "" {
		return "unknown"
	}
	return stdout[0]
}

func (s *sshClient) GetHostArch() (hostArch string) {
	cmd := "arch"
	stdout, err := s.RemoteRunCmd(cmd)
	if err != nil {
		logrus.Error(err)
		return "x86_64"
	}
	if len(stdout) == 0 {
		return "x86_64"
	}
	if stdout[0] == "" {
		return "x86_64"
	}
	return stdout[0]
}

func (s *sshClient) GetHostOs() (hostOs string) {
	cmd := "cat /etc/redhat-release"
	stdout, err := s.RemoteRunCmd(cmd)
	if err != nil {
		logrus.Error(err)
		return "Ubuntu"
	}
	if len(stdout) == 0 {
		return "未知"
	}
	if stdout[0] == "" {
		return "未知"
	}
	return stdout[0]
}

func (s *sshClient) Mount(sharePath, mountPoint string) (err error) {
	// 检查目标挂载目录是否存在，不存在则创建
	lsCmd := fmt.Sprintf("ls %s", mountPoint)
	_, err1 := s.RemoteRunCmd(lsCmd)
	if err1 != nil {
		logrus.Warn(err1)
		mkdirCmd := fmt.Sprintf("mkdir -p %s", mountPoint)
		_, err = s.RemoteRunCmd(mkdirCmd)
		if err != nil {
			err = errors.New(msg.ERROR_REMOTE_HOST_MOUNT_FAILED,
				msg.GetMsg(msg.ERROR_REMOTE_HOST_MOUNT_FAILED, s.loginInfo.Ip, err.Error()))
			logrus.Error(err)
			return
		}
		_, err = s.RemoteRunCmd(lsCmd)
		if err != nil {
			err = errors.New(msg.ERROR_REMOTE_HOST_MOUNT_FAILED,
				msg.GetMsg(msg.ERROR_REMOTE_HOST_MOUNT_FAILED, s.loginInfo.Ip, err.Error()))
			logrus.Error(err)
			return
		}
	}
	// todo: 不检查目标挂载目录

	// 挂载NFS
	mountCmd := fmt.Sprintf("mount -t nfs %s %s", sharePath, mountPoint)
	_, err = s.RemoteRunCmd(mountCmd)
	if err != nil {
		err = errors.New(msg.ERROR_REMOTE_HOST_MOUNT_FAILED,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_MOUNT_FAILED, s.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (s *sshClient) UnMount(mountPoint string, delMp bool) (err error) {
	umountCmd := fmt.Sprintf("umount -l %s", mountPoint)
	_, err = s.RemoteRunCmd(umountCmd)
	if err != nil {
		err = errors.New(msg.ERROR_REMOTE_HOST_UN_MOUNT_FAILED,
			msg.GetMsg(msg.ERROR_REMOTE_HOST_UN_MOUNT_FAILED, s.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return
	}
	// 删除挂载目录
	if delMp {
		rmdirCmd := fmt.Sprintf("rm -rf %s", mountPoint)
		_, _ = s.RemoteRunCmd(rmdirCmd)
	}
	return
}
