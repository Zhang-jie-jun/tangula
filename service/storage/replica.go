package storage

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/host"
	"github.com/Zhang-jie-jun/tangula/internal/dao/image"
	"github.com/Zhang-jie-jun/tangula/internal/dao/platform"
	"github.com/Zhang-jie-jun/tangula/internal/dao/pool"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/internal/dao/snapshot"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/Zhang-jie-jun/tangula/service/websockets"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CreateReplica(param *view.ReplicaCreateParam, user *auth.User) (result map[string]interface{}, err error) {
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_REPLICA_FAILED, param.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_REPLICA, param.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.CREATE_REPLICA_SUCCESS, param.Name)
			service.CreateLogRecord(msg.CREATE_REPLICA, param.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	replicaName := util.Strips(param.Name, " ")
	if replica.ReplicaMgm.CheckIsExistByNameAndCreateUser(replicaName, user.Account) {
		err = errors.New(msg.ERROR_REPLICA_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_REPLICA_NAME_IS_EXIST))
		return
	}
	storagePool, err := pool.StorePoolMgm.FindByUuid(param.StorePoolId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// size单位转换[GB转为字节]
	size := param.Size * 1024 * 1024 * 1024
	// 获取存储信息
	clusterInfo, err := ceph.Client.GetClusterInfo()
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_REPLICA, msg.GetMsg(msg.ERROR_CREATE_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info("获取存储信息===========>", clusterInfo)
	// 检查存储可用空间
	if clusterInfo == nil || size > clusterInfo.AvailSize {
		err = errors.New(msg.ERROR_CREATE_REPLICA, msg.GetMsg(msg.ERROR_CREATE_REPLICA,
			msg.GetMsg(msg.ERROR_STORE_INSUFFICIENT_SPACE, clusterInfo.AvailSize)))
		logrus.Error(err)
		return
	}
	logrus.Info("检查存储可用空间===========>")
	// 生成一个随机UUID作为Ceph存储上的镜像名称
	uuid := util.GenerateGuid()
	err = ceph.Client.CreateImage(storagePool.Uuid, uuid, size)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_REPLICA, msg.GetMsg(msg.ERROR_CREATE_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info("ceph创建镜像===========>", uuid)
	var rep replica.Replica
	rep.Name = replicaName
	rep.Desc = param.Desc
	rep.Uuid = uuid
	rep.Size = size
	rep.Type = param.Type
	rep.Status = contants.NOT_MOUNT
	rep.Export = "-"
	rep.Pool = storagePool
	rep.PoolId = storagePool.Id
	rep.CreateUser = user.Account
	rep, err = replica.ReplicaMgm.CreateReplica(rep)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_REPLICA, msg.GetMsg(msg.ERROR_CREATE_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	result = rep.TransformMap()
	return
}

func EditReplica(id uint, editParam *view.ReplicaEditParam, user *auth.User) (map[string]interface{}, error) {
	replicaObj, err := replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, contants.PRIVATE, replicaObj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	// size单位转换[GB转为字节]
	size := editParam.Size * 1024 * 1024 * 1024

	// 检查状态
	if size > 0 && size != replicaObj.Size {
		if replicaObj.Status != contants.NOT_MOUNT {
			err = errors.New(msg.ERROR_REPLICAT_STATUS, msg.GetMsg(msg.ERROR_REPLICAT_STATUS))
			logrus.Error(err)
			return nil, err
		}
	}

	if editParam.Type != 0 && editParam.Type != replicaObj.Type {
		if replicaObj.Status != contants.NOT_MOUNT {
			err = errors.New(msg.ERROR_REPLICAT_STATUS, msg.GetMsg(msg.ERROR_REPLICAT_STATUS))
			logrus.Error(err)
			return nil, err
		}
	}
	if size > 0 && size != replicaObj.Size {
		if size <= replicaObj.Size {
			err = errors.New(msg.ERROR_EDIT_REPLICA,
				msg.GetMsg(msg.ERROR_EDIT_REPLICA, msg.GetMsg(msg.ERROR_REPLICA_CAPACITY_EXPANSION_ONLY)))
			logrus.Error(err)
			return nil, err
		}
		// 获取存储信息
		clusterInfo, err := ceph.Client.GetClusterInfo()
		if err != nil {
			err = errors.New(msg.ERROR_EDIT_REPLICA, msg.GetMsg(msg.ERROR_EDIT_REPLICA, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		// 检查存储可用空间
		if clusterInfo == nil || size > clusterInfo.AvailSize {
			err = errors.New(msg.ERROR_EDIT_REPLICA, msg.GetMsg(msg.ERROR_EDIT_REPLICA,
				msg.GetMsg(msg.ERROR_STORE_INSUFFICIENT_SPACE, clusterInfo.AvailSize)))
			logrus.Error(err)
			return nil, err
		}
		err = ceph.Client.ReSizeImage(replicaObj.Pool.Uuid, replicaObj.Uuid, size)
		if err != nil {
			err = errors.New(msg.ERROR_EDIT_REPLICA, msg.GetMsg(msg.ERROR_EDIT_REPLICA, err.Error()))
			logrus.Error(err)
			return nil, err
		}
	} else {
		size = replicaObj.Size
	}

	replicaName := util.Strips(editParam.Name, " ")
	if replicaName != replicaObj.Name { //同名就不用改了
		if replica.ReplicaMgm.CheckIsExistByNameAndCreateUser(replicaName, user.Account) {
			err = errors.New(msg.ERROR_REPLICA_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_REPLICA_NAME_IS_EXIST))
			return nil, err
		}
		replicaObj.Name = replicaName
	}

	replicaObj.Desc = editParam.Desc
	replicaObj.Size = size
	if editParam.Type != 0 {
		replicaObj.Type = editParam.Type
	}
	replicaObj, err = replica.ReplicaMgm.UpdateReplica(replicaObj)
	if err != nil {
		err = errors.New(msg.ERROR_EDIT_REPLICA, msg.GetMsg(msg.ERROR_EDIT_REPLICA, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	result := replicaObj.TransformMap()
	return result, nil
}

func GetReplica(queryParam *view.ReplicaQueryParam, user *auth.User) (totalNum int64,
	result []map[string]interface{}, err error) {
	totalNum, replicas, err := replica.ReplicaMgm.
		GetReplica(queryParam.Index, queryParam.Count, user.Account, queryParam.Status, queryParam.Type, queryParam.Filter, queryParam.Uuid)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	for _, iter := range replicas {
		var res map[string]interface{}
		res = iter.TransformMap()
		if iter.Status == contants.MOUNTED {
			// 获取挂载信息
			mountInfo, err := getMountInfo(iter.Id)
			if err != nil {
				logrus.Error(err)
			} else {
				res["mount_info"] = mountInfo
			}
		}
		result = append(result, res)
	}

	return
}

func GetReplicaById(id uint, user *auth.User) (result map[string]interface{}, err error) {
	obj, err := replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	result["isUpload"] = false
	if obj.Status == contants.MOUNTED {
		// 获取挂载信息
		mountInfo, err := getMountInfo(obj.Id)
		if err != nil {
			logrus.Error(err)
		} else {
			result["mount_info"] = mountInfo
		}

		//是否已上传配置文件
		if obj.Type == contants.CASVM {
			configFleName := contants.AppCfg.System.MountPath + "/" + obj.Pool.Uuid + "/" + obj.Uuid + "/cas.json"
			if util.IsFileExists(configFleName) {
				result["isUpload"] = true
			}

			//解析磁盘文件
			diskPath := contants.AppCfg.System.MountPath + "/" + obj.Pool.Uuid + "/" + obj.Uuid
			cmdRes, cmdErr := util.GetMountDiskInfo(diskPath)
			if cmdErr != nil {
				logrus.Error(fmt.Sprintf("解析磁盘错误:%s", cmdErr))
			} else {
				diskInfoList := []view.DiskInfo{}
				for _, line := range cmdRes {
					logrus.Info(fmt.Sprintf("获取磁盘信息:%s", line))
					splitRet := strings.Split(line, " ")
					if len(splitRet) != 2 {
						logrus.Error("split line ret error.")
						continue
					}
					diskInfo := view.DiskInfo{}
					diskInfo.Name = splitRet[1]
					diskInfo.SizeByte, _ = strconv.ParseInt(splitRet[0], 10, 64)
					diskInfoList = append(diskInfoList, diskInfo)
				}
				result["diskInfo"] = diskInfoList
			}

		}

	}

	return
}

func getMountInfo(replicaId uint) (result map[string]interface{}, err error) {
	mountInfo, err := replica.MountInfoMgm.FindMountInfoByReplicaId(replicaId)
	if err != nil {
		logrus.Error(err)
		return
	}
	result = mountInfo.TransformMap()
	if mountInfo.TargetType == contants.LINUX || mountInfo.TargetType == contants.WINDOWS {
		targetInfo, err := host.HostMgm.FindById(mountInfo.TargetId)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result["targetInfo"] = targetInfo.TransformMap()
	} else {
		targetInfo, err := platform.PlatformMgm.FindById(mountInfo.TargetId)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result["targetInfo"] = targetInfo.TransformMap()
	}
	return
}

func DeleteReplica(id uint, user *auth.User) (err error) {
	var obj replica.Replica
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DELETE_REPLICA_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.DELETE_REPLICA, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DELETE_REPLICA_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.DELETE_REPLICA, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 检查资源权限
	err = service.CheckResource(service.DELETE_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 检查状态
	if obj.Status != contants.NOT_MOUNT {
		err = errors.New(msg.ERROR_REPLICAT_STATUS, msg.GetMsg(msg.ERROR_REPLICAT_STATUS))
		logrus.Error(err)
		return
	}
	//删除快照
	snap, snapShotErr := snapshot.SnapshotMgm.FindByReplicaId(id)
	if snapShotErr != nil {
		logrus.Error("查询副本快照错误:", snapShotErr)
	} else {
		if len(snap) > 0 {
			for _, rep := range snap {
				logrus.Info("删除副本快照:", rep.Id)
				deleteSnapErr := DeleteSnapshot(rep.Id, user)
				if deleteSnapErr != nil {
					logrus.Error("删除副本快照错误:", deleteSnapErr)
					continue
				}
			}
		}
	}

	err = ceph.Client.DeleteImage(obj.Pool.Uuid, obj.Uuid)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_REPLICA, msg.GetMsg(msg.ERROR_DELETE_REPLICA, err.Error()))
		logrus.Error(err)
		//return
	}
	err = replica.ReplicaMgm.DeleteReplica(obj.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_REPLICA, msg.GetMsg(msg.ERROR_DELETE_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func CreateImageByReplica(id uint, param *view.ImageCreateParam, user *auth.User) (err error) {
	var obj replica.Replica
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_IMAGE_BY_REPLICA_FAILED, obj.Name, param.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_IMAGE_BY_REPLICA, obj.Name, detail, user.Account, contants.LOG_FAILED)
		}
	}()
	replicaName := util.Strips(param.Name, " ")
	if image.ImageMgm.CheckIsExistByName(replicaName) {
		err = errors.New(msg.ERROR_IMAGE_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_IMAGE_NAME_IS_EXIST))
		return
	}
	obj, err = replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 检查状态
	if obj.Status != contants.NOT_MOUNT {
		err = errors.New(msg.ERROR_REPLICAT_STATUS, msg.GetMsg(msg.ERROR_REPLICAT_STATUS))
		logrus.Error(err)
		return
	}
	err = replica.ReplicaMgm.UpdateReplicaStatus(obj.Id, contants.CLONEING)
	if err != nil {
		err = errors.New(msg.ERROR_REPLICA_TO_IMAGE, msg.GetMsg(msg.ERROR_REPLICA_TO_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	go func() {
		time.Sleep(2 * time.Second)
		// 生成一个随机UUID作为Ceph存储上的镜像名称
		uuid := util.GenerateGuid()
		logrus.Info("============开始生成镜像==============", obj.Uuid, "=>", uuid)
		err = ceph.Client.CopyImage(obj.Pool.Uuid, obj.Uuid, uuid)
		if err != nil {
			err = errors.New(msg.ERROR_REPLICA_TO_IMAGE, msg.GetMsg(msg.ERROR_REPLICA_TO_IMAGE, err.Error()))
			logrus.Error(err)
		}
		logrus.Info("============成功生成镜像==============", obj.Uuid, "=>", uuid)

		//创建image记录
		var ima image.Image
		ima.Name = replicaName
		ima.Desc = param.Desc
		ima.Uuid = uuid
		ima.Size = obj.Size
		ima.Type = obj.Type
		ima.Status = contants.NOT_MOUNT
		ima.Pool = obj.Pool
		ima.PoolId = obj.PoolId
		ima.CreateUser = user.Account
		ima.AuthType = contants.PRIVATE
		ima, err = image.ImageMgm.CreateImage(ima)
		if err != nil {
			err = errors.New(msg.ERROR_REPLICA_TO_IMAGE, msg.GetMsg(msg.ERROR_REPLICA_TO_IMAGE, err.Error()))
			logrus.Error(err)
			return
		}
		defer func() {
			if err != nil {
				detail := msg.GetOperation(msg.CREATE_IMAGE_BY_REPLICA_FAILED, obj.Name, param.Name, "未知错误")
				service.CreateLogRecord(msg.CREATE_IMAGE_BY_REPLICA, obj.Name, detail, user.Account, contants.LOG_FAILED)
				return
			}
		}()
		detail := msg.GetOperation(msg.CREATE_IMAGE_BY_REPLICA_SUCCESS, obj.Name, param.Name)
		service.CreateLogRecord(msg.CREATE_IMAGE_BY_REPLICA, obj.Name, detail, user.Account, contants.LOG_SUCCESS)

		//更新副本状态为空闲
		err = replica.ReplicaMgm.UpdateReplicaStatus(obj.Id, contants.NOT_MOUNT)
		if err != nil {
			err = errors.New(msg.ERROR_REPLICA_TO_IMAGE, msg.GetMsg(msg.ERROR_REPLICA_TO_IMAGE, err.Error()))
			logrus.Error(err)
			return
		}
		websockets.UpdateMsg(user.Account)
	}()
	return
}

func CreateSnapshotByReplica(id uint, param *view.ImageCreateParam, user *auth.User) (result map[string]interface{}, err error) {
	var obj replica.Replica
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_SNAPSHOT_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_SNAPSHOT, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.CREATE_SNAPSHOT_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.CREATE_SNAPSHOT, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	if snapshot.SnapshotMgm.CheckIsExistByNameAndReplicaId(param.Name, id) {
		err = errors.New(msg.ERROR_SNAPSHOT_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_SNAPSHOT_NAME_IS_EXIST))
		return
	}
	obj, err = replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 生成一个随机UUID作为Ceph存储上的快照名称
	uuid := util.GenerateGuid()
	err = ceph.Client.CreateSnapshot(obj.Pool.Uuid, obj.Uuid, uuid)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_SNAPSHOT, msg.GetMsg(msg.ERROR_CREATE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return
	}
	var snap snapshot.Snapshot
	snap.Name = param.Name
	snap.Desc = param.Desc
	snap.Uuid = uuid
	snap.Replica = obj
	snap.ReplicaId = obj.Id
	snap, err = snapshot.SnapshotMgm.CreateSnapshot(snap)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_SNAPSHOT, msg.GetMsg(msg.ERROR_CREATE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return
	}
	result = snap.TransformMap()
	result["replica"] = obj.TransformMap()
	return
}

func GetSnapshotByReplica(id uint, param *view.QueryParam, user *auth.User) (
	totalNum int64, result []map[string]interface{}, err error) {
	obj, err := replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	totalNum, snapshots, err := snapshot.SnapshotMgm.GetSnapshots(param.Index, param.Count, id, param.Filter)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SNAPSHOT_INFO, msg.GetMsg(msg.ERROR_GET_SNAPSHOT_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	logrus.Info(snapshots)
	for _, iter := range snapshots {
		snapMap := iter.TransformMap()
		snapMap["replica"] = obj.TransformMap()
		result = append(result, snapMap)
	}
	return
}

func DeleteSnapshotByReplica(id uint, user *auth.User) (err error) {
	var obj replica.Replica
	var snapshots []snapshot.Snapshot
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			var snapArr string
			for _, it := range snapshots {
				if snapArr == "" {
					snapArr = it.Name
				} else {
					snapArr = snapArr + "," + it.Name
				}
			}
			detail := msg.GetOperation(msg.DELETE_SNAPSHOT_BY_REPLICA_FAILED, obj.Name, snapArr, err.Error())
			service.CreateLogRecord(msg.DELETE_SNAPSHOT_BY_REPLICA, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			var snapArr string
			for _, it := range snapshots {
				if snapArr == "" {
					snapArr = it.Name
				} else {
					snapArr = snapArr + "," + it.Name
				}
			}
			detail := msg.GetOperation(msg.DELETE_SNAPSHOT_BY_REPLICA_SUCCESS, obj.Name, snapArr)
			service.CreateLogRecord(msg.DELETE_SNAPSHOT_BY_REPLICA, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 检查资源权限
	err = service.CheckResource(service.DELETE_RESOURCE, user, contants.PRIVATE, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	index := 0
	count := 15
	filter := ""
	for {
		totalNum, temp, err := snapshot.SnapshotMgm.GetSnapshots(index, count, id, filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_SNAPSHOT_INFO, msg.GetMsg(msg.ERROR_GET_SNAPSHOT_INFO, err.Error()))
			logrus.Error(err)
			break
		}
		index = index + count
		snapshots = append(snapshots, temp...)
		if totalNum < int64(index) {
			break
		}
	}
	for _, iter := range snapshots {
		err = ceph.Client.DeleteSnapshot(obj.Pool.Uuid, obj.Uuid, iter.Uuid)
		if err != nil {
			err = errors.New(msg.ERROR_DELETE_SNAPSHOT, msg.GetMsg(msg.ERROR_DELETE_SNAPSHOT, err.Error()))
			logrus.Error(err)
			continue
		}
		err = snapshot.SnapshotMgm.DeleteSnapshot(iter.Id)
		if err != nil {
			err = errors.New(msg.ERROR_DELETE_SNAPSHOT, msg.GetMsg(msg.ERROR_DELETE_SNAPSHOT, err.Error()))
			logrus.Error(err)
			continue
		}
	}
	return
}

func BatchDeleteReplica(param *view.BatchDeleteReplicaParam, user *auth.User) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	var res replica.Replica
	var err error
	for _, replicaId := range param.ReplicaIdList {
		//先检查副本状态
		res, err = replica.ReplicaMgm.FindById(replicaId)
		if err != nil {
			logrus.Error(err)
			tmpMap := make(map[string]interface{})
			tmpMap[strconv.Itoa(int(replicaId))] = err.Error()
			result = append(result, tmpMap)
			continue
		}
		if res.Status != contants.NOT_MOUNT {
			logrus.Error(err)
			tmpMap := make(map[string]interface{})
			tmpMap[strconv.Itoa(int(replicaId))] = "副本:" + strconv.Itoa(int(replicaId)) + "不是空闲状态"
			result = append(result, tmpMap)
			continue
		}
		err := DeleteReplica(replicaId, user)
		if err != nil {
			tmpMap := make(map[string]interface{})
			tmpMap[strconv.Itoa(int(replicaId))] = err.Error()
			result = append(result, tmpMap)
		}
	}
	return result, nil
}

func UploadJson(id uint, context *gin.Context, user *auth.User) error {
	var err error
	var file *multipart.FileHeader

	replicaObj, err := replica.ReplicaMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return err
	}

	//if replicaObj.Type != contants.CASVM {
	//	err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, "只有CAS类型副本才能上传配置文件"))
	//	return err
	//}
	if replicaObj.Status != contants.MOUNTED {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, "只有已挂载的副本才能上传配置文件"))
		return err
	}

	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.UPLOAD_SCRIPT_FILE_FAILED, cast.ToString(replicaObj.Id)+":"+file.Filename, err.Error())
			service.CreateLogRecord(msg.UPLOAD_SCRIPT_FILE, cast.ToString(replicaObj.Id)+":"+file.Filename, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.UPLOAD_SCRIPT_FILE_SUCCESS, cast.ToString(replicaObj.Id)+":"+file.Filename)
			service.CreateLogRecord(msg.UPLOAD_SCRIPT_FILE, cast.ToString(replicaObj.Id)+":"+file.Filename, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()

	// 获取文件信息
	file, err = context.FormFile("file")
	if err != nil {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, err.Error()))
		logrus.Error(err)
		return err
	}

	//if file.Filename != "cas.json" {
	//	err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, "配置文件名称必须是cas.json"))
	//	logrus.Error(err)
	//	return err
	//}

	// 脚本文件不能大于10M
	if file.Size > 10*1024*1024 {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE,
			msg.GetMsg(msg.ERROR_UPLOAD_FILE_SZIE_IS_TOO_BIG)))
		return err
	}
	//检查目录
	//cacheDirName:="D:\\dist"
	cacheDirName := contants.AppCfg.System.MountPath + "/" + replicaObj.Pool.Uuid + "/" + replicaObj.Uuid
	if _, cachErr := os.Stat(cacheDirName); cachErr != nil {
		// not exists
		logrus.Error(cachErr)
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, "副本目录不存在:"+cacheDirName))
		return err
	}

	// 保存文件
	dst := filepath.Join(cacheDirName, file.Filename)
	if err = context.SaveUploadedFile(file, dst); err != nil {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, err.Error()))
		logrus.Error(err)
		return err
	}
	return nil
}

