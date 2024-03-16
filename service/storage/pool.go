package storage

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/image"
	"github.com/Zhang-jie-jun/tangula/internal/dao/pool"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/sirupsen/logrus"
)

func CreateStorePool(param *view.NameParam, user *auth.User) (result map[string]interface{}, err error) {
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_POOL_FAILED, param.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_POOL, param.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.CREATE_POOL_SUCCESS, param.Name)
			service.CreateLogRecord(msg.CREATE_POOL, param.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	if pool.StorePoolMgm.CheckIsExistByNameAndCreateUser(param.Name, user.Account) {
		err = errors.New(msg.ERROR_STORE_POOL_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_STORE_POOL_NAME_IS_EXIST))
		return
	}
	// 生成一个随机UUID作为Ceph存储上的存储池名称
	uuid := util.GenerateGuid()
	err = ceph.Client.CreatePool(uuid)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_STORE_POOL, msg.GetMsg(msg.ERROR_CREATE_STORE_POOL, err.Error()))
		logrus.Error(err)
		return
	}

	var storePool pool.StorePool
	storePool.Name = param.Name
	storePool.Desc = param.Desc
	storePool.Uuid = uuid
	storePool.CreateUser = user.Account
	storePool, err = pool.StorePoolMgm.CreateStorePool(storePool)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_STORE_POOL, msg.GetMsg(msg.ERROR_CREATE_STORE_POOL, err.Error()))
		logrus.Error(err)
		return
	}
	result = storePool.TransformMap()
	return
}

func DeleteStorePool(id uint, user *auth.User) (err error) {
	var storePool pool.StorePool
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DELETE_POOL_FAILED, storePool.Name, err.Error())
			service.CreateLogRecord(msg.DELETE_POOL, storePool.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DELETE_POOL_SUCCESS, storePool.Name)
			service.CreateLogRecord(msg.DELETE_POOL, storePool.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	storePool, err = pool.StorePoolMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 获取存储池内的副本，检查是否存在挂载状态的镜像
	replicaArr, err := replica.ReplicaMgm.GetReplicaByPoolId(storePool.Id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	for _, rep := range replicaArr {
		if rep.Status != contants.NOT_MOUNT {
			err = errors.New(msg.ERROR_DELETE_STORE_POOL,
				msg.GetMsg(msg.ERROR_DELETE_STORE_POOL, msg.GetMsg(msg.ERROR_EXIST_MOUNT_STATUS_REPLICA, rep.Name)))
			logrus.Error(err)
		}
	}
	// 先删除镜像与副本再删存储池
	err = replica.ReplicaMgm.DeleteReplicaByPoolId(storePool.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_REPLICA, msg.GetMsg(msg.ERROR_DELETE_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	err = image.ImageMgm.DeleteImageByPoolId(storePool.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_IMAGE, msg.GetMsg(msg.ERROR_DELETE_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}

	err = ceph.Client.DeletePool(storePool.Uuid)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = pool.StorePoolMgm.DeleteStorePool(storePool.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_STORE_POOL, msg.GetMsg(msg.ERROR_DELETE_STORE_POOL, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func GetStorePools(queryParam *view.QueryParam) (totalNum int64, result []map[string]interface{}, err error) {
	var storePools []pool.StorePool
	totalNum, storePools, err = pool.StorePoolMgm.GetStorePools(queryParam.Index, queryParam.Count, queryParam.Filter)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	for _, iter := range storePools {
		result = append(result, iter.TransformMap())
	}
	return
}

func GetStorePoolById(id uint) (result map[string]interface{}, err error) {
	storePool, err := pool.StorePoolMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_STORE_POOL_INFO, msg.GetMsg(msg.ERROR_GET_STORE_POOL_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	result = storePool.TransformMap()
	return
}
