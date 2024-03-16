package main

import (
	"fmt"
)

var (
	Address   = "192.168.212.32"
	Username  = "jack.zhang"
	Password  = "zhangjiejun520..."
	storePool = ""
)

func GetUserInfo(client *ApiClient) {
	// 获取用户信息
	url := "/auth/user"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserInfo:%v\n", response)
}

func GetUserInfoById(client *ApiClient) {
	// 获取用户信息
	url := "/auth/user/20"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserInfoById:%v\n", response)
}

func GetUserList(client *ApiClient) {
	// 获取所有用户与管理员列表
	url := "/auth/users"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserList:%v\n", response)
	// 分页获取所有用户与管理员列表
	param := map[string]interface{}{"index": 0, "count": 20}
	response, err = client.SendRequest(GET, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserList:%v\n", response)
	// 获取所有用户信息
	param1 := map[string]interface{}{"index": 0, "count": 20, "roleName": "普通用户"}
	response, err = client.SendRequest(GET, param1, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserList:%v\n", response)
	// 获取所有管理员信息
	param2 := map[string]interface{}{"index": 0, "count": 20, "roleName": "系统管理员"}
	response, err = client.SendRequest(GET, param2, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserList:%v\n", response)
	// 获取指定用户信息
	param3 := map[string]interface{}{"index": 0, "count": 20, "filter": "jack"}
	response, err = client.SendRequest(GET, param3, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserList:%v\n", response)
	param4 := map[string]interface{}{"index": 0, "count": 20, "filter": "zhang.jiejun@outlook.com"}
	response, err = client.SendRequest(GET, param4, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserList:%v\n", response)
	// 获取指定用户信息
	param5 := map[string]interface{}{"index": 0, "count": 20, "roleName": "超级管理员", "filter": "jack"}
	response, err = client.SendRequest(GET, param5, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetUserList:%v\n", response)
}

func GetRoleInfo(client *ApiClient) {
	// 获取所有角色列表
	url := "/auth/role"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetRoleInfo:%v\n", response)
}

func GetHostSupportType(client *ApiClient) {
	// 获取已支持主机类型
	url := "/resource/host/support"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHostSupportType:%v\n", response)
}

func CreateHost(client *ApiClient) {
	// 添加主机资源
	url := "/resource/host"
	param := map[string]interface{}{"name": "212.54", "desc": "this is test host!", "type": 51, "ip": "192.168.212.54",
		"port": 22, "username": "root", "password": "passwd.com123"}
	response, err := client.SendRequest(POST, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("CreateHost:%v\n", response)
}

func EditHost(client *ApiClient) {
	url := "/resource/host/1"
	param := map[string]interface{}{"desc": "this is edit host!", "ip": "192.168.212.54",
		"port": 22, "username": "root", "password": "passwd.com123"}
	response, err := client.SendRequest(PUT, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("EditHost:%v\n", response)
}

func GetHostById(client *ApiClient) {
	url := "/resource/host/1"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHostById:%v\n", response)
}

func GetHosts(client *ApiClient) {
	url := "/resource/host"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHosts:%v\n", response)
	// 获取私有主机
	param := map[string]interface{}{"type": 111}
	response, err = client.SendRequest(GET, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHosts:%v\n", response)
	// 获取公共主机
	param1 := map[string]interface{}{"type": 222}
	response, err = client.SendRequest(GET, param1, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHosts:%v\n", response)
	// 过滤主机IP
	param2 := map[string]interface{}{"filter": "192."}
	response, err = client.SendRequest(GET, param2, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHosts:%v\n", response)
	// 过滤主机名称
	param3 := map[string]interface{}{"filter": "test"}
	response, err = client.SendRequest(GET, param3, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHosts:%v\n", response)
}

func PublishHost(client *ApiClient) {
	url := "/resource/host/1/publish"
	response, err := client.SendRequest(POST, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("PublishHost:%v\n", response)
}

func DeleteHost(client *ApiClient) {
	url := "/resource/host/1"
	param := map[string]interface{}{}
	response, err := client.SendRequest(DELETE, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("CreateHost:%v\n", response)
}

func GetPlatformSupportType(client *ApiClient) {
	// 获取已支持主机类型
	url := "/resource/platform/support"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetHostSupportType:%v\n", response)
}

func CreatePlatform(client *ApiClient) {
	// 添加主机资源
	url := "/resource/platform"
	param := map[string]interface{}{"name": "platform_vmware_esxi", "desc": "this is test platform!", "type": 11,
		"ip": "192.168.212.52", "port": 902, "username": "administrator@vsphere.local", "password": "passwd.com123"}
	response, err := client.SendRequest(POST, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("CreatePlatform:%v\n", response)
}

func EditPlatform(client *ApiClient) {
	url := "/resource/platform/1"
	param := map[string]interface{}{"desc": "this is edit platform!", "ip": "192.168.212.51",
		"port": 902, "username": "root", "password": "passwd.com123"}
	response, err := client.SendRequest(PUT, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("EditPlatform:%v\n", response)
}

func GetPlatformById(client *ApiClient) {
	url := "/resource/platform/1"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetPlatformById:%v\n", response)
}

func GetPlatforms(client *ApiClient) {
	url := "/resource/platform"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetPlatforms:%v\n", response)
	// 获取私有主机
	param := map[string]interface{}{"auth": 111}
	response, err = client.SendRequest(GET, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetPlatforms:%v\n", response)
	// 获取公共主机
	param1 := map[string]interface{}{"type": 11}
	response, err = client.SendRequest(GET, param1, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetPlatforms:%v\n", response)
	// 过滤主机IP
	param2 := map[string]interface{}{"filter": "192."}
	response, err = client.SendRequest(GET, param2, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetPlatforms:%v\n", response)
	// 过滤主机名称
	param3 := map[string]interface{}{"filter": "test"}
	response, err = client.SendRequest(GET, param3, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetPlatforms:%v\n", response)
}

func PublishPlatform(client *ApiClient) {
	url := "/resource/platform/1/publish"
	response, err := client.SendRequest(POST, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("PublishPlatform:%v\n", response)
}

func DeletePlatform(client *ApiClient) {
	url := "/resource/platform/1"
	param := map[string]interface{}{}
	response, err := client.SendRequest(DELETE, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("DeletePlatform:%v\n", response)
}

func CreateStorePool(client *ApiClient) {
	// 创建存储池
	url := "/resource/store_pool"
	param := map[string]interface{}{"name": "Linux_Pool", "desc": "this is test linux store pool!"}
	response, err := client.SendRequest(POST, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	storePool = response["uuid"].(string)
	fmt.Printf("CreateStorePool:%v\n", response)
}

func GetStorePools(client *ApiClient) {
	url := "/resource/store_pool"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetStorePools:%v\n", response)
}

func DeleteStorePool(client *ApiClient) {
	url := "/resource/store_pool/1"
	param := map[string]interface{}{}
	response, err := client.SendRequest(DELETE, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("store_pool:%v\n", response)
}

func GetLogRecord(client *ApiClient) {
	url := "/log/record"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetLogRecord:%v\n", response)
}

func GetLogRecordDetail(client *ApiClient) {
	url := "/log/record/15/detail"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetLogRecordDetail:%v\n", response)
}

func CreateReplica(client *ApiClient) {
	url := "/store_pool/replica"
	param := map[string]interface{}{"name": "test_Replica", "type": 1003, "size": 100, "storePoolId": storePool}
	response, err := client.SendRequest(POST, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("CreateReplica:%v\n", response)
}

func GetReplica(client *ApiClient) {
	url := "/store_pool/replica"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetReplica:%v\n", response)
}

func GetReplicaById(client *ApiClient) {
	url := "/store_pool/replica/1"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetReplicaById:%v\n", response)
}

func CreateSnapshotByReplica(client *ApiClient) {
	url := "/store_pool/replica/1/snapshot"
	param1 := map[string]interface{}{"name": "snapshot1", "desc": "this is snapshot1"}
	response, err := client.SendRequest(POST, param1, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("CreateSnapshotByReplica:%v\n", response)
	param2 := map[string]interface{}{"name": "snapshot2", "desc": "this is snapshot1"}
	response, err = client.SendRequest(POST, param2, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("CreateSnapshotByReplica:%v\n", response)
}

func GetSnapshotByReplica(client *ApiClient) {
	url := "/store_pool/replica/2/snapshot"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetSnapshotByReplica:%v\n", response)
}

func CreateImageByReplica(client *ApiClient) {
	url := "/store_pool/replica/1/image"
	param := map[string]interface{}{"name": "image1", "desc": "this is image1"}
	response, err := client.SendRequest(POST, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("CreateImageByReplica:%v\n", response)
}

func DeleteSnapshotByReplica(client *ApiClient) {
	url := "/store_pool/replica/1/snapshot"
	response, err := client.SendRequest(DELETE, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("DeleteSnapshotByReplica:%v\n", response)
}

func DeleteReplica(client *ApiClient) {
	url := "/store_pool/replica/1"
	response, err := client.SendRequest(DELETE, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("DeleteReplica:%v\n", response)
}

func MountReplica(client *ApiClient) {
	url := "/store_pool/replica/mount"
	//appConfig := map[string]interface{}{"locationPath": "VMware开发测试环境/vm/122",
	//	"computeResource": "VMware开发测试环境/host/cluster/Resources/Jack(开发环境)", "isRegisterVM": false}
	appConfig := map[string]interface{}{"mountPoint": "/test"}
	mountInfo := map[string]interface{}{"replicaId": 1, "targetType": 51, "targetId": 1, "mountType": 1,
		"isExecScript": true, "scriptName": "default.sh", "appConfig": appConfig}
	param := map[string]interface{}{"mountInfo": mountInfo}
	response, err := client.SendRequest(POST, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("MountReplica:%v\n", response)
}

func UnMountReplica(client *ApiClient) {
	url := "/store_pool/replica/1/unmount"
	response, err := client.SendRequest(POST, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("UnMountReplica:%v\n", response)
}

func GetSnapshotById(client *ApiClient) {
	url := "/store_pool/snapshot/3"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetSnapshotById:%v\n", response)
}

func DeleteSnapshotById(client *ApiClient) {
	url := "/store_pool/snapshot/3"
	response, err := client.SendRequest(DELETE, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("DeleteSnapshotById:%v\n", response)
}

func GetImage(client *ApiClient) {
	url := "/store_pool/image"
	param := map[string]interface{}{"auth": 111}
	response, err := client.SendRequest(GET, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetReplica:%v\n", response)
}

func GetDataSources(path string, isGetAll bool, client *ApiClient) {
	url := "/resource/platform/1/vmware/datasources"
	param := map[string]interface{}{"fullPath": path, "showType": 1, "isGetAll": isGetAll, "isGetVm": true}
	response, err := client.SendRequest(GET, param, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	if !isGetAll {
		if response["data"] != nil {
			datasources := response["data"].([]interface{})
			for _, datasource := range datasources {
				fmt.Printf("datasources:%v\n", datasources)
				GetDataSources(datasource.(map[string]interface{})["path"].(string), isGetAll, client)
			}
		}
	} else {
		fmt.Printf("GetDataSources:%v\n", response)
	}
}

func GetFiles(client *ApiClient) {
	url := "/file/tangula/ftp"
	response, err := client.SendRequest(GET, nil, url)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("GetFiles:%v\n", response)
}

func main() {
	client, err := NewApiClient(Address, Username, Password)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	// 用户鉴权
	//GetUserInfo(client)
	//GetUserInfoById(client)
	//GetUserList(client)
	//GetRoleInfo(client)

	// 主机管理
	//GetHostSupportType(client)
	//CreateHost(client)
	//GetHostById(client)
	//EditHost(client)
	//PublishHost(client)
	//GetHosts(client)
	//DeleteHost(client)

	// 平台管理
	//GetPlatformSupportType(client)
	//CreatePlatform(client)
	//GetPlatformById(client)
	//EditPlatform(client)
	//PublishPlatform(client)
	//GetPlatforms(client)
	//DeletePlatform(client)
	//GetDataSources("", true, client)

	// 存储池管理
	//CreateStorePool(client)
	//GetStorePools(client)
	//DeleteStorePool(client)

	// 副本管理
	//CreateReplica(client)
	//GetReplica(client)
	//GetReplicaById(client)
	//CreateSnapshotByReplica(client)
	//GetSnapshotByReplica(client)
	//CreateImageByReplica(client)
	//DeleteSnapshotByReplica(client)
	//DeleteReplica(client)
	//MountReplica(client)
	//UnMountReplica(client)

	// 快照管理
	//GetSnapshotById(client)
	//DeleteSnapshotById(client)

	// 镜像管理
	//GetImage(client)

	// 日志管理
	GetLogRecord(client)

	// 文件服务器
	GetFiles(client)

}
