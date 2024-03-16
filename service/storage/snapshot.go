package storage

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/image"
	"github.com/Zhang-jie-jun/tangula/internal/dao/pool"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/internal/dao/snapshot"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/sirupsen/logrus"
)

func GetSnapshotById(id uint, user *auth.User) (result map[string]interface{}, err error) {
	snap, err := snapshot.SnapshotMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SNAPSHOT_INFO, msg.GetMsg(msg.ERROR_GET_SNAPSHOT_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	rep, err := replica.ReplicaMgm.FindById(snap.ReplicaId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	result = snap.TransformMap()
	result["replica"] = rep.TransformMap()
	return
}

func RollbackSnapshot(id uint, user *auth.User) (err error) {
	var snap snapshot.Snapshot
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.ROLLBACK_SNAPSHOT_FAILED, snap.Name, err.Error())
			service.CreateLogRecord(msg.ROLLBACK_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.ROLLBACK_SNAPSHOT_SUCCESS, snap.Name)
			service.CreateLogRecord(msg.ROLLBACK_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	snap, err = snapshot.SnapshotMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SNAPSHOT_INFO, msg.GetMsg(msg.ERROR_GET_SNAPSHOT_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 获取存储池信息
	var storePool pool.StorePool
	storePool, err = pool.StorePoolMgm.FindById(snap.Replica.PoolId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	if snap.Replica.Status != contants.NOT_MOUNT {
		err = errors.New(msg.ERROR_ROLL_BACK_SNAPSHOT,
			msg.GetMsg(msg.ERROR_ROLL_BACK_SNAPSHOT, msg.GetMsg(msg.ERROR_REPLICAT_STATUS)))
		logrus.Error(err)
		return
	}
	err = ceph.Client.RollbackSnapshot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid)
	if err != nil {
		err = errors.New(msg.ERROR_ROLL_BACK_SNAPSHOT, msg.GetMsg(msg.ERROR_ROLL_BACK_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func DeleteSnapshot(id uint, user *auth.User) (err error) {
	var snap snapshot.Snapshot
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DELETE_SNAPSHOT_FAILED, snap.Name, err.Error())
			service.CreateLogRecord(msg.DELETE_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DELETE_SNAPSHOT_SUCCESS, snap.Name)
			service.CreateLogRecord(msg.DELETE_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	snap, err = snapshot.SnapshotMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SNAPSHOT_INFO, msg.GetMsg(msg.ERROR_GET_SNAPSHOT_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 获取存储池信息
	var storePool pool.StorePool
	storePool, err = pool.StorePoolMgm.FindById(snap.Replica.PoolId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = ceph.Client.DeleteSnapshot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_SNAPSHOT, msg.GetMsg(msg.ERROR_DELETE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return
	}
	err = snapshot.SnapshotMgm.DeleteSnapshot(snap.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_SNAPSHOT, msg.GetMsg(msg.ERROR_DELETE_SNAPSHOT, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func CreateImageBySnapshot(id uint, param *view.ImageCreateParam, user *auth.User) (err error) {
	var snap snapshot.Snapshot
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_IMAGE_BY_SNAPSHOT_FAILED, snap.Name, param.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_IMAGE_BY_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_FAILED)
		}
	}()
	if image.ImageMgm.CheckIsExistByName(param.Name) {
		err = errors.New(msg.ERROR_IMAGE_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_IMAGE_NAME_IS_EXIST))
		return
	}
	snap, err = snapshot.SnapshotMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SNAPSHOT_INFO, msg.GetMsg(msg.ERROR_GET_SNAPSHOT_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 获取存储池信息
	var storePool pool.StorePool
	storePool, err = pool.StorePoolMgm.FindById(snap.Replica.PoolId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}

	//先创建克隆中镜像记录
	var ima image.Image
	ima.Name = param.Name
	ima.Desc = param.Desc
	ima.Size = snap.Replica.Size
	ima.Type = snap.Replica.Type
	ima.Status = contants.CLONEING
	ima.Pool = snap.Replica.Pool
	ima.PoolId = snap.Replica.PoolId
	ima.CreateUser = user.Account
	ima.AuthType = contants.PRIVATE
	ima, err = image.ImageMgm.CreateImage(ima)
	if err != nil {
		err = errors.New(msg.ERROR_SNAPSHOT_TO_IMAGE, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	detail := msg.GetOperation(msg.CREATE_IMAGE_BY_SNAPSHOT_SUCCESS, snap.Name, param.Name)
	service.CreateLogRecord(msg.CREATE_IMAGE_BY_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_SUCCESS)

	go func() {
		defer func() {
			if err != nil {
				detail := msg.GetOperation(msg.CREATE_IMAGE_BY_SNAPSHOT_FAILED, snap.Name, param.Name, err.Error())
				service.CreateLogRecord(msg.CREATE_IMAGE_BY_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_FAILED)
				//删除镜像记录
				logrus.Error("删除镜像记录==>", ima.Name)
				delErr := image.ImageMgm.DeleteImage(ima.Id)
				logrus.Error(delErr)
			}
		}()
		// 生成一个随机UUID作为Ceph存储上的镜像名称
		uuid := util.GenerateGuid()
		// 1.保护快照
		logrus.Info("============保护快照==============", uuid)
		err = ceph.Client.ProtectSnapShot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_IMAGE, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_IMAGE, err.Error()))
			logrus.Error(err)
			return
		}
		// 2.从快照克隆镜像
		logrus.Info("============从快照克隆镜像==============", uuid)
		err = ceph.Client.CloneImageBySnapshot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid, uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_IMAGE, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_IMAGE, err.Error()))
			logrus.Error(err)
			return
		}
		// 3.分离快照依赖
		logrus.Info("============分离快照依赖==============", uuid)
		err = ceph.Client.FlattenImage(storePool.Uuid, uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_IMAGE, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_IMAGE, err.Error()))
			logrus.Error(err)
			return
		}
		// 4.解除快照保护
		logrus.Info("============解除快照保护==============", uuid)
		err = ceph.Client.UnProtectSnapShot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_IMAGE, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_IMAGE, err.Error()))
			logrus.Error(err)
			return
		}

		updErr := image.ImageMgm.UpdateImageStatusAndUuid(ima.Id, uuid, contants.NOT_MOUNT)
		if updErr != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_IMAGE, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_IMAGE, updErr.Error()))
			logrus.Error(err)
			return
		}
		detail := msg.GetOperation(msg.CREATE_IMAGE_BY_SNAPSHOT_SUCCESS, snap.Name, param.Name)
		service.CreateLogRecord(msg.CREATE_IMAGE_BY_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_SUCCESS)
	}()
	return
}

