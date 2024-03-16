package vsphere

import (
	"context"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"net/url"
	"time"
)

func NewClient(login *contants.LoginInfo) (client *VsphereClient, err error) {
	v := &VsphereClient{loginInfo: login, Ctx: context.Background()}
	err = v.Login()
	if err == nil {
		client = v
	}
	return
}

func (v *VsphereClient) Login() (err error) {
	userInfo := &url.URL{
		Scheme: "https",
		Host:   v.loginInfo.Ip,
		Path:   "/sdk",
	}
	userInfo.User = url.UserPassword(v.loginInfo.UserName, v.loginInfo.PassWord)
	client, err := govmomi.NewClient(v.Ctx, userInfo, true)
	if err != nil {
		err = errors.New(msg.ERROR_LOGIN_VSPHERE, msg.GetMsg(msg.ERROR_LOGIN_VSPHERE, v.loginInfo.Ip, err.Error()))
		logrus.Error(err)
		return err
	}
	if client == nil {
		err = errors.New(msg.ERROR_LOGIN_VSPHERE,
			msg.GetMsg(msg.ERROR_LOGIN_VSPHERE, v.loginInfo.Ip, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return
	}
	v.Vmomi = client
	return
}

func (v *VsphereClient) ReLogin() {
	if v.Vmomi == nil {
		err := v.Login()
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (v *VsphereClient) Logout() {
	err := v.Vmomi.Logout(v.Ctx)
	if err != nil {
		logrus.Error(err)
	}
}

// 平台类型
func (v *VsphereClient) IsvCenter() bool {
	return v.Vmomi.Client.IsVC()
}

// 获取平台版本
func (v *VsphereClient) GetVersion() string {
	version := v.Vmomi.ServiceContent.About.ApiVersion
	if v.Vmomi.Client.IsVC() {
		return fmt.Sprintf("vCenter %s", version)
	} else {
		return fmt.Sprintf("ESXi %s", version)
	}
}

// 获取所有虚拟机
func (v *VsphereClient) GetVMs() ([]DataSource, error) {
	var rootObj []mo.Folder
	objs := []types.ManagedObjectReference{v.getRootPath()}
	err := v.Vmomi.Retrieve(v.Ctx, objs, []string{}, &rootObj)
	if err != nil {
		err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	m := view.NewManager(v.Vmomi.Client)

	cv, err := m.CreateContainerView(v.Ctx, objs[0], []string{"VirtualMachine"}, true)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_VSPHERE_CONTAINER_VIEW,
			msg.GetMsg(msg.ERROR_CREATE_VSPHERE_CONTAINER_VIEW, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	var vms []mo.VirtualMachine
	err = cv.Retrieve(v.Ctx, []string{"VirtualMachine"}, []string{}, &vms)
	if err != nil {
		err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
		return nil, err
	}
	var dataSources []DataSource
	for _, iter := range vms {
		var dataSource DataSource
		dataSource.Type = INVENTORY_VIRTALMACHINE
		dataSource.Name = iter.Name
		dataSource.Expanded = false
		dataSource.Ip = iter.Summary.Guest.IpAddress
		if iter.Config != nil {
			dataSource.Uuid = iter.Config.InstanceUuid
		}
		dataSources = append(dataSources, dataSource)
	}
	return dataSources, nil
}

// 获取指定路径下的子对象
func (v *VsphereClient) GetPathSubObjects(path string) ([]*DataSource, error) {
	v.ReLogin()
	var dataSources []*DataSource
	if path == "" {
		var rootObj []mo.Folder
		objs := []types.ManagedObjectReference{v.getRootPath()}
		err := v.Vmomi.Retrieve(v.Ctx, objs, []string{}, &rootObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		for _, obj := range rootObj[0].ChildEntity {
			dataSource, err := v.getObjectProperty("", &obj)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, dataSource)
		}
	} else {
		obj, err := v.findByInventoryPath(path)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		switch obj.Type {
		case INVENTORY_FOLDER:
			var folder mo.Folder
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &folder)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range folder.ChildEntity {
				dataSource, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, dataSource)
			}
		case INVENTORY_DATACENTER:
			var dc mo.Datacenter
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &dc)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			// 获取计算资源
			hostFolder, err := v.getObjectProperty(path, &dc.HostFolder)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, hostFolder)
			// 获取存储
			dataStoreFolder, err := v.getObjectProperty(path, &dc.DatastoreFolder)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, dataStoreFolder)
			// 获取网卡
			networkFolder, err := v.getObjectProperty(path, &dc.NetworkFolder)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, networkFolder)
			// 获取vmFolder
			vmFolder, err := v.getObjectProperty(path, &dc.VmFolder)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, vmFolder)
		case INVENTORY_CLUSTERCOMPUTERESOURCE:
			var ccr mo.ClusterComputeResource
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &ccr)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			// 获取主机
			for _, temp := range ccr.Host {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
			resourcesPool, err := v.getObjectProperty(path, ccr.ResourcePool)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, resourcesPool)
		case INVENTORY_COMPUTERESOURCE:
			var cr mo.ComputeResource
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &cr)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			// 获取主机
			for _, temp := range cr.Host {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
			data, err := v.getObjectProperty(path, cr.ResourcePool)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, data)
		case INVENTORY_HOSTSYSTEM:
		case INVENTORY_DATASTORE:
			var ds mo.Datastore
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &ds)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range ds.Vm {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
		case INVENTORY_NETWORK:
			var net mo.Network
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &net)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range net.Vm {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
		case INVENTORY_RESOURCEPOOL:
			var rs mo.ResourcePool
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &rs)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range rs.ResourcePool {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
			for _, temp := range rs.Vm {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
		case INVENTORY_VIRTUALAPP:
			var vApp mo.VirtualApp
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &vApp)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range vApp.ResourcePool.ResourcePool {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
			for _, temp := range vApp.Vm {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
		default:
			err = errors.New(msg.ERROR_GET_SUB_OBJECT_UN_DEFINE,
				msg.GetMsg(msg.ERROR_GET_SUB_OBJECT_UN_DEFINE))
			logrus.Error(err)
			return nil, err
		}
	}
	return dataSources, nil
}

