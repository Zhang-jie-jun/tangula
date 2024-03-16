package resource

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/platform"
	"github.com/Zhang-jie-jun/tangula/internal/vsphere"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	"github.com/Zhang-jie-jun/tangula/service/app/cas"
	"github.com/Zhang-jie-jun/tangula/service/app/fusioncompute"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func GetVMwareDataSources(id uint, param *view.VMwareDataSources, user *auth.User) (int64, []map[string]interface{}, error) {
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error())))
		logrus.Error(err)
		return 0, nil, err
	}
	if obj.Type != contants.VMWARE {
		err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_PLATFORM_TYPE_NO_MATCH)))
		logrus.Error(err)
		return 0, nil, err
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, err.Error()))
		logrus.Error(err)
		return 0, nil, err
	}
	passWord, e := util.AesDecrypt(obj.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return 0, nil, err
	}
	// 登录VMware平台
	var login contants.LoginInfo
	login.Ip = obj.Ip
	login.Port = obj.Port
	login.UserName = obj.UserName
	login.PassWord = passWord
	client, err := vsphere.NewClient(&login)
	if err != nil {
		err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, err.Error()))
		logrus.Error(err)
		return 0, nil, err
	}
	var dataSources []map[string]interface{}
	if param.ShowType == 1 {
		//先查询第一级子目录，检查是否cluster的场景
		var dataSourcesTmp []map[string]interface{}
		var fullPathTmp = param.FullPath
		dataSourcesTmp, err = getDataSourceByComputeResource(client, param.FullPath, false, false)
		if err != nil {
			err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
				msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, err.Error()))
			logrus.Error(err)
			return 0, nil, err
		}
		for _, ds := range dataSourcesTmp {
			if ds["name"] == param.HostName {
				if ds["type"] == vsphere.INVENTORY_COMPUTERESOURCE { //子目录下有主机名称，且不为clust，则计算资源路径要拼上主机名称
					fullPathTmp += "/host/" + param.HostName
				}
			}
		}
		logrus.Info(fmt.Sprintf("根据路径查询计算资源:%s", fullPathTmp))
		//查询所有计算资源
		dataSources, err = getDataSourceByComputeResource(client, fullPathTmp, param.ISGetAll, param.IsGetVm)
		if err != nil {
			err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
				msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, err.Error()))
			logrus.Error(err)
			return 0, nil, err
		}
	} else if param.ShowType == 2 {
		dataSources, err = getDataSourceByVmTemplate(client, param.FullPath, param.ISGetAll, param.IsGetVm)
		if err != nil {
			err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
				msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, err.Error()))
			logrus.Error(err)
			return 0, nil, err
		}
	} else {
		err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_SHOW_TYPE_NON_SUPPORT)))
		logrus.Error(err)
		return 0, nil, err
	}
	length := len(dataSources)
	return int64(length), dataSources, nil
}

func getDataSourceByComputeResource(client *vsphere.VsphereClient, path string, isGetAll, isGetVm bool) ([]map[string]interface{}, error) {
	var dataSourceMaps []map[string]interface{}
	data, err := client.GetDataSourcesByComputeResource(path, isGetVm)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for _, iter := range data {
		if iter.Type == vsphere.INVENTORY_DATACENTER {
			iter.CheckAble = false
		}
		dataSourceMaps = append(dataSourceMaps, iter.TransformMap())
	}
	if isGetAll {
		for _, dataSourceMap := range dataSourceMaps {
			if dataSourceMap["expanded"].(bool) {
				tempArr, err := getDataSourceByComputeResource(client, dataSourceMap["path"].(string), isGetAll, isGetVm)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSourceMap["subObject"] = tempArr
				dataSourceMap["subObjectNum"] = len(tempArr)
			}
		}
	} else {
		for _, dataSourceMap := range dataSourceMaps {
			dataSourceMap["subObject"] = nil
			dataSourceMap["subObjectNum"] = 0
		}
	}
	return dataSourceMaps, nil
}