func CreateReplicaBySnapshot(id uint, param *view.ImageCreateParam, user *auth.User) (err error) {
	var snap snapshot.Snapshot
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_REPLICA_BY_SNAPSHOT_FAILED, snap.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_REPLICA_BY_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_FAILED)
		}
	}()
	if image.ImageMgm.CheckIsExistByName(param.Name) {
		err = errors.New(msg.ERROR_IMAGE_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_IMAGE_NAME_IS_EXIST))
		return
	}
	if replica.ReplicaMgm.CheckIsExistByNameAndCreateUser(param.Name, user.Account) {
		err = errors.New(msg.ERROR_REPLICA_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_REPLICA_NAME_IS_EXIST))
		return
	}
	snap, err = snapshot.SnapshotMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SNAPSHOT_INFO, msg.GetMsg(msg.ERROR_GET_SNAPSHOT_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 获取存储池信息
	var storePool pool.StorePool
	storePool, err = pool.StorePoolMgm.FindById(snap.Replica.PoolId)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	//先创建一个克隆中副本
	var rep replica.Replica
	rep.Name = param.Name
	rep.Desc = param.Desc
	rep.Size = snap.Replica.Size
	rep.Type = snap.Replica.Type
	rep.Status = contants.CLONEING
	rep.Export = "-"
	rep.Pool = snap.Replica.Pool
	rep.PoolId = snap.Replica.PoolId
	rep.CreateUser = user.Account
	rep, err = replica.ReplicaMgm.CreateReplica(rep)
	if err != nil {
		err = errors.New(msg.ERROR_SNAPSHOT_TO_REPLICA, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	go func() {
		defer func() {
			if err != nil {
				detail := msg.GetOperation(msg.CREATE_REPLICA_BY_SNAPSHOT_FAILED, snap.Name, param.Name, err.Error())
				service.CreateLogRecord(msg.CREATE_REPLICA_BY_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_FAILED)
				//删除副本记录
				delerr := replica.ReplicaMgm.DeleteReplica(rep.Id)
				if delerr != nil {
					logrus.Error(delerr)
				}
			}
		}()
		// 生成一个随机UUID作为Ceph存储上的镜像名称
		uuid := util.GenerateGuid()
		logrus.Info("============保护快照==============", uuid)
		// 1.保护快照
		err = ceph.Client.ProtectSnapShot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_REPLICA, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_REPLICA, err.Error()))
			logrus.Error(err)
			return
		}
		// 2.从快照克隆镜像
		logrus.Info("============从快照克隆镜像==============", uuid)
		err = ceph.Client.CloneImageBySnapshot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid, uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_REPLICA, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_REPLICA, err.Error()))
			logrus.Error(err)
			return
		}
		// 3.分离快照依赖
		logrus.Info("============分离快照依赖==============", uuid)
		err = ceph.Client.FlattenImage(storePool.Uuid, uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_REPLICA, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_REPLICA, err.Error()))
			logrus.Error(err)
			return
		}
		// 4.解除快照保护
		logrus.Info("============解除快照保护==============", uuid)
		err = ceph.Client.UnProtectSnapShot(storePool.Uuid, snap.Replica.Uuid, snap.Uuid)
		if err != nil {
			err = errors.New(msg.ERROR_SNAPSHOT_TO_REPLICA, msg.GetMsg(msg.ERROR_SNAPSHOT_TO_REPLICA, err.Error()))
			logrus.Error(err)
			return
		}
		//更新副本状态
		_ = replica.ReplicaMgm.UpdateReplicaStatusAndUuid(rep.Id, uuid, contants.NOT_MOUNT)
		detail := msg.GetOperation(msg.CREATE_REPLICA_BY_SNAPSHOT_SUCCESS, snap.Name, param.Name)
		service.CreateLogRecord(msg.CREATE_REPLICA_BY_SNAPSHOT, snap.Name, detail, user.Account, contants.LOG_SUCCESS)
	}()
	return
}
