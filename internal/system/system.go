package system

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/sirupsen/logrus"
)

type System struct{}

var SysManage = System{}

func (sys *System) GetLocalIp() string {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		err = errors.New(msg.ERROR_GET_LOCAL_ADDR, msg.GetMsg(msg.ERROR_GET_LOCAL_ADDR, err.Error()))
		logrus.Error(err)
		return contants.AppCfg.Server.Host
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if ipnet.IP.String() == contants.AppCfg.Server.Host {
					ip = ipnet.IP.String()
					break
				}
				ip = ipnet.IP.String()
			}

		}
	}
	return ip
}

func (sys *System) FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (sys *System) GetAllDevInfo() (infos []*DevInfo, err error) {
	cmd := "blkid"
	result, err := sys.RunCommand(cmd)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, line := range result {
		temp := strings.Fields(line)
		if len(temp) != 3 {
			continue
		}
		var info DevInfo
		info.DevPath = strings.Split(temp[0], ":")[0]
		info.Uuid = strings.Split(temp[1], "\"")[1]
		info.Type = strings.Split(temp[2], "\"")[1]
		logrus.Info(info)
		infos = append(infos, &info)
	}
	return
}

func (sys *System) CheckDevFormat(devPath string) bool {
	cmd := fmt.Sprintf("blkid %s", devPath)
	result, err := sys.RunCommand(cmd)
	if err != nil {
		logrus.Error(err)
		return false
	}
	if len(result) == 0 {
		return false
	}
	return true
}

func (sys *System) CheckDevMounted(devPath string) bool {
	cmd := fmt.Sprintf("lsblk %s", devPath)
	result, err := sys.RunCommand(cmd)
	if err != nil {
		logrus.Error(err)
		return false
	}
	if len(result) != 2 {
		return false
	}
	tempArr := strings.Fields(result[1])
	if len(tempArr) == 7 && tempArr[6] != "" {
		return true
	}
	return false
}

func (sys *System) GetDevMountPoint(devPath string) (mountPoint string) {
	cmd := fmt.Sprintf("lsblk %s", devPath)
	result, err := sys.RunCommand(cmd)
	if err != nil {
		logrus.Error(err)
		return
	}
	if len(result) != 2 {
		return
	}
	tempArr := strings.Fields(result[1])
	if len(tempArr) == 7 && tempArr[6] != "" {
		mountPoint = tempArr[6]
		return
	}
	return
}

