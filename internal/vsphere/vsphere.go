package vsphere

import (
	"fmt"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// 获取根对象
func (v *VsphereClient) getRootPath() types.ManagedObjectReference {
	return v.Vmomi.ServiceContent.RootFolder
}

// 根据清单路径获取对象
func (v VsphereClient) findByInventoryPath(path string) (*types.ManagedObjectReference, error) {
	req := types.FindByInventoryPath{
		This:          *v.Vmomi.Client.ServiceContent.SearchIndex,
		InventoryPath: path,
	}

	res, err := methods.FindByInventoryPath(v.Ctx, v.Vmomi.Client, &req)
	if err != nil {
		err = errors.New(msg.ERROR_FIND_BY_INVENTORY_PATH,
			msg.GetMsg(msg.ERROR_FIND_BY_INVENTORY_PATH, path, err.Error()))
		return nil, err
	}

	if res.Returnval == nil {
		err = errors.New(msg.ERROR_FIND_BY_INVENTORY_PATH,
			msg.GetMsg(msg.ERROR_FIND_BY_INVENTORY_PATH, path, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL)))
		return nil, err
	}
	obj := object.NewReference(v.Vmomi.Client, *res.Returnval).Reference()
	return &obj, err
}

// 根据uuid获取虚拟机对象
func (v *VsphereClient) findByUUid(uuid string) (*types.ManagedObjectReference, error) {
	ins := true
	req := types.FindByUuid{
		This:         *v.Vmomi.Client.ServiceContent.SearchIndex,
		Uuid:         uuid,
		VmSearch:     true,
		InstanceUuid: &ins,
	}
	res, err := methods.FindByUuid(v.Ctx, v.Vmomi.Client, &req)
	if err != nil {
		err = errors.New(msg.ERROR_FIND_BY_INVENTORY_UUID,
			msg.GetMsg(msg.ERROR_FIND_BY_INVENTORY_UUID, uuid, err.Error()))
		return nil, err
	}
	if res.Returnval == nil {
		err = errors.New(msg.ERROR_FIND_BY_INVENTORY_UUID,
			msg.GetMsg(msg.ERROR_FIND_BY_INVENTORY_UUID, uuid, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL)))
		return nil, err
	}
	obj := object.NewReference(v.Vmomi.Client, *res.Returnval).Reference()
	return &obj, nil
}

// 获取指定对象属性
func (v *VsphereClient) getObjectProperty(parentPath string, objRef *types.ManagedObjectReference) (*DataSource, error) {
	var dataSources DataSource
	if objRef == nil {
		err := errors.New(msg.ERROR_OBJECT_IS_NULL, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL))
		logrus.Error(err)
		return nil, err
	}
	switch objRef.Type {
	case INVENTORY_FOLDER:
		var inventoryObj mo.Folder
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_FOLDER
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_DATACENTER:
		var inventoryObj mo.Datacenter
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_DATACENTER
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_CLUSTERCOMPUTERESOURCE:
		var inventoryObj mo.ClusterComputeResource
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_CLUSTERCOMPUTERESOURCE
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_COMPUTERESOURCE:
		var inventoryObj mo.ComputeResource
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_COMPUTERESOURCE
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_DATASTORE:
		var inventoryObj mo.Datastore
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_DATASTORE
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_DISTRIBUTEDVIRTUALSWITCH:
		var inventoryObj mo.DistributedVirtualSwitch
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Expanded = false
		dataSources.CheckAble = true
		dataSources.Type = INVENTORY_DISTRIBUTEDVIRTUALSWITCH
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_HOSTSYSTEM:
		var inventoryObj mo.HostSystem
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_HOSTSYSTEM
		dataSources.Expanded = false
		dataSources.CheckAble = false
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_NETWORK:
		var inventoryObj mo.Network
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_NETWORK
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_RESOURCEPOOL:
		var inventoryObj mo.ResourcePool
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_RESOURCEPOOL
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_VIRTUALAPP:
		var inventoryObj mo.VirtualApp
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_VIRTUALAPP
		dataSources.Expanded = true
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	case INVENTORY_VIRTALMACHINE:
		var inventoryObj mo.VirtualMachine
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*objRef}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		dataSources.Name = inventoryObj.Name
		if inventoryObj.Config != nil {
			dataSources.Uuid = inventoryObj.Config.InstanceUuid
		}
		dataSources.Type = INVENTORY_VIRTALMACHINE
		dataSources.Expanded = false
		dataSources.CheckAble = true
		if parentPath == "" {
			dataSources.Path = inventoryObj.Name
		} else {
			dataSources.Path = fmt.Sprintf("%s/%s", parentPath, inventoryObj.Name)
		}
	default:
		dataSources.Name = ""
		dataSources.Uuid = ""
		dataSources.Type = INVENTORY_UNDEFINE
		dataSources.Expanded = false
		dataSources.CheckAble = false
		dataSources.Path = parentPath
	}
	return &dataSources, nil
}