// 按计算资源的方式获取数据源
func (v *VsphereClient) GetDataSourcesByComputeResource(path string, isGetVm bool) ([]*DataSource, error) {
	v.ReLogin()
	var dataSources []*DataSource
	if path == "" {
		var rootObj []mo.Folder
		objs := []types.ManagedObjectReference{v.getRootPath()}
		err := v.Vmomi.Retrieve(v.Ctx, objs, []string{}, &rootObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		for _, obj := range rootObj[0].ChildEntity {
			dataSource, err := v.getObjectProperty("", &obj)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, dataSource)
		}
	} else {
		obj, err := v.findByInventoryPath(path)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		//logrus.Info(fmt.Sprintf("获取计算资源类型:%s", obj.Type))
		switch obj.Type {
		case INVENTORY_FOLDER:
			var folder mo.Folder
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &folder)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range folder.ChildEntity {
				dataSource, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, dataSource)
			}
		case INVENTORY_DATACENTER:
			var dc mo.Datacenter
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &dc)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			// 获取计算资源
			hostFolder, err := v.getObjectProperty(path, &dc.HostFolder)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			// 再次获取，过滤掉隐藏路径[host]
			hostFolderSub, err := v.GetPathSubObjects(hostFolder.Path)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, hostFolderSub...)
		case INVENTORY_CLUSTERCOMPUTERESOURCE:
			var ccr mo.ClusterComputeResource
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &ccr)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			temp, err := v.getObjectProperty(path, ccr.ResourcePool)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			// 再次获取，过滤掉根资源池[Resources]
			subObjects, err := v.GetPathSubObjects(temp.Path)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			for _, subObject := range subObjects {
				// 是否获取虚拟机
				if !isGetVm && (subObject.Type == INVENTORY_VIRTALMACHINE) {
					continue
				}
				dataSources = append(dataSources, subObject)
			}
		case INVENTORY_COMPUTERESOURCE:
			var cr mo.ComputeResource
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &cr)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			temp, err := v.getObjectProperty(path, cr.ResourcePool)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			// 再次获取，过滤掉根资源池[Resources]
			subObjects, err := v.GetPathSubObjects(temp.Path)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			for _, subObject := range subObjects {
				// 是否获取虚拟机
				if !isGetVm && (subObject.Type == INVENTORY_VIRTALMACHINE) {
					continue
				}
				dataSources = append(dataSources, subObject)
			}
		case INVENTORY_RESOURCEPOOL:
			var rs mo.ResourcePool
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &rs)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range rs.ResourcePool {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
			// 是否获取虚拟机
			if isGetVm {
				for _, temp := range rs.Vm {
					data, err := v.getObjectProperty(path, &temp)
					if err != nil {
						logrus.Error(err)
						return nil, err
					}
					dataSources = append(dataSources, data)
				}
			}
		case INVENTORY_VIRTUALAPP:
			var vApp mo.VirtualApp
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &vApp)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range vApp.ResourcePool.ResourcePool {
				data, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				dataSources = append(dataSources, data)
			}
			// 是否获取虚拟机
			if isGetVm {
				for _, temp := range vApp.Vm {
					data, err := v.getObjectProperty(path, &temp)
					if err != nil {
						logrus.Error(err)
						return nil, err
					}
					dataSources = append(dataSources, data)
				}
			}
		default:
			err = errors.New(msg.ERROR_GET_SUB_OBJECT_UN_DEFINE,
				msg.GetMsg(msg.ERROR_GET_SUB_OBJECT_UN_DEFINE))
			logrus.Error(err)
			return nil, err
		}
	}
	return dataSources, nil
}

