package storage

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/ceph"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/image"
	"github.com/Zhang-jie-jun/tangula/internal/dao/replica"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/sirupsen/logrus"
)

func GetImage(queryParam *view.ImageQueryParam, user *auth.User) (totalNum int64,
	result []map[string]interface{}, err error) {
	var images []image.Image
	if queryParam.Auth == contants.PRIVATE {
		totalNum, images, err = image.ImageMgm.
			GetPrivateImageList(queryParam.Index, queryParam.Count, user.Account, queryParam.Type, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_IMAGE_INFO, msg.GetMsg(msg.ERROR_GET_IMAGE_INFO, err.Error()))
			logrus.Error(err)
			return
		}
	} else if queryParam.Auth == contants.PUBLIC {
		totalNum, images, err = image.ImageMgm.
			GetPublicImageList(queryParam.Index, queryParam.Count, queryParam.Type, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_IMAGE_INFO, msg.GetMsg(msg.ERROR_GET_IMAGE_INFO, err.Error()))
			logrus.Error(err)
			return
		}
	} else {
		totalNum, images, err = image.ImageMgm.
			GetImageList(queryParam.Index, queryParam.Count, queryParam.Type, queryParam.Filter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_IMAGE_INFO, msg.GetMsg(msg.ERROR_GET_IMAGE_INFO, err.Error()))
			logrus.Error(err)
			return
		}
		// 所有用户只能获取公共镜像与自身创建的私有镜像资源
		for i := 0; i < len(images); {
			if e := service.CheckResource(service.QUERY_RESOURCE, user, images[i].AuthType, images[i].CreateUser); e != nil {
				images = append(images[:i], images[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, iter := range images {
		result = append(result, iter.TransformMap())
	}
	return
}

func GetImageById(id uint, user *auth.User) (result map[string]interface{}, err error) {
	obj, err := image.ImageMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_IMAGE_INFO, msg.GetMsg(msg.ERROR_GET_IMAGE_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func DeleteImage(id uint, user *auth.User) (err error) {
	var obj image.Image
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DELETE_IMAGE_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.DELETE_IMAGE, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DELETE_IMAGE_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.DELETE_IMAGE, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = image.ImageMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_IMAGE_INFO, msg.GetMsg(msg.ERROR_GET_IMAGE_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	// 检查资源权限
	err = service.CheckResource(service.DELETE_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = ceph.Client.DeleteImage(obj.Pool.Uuid, obj.Uuid)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_IMAGE, msg.GetMsg(msg.ERROR_DELETE_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	err = image.ImageMgm.DeleteImage(obj.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_IMAGE, msg.GetMsg(msg.ERROR_DELETE_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	return
}

func PublishImage(id uint, user *auth.User) (result map[string]interface{}, err error) {
	var obj image.Image
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.PUBLISH_IMAGE_FAILED, obj.Name, err.Error())
			service.CreateLogRecord(msg.PUBLISH_IMAGE, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.PUBLISH_IMAGE_SUCCESS, obj.Name)
			service.CreateLogRecord(msg.PUBLISH_IMAGE, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	obj, err = image.ImageMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_IMAGE_INFO, msg.GetMsg(msg.ERROR_GET_IMAGE_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.PUBLISH_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	obj.AuthType = contants.PUBLIC
	obj, err = image.ImageMgm.UpdateImage(obj)
	if err != nil {
		err = errors.New(msg.ERROR_PUBLISH_IMAGE, msg.GetMsg(msg.ERROR_PUBLISH_IMAGE, err.Error()))
		logrus.Error(err)
		return
	}
	result = obj.TransformMap()
	return
}

func CreateReplicaByImage(id uint, param *view.ImageCreateParam, user *auth.User) (result map[string]interface{}, err error) {
	var obj image.Image
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.CREATE_REPLICA_BY_IMAGE_FAILED, obj.Name, param.Name, err.Error())
			service.CreateLogRecord(msg.CREATE_REPLICA_BY_IMAGE, obj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.CREATE_REPLICA_BY_IMAGE_BEGIN, obj.Name, param.Name)
			service.CreateLogRecord(msg.CREATE_REPLICA_BY_IMAGE, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	if replica.ReplicaMgm.CheckIsExistByNameAndCreateUser(param.Name, user.Account) {
		err = errors.New(msg.ERROR_REPLICA_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_REPLICA_NAME_IS_EXIST))
		return
	}
	obj, err = image.ImageMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_IMAGE_INFO, msg.GetMsg(msg.ERROR_GET_IMAGE_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = image.ImageMgm.UpdateImageStatus(obj.Id, contants.CLONEING)
	if err != nil {
		err = errors.New(msg.ERROR_IMAGE_TO_REPLICA, msg.GetMsg(msg.ERROR_IMAGE_TO_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	var rep replica.Replica
	rep.Name = param.Name
	rep.Desc = param.Desc
	rep.Size = obj.Size
	rep.Type = obj.Type
	rep.Status = contants.CLONEING
	rep.Export = "-"
	rep.Pool = obj.Pool
	rep.PoolId = obj.PoolId
	rep.CreateUser = user.Account
	rep, err = replica.ReplicaMgm.CreateReplica(rep)
	if err != nil {
		err = errors.New(msg.ERROR_IMAGE_TO_REPLICA, msg.GetMsg(msg.ERROR_IMAGE_TO_REPLICA, err.Error()))
		logrus.Error(err)
		return
	}
	result = rep.TransformMap()
	go func() {
		defer func() {
			if err != nil {
				detail := msg.GetOperation(msg.CREATE_REPLICA_BY_IMAGE_FAILED, obj.Name, param.Name, err.Error())
				service.CreateLogRecord(msg.CREATE_REPLICA_BY_IMAGE, obj.Name, detail, user.Account, contants.LOG_FAILED)
				err3 := image.ImageMgm.UpdateImageStatus(obj.Id, contants.NOT_MOUNT)
				if err3 != nil {
					logrus.Error(err3)
				}
				err4 := replica.ReplicaMgm.DeleteReplica(rep.Id)
				if err3 != nil {
					logrus.Error(err4)
				}
			} else {
				detail := msg.GetOperation(msg.CREATE_REPLICA_BY_IMAGE_SUCCESS, obj.Name, param.Name)
				service.CreateLogRecord(msg.CREATE_REPLICA_BY_IMAGE, obj.Name, detail, user.Account, contants.LOG_SUCCESS)
			}
		}()
		// 生成一个随机UUID作为Ceph存储上的镜像名称
		uuid := util.GenerateGuid()
		logrus.Info("============开始由镜像生成副本==============", obj.Uuid, "=>", uuid)
		err1 := ceph.Client.CopyImage(obj.Pool.Uuid, obj.Uuid, uuid)
		if err1 != nil {
			err = errors.New(msg.ERROR_IMAGE_TO_REPLICA, msg.GetMsg(msg.ERROR_IMAGE_TO_REPLICA, err1.Error()))
			logrus.Error(err1)
			return
		}
		logrus.Info("============镜像生成副本成功==============", obj.Uuid, "=>", uuid)
		//更新镜像状态
		err2 := image.ImageMgm.UpdateImageStatus(obj.Id, contants.NOT_MOUNT)
		if err2 != nil {
			err = errors.New(msg.ERROR_IMAGE_TO_REPLICA, msg.GetMsg(msg.ERROR_IMAGE_TO_REPLICA, err2.Error()))
			logrus.Error(err2)
			return
		}
		//更新副本状态
		_ = replica.ReplicaMgm.UpdateReplicaStatusAndUuid(rep.Id, uuid, contants.NOT_MOUNT)
	}()
	return
}
func EditImage(id uint, editParam *view.ImageEditParam, user *auth.User) (map[string]interface{}, error) {
	imgObj, err := image.ImageMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_REPLICA_INFO, msg.GetMsg(msg.ERROR_GET_REPLICA_INFO, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, imgObj.AuthType, imgObj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	imgName := util.Strips(editParam.Name, " ")
	if imgName != imgObj.Name { //同名就不用改了
		if image.ImageMgm.CheckIsExistByName(imgName) {
			err = errors.New(msg.ERROR_IMAGE_NAME_IS_EXIST, msg.GetMsg(msg.ERROR_IMAGE_NAME_IS_EXIST))
			return nil, err
		}
		imgObj.Name = imgName
	}
	imgObj.Desc = editParam.Desc
	imgObj, err = image.ImageMgm.UpdateImage(imgObj)
	if err != nil {
		err = errors.New(msg.ERROR_EDIT_IMAGE, msg.GetMsg(msg.ERROR_EDIT_IMAGE, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	result := imgObj.TransformMap()
	return result, nil
}
