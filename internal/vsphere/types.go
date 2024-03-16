package vsphere

import (
	"context"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/vmware/govmomi"
	"reflect"
)

type VsphereClient struct {
	Vmomi     *govmomi.Client
	loginInfo *contants.LoginInfo
	Ctx       context.Context
}

const (
	INVENTORY_UNDEFINE                 = ""                         // 未定义
	INVENTORY_CLUSTERCOMPUTERESOURCE   = "ClusterComputeResource"   // 集群计算资源
	INVENTORY_COMPUTERESOURCE          = "ComputeResource"          // 计算资源
	INVENTORY_DATACENTER               = "Datacenter"               // 数据中心
	INVENTORY_DATASTORE                = "Datastore"                // 数据存储
	INVENTORY_DISTRIBUTEDVIRTUALSWITCH = "DistributedVirtualSwitch" // 分布式交换机
	INVENTORY_FOLDER                   = "Folder"                   // 文件夹
	INVENTORY_HOSTSYSTEM               = "HostSystem"               // 主机
	INVENTORY_NETWORK                  = "Network"                  // 网络
	INVENTORY_RESOURCEPOOL             = "ResourcePool"             // 资源池
	INVENTORY_VIRTUALAPP               = "VirtualApp"               // vApp
	INVENTORY_VIRTALMACHINE            = "VirtualMachine"           // 虚拟机
)

type DataSource struct {
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	Uuid      string `json:"uuid"`
	Path      string `json:"path"`
	Type      string `json:"type"`
	Expanded  bool   `json:"expanded"`
	CheckAble bool   `json:"checkable"`
}

func (d *DataSource) TransformMap() map[string]interface{} {
	t := reflect.TypeOf(d)
	v := reflect.ValueOf(d)
	data := make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		k := t.Elem().Field(i).Tag.Get("json")
		if len(k) == 0 {
			continue
		}
		data[k] = v.Elem().Field(i).Interface()
	}
	return data
}

type IpAddr struct {
	Ip       string
	Netmask  string
	Gateway  string
	Hostname string
}

type NasDataStoreParam struct {
	// HostPath   string
	HostName   string
	AccessMode string
	LocalPath  string
	RemoteHost string
	RemotePath string
	StoreType  string
}

type RegisterVmParam struct {
	VmxPath          string
	VmName           string
	FolderPath       string
	HostPath         string
	ResourcePoolPath string
}