// 按虚拟机模板的方式获取数据源
func (v *VsphereClient) GetDataSourcesByVmTemplate(path string, isGetVm bool) ([]*DataSource, error) {
	v.ReLogin()
	var dataSources []*DataSource
	if path == "" {
		var rootObj []mo.Folder
		objs := []types.ManagedObjectReference{v.getRootPath()}
		err := v.Vmomi.Retrieve(v.Ctx, objs, []string{}, &rootObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		for _, obj := range rootObj[0].ChildEntity {
			dataSource, err := v.getObjectProperty("", &obj)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			dataSources = append(dataSources, dataSource)
		}
	} else {
		obj, err := v.findByInventoryPath(path)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		switch obj.Type {
		case INVENTORY_FOLDER:
			var folder mo.Folder
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &folder)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			for _, temp := range folder.ChildEntity {
				dataSource, err := v.getObjectProperty(path, &temp)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				// 是否获取虚拟机
				if !isGetVm && (dataSource.Type == INVENTORY_VIRTALMACHINE) {
					continue
				}
				dataSources = append(dataSources, dataSource)
			}
		case INVENTORY_DATACENTER:
			var dc mo.Datacenter
			err = v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*obj}, []string{}, &dc)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			// 获取虚拟机文件夹
			vmFolder, err := v.getObjectProperty(path, &dc.VmFolder)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			// 再次获取，过滤掉隐藏路径[vm]
			vmFolderSub, err := v.GetPathSubObjects(vmFolder.Path)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			for _, subObject := range vmFolderSub {
				// 是否获取虚拟机
				if !isGetVm && (subObject.Type == INVENTORY_VIRTALMACHINE) {
					continue
				}
				dataSources = append(dataSources, subObject)
			}
		case INVENTORY_VIRTUALAPP:
			// todo: 过滤不展示
		default:
			err = errors.New(msg.ERROR_GET_SUB_OBJECT_UN_DEFINE,
				msg.GetMsg(msg.ERROR_GET_SUB_OBJECT_UN_DEFINE))
			logrus.Error(err)
			return nil, err
		}
	}
	return dataSources, nil
}

// 配置虚拟机
func (v *VsphereClient) Customize(uuid string, addr *IpAddr) error {
	v.ReLogin()
	vmRef, err := v.findByUUid(uuid)
	if err != nil {
		err = errors.New(msg.ERROR_CONFIG_VM_FAILED, msg.GetMsg(msg.ERROR_CONFIG_VM_FAILED, err.Error()))
		return err
	}
	if addr.Ip == "" || addr.Netmask == "" || addr.Gateway == "" {
		return nil
	}
	vmObj := object.NewVirtualMachine(v.Vmomi.Client, vmRef.Reference())
	cam := types.CustomizationAdapterMapping{
		Adapter: types.CustomizationIPSettings{
			Ip:         &types.CustomizationFixedIp{IpAddress: addr.Ip},
			SubnetMask: addr.Netmask,
			Gateway:    []string{addr.Gateway},
		},
	}
	customSpec := types.CustomizationSpec{
		NicSettingMap: []types.CustomizationAdapterMapping{cam},
		Identity:      &types.CustomizationLinuxPrep{HostName: &types.CustomizationFixedName{Name: addr.Hostname}},
	}
	task, err := vmObj.Customize(v.Ctx, customSpec)
	if err != nil {
		err = errors.New(msg.ERROR_CONFIG_VM_FAILED, msg.GetMsg(msg.ERROR_CONFIG_VM_FAILED, err.Error()))
		return err
	}
	err = task.Wait(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_CONFIG_VM_FAILED, msg.GetMsg(msg.ERROR_CONFIG_VM_FAILED, err.Error()))
		return err
	}
	return nil
}