func (sys *System) FormatXFS(devPath string) (err error) {
	// 检查设备文件是否存在
	if !sys.FileExists(devPath) {
		err = errors.New(msg.ERROR_DEV_NOT_EXISTS, msg.GetMsg(msg.ERROR_DEV_NOT_EXISTS, devPath))
		logrus.Error(err)
		return
	}
	// 检查设备是否格式化
	//devInfos, err := sys.GetAllDevInfo()
	//for _, iter := range devInfos {
	//	if devPath == iter.DevPath {
	//		return
	//	}
	//}
	if sys.CheckDevFormat(devPath) {
		return
	}
	// 格式化设备
	cmd := fmt.Sprintf("mkfs.xfs -f %s", devPath)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) GrowfsXFS(devPath string) (err error) {
	cmd := fmt.Sprintf("xfs_growfs %s", devPath)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		logrus.Error(err)
		return
	}
	return

}
func (sys *System) Mount(devPath, mountPoint string) (err error) {
	// 检查设备文件是否存在
	logrus.Info("************************检查设备文件是否存在************************")
	if !sys.FileExists(devPath) {
		err = errors.New(msg.ERROR_DEV_NOT_EXISTS, msg.GetMsg(msg.ERROR_DEV_NOT_EXISTS, devPath))
		logrus.Error(err)
		return
	}
	// 检查设备是否已经挂载
	logrus.Info("************************检查设备是否已经挂载************************")
	if sys.CheckDevMounted(devPath) {
		// 已经挂载，且挂载到指定的挂载点上
		if sys.GetDevMountPoint(devPath) == mountPoint {
			return
		} else {
			// 已经挂载， 挂载到非指定的挂载点上
			err = errors.New(msg.ERROR_MOUNT_POINT_NOT_INCORRECT,
				msg.GetMsg(msg.ERROR_MOUNT_POINT_NOT_INCORRECT, devPath, mountPoint))
			logrus.Error(err)
			return
		}
	}
	// 检查挂载目录是否存在, 存在需要检查是否为空目录，不存在则创建
	logrus.Info("************************检查挂载目录是否存在,存在需要检查是否为空目录,不存在则创建************************")
	if !sys.FileExists(mountPoint) {
		err = os.MkdirAll(mountPoint, 0666)
		logrus.Info(fmt.Sprintf("创建挂载目录：%s", mountPoint))
		if err != nil {
			err = errors.New(msg.ERROR_CREATE_DIR, msg.GetMsg(msg.ERROR_CREATE_DIR, mountPoint, err.Error()))
			logrus.Error(err)
			return
		}
	} else {
		subList, _ := ioutil.ReadDir(mountPoint)
		if len(subList) != 0 {
			err = errors.New(msg.ERROR_MOUNT_POINT_OCCUPY, msg.GetMsg(msg.ERROR_MOUNT_POINT_OCCUPY, mountPoint))
			logrus.Error(err)
			return
		}
		logrus.Info(fmt.Sprintf("挂载目录已存在: %s", mountPoint))
	}

	time.Sleep(time.Second * 3)

	// 执行挂载命令
	cmd := fmt.Sprintf("mount -t xfs -o discard,nouuid %s %s", devPath, mountPoint)
	logrus.Info("************************执行挂载命令************************", cmd)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) UnMount(devPath, mountPoint string) (err error) {
	logrus.Info("开始unmount===================>", mountPoint)

	// 检查设备文件是否存在
	if !sys.FileExists(devPath) {
		err = errors.New(msg.ERROR_DEV_NOT_EXISTS, msg.GetMsg(msg.ERROR_DEV_NOT_EXISTS, devPath))
		logrus.Error(err)
		return
	}
	// 检查挂载目录是否存在
	if !sys.FileExists(mountPoint) {
		return
	}
	// 检查设备是否已经挂载
	if !sys.CheckDevMounted(devPath) {
		// 未挂载
		return
	}
	// 检查设备挂载点是否匹配
	if sys.GetDevMountPoint(devPath) != mountPoint {
		// 挂载到非指定的挂载点上
		err = errors.New(msg.ERROR_MOUNT_POINT_NOT_INCORRECT, msg.GetMsg(msg.ERROR_MOUNT_POINT_NOT_INCORRECT, devPath, mountPoint))
		logrus.Error(err)
		return
	}
	// 检查目录是否为挂载点
	cmd1 := fmt.Sprintf("mountpoint %s", mountPoint)
	result, err := sys.RunCommand(cmd1)
	if err != nil {
		err = errors.New(msg.ERROR_UNMOUNT_FILED,
			msg.GetMsg(msg.ERROR_UNMOUNT_FILED, msg.GetMsg(msg.ERROR_DIR_NOT_MOUNT_POINT, mountPoint)))
		logrus.Error(err)
		return
	}
	logrus.Info(result)
	// 执行卸载命令
	cmd2 := fmt.Sprintf("umount -l %s", devPath)
	_, err = sys.RunCommand(cmd2)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 睡眠3S防止文件占用未解除
	time.Sleep(time.Second * 3)
	// 执行删除挂载目录命令
	cmd3 := fmt.Sprintf("rm -rf %s", mountPoint)
	_, err = sys.RunCommand(cmd3)
	if err != nil {
		logrus.Error(err)
		return
	}
	time.Sleep(time.Second * 3)
	return
}

