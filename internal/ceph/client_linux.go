//go:build linux
// +build linux

package ceph

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/internal/system"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	"github.com/sirupsen/logrus"
	"strings"
)

type client struct{}

var Client = client{}

// 连接ceph集群
func (c *client) loginCeph() (err error) {
	Connect, err = rados.NewConn()
	if err != nil {
		logrus.Errorf("Connect: new connect error:%v\n", err)
		err = errors.New(msg.ERROR_CEPH_NOT_CONNECT, msg.GetMsg(msg.ERROR_CEPH_NOT_CONNECT))
		return err
	}

	err = Connect.ReadDefaultConfigFile()
	if err != nil {
		logrus.Errorf("Connect: read config file error:%v\n", err)
		err = errors.New(msg.ERROR_CEPH_NOT_CONNECT, msg.GetMsg(msg.ERROR_CEPH_NOT_CONNECT))
		return err
	}

	err = Connect.Connect()
	if err != nil {
		logrus.Errorf("Connect: connect error:%v\n", err)
		err = errors.New(msg.ERROR_CEPH_NOT_CONNECT, msg.GetMsg(msg.ERROR_CEPH_NOT_CONNECT))
		return err
	}
	if Connect == nil {
		err = errors.New(msg.ERROR_CEPH_NOT_CONNECT, msg.GetMsg(msg.ERROR_CEPH_NOT_CONNECT))
		return err
	}

	logrus.Info("-------------连接ceph成功!-----------------")
	return err
}

// 执行ceph命令
func (c *client) MonCommand(command []byte) (buf []byte, err error) {
	if Connect == nil {
		err = Client.loginCeph()
		if err != nil {
			logrus.Errorf("Connect: connect error:%v\n", err)
			return
		}
	}
	buf, _, err = Connect.MonCommand(command)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 获取集群信息
func (c *client) GetClusterInfo() (info *ClusterInfo, err error) {
	if Connect == nil {
		err = Client.loginCeph()
		if err != nil {
			logrus.Errorf("Connect: connect error:%v\n", err)
			return
		}
	}
	stats, err := Connect.GetClusterStats()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_GET_CLUSTER_STATUS, msg.GetMsg(msg.ERROR_CEPH_GET_CLUSTER_STATUS, err.Error()))
		logrus.Error(err)
		return
	}
	info = &ClusterInfo{TotalSize: stats.Kb * 1024, UsedSize: stats.Kb_used * 1024,
		AvailSize: stats.Kb_avail * 1024, TotalObject: stats.Num_objects}
	return
}

// 创建存储池
func (c *client) CreatePool(poolName string) (err error) {
	if Connect == nil {
		err = Client.loginCeph()
		if err != nil {
			logrus.Errorf("Connect: connect error:%v\n", err)
			return err
		}
	}

	logrus.Info("------------开始创建存储池:-----------", poolName)
	err = Connect.MakePool(poolName)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_CREATE_STORE_POOL, msg.GetMsg(msg.ERROR_CEPH_CREATE_STORE_POOL, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 删除存储池
func (c *client) DeletePool(poolName string) (err error) {
	if Connect == nil {
		err = Client.loginCeph()
		if err != nil {
			logrus.Errorf("Connect: connect error:%v\n", err)
			return err
		}
	}
	err = Connect.DeletePool(poolName)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_DELETE_STORE_POOL, msg.GetMsg(msg.ERROR_CEPH_DELETE_STORE_POOL, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

func (c *client) openPool(poolName string) (ctx *rados.IOContext, err error) {
	if Connect == nil {
		err = Client.loginCeph()
		if err != nil {
			logrus.Errorf("Connect: connect error:%v\n", err)
			return
		}
	}
	ctx, err = Connect.OpenIOContext(poolName)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_POOL_NOT_OPEN, msg.GetMsg(msg.ERROR_CEPH_POOL_NOT_OPEN, err.Error()))
		logrus.Error(err)
		return ctx, err
	}
	if ctx == nil {
		err = errors.New(msg.ERROR_CEPH_NOT_CONNECT, msg.GetMsg(msg.ERROR_CEPH_NOT_CONNECT))
		return ctx, err
	}
	return ctx, err
}

// 创建image
func (c *client) CreateImage(poolName, imageName string, size uint64) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	options := rbd.NewRbdImageOptions()
	err = rbd.CreateImage(ctx, imageName, size, options)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_CREATE_IMAGE, msg.GetMsg(msg.ERROR_CEPH_CREATE_IMAGE, err.Error()))
		return err
	}
	return err
}