// 开机
func (v *VsphereClient) PowerOn(uuid string) error {
	v.ReLogin()
	vmRef, err := v.findByUUid(uuid)
	if err != nil {
		err = errors.New(msg.ERROR_POWER_ON_VM_FAILED, msg.GetMsg(msg.ERROR_POWER_ON_VM_FAILED, err.Error()))
		return err
	}
	vmObj := object.NewVirtualMachine(v.Vmomi.Client, vmRef.Reference())
	task, err := vmObj.PowerOn(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_POWER_ON_VM_FAILED, msg.GetMsg(msg.ERROR_POWER_ON_VM_FAILED, err.Error()))
		return err
	}
	err = task.Wait(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_POWER_ON_VM_FAILED, msg.GetMsg(msg.ERROR_POWER_ON_VM_FAILED, err.Error()))
		return err
	}
	return nil
}

// 关机
func (v *VsphereClient) PowerOff(uuid string) error {
	v.ReLogin()
	vmRef, err := v.findByUUid(uuid)
	if err != nil {
		err = errors.New(msg.ERROR_POWER_OFF_VM_FAILED, msg.GetMsg(msg.ERROR_POWER_OFF_VM_FAILED, err.Error()))
		return err
	}
	vmObj := object.NewVirtualMachine(v.Vmomi.Client, vmRef.Reference())
	task, err := vmObj.PowerOff(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_POWER_OFF_VM_FAILED, msg.GetMsg(msg.ERROR_POWER_OFF_VM_FAILED, err.Error()))
		return err
	}
	err = task.Wait(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_POWER_OFF_VM_FAILED, msg.GetMsg(msg.ERROR_POWER_OFF_VM_FAILED, err.Error()))
		return err
	}
	return nil
}

// 创建快照
func (v *VsphereClient) CreateSnapShot(uuid string, name, desc string) error {
	v.ReLogin()
	vmRef, err := v.findByUUid(uuid)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_VM_SNAPSHOT_FAILED, msg.GetMsg(msg.ERROR_CREATE_VM_SNAPSHOT_FAILED, err.Error()))
		return err
	}
	vmObj := object.NewVirtualMachine(v.Vmomi.Client, vmRef.Reference())
	task, err := vmObj.CreateSnapshot(v.Ctx, name, desc, false, false)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_VM_SNAPSHOT_FAILED, msg.GetMsg(msg.ERROR_CREATE_VM_SNAPSHOT_FAILED, err.Error()))
		return err
	}
	err = task.Wait(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_VM_SNAPSHOT_FAILED, msg.GetMsg(msg.ERROR_CREATE_VM_SNAPSHOT_FAILED, err.Error()))
		return err
	}
	return nil
}

// 恢复快照
func (v *VsphereClient) RecoverToSnapshot(uuid string, name string, suppressPowerOn bool) error {
	v.ReLogin()
	vmRef, err := v.findByUUid(uuid)
	if err != nil {
		err = errors.New(msg.ERROR_RECOVER_VM_SNAPSHOT_FAILED,
			msg.GetMsg(msg.ERROR_RECOVER_VM_SNAPSHOT_FAILED, err.Error()))
		return err
	}
	vmObj := object.NewVirtualMachine(v.Vmomi.Client, vmRef.Reference())
	task, err := vmObj.RevertToSnapshot(v.Ctx, name, suppressPowerOn)
	if err != nil {
		err = errors.New(msg.ERROR_RECOVER_VM_SNAPSHOT_FAILED,
			msg.GetMsg(msg.ERROR_RECOVER_VM_SNAPSHOT_FAILED, err.Error()))
		return err
	}
	err = task.Wait(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_RECOVER_VM_SNAPSHOT_FAILED,
			msg.GetMsg(msg.ERROR_RECOVER_VM_SNAPSHOT_FAILED, err.Error()))
		return err
	}
	return nil
}