func (sys *System) ServerStatus(serverName string) (status ServerStatus, err error) {
	cmd := fmt.Sprintf("systemctl status %s", serverName)
	result, err := sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SERVER_STATUS, msg.GetMsg(msg.ERROR_GET_SERVER_STATUS, serverName, err.Error()))
		logrus.Error(err)
		return
	}
	for _, line := range result {
		if !strings.Contains(line, "Active:") {
			continue
		} else {
			if strings.Contains(strings.Split(line, ":")[1], "active") {
				status = ACTIVE
			} else if strings.Contains(strings.Split(line, ":")[1], "inactive") {
				status = INACTIVE
			} else if strings.Contains(strings.Split(line, ":")[1], "failed") {
				status = FAILED
			}
		}
	}
	return
}

func (sys *System) ServerStart(serverName string) (err error) {
	cmd := fmt.Sprintf("systemctl start %s", serverName)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_START_SERVER, msg.GetMsg(msg.ERROR_START_SERVER, serverName, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) ServerRestart(serverName string) (err error) {
	cmd := fmt.Sprintf("systemctl restart %s", serverName)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_RESTART_SERVER, msg.GetMsg(msg.ERROR_RESTART_SERVER, serverName, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) ServerReload(serverName string) (err error) {
	cmd := fmt.Sprintf("systemctl reload %s", serverName)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_RELOAD_SERVER, msg.GetMsg(msg.ERROR_RELOAD_SERVER, serverName, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) ServerStop(serverName string) (err error) {
	cmd := fmt.Sprintf("systemctl stop %s", serverName)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_STOP_SERVER, msg.GetMsg(msg.ERROR_STOP_SERVER, serverName, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) ServerEnable(serverName string) (err error) {
	cmd := fmt.Sprintf("systemctl enable %s", serverName)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_ENABLE_SERVER, msg.GetMsg(msg.ERROR_ENABLE_SERVER, serverName, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) ServerDisable(serverName string) (err error) {
	cmd := fmt.Sprintf("systemctl disable %s", serverName)
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_DISABLE_SERVER, msg.GetMsg(msg.ERROR_DISABLE_SERVER, serverName, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (sys *System) AddNFS(mountPoint, targetIp string) (err error) {
	// 检查NFS服务状态
	logrus.Info("*******************************检查NFS服务状态*******************************")
	status, err := sys.ServerStatus("nfs-server")
	if status != ACTIVE {
		// 重启服务
		err = sys.ServerRestart("nfs-server")
		if err != nil {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err.Error()))
			logrus.Error(err)
			return
		}
		if status, err = sys.ServerStatus("nfs-server"); status != ACTIVE {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE,
				msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_SERVER_NOT_START, "nfs-server")))
			logrus.Error(err)
		}
	}
	// 检查rpcbind服务状态
	logrus.Info("*******************************检查rpcbind服务状态*******************************")
	status, err = sys.ServerStatus("rpcbind")
	if status != ACTIVE {
		// 重启服务
		err = sys.ServerRestart("rpcbind")
		if err != nil {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err.Error()))
			logrus.Error(err)
			return
		}
		if status, err = sys.ServerStatus("rpcbind"); status != ACTIVE {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE,
				msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_SERVER_NOT_START, "rpcbind")))
			logrus.Error(err)
		}
	}
	// 检查挂载目录是否存在
	logrus.Info("*******************************检查挂载目录是否存在*******************************")
	if !sys.FileExists(mountPoint) {
		err = errors.New(msg.ERROR_DIR_NOT_EXISTS, msg.GetMsg(msg.ERROR_DIR_NOT_EXISTS, mountPoint))
		logrus.Error(err)
		return
	}
	if targetIp == "" {
		targetIp = "*"
	}
	/**
	// 检查exports文件是否存在，不存在则创建
	logrus.Info("*******************************检查exports文件是否存在，不存在则创建*******************************")
	exports := "/etc/exports"
	if !sys.FileExists(exports) {
		_, err = os.Create(exports)
		if err != nil {
			err = errors.New(msg.ERROR_CREATE_DIR, msg.GetMsg(msg.ERROR_CREATE_DIR, exports, err.Error()))
			logrus.Error(err)
			return
		}
	}

	fp, err := os.OpenFile(exports, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err.Error()))
		logrus.Error(err)
		return
	}
	defer func() {
		e := fp.Close()
		if e != nil {
			logrus.Error(e)
		}
	}()

	isWrite := true
	br := bufio.NewReader(fp)
	for {
		line, _, err2 := br.ReadLine()
		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err2.Error()))
			logrus.Error(err)
			return
		}
		// 已经存在不需要重复添加
		if strings.Contains(string(line), mountPoint) && strings.Contains(string(line), targetIp) {
			isWrite = false
		}
	}

	if isWrite {
		logrus.Info("*******************************构造NFS共享路径*******************************")
		// 构造NFS共享路径
		mountPath := fmt.Sprintf("%s %s(rw,sync,no_root_squash,fsid=%d)\n",
			mountPoint, targetIp, time.Now().UnixNano())
		// 追加写入
		w := bufio.NewWriter(fp)
		_, err = w.WriteString(mountPath)
		if err != nil {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err.Error()))
			logrus.Error(err)
			return
		}
		_ = w.Flush()
		_ = fp.Sync()
	}
	**/
	splitList := strings.Split(mountPoint, "/")
	replicaId := splitList[len(splitList)-1]
	exportsFile := fmt.Sprintf("/etc/exports.d/%s.exports", replicaId)
	logrus.Infof("副本对应的配置文件:%s", exportsFile)
	if !sys.FileExists(exportsFile) {
		_, err = os.Create(exportsFile)
		if err != nil {
			err = errors.New(msg.ERROR_CREATE_DIR, msg.GetMsg(msg.ERROR_CREATE_DIR, exportsFile, err.Error()))
			logrus.Error(err)
			return
		}
	}
	fpNew, err := os.OpenFile(exportsFile, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err.Error()))
		logrus.Error(err)
		return
	}
	defer func() {
		e := fpNew.Close()
		if e != nil {
			logrus.Error(e)
		}
	}()

	isWriteNew := true
	brNew := bufio.NewReader(fpNew)
	for {
		line, _, err2 := brNew.ReadLine()
		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err2.Error()))
			logrus.Error(err)
			return
		}
		// 已经存在不需要重复添加
		if strings.Contains(string(line), mountPoint) && strings.Contains(string(line), targetIp) {
			isWriteNew = false
		}
	}

	if isWriteNew {
		logrus.Info("*******************************为每个副本单独构造NFS共享路径*******************************")
		// 构造NFS共享路径
		mountPath := fmt.Sprintf("%s %s(rw,sync,no_root_squash,fsid=%d)\n",
			mountPoint, targetIp, time.Now().UnixNano())
		// 追加写入
		w := bufio.NewWriter(fpNew)
		_, err = w.WriteString(mountPath)
		if err != nil {
			err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err.Error()))
			logrus.Error(err)
			return
		}
		_ = w.Flush()
		_ = fpNew.Sync()
	}

	// 刷新配置文件
	logrus.Info("*******************************刷新配置文件*******************************")
	cmd := fmt.Sprintf("exportfs -r")
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_ADD_NFS_SHARE, msg.GetMsg(msg.ERROR_ADD_NFS_SHARE, err.Error()))
		logrus.Error(err)
		//刷新有问题暂不做处理
	}
	return
}