// 删除image
func (c *client) DeleteImage(poolName, imageName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()

	//强制解除映射
	cmd := fmt.Sprintf("rbd unmap -o force /dev/rbd/%s/%s", poolName, imageName)
	logrus.Info("解除rdb映射:", cmd)
	_, rbdErr := system.SysManage.RunCommand(cmd)
	if rbdErr != nil {
		logrus.Error("解除rdb映射异常:", rbdErr)
	}

	err = rbd.RemoveImage(ctx, imageName)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_DELETE_IMAGE, msg.GetMsg(msg.ERROR_CEPH_DELETE_IMAGE, err.Error()+"，请确定副本在未被进程占用的情况下再进行卸载"))
		logrus.Error(err)
		return err
	}

	return err
}

func (c *client) getImage(ctx *rados.IOContext, imageName string) (image *rbd.Image, err error) {
	image, err = rbd.OpenImage(ctx, imageName, rbd.NoSnapshot)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_IMAGE_NOT_OPEN, msg.GetMsg(msg.ERROR_CEPH_IMAGE_NOT_OPEN, err.Error()))
		logrus.Error(err)
		return image, err
	}
	return image, err
}

// 复制image
func (c *client) CopyImage(poolName, imageName, destName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	err = image.Copy(ctx, destName)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_COPY_IMAGE, msg.GetMsg(msg.ERROR_CEPH_COPY_IMAGE, err.Error()))
		logrus.Error(err)
		return err
	}
	logrus.Info("ceph拷贝image====>:", destName)
	return err
}