func (v *VsphereClient) GetHostNameByPath(hostPath string) (string, error) {
	hostObj, err := v.getHostObjectByPath(hostPath)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return hostObj.Name, nil
}

// 获取NFS存储列表
func (v *VsphereClient) GetNasDatastore() ([]string, error) {
	var dataStoreNames []string
	hostVec, err := v.GetAllHostSystem("")
	if err != nil {
		err = errors.New(msg.ERROR_GET_HOST_OBJECT_BY_NAME, msg.GetMsg(msg.ERROR_GET_HOST_OBJECT_BY_NAME, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	for _, host := range hostVec {
		var inventoryObj mo.HostSystem
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*host}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
			logrus.Error(err)
			return nil, err
		}
		for _, storeObj := range inventoryObj.Datastore {
			var dataStore mo.Datastore
			err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{storeObj}, []string{}, &dataStore)
			if err != nil {
				err = errors.New(msg.ERROR_GET_INVENTORY_OBJECT, msg.GetMsg(msg.ERROR_GET_INVENTORY_OBJECT, err.Error()))
				logrus.Error(err)
				return nil, err
			}
			if dataStore.Summary.Type == "NFS" || dataStore.Summary.Type == "NFS41" {
				dataStoreNames = append(dataStoreNames, dataStore.Name)
			}
		}
	}
	return dataStoreNames, nil
}

// 创建nas存储
func (v *VsphereClient) CreateNasDatastore(storeInfo *NasDataStoreParam) error {
	logrus.Info(fmt.Sprintf("创建nas存储参数:"))
	logrus.Info(fmt.Sprintf("【StoreType】:%s", storeInfo.StoreType))
	logrus.Info(fmt.Sprintf("【LocalPath】:%s", storeInfo.LocalPath))
	// logrus.Info(fmt.Sprintf("【HostPath】:%s", storeInfo.HostPath))
	logrus.Info(fmt.Sprintf("【HostName】:%s", storeInfo.HostName))
	logrus.Info(fmt.Sprintf("【RemoteHost】:%s", storeInfo.RemoteHost))
	logrus.Info(fmt.Sprintf("【RemotePath】:%s", storeInfo.RemotePath))
	logrus.Info(fmt.Sprintf("【AccessMode】:%s", storeInfo.AccessMode))
	v.ReLogin()
	// sysObj, err := v.GetHostDataStoreSystem(storeInfo.HostPath)
	sysObj, err := v.GetHostDataStoreSystem(storeInfo.HostName)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_NFS_DATASTORE_FAILED,
			msg.GetMsg(msg.ERROR_CREATE_NFS_DATASTORE_FAILED, err.Error()))
		logrus.Error(err)
		return err
	}
	hostNasSpec := types.HostNasVolumeSpec{
		RemoteHost:      storeInfo.RemoteHost,
		RemotePath:      storeInfo.RemotePath,
		LocalPath:       storeInfo.LocalPath,
		AccessMode:      storeInfo.AccessMode,
		Type:            storeInfo.StoreType,
		UserName:        "",
		Password:        "",
		RemoteHostNames: []string{storeInfo.RemoteHost},
		SecurityType:    "AUTH_SYS"}
	_, err = sysObj.CreateNasDatastore(v.Ctx, hostNasSpec)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_NFS_DATASTORE_FAILED,
			msg.GetMsg(msg.ERROR_CREATE_NFS_DATASTORE_FAILED, err.Error()))
		return err
	}
	return nil
}