func (sys *System) RemoveNFS(mountPoint, targetIp string) error {
	logrus.Info("开始卸载NFS存储===================>", mountPoint)
	// 弃用/etc/exports配置文件模式
	/**
	exports := "/etc/exports"
	if !sys.FileExists(exports) {
		err := errors.New(msg.ERROR_REMOVE_NFS_SHARE,
			msg.GetMsg(msg.ERROR_REMOVE_NFS_SHARE, msg.GetMsg(msg.ERROR_FILE_NOT_EXISTS, exports)))
		logrus.Error(err)
		return err
	}
	if targetIp == "" {
		targetIp = "*"
	}
	fp1, err := os.OpenFile(exports, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		err = errors.New(msg.ERROR_REMOVE_NFS_SHARE, msg.GetMsg(msg.ERROR_REMOVE_NFS_SHARE, err.Error()))
		logrus.Error(err)
		return err
	}
	defer func() {
		e := fp1.Close()
		if e != nil {
			logrus.Error(e)
		}
	}()
	var items []string
	br := bufio.NewReader(fp1)
	logrus.Info("读取/etc/exports配置文件")
	for {
		line, _, err2 := br.ReadLine()
		if err2 == io.EOF {
			logrus.Error("读取/etc/exports配置文件失败:", line)
			break
		}
		if err2 != nil {
			err = errors.New(msg.ERROR_REMOVE_NFS_SHARE, msg.GetMsg(msg.ERROR_REMOVE_NFS_SHARE, err2.Error()))
			logrus.Error(err)
			return err
		}
		// 找到存在的目标项并过滤掉
		if strings.Contains(string(line), mountPoint) && strings.Contains(string(line), targetIp) {
			logrus.Info(fmt.Sprintf("过滤配置文件行=>%s", line))
			continue
		}
		items = append(items, string(line))
	}
	// e := fp1.Close()
	// if e != nil {
	// 	logrus.Error(e)
	// 	return err
	// }
	// 清空文件重新写入
	fp2, err := os.OpenFile(exports, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		err = errors.New(msg.ERROR_REMOVE_NFS_SHARE, msg.GetMsg(msg.ERROR_REMOVE_NFS_SHARE, err.Error()))
		logrus.Error(err)
		return err
	}
	defer func() {
		e := fp2.Close()
		if e != nil {
			logrus.Error(e)
		}
	}()
	w := bufio.NewWriter(fp2)
	logrus.Info("NFS配置文件：", items)
	for _, item := range items {
		// 加个换行符
		_, err = w.WriteString(item + "\n")
		if err != nil {
			err = errors.New(msg.ERROR_REMOVE_NFS_SHARE, msg.GetMsg(msg.ERROR_REMOVE_NFS_SHARE, err.Error()))
			logrus.Error(err)
			return err
		}
	}
	_ = w.Flush()
	_ = fp2.Sync()
	**/
	//删除每个副本对应的配置文件
	splitList := strings.Split(mountPoint, "/")
	replicaId := splitList[len(splitList)-1]
	exportsFile := fmt.Sprintf("/etc/exports.d/%s.exports", replicaId)
	err := os.Remove(exportsFile)
	logrus.Info(fmt.Sprintf("删除每个副本对应的配置文件:%s", exportsFile))
	if err != nil {
		err = errors.New(msg.ERROR_REMOVE_NFS_SHARE, msg.GetMsg(msg.ERROR_REMOVE_NFS_SHARE, err.Error()))
		logrus.Error(err)
	}

	// 刷新配置文件
	cmd := fmt.Sprintf("exportfs -r")
	_, err = sys.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_REMOVE_NFS_SHARE, msg.GetMsg(msg.ERROR_REMOVE_NFS_SHARE, err.Error()))
		logrus.Error(err)
		return err
	}
	time.Sleep(3 * time.Second)
	return nil
}