// 从快照克隆image[克隆前需要检查快照是否受保护]
func (c *client) CloneImageBySnapshot(poolName, imageName, snapName, destName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	features, err := image.GetFeatures()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_GET_IMAGE_FEATURES, msg.GetMsg(msg.ERROR_CEPH_GET_IMAGE_FEATURES, err.Error()))
		logrus.Error(err)
		return err
	}
	_, err = image.Clone(snapName, ctx, destName, features, 22)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_CLONE_IMAGE, msg.GetMsg(msg.ERROR_CEPH_CLONE_IMAGE, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 分离image快照依赖
func (c *client) FlattenImage(poolName, imageName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	err = image.Flatten()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_FLATTEN_IMAGE, msg.GetMsg(msg.ERROR_CEPH_FLATTEN_IMAGE, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 重置image大小
func (c *client) ReSizeImage(poolName, imageName string, size uint64) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	err = image.Resize(size)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_RESIZE_IMAGE, msg.GetMsg(msg.ERROR_CEPH_RESIZE_IMAGE, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 重命名镜像
func (c *client) ReNameImage(poolName, imageName, destImageName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	err = image.Rename(destImageName)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_RENAME_IMAGE, msg.GetMsg(msg.ERROR_CEPH_RENAME_IMAGE, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 刷新缓存数据到镜像
func (c *client) FlushImage(poolName, imageName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	err = image.Flush()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_FLUSH_IMAGE, msg.GetMsg(msg.ERROR_CEPH_FLUSH_IMAGE, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 创建image快照
func (c *client) CreateSnapshot(poolName, imageName, snapName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	_, err = image.CreateSnapshot(snapName)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_CREATE_IMAGE_SNAPSHOT, msg.GetMsg(msg.ERROR_CEPH_CREATE_IMAGE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 删除image快照
func (c *client) DeleteSnapshot(poolName, imageName, snapName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	snapshot := image.GetSnapshot(snapName)
	err = snapshot.Remove()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_DELETE_IMAGE_SNAPSHOT, msg.GetMsg(msg.ERROR_CEPH_DELETE_IMAGE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 回滚image快照
func (c *client) RollbackSnapshot(poolName, imageName, snapName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	snapshot := image.GetSnapshot(snapName)
	err = snapshot.Rollback()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_ROLLBACK_IMAGE_SNAPSHOT, msg.GetMsg(msg.ERROR_CEPH_ROLLBACK_IMAGE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 保护快照
func (c *client) ProtectSnapShot(poolName, imageName, snapName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	snapshot := image.GetSnapshot(snapName)
	err = snapshot.Protect()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_PROTECT_IMAGE_SNAPSHOT, msg.GetMsg(msg.ERROR_CEPH_PROTECT_IMAGE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

// 解除快照保护
func (c *client) UnProtectSnapShot(poolName, imageName, snapName string) (err error) {
	ctx, err := c.openPool(poolName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ctx.Destroy()
	image, err := c.getImage(ctx, imageName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		w := image.Close()
		if w != nil {
			logrus.Warn(w)
		}
	}()
	snapshot := image.GetSnapshot(snapName)
	err = snapshot.Unprotect()
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_UN_PROTECT_IMAGE_SNAPSHOT, msg.GetMsg(msg.ERROR_CEPH_UN_PROTECT_IMAGE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return err
	}
	return err
}

func (c *client) MapRBDImage(poolName, imageName string) (devPath string, err error) {
	cmd := fmt.Sprintf("rbd map %s --pool %s", imageName, poolName)
	result, err := system.SysManage.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_MAP_IMAGE, msg.GetMsg(msg.ERROR_CEPH_MAP_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	if len(result) != 1 {
		err = errors.New(msg.ERROR_CEPH_MAP_IMAGE,
			msg.GetMsg(msg.ERROR_CEPH_MAP_IMAGE, msg.GetMsg(msg.ERROR_EXPECTED_RESULT_IS_INCORRECT)))
		logrus.Error(err)
		return
	}
	devPath = result[0]
	return
}

func (c *client) UnMapRBDImage(poolName, imageName string) (err error) {
	logrus.Info("开始unmapRBD===================>", imageName)

	cmd := fmt.Sprintf("rbd unmap  -o force /dev/rbd/%s/%s", poolName, imageName)
	_, err = system.SysManage.RunCommand(cmd)

	if err != nil {
		err = errors.New(msg.ERROR_CEPH_UN_MAP_IMAGE, msg.GetMsg(msg.ERROR_CEPH_UN_MAP_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func (c *client) ShowMapRBDImage() (mapInfos []MapInfo, err error) {
	cmd := "rbd showmapped"
	result, err := system.SysManage.RunCommand(cmd)
	if err != nil {
		err = errors.New(msg.ERROR_CEPH_MAP_IMAGE, msg.GetMsg(msg.ERROR_CEPH_MAP_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	for _, line := range result {
		temp := strings.Fields(line)
		if len(temp) != 5 {
			err = errors.New(msg.ERROR_CEPH_GET_MAP_INFO,
				msg.GetMsg(msg.ERROR_CEPH_GET_MAP_INFO, msg.GetMsg(msg.ERROR_EXPECTED_RESULT_IS_INCORRECT)))
			logrus.Error(err)
			return
		}
		// 过滤标签行
		if temp[0] == "id" {
			continue
		}
		var info MapInfo
		info.PoolName = temp[1]
		info.ImageName = temp[2]
		info.DevPath = temp[4]
		mapInfos = append(mapInfos, info)
	}
	return
}