// 卸载NFS存储
func (v *VsphereClient) RemoveNasDatastore(storeName string) error {
	hostVec, err := v.GetAllHostSystem("")
	if err != nil {
		err = errors.New(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED,
			msg.GetMsg(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED, err.Error()))
		logrus.Error(err)
		return err
	}
	for _, host := range hostVec {
		var inventoryObj mo.HostSystem
		err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{*host}, []string{}, &inventoryObj)
		if err != nil {
			err = errors.New(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED,
				msg.GetMsg(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED, err.Error()))
			logrus.Error(err)
			return err
		}
		for _, storeObj := range inventoryObj.Datastore {
			var dataStore mo.Datastore
			err := v.Vmomi.Retrieve(v.Ctx, []types.ManagedObjectReference{storeObj}, []string{}, &dataStore)
			if err != nil {
				err = errors.New(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED,
					msg.GetMsg(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED, err.Error()+"---"+dataStore.Name))
				logrus.Error(err)
				return err
			}
			time.Sleep(time.Second * 2)
			if dataStore.Name == storeName {
				logrus.Info(fmt.Sprintf("执行卸载存储【%s】", dataStore.Name))
				ds := object.NewDatastore(v.Vmomi.Client, dataStore.Reference())
				datastoreSystem := object.NewHostDatastoreSystem(v.Vmomi.Client, inventoryObj.ConfigManager.DatastoreSystem.Reference())
				err := datastoreSystem.Remove(v.Ctx, ds)
				if err != nil {
					err = errors.New(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED,
						msg.GetMsg(msg.ERROR_REMOVE_NFS_DATASTORE_FAILED, err.Error()))
					logrus.Error(err)
					return err
				}
			}
		}
	}
	return nil
}

// 注册虚拟
func (v *VsphereClient) RegisterVm(param *RegisterVmParam) error {
	v.ReLogin()
	// 获取folder对象
	var folder *object.Folder
	if param.FolderPath != "" {
		folderObj, err := v.findByInventoryPath(param.FolderPath)
		if err != nil {
			err = errors.New(msg.ERROR_REGISTER_VM_FAILED,
				msg.GetMsg(msg.ERROR_REGISTER_VM_FAILED, err.Error()))
			logrus.Error(err)
			return err
		}
		folder = object.NewFolder(v.Vmomi.Client, *folderObj)
	} else {
		err := errors.New(msg.ERROR_REGISTER_VM_FAILED,
			msg.GetMsg(msg.ERROR_REGISTER_VM_FAILED, msg.GetMsg(msg.ERROR_PARAM_IS_NULL)))
		logrus.Error(err)
		return err
	}
	// 获取主机对象
	var hostSystem *object.HostSystem
	if param.HostPath != "" {
		hostObj, err := v.getHostObjectByPath(param.HostPath)
		if err != nil {
			logrus.Error(err)
		} else {
			hostSystem = object.NewHostSystem(v.Vmomi.Client, hostObj.Reference())
		}
	}
	// 获取资源池对象
	var resourcePool *object.ResourcePool
	if param.ResourcePoolPath != "" {
		resourcePoolObj, err := v.findByInventoryPath(param.ResourcePoolPath)
		if err != nil {
			logrus.Error(err)
		} else {
			resourcePool = object.NewResourcePool(v.Vmomi.Client, *resourcePoolObj)
		}
	}
	task, err := folder.RegisterVM(v.Ctx, param.VmxPath, param.VmName, false, resourcePool, hostSystem)
	if err != nil {
		err = errors.New(msg.ERROR_REGISTER_VM_FAILED,
			msg.GetMsg(msg.ERROR_REGISTER_VM_FAILED, err.Error()))
		return err
	}
	err = task.Wait(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_REGISTER_VM_FAILED,
			msg.GetMsg(msg.ERROR_REGISTER_VM_FAILED, err.Error()))
		return err
	}
	return nil
}

// 取消注册
func (v *VsphereClient) UnRegisterVm(uuid string) error {
	v.ReLogin()
	vmRef, err := v.findByUUid(uuid)
	if err != nil {
		err = errors.New(msg.ERROR_UN_REGISTER_VM_FAILED,
			msg.GetMsg(msg.ERROR_UN_REGISTER_VM_FAILED, err.Error()))
		return err
	}
	vmObj := object.NewVirtualMachine(v.Vmomi.Client, vmRef.Reference())
	err = vmObj.Unregister(v.Ctx)
	if err != nil {
		err = errors.New(msg.ERROR_UN_REGISTER_VM_FAILED,
			msg.GetMsg(msg.ERROR_UN_REGISTER_VM_FAILED, err.Error()))
		logrus.Error(err)
		return err
	}
	return nil
}