// 运行shell脚本
func (sys *System) RunShellScript(scriptPath, operation, sharaPath string) (string, error) {
	if scriptPath == "" {
		err := errors.New(msg.ERROR_RUN_SHELL_SCRIPT,
			msg.GetMsg(msg.ERROR_RUN_SHELL_SCRIPT, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return "", err
	}
	if string([]byte(scriptPath)[len(scriptPath)-3:]) != ".sh" {
		err := errors.New(msg.ERROR_RUN_SHELL_SCRIPT,
			msg.GetMsg(msg.ERROR_RUN_SHELL_SCRIPT, msg.GetMsg(msg.ERROR_NON_SUPPORT_NOT_SHELL)))
		logrus.Error(err)
		return "", err
	}
	err := os.Chmod(scriptPath, 0777)
	if err != nil {
		err := errors.New(msg.ERROR_RUN_SHELL_SCRIPT, msg.GetMsg(msg.ERROR_RUN_SHELL_SCRIPT, err.Error()))
		logrus.Error(err)
		return "", err
	}
	command := fmt.Sprintf("%s %s %s", scriptPath, operation, sharaPath)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		err := errors.New(msg.ERROR_RUN_SHELL_SCRIPT, msg.GetMsg(msg.ERROR_RUN_SHELL_SCRIPT, err.Error()))
		logrus.Error(err)
		return "", err
	}
	return string(output), nil
}

// todo: 此方法禁止执行交互式命令，使用管道连接的方式，无法获取响应，超时后进程会被kill掉
func (sys *System) RunCommand(command string) (lines []string, err error) {
	logrus.Info(command)
	// 参数检查
	if len(command) == 0 {
		err = errors.New(msg.ERROR_RUN_SYSTEM_COMMAND,
			msg.GetMsg(msg.ERROR_RUN_SYSTEM_COMMAND, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return
	}
	// 检查命令是否可执行
	path, err := exec.LookPath(strings.Split(command, " ")[0])
	if err != nil {
		err = errors.New(msg.ERROR_RUN_SYSTEM_COMMAND, msg.GetMsg(msg.ERROR_RUN_SYSTEM_COMMAND, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Infof(path)
	// 设置5分钟超时
	timeout := time.Duration(5) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	name := "/bin/bash"
	args := []string{"-c", command}
	cmd := exec.CommandContext(ctx, name, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// 创建命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.New(msg.ERROR_RUN_SYSTEM_COMMAND, msg.GetMsg(msg.ERROR_RUN_SYSTEM_COMMAND, err.Error()))
		logrus.Error(err)
		return
	}
	// 执行命令
	logrus.Info(fmt.Sprintf("开始执行命令:%s", cmd))
	err = cmd.Start()
	if err != nil {
		err = errors.New(msg.ERROR_RUN_SYSTEM_COMMAND, msg.GetMsg(msg.ERROR_RUN_SYSTEM_COMMAND, err.Error()))
		logrus.Error(err)
		return
	}
	// 按行获取执行结果
	reader := bufio.NewReader(stdout)
	for {
		lineByte, _, err1 := reader.ReadLine()
		if err1 != nil {
			if err1 == io.EOF {
				break
			}
			err = errors.New(msg.ERROR_RUN_SYSTEM_COMMAND, msg.GetMsg(msg.ERROR_RUN_SYSTEM_COMMAND, err1.Error()))
			logrus.Error(err)
			break
		}
		// 去掉多余的换行符
		line := strings.Replace(string(lineByte), "\n", "", -1)
		if line == "" || line == " " {
			continue
		} else {
			lines = append(lines, line)
		}
	}
	// 等待执行结束
	err1 := cmd.Wait()
	logrus.Info(fmt.Sprintf("【%s】命令执行结果:%s", cmd, lines))
	if err1 != nil {
		if stderr.Len() == 0 && lines == nil {
			// 文件描述符链接到os.DevNul上，忽略错误
			return
		}
		errInfo := err1.Error() + ":" + stderr.String()
		err = errors.New(msg.ERROR_RUN_SYSTEM_COMMAND, msg.GetMsg(msg.ERROR_RUN_SYSTEM_COMMAND, errInfo))
		logrus.Error(err)
		return
	}
	return
}