func ViewJson(context *gin.Context) error {
	var err error
	var idparam view.UpdIdParam
	if paraErr := context.ShouldBind(&idparam); paraErr != nil {
		logrus.Error(paraErr)
		return err
	}

	replicaObj, err := replica.ReplicaMgm.FindById(idparam.Id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SCRIPT_FILE_CONTENT, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_CONTENT, err.Error()))
		logrus.Error(err)
		return err
	}

	//检查目录
	downloadFile := ""
	if replicaObj.Type == contants.VMWAREVM {
		downloadFile = contants.AppCfg.System.MountPath + "/" + replicaObj.Pool.Uuid + "/" + replicaObj.Uuid + "/meta.json"
	} else if replicaObj.Type == contants.CASVM {
		downloadFile = contants.AppCfg.System.MountPath + "/" + replicaObj.Pool.Uuid + "/" + replicaObj.Uuid + "/cas.json"
	}

	//context.Writer.Header().Add("Content-Type", "application/octet-stream")
	context.Writer.Header().Add("Content-Type", "application/json")
	//浏览器下载或预览
	context.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", "cas.json"))
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Cache-Control", "no-cache")
	context.File(downloadFile)
	return nil
}

func DoCasJob(jobParam *view.DoCasJobParam, context *gin.Context) (result map[string]interface{}, err error) {

	replicaObj, replicaErr := replica.ReplicaMgm.FindById(jobParam.Id)
	if replicaErr != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, replicaErr.Error()))
		logrus.Error(err)
		return
	}

	logrus.Info("-------------------获取挂载信息---------------------", replicaObj.Name)
	mountInfo, mountErr := replica.MountInfoMgm.FindMountInfoByReplicaId(replicaObj.Id)
	if mountErr != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, mountErr.Error()))
		logrus.Error(err)
		return
	}

	platformOrm, platformErr := platform.PlatformMgm.FindById(mountInfo.TargetId)
	if platformErr != nil {
		err = errors.New(msg.ERROR_GET_PLATFORM_INFO, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, platformErr.Error()))
		logrus.Error(err)
		return
	}
	passWord, e := util.AesDecrypt(platformOrm.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return
	}
	paramMap := map[string]string{
		"vcplatformip":     platformOrm.Ip,
		"vcplatformuser":   platformOrm.UserName,
		"vcplatformpasswd": passWord,
		"jsonname":         "cas.json",
		"SelfIP":           jobParam.ConsoleIp,
		"VERSION":          "AB7.0",
		"AT_tag":           jobParam.Tag,
	}
	buildId, buildErr := util.BuildJenkins(context, contants.AppCfg.JENKINS.JobNameCas, paramMap)
	if buildErr != nil {
		err = errors.New(msg.ERROR_GET_SCRIPT_FILE_CONTENT, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_CONTENT, buildErr.Error()))
		logrus.Error(buildErr)
	}

	res := map[string]interface{}{
		"buildId": buildId,
	}
	return res, err
}