// 获取所有主机对象
func (v *VsphereClient) GetAllHostSystem(path string) ([]*types.ManagedObjectReference, error) {
	var hosts []*types.ManagedObjectReference
	datasources, err := v.GetPathSubObjects(path)
	if err != nil {
		err = errors.New(msg.ERROR_GET_ALL_HOST_OBJECT, msg.GetMsg(msg.ERROR_GET_ALL_HOST_OBJECT, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	for _, datasource := range datasources {
		if datasource.Type == INVENTORY_HOSTSYSTEM {
			host, err := v.findByInventoryPath(datasource.Path)
			if err != nil {
				err = errors.New(msg.ERROR_GET_ALL_HOST_OBJECT, msg.GetMsg(msg.ERROR_GET_ALL_HOST_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			hosts = append(hosts, host)
			continue
		}
		if datasource.Type == INVENTORY_RESOURCEPOOL || datasource.Type == INVENTORY_DATASTORE ||
			datasource.Type == INVENTORY_NETWORK || datasource.Type == INVENTORY_VIRTUALAPP ||
			datasource.Type == INVENTORY_DISTRIBUTEDVIRTUALSWITCH || datasource.Type == INVENTORY_VIRTALMACHINE {
			continue
		}
		if datasource.Expanded {
			tempHosts, err := v.GetAllHostSystem(datasource.Path)
			if err != nil {
				err = errors.New(msg.ERROR_GET_ALL_HOST_OBJECT, msg.GetMsg(msg.ERROR_GET_ALL_HOST_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			hosts = append(hosts, tempHosts...)
		}
	}
	return hosts, nil
}

// 根据名称获取主机对象
func (v *VsphereClient) GetHostSystemByName(hostName string) (*mo.HostSystem, error) {
	hostVec, err := v.GetAllHostSystem("")
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_NAME, msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_NAME, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	for _, iter := range hostVec {
		hostInfo, err := v.getObjectProperty("", iter)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_NAME, msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_NAME, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		var inventoryObj mo.HostSystem
		if hostInfo.Name == hostName {
			HostSystemErr := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*iter}, []string{}, &inventoryObj)
			if HostSystemErr != nil {
				err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_NAME,
					msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_NAME, HostSystemErr.Error()))
				logrus.Error(HostSystemErr)
				return nil, err
			}
			return &inventoryObj, nil
		}
	}
	err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_NAME,
		msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_NAME, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL)))
	return nil, err
}

// 根据主机路径获取主机对象
func (v *VsphereClient) getHostObjectByPath(hostPath string) (*mo.HostSystem, error) {
	computeResourceRef, err := v.findByInventoryPath(hostPath)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
			msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	result, err := v.getObjectProperty("", computeResourceRef)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var inventoryObj mo.HostSystem
	if result.Type == INVENTORY_CLUSTERCOMPUTERESOURCE {
		var ccrObj mo.ClusterComputeResource
		err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*computeResourceRef}, []string{}, &ccrObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		logrus.Info(fmt.Sprintf("【%s】 集群下主机数量: %d", hostPath, len(ccrObj.Host)))
		if len(ccrObj.Host) == 0 {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL)))
			logrus.Error(err)
			return nil, err
		}
		err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{ccrObj.Host[0]}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
			logrus.Error(err)
			return nil, err
		}
	} else if result.Type == INVENTORY_COMPUTERESOURCE {
		var crObj mo.ComputeResource
		err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*computeResourceRef}, []string{}, &crObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
			logrus.Error(err)
			return nil, err
		}

		logrus.Info(fmt.Sprintf("【%s】 路径下主机数量: %d", hostPath, len(crObj.Host)))
		if len(crObj.Host) == 0 {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL)))
			logrus.Error(err)
			return nil, err
		}
		err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{crObj.Host[0]}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
			logrus.Error(err)
			return nil, err
		}
	} else {
		err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
			msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL)))
		logrus.Error(err)
		return nil, err
	}
	return &inventoryObj, nil
}

// 根据主机路径获取主机存储管理对象
func (v *VsphereClient) GetHostDataStoreSystem(hosName string) (*object.HostDatastoreSystem, error) {
	// hostObj, err := v.getHostObjectByPath(hostPath)
	hostObj, err := v.GetHostSystemByName(hosName)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_MANAGE_OBJECT_BY_NAME,
			msg.GetMsg(msg.ERROR_GET_HOST_MANAGE_OBJECT_BY_NAME, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	// configManager
	datastoreSystem := object.NewHostDatastoreSystem(v.Vmomi.Client, hostObj.ConfigManager.DatastoreSystem.Reference())
	return datastoreSystem, nil
}

// 根据路径获取主机列表
func (v *VsphereClient) GetHostObjectListByPath(hostPath string) ([]string, error) {
	var hostList []string
	computeResourceRef, err := v.findByInventoryPath(hostPath)
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
			msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	result, err := v.getObjectProperty("", computeResourceRef)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if result.Type == INVENTORY_CLUSTERCOMPUTERESOURCE {
		var ccrObj mo.ClusterComputeResource
		err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*computeResourceRef}, []string{}, &ccrObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
			logrus.Error(err)
		}
		logrus.Info(fmt.Sprintf("【%s】 集群下主机数量: %d", hostPath, len(ccrObj.Host)))
		for _, iter := range ccrObj.Host {
			hostInfo, getErr := v.getObjectProperty("", &iter)
			if getErr != nil {
				logrus.Error(getErr)
				continue
			}
			logrus.Info(hostInfo.Name)
			hostList = append(hostList, hostInfo.Name)
		}

	} else if result.Type == INVENTORY_COMPUTERESOURCE {
		var crObj mo.ComputeResource
		err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*computeResourceRef}, []string{}, &crObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
				msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, err.Error()))
			logrus.Error(err)
		}

		logrus.Info(fmt.Sprintf("【%s】 路径下主机数量: %d", hostPath, len(crObj.Host)))
		for _, iter := range crObj.Host {
			hostInfo, getErr := v.getObjectProperty("", &iter)
			if getErr != nil {
				logrus.Error(getErr)
				continue
			}
			logrus.Info(hostInfo.Name)
			hostList = append(hostList, hostInfo.Name)

		}

	} else {
		err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_PATH,
			msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_PATH, msg.GetMsg(msg.ERROR_OBJECT_IS_NULL)))
		logrus.Error(err)
		return nil, err
	}
	logrus.Info(fmt.Sprintf("主机列表:%s", hostList))
	return hostList, err

}