func getDataSourceByVmTemplate(client *vsphere.VsphereClient, path string, isGetAll, isGetVm bool) ([]map[string]interface{}, error) {
	var dataSourceMaps []map[string]interface{}
	data, err := client.GetDataSourcesByVmTemplate(path, isGetVm)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for _, iter := range data {
		dataSourceMaps = append(dataSourceMaps, iter.TransformMap())
	}
	if isGetAll {
		for _, dataSourceMap := range dataSourceMaps {
			if dataSourceMap["expanded"].(bool) {
				tempArr, err := getDataSourceByVmTemplate(client, dataSourceMap["path"].(string), isGetAll, isGetVm)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSourceMap["subObject"] = tempArr
				dataSourceMap["subObjectNum"] = len(tempArr)
			}
		}
	} else {
		for _, dataSourceMap := range dataSourceMaps {
			dataSourceMap["subObject"] = nil
			dataSourceMap["subObjectNum"] = 0
		}
	}
	return dataSourceMaps, nil
}

func GetVMwareHostsByPath(id uint, path string, user *auth.User) (map[string]interface{}, error) {

	hostsmap := map[string]interface{}{
		"data": nil,
	}
	var allHostsList []string
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_VMWARE_HOSTS_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_HOSTS_FAILED, msg.GetMsg(msg.ERROR_GET_PLATFORM_INFO, err.Error())))
		logrus.Error(err)
		return nil, err
	}
	if obj.Type != contants.VMWARE {
		err = errors.New(msg.ERROR_GET_VMWARE_HOSTS_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_HOSTS_FAILED, msg.GetMsg(msg.ERROR_PLATFORM_TYPE_NO_MATCH)))
		logrus.Error(err)
		return nil, err
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		err = errors.New(msg.ERROR_GET_VMWARE_HOSTS_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_HOSTS_FAILED, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	passWord, e := util.AesDecrypt(obj.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return nil, err
	}
	// 登录VMware平台
	var login contants.LoginInfo
	login.Ip = obj.Ip
	login.Port = obj.Port
	login.UserName = obj.UserName
	login.PassWord = passWord
	client, err := vsphere.NewClient(&login)
	if err != nil {
		err = errors.New(msg.ERROR_GET_VMWARE_HOSTS_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_HOSTS_FAILED, err.Error()))
		logrus.Error(err)
		return nil, err
	}

	dataSources, dataSourcesErr := getDataSourceByComputeResource(client, path, false, false)
	if dataSourcesErr != nil {
		err = errors.New(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_VMWARE_DATASOURCE_FAILED, dataSourcesErr.Error()))
		logrus.Error(dataSourcesErr)
		return nil, err
	}
	for _, iter := range dataSources {
		if iter["type"] == "ClusterComputeResource" || iter["type"] == "ComputeResource" {
			hostList, hostErr := client.GetHostObjectListByPath(cast.ToString(iter["path"]))
			if hostErr != nil {
				err = errors.New(msg.ERROR_GET_VMWARE_HOSTS_FAILED,
					msg.GetMsg(msg.ERROR_GET_VMWARE_HOSTS_FAILED, hostErr.Error()))
				logrus.Error(hostErr)
				continue
			} else {
				for _, h := range hostList {
					allHostsList = append(allHostsList, h)
				}
			}
		}
	}
	hostsmap["data"] = allHostsList
	return hostsmap, err

}

func GetCasHosts(id uint) (res map[string]interface{}, err error) {
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, err.Error())))
		logrus.Error(err)
		return nil, err
	}
	if obj.Type != contants.CAS {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_PLATFORM_TYPE_NO_MATCH)))
		logrus.Error(err)
		return nil, err
	}

	realPassWord, e := util.AesDecrypt(obj.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return nil, err
	}

	//查询主机列表
	var loginInfo contants.LoginInfo
	var casIp = obj.Ip
	fmt.Println(casIp)
	loginInfo.Ip = obj.Ip
	loginInfo.Port = obj.Port
	loginInfo.UserName = obj.UserName
	loginInfo.PassWord = realPassWord
	result, err := cas.GetCasHosts(loginInfo)
	if err != nil {
		return nil, err
	}
	/**
	reqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/host", obj.Ip)
	queryRes, queryErr := util.AuthRequest(obj.UserName, realPassWord, "GET", reqUrl, nil)
	if err != nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, obj.Ip, queryErr.Error()))
		logrus.Error(err)
		return nil, err
	}
	**/
	return result, err
}

func GetCasHostDetails(id uint, hostId uint) (res map[string]interface{}, err error) {
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, err.Error())))
		logrus.Error(err)
		return nil, err
	}
	if obj.Type != contants.CAS {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_PLATFORM_TYPE_NO_MATCH)))
		logrus.Error(err)
		return nil, err
	}

	realPassWord, e := util.AesDecrypt(obj.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return nil, err
	}

	//查询主机详情
	var loginInfo contants.LoginInfo
	var casIp = obj.Ip
	fmt.Println(casIp)
	loginInfo.Ip = obj.Ip
	loginInfo.Port = obj.Port
	loginInfo.UserName = obj.UserName
	loginInfo.PassWord = realPassWord
	result, err := cas.GetHostInfo(loginInfo, hostId)
	if err != nil {
		return nil, err
	}
	return result, err
}

func GetCasVmList(id uint, hostId uint, user *auth.User) (map[string]interface{}, error) {
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, err.Error())))
		logrus.Error(err)
		return nil, err
	}
	if obj.Type != contants.CAS {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_PLATFORM_TYPE_NO_MATCH)))
		logrus.Error(err)
		return nil, err
	}
	err = service.CheckResource(service.QUERY_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	realPassWord, e := util.AesDecrypt(obj.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return nil, err
	}

	//查询虚拟机列表
	var loginInfo contants.LoginInfo
	loginInfo.Ip = obj.Ip
	loginInfo.Port = obj.Port
	loginInfo.UserName = obj.UserName
	loginInfo.PassWord = realPassWord
	result, err := cas.GetCasVms(loginInfo, hostId)
	if err != nil {
		return nil, err
	}
	/**
	reqUrl := fmt.Sprintf("http://%s:8080/cas/casrs/host/id/%d/vm", obj.Ip, hostId)
	queryRes, queryErr := util.AuthRequest(obj.UserName, realPassWord, "GET", reqUrl, nil)
	if err != nil {
		err = errors.New(msg.ERRIR_QUERY_CAS, msg.GetMsg(msg.ERRIR_QUERY_CAS, obj.Ip, queryErr.Error()))
		logrus.Error(err)
		return nil, err
	}
	**/
	return result, err

}

func GetCasStorageList(id uint, hostId uint, user *auth.User) (res map[string]interface{}, err error) {
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, err.Error())))
		logrus.Error(err)
		return nil, err
	}
	if obj.Type != contants.CAS {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, msg.GetMsg(msg.ERROR_PLATFORM_TYPE_NO_MATCH)))
		logrus.Error(err)
		return nil, err
	}

	err = service.CheckResource(service.QUERY_RESOURCE, user, obj.AuthType, obj.CreateUser)
	if err != nil {
		err = errors.New(msg.ERROR_GET_CAS_DATASOURCE_FAILED,
			msg.GetMsg(msg.ERROR_GET_CAS_DATASOURCE_FAILED, err.Error()))
		logrus.Error(err)
		return nil, err
	}

	realPassWord, e := util.AesDecrypt(obj.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return nil, err
	}

	//查询存储
	var loginInfo contants.LoginInfo
	loginInfo.Ip = obj.Ip
	loginInfo.Port = obj.Port
	loginInfo.UserName = obj.UserName
	loginInfo.PassWord = realPassWord
	result, err := cas.GetCasStorages(loginInfo, hostId)
	if err != nil {
		return nil, err
	}
	return result, err
}

func GetFcHosts(id uint) (res map[string]interface{}, err error) {
	obj, err := platform.PlatformMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_INFO_FAILED,
			msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_INFO_FAILED, msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_INFO_FAILED, err.Error())))
		logrus.Error(err)
		return nil, err
	}
	if obj.Type != contants.FUSIONCOMPUTE {
		err = errors.New(msg.ERROR_GET_FUSIONCOMPUTE_INFO_FAILED,
			msg.GetMsg(msg.ERROR_GET_FUSIONCOMPUTE_INFO_FAILED, msg.GetMsg(msg.ERROR_PLATFORM_TYPE_NO_MATCH)))
		logrus.Error(err)
		return nil, err
	}

	realPassWord, e := util.AesDecrypt(obj.PassWord)
	if e != nil {
		err = errors.New(msg.ERROR_AES_DECRYPT, msg.GetMsg(msg.ERROR_AES_DECRYPT, e.Error()))
		logrus.Error(err)
		return nil, err
	}

	//查询主机列表
	var loginInfo contants.LoginInfo
	loginInfo.Ip = obj.Ip
	loginInfo.Port = obj.Port
	loginInfo.UserName = obj.UserName
	loginInfo.PassWord = realPassWord
	result, err := fusioncompute.GetHosts(loginInfo)
	if err != nil {
		return nil, err
	}
	return result, err
}
