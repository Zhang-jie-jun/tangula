package vmware

import (
	"context"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/host"
	"github.com/Zhang-jie-jun/tangula/internal/dao/platform"
	"github.com/Zhang-jie-jun/tangula/internal/vsphere"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/service/app/linux"
	"github.com/Zhang-jie-jun/tangula/service/app/windows"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"net/url"
	"os"
	"strings"
)

type VmWare struct {
	IP     string
	User   string
	Pwd    string
	Client *govmomi.Client
	Ctx    context.Context
}

func NewVmWare(IP, User, Pwd string) *VmWare {
	u := &url.URL{
		Scheme: "https",
		Host:   IP,
		Path:   "/sdk",
	}
	ctx := context.Background()
	u.User = url.UserPassword(User, Pwd)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		panic(err)
	}
	return &VmWare{
		IP:     IP,
		User:   User,
		Pwd:    Pwd,
		Client: client,
		Ctx:    ctx,
	}
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func (vw *VmWare) getBase(tp string) (v *view.ContainerView, error error) {
	m := view.NewManager(vw.Client.Client)

	v, err := m.CreateContainerView(vw.Ctx, vw.Client.Client.ServiceContent.RootFolder, []string{tp}, true)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (vw *VmWare) GetAllVmClient() (vmList []VirtualMachines, templateList []TemplateInfo, err error) {
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, nil, err
	}
	defer v.Destroy(vw.Ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.Ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, nil, err
	}
	for _, vm := range vms {
		//if vm.Summary.Config.Name == "测试机器" {
		//	v := object.NewVirtualMachine(vw.Client.Client, vm.Self)
		//	vw.SetIP(v)
		//}
		if vm.Summary.Config.Template {
			templateList = append(templateList, TemplateInfo{
				Name:   vm.Summary.Config.Name,
				System: vm.Summary.Config.GuestFullName,
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		} else {
			vmList = append(vmList, VirtualMachines{
				Uuid:       vm.Summary.Config.InstanceUuid,
				Name:       vm.Summary.Config.Name,
				System:     vm.Summary.Config.GuestFullName,
				Ip:         vm.Summary.Guest.IpAddress,
				Hostname:   vm.Summary.Guest.HostName,
				PowerState: string(vm.Summary.Runtime.PowerState),
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		}
	}
	fmt.Println(vmList)
	return vmList, templateList, nil
}

func (vw *VmWare) GetAllHost() (hostList []*HostSummary, err error) {
	v, err := vw.getBase("HostSystem")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var hss []mo.HostSystem
	err = v.Retrieve(vw.Ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}
	for _, hs := range hss {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		hostList = append(hostList, &HostSummary{
			Host: Host{
				Type:  hs.Summary.Host.Type,
				Value: hs.Summary.Host.Value,
			},
			Name:        hs.Summary.Config.Name,
			UsedCPU:     int64(hs.Summary.QuickStats.OverallCpuUsage),
			TotalCPU:    totalCPU,
			FreeCPU:     freeCPU,
			UsedMemory:  int64((units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage)) * 1024 * 1024),
			TotalMemory: int64(units.ByteSize(hs.Summary.Hardware.MemorySize)),
			FreeMemory:  freeMemory,
			HostSelf:    hs.Self,
		})
	}
	return hostList, err
}

func (vw *VmWare) GetAllNetwork() (networkList []map[string]string, err error) {
	v, err := vw.getBase("Network")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var networks []mo.Network
	err = v.Retrieve(vw.Ctx, []string{"Network"}, nil, &networks)
	if err != nil {
		return nil, err
	}
	for _, net := range networks {
		networkList = append(networkList, map[string]string{
			"Vlan":      net.Name,
			"NetworkID": strings.Split(net.Reference().String(), ":")[1],
		})
	}
	return networkList, nil
}

func (vw *VmWare) GetAllDatastore() (datastoreList []DatastoreSummary, err error) {
	v, err := vw.getBase("Datastore")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var dss []mo.Datastore
	err = v.Retrieve(vw.Ctx, []string{"Datastore"}, []string{"summary"}, &dss)
	if err != nil {
		return nil, err
	}
	for _, ds := range dss {
		datastoreList = append(datastoreList, DatastoreSummary{
			Name: ds.Summary.Name,
			Datastore: Datastore{
				Type:  ds.Summary.Datastore.Type,
				Value: ds.Summary.Datastore.Value,
			},
			Type:          ds.Summary.Type,
			Capacity:      int64(units.ByteSize(ds.Summary.Capacity)),
			FreeSpace:     int64(units.ByteSize(ds.Summary.FreeSpace)),
			DatastoreSelf: ds.Self,
		})
	}
	return
}

func (vw *VmWare) GetHostVm() (hostVm map[string][]VMS, err error) {
	hostList, err := vw.GetAllHost() //
	if err != nil {
		return
	}
	var hostIDList []string
	hostVm = make(map[string][]VMS)
	for _, host := range hostList {
		hostIDList = append(hostIDList, host.Host.Value)
		hostVm[host.Host.Value] = []VMS{}
	}
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.Ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, err
	}
	for _, vm := range vms {
		if IsContain(hostIDList, vm.Summary.Runtime.Host.Value) {
			hostVm[vm.Summary.Runtime.Host.Value] = append(hostVm[vm.Summary.Runtime.Host.Value], VMS{
				Name:  vm.Summary.Config.Name,
				Value: vm.Summary.Vm.Value,
			})
		}
		//s, _ := json.Marshal(vm.Summary)
		//fmt.Println(string(s))
	}
	//fmt.Println(hostVm)
	return
}

func (vw *VmWare) GetAllCluster() (clusterList []ClusterInfo, err error) {
	v, err := vw.getBase("ClusterComputeResource")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var crs []mo.ClusterComputeResource
	err = v.Retrieve(vw.Ctx, []string{"ClusterComputeResource"}, []string{}, &crs)
	if err != nil {
		return nil, err
	}
	for _, cr := range crs {
		clusterList = append(clusterList, ClusterInfo{
			Cluster: Self{
				Type:  cr.Self.Type,
				Value: cr.Self.Value,
			},
			Name: cr.Name,
			Parent: Self{
				Type:  cr.Parent.Type,
				Value: cr.Parent.Value,
			},
			ResourcePool: Self{
				Type:  cr.ResourcePool.Type,
				Value: cr.ResourcePool.Value,
			},
			Hosts:     cr.Host,
			Datastore: cr.Datastore,
		})
	}
	fmt.Println(clusterList)
	return
}

func (vw *VmWare) GetAllDatacenter() (dataCenterList []DataCenter, err error) {
	v, err := vw.getBase("Datacenter")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var dcs []mo.Datacenter
	err = v.Retrieve(vw.Ctx, []string{"Datacenter"}, []string{}, &dcs)
	if err != nil {
		return nil, err
	}
	for _, dc := range dcs {
		dataCenterList = append(dataCenterList, DataCenter{
			Datacenter: Self{
				Type:  dc.Self.Type,
				Value: dc.Self.Value,
			},
			Name: dc.Name,
			VmFolder: Self{
				Type:  dc.VmFolder.Type,
				Value: dc.VmFolder.Value,
			},
			HostFolder: Self{
				Type:  dc.HostFolder.Type,
				Value: dc.HostFolder.Value,
			},
			DatastoreFolder: Self{
				Type:  dc.DatastoreFolder.Type,
				Value: dc.DatastoreFolder.Value,
			},
		})
	}
	fmt.Println(dataCenterList)
	return
}

func (vw *VmWare) GetAllResourcePool() (resourceList []ResourcePoolInfo, err error) {
	v, err := vw.getBase("ResourcePool")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var rps []mo.ResourcePool
	err = v.Retrieve(vw.Ctx, []string{"ResourcePool"}, []string{}, &rps)
	for _, rp := range rps {
		//if rp.Name == "测试虚机" {
		//	s, _ := json.Marshal(rp)
		//	fmt.Println(string(s))
		//}
		resourceList = append(resourceList, ResourcePoolInfo{
			ResourcePool: Self{
				Type:  rp.Self.Type,
				Value: rp.Self.Value,
			},
			Name: rp.Name,
			Parent: Self{
				Type:  rp.Parent.Type,
				Value: rp.Parent.Value,
			},
			ResourcePoolList: rp.ResourcePool,
			Resource:         rp.Self,
		})
	}
	return
}

func (vw *VmWare) GetFolder() (folderList []FolderInfo, err error) {
	v, err := vw.getBase("Folder")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.Ctx)
	var folders []mo.Folder
	err = v.Retrieve(vw.Ctx, []string{"Folder"}, []string{}, &folders)
	for _, folder := range folders {
		//newFolder := object.NewFolder(vw.client.Client, folder.Self)
		//fmt.Println(newFolder)
		folderList = append(folderList, FolderInfo{
			Folder: Self{
				Type:  folder.Self.Type,
				Value: folder.Self.Value,
			},
			Name:        folder.Name,
			ChildEntity: folder.ChildEntity,
			Parent: Self{
				Type:  folder.Parent.Type,
				Value: folder.Parent.Value,
			},
			FolderSelf: folder.Self,
		})
		//break
	}
	return folderList, nil
}

func (vw *VmWare) getLibraryItem(ctx context.Context, rc *rest.Client) (*library.Item, error) {
	const (
		libraryName     = "模板"
		libraryItemName = "template-rehl7.7"
		libraryItemType = "ovf"
	)

	m := library.NewManager(rc)
	libraries, err := m.FindLibrary(ctx, library.Find{Name: libraryName})
	if err != nil {
		fmt.Printf("Find library by name %s failed, %v", libraryName, err)
		return nil, err
	}

	if len(libraries) == 0 {
		fmt.Printf("Library %s was not found", libraryName)
		return nil, fmt.Errorf("library %s was not found", libraryName)
	}

	if len(libraries) > 1 {
		fmt.Printf("There are multiple libraries with the name %s", libraryName)
		return nil, fmt.Errorf("there are multiple libraries with the name %s", libraryName)
	}

	items, err := m.FindLibraryItems(ctx, library.FindItem{Name: libraryItemName,
		Type: libraryItemType, LibraryID: libraries[0]})

	if err != nil {
		fmt.Printf("Find library item by name %s failed", libraryItemName)
		return nil, fmt.Errorf("find library item by name %s failed", libraryItemName)
	}

	if len(items) == 0 {
		fmt.Printf("Library item %s was not found", libraryItemName)
		return nil, fmt.Errorf("library item %s was not found", libraryItemName)
	}

	if len(items) > 1 {
		fmt.Printf("There are multiple library items with the name %s", libraryItemName)
		return nil, fmt.Errorf("there are multiple library items with the name %s", libraryItemName)
	}

	item, err := m.GetLibraryItem(ctx, items[0])
	if err != nil {
		fmt.Printf("Get library item by %s failed, %v", items[0], err)
		return nil, err
	}
	return item, nil
}

func (vw *VmWare) CreateVM() {
	createData := CreateMap{
		TempName:    "xxx",
		Datacenter:  "xxx",
		Cluster:     "xxx",
		Host:        "xxx",
		Resources:   "xxx",
		Storage:     "xxx",
		VmName:      "xxx",
		SysHostName: "xxx",
		Network:     "xxx",
	}
	_, templateList, err := vw.GetAllVmClient()
	if err != nil {
		panic(err)
	}
	var templateNameList []string
	for _, template := range templateList {
		templateNameList = append(templateNameList, template.Name)
	}
	if !IsContain(templateNameList, createData.TempName) {
		fmt.Fprintf(os.Stderr, "模版不存在，虚拟机创建失败")
		return
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	for _, resource := range resourceList {
		if resource.Name == createData.Resources {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机创建失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	datastoreList, err := vw.GetAllDatastore()
	if err != nil {
		panic(err)
	}
	var datastoreID, datastoreStr string
	for _, datastore := range datastoreList {
		if datastore.Name == createData.Storage {
			datastoreID = datastore.Datastore.Value
			datastoreStr = datastore.Name
		}
	}
	if datastoreStr == "" {
		fmt.Fprintf(os.Stderr, "存储中心不存在，虚拟机创建失败")
		return
	}
	fmt.Println("DatastoreID", datastoreID)
	networkList, err := vw.GetAllNetwork()
	if err != nil {
		panic(err)
	}
	var networkID, networkStr string
	for _, network := range networkList {
		if network["Vlan"] == createData.Network {
			networkStr = network["Vlan"]
			networkID = network["NetworkID"]
		}
	}

	if networkStr == "" {
		fmt.Fprintf(os.Stderr, "网络不存在，虚拟机创建失败")
		return
	}
	fmt.Println("NetworkID", networkID)
	finder := find.NewFinder(vw.Client.Client)
	//resourcePools, err := finder.DatacenterList(vw.ctx, "*")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to list resource pool at vc %v", err)
	//	os.Exit(1)
	//}
	//fmt.Println(reflect.TypeOf(resourcePools[0].Reference().Value), resourcePools)
	folders, err := finder.FolderList(vw.Ctx, "*")
	var folderID string
	for _, folder := range folders {
		if folder.InventoryPath == "/"+createData.Datacenter+"/vm" {
			folderID = folder.Reference().Value
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list folder at vc  %v", err)
		return
	}
	rc := rest.NewClient(vw.Client.Client)
	if err := rc.Login(vw.Ctx, url.UserPassword(vw.User, vw.Pwd)); err != nil {
		fmt.Fprintf(os.Stderr, "rc Login filed, %v", err)
		return
	}
	item, err := vw.getLibraryItem(vw.Ctx, rc)
	if err != nil {
		panic(err)
	}
	//cloneSpec := &types.VirtualMachineCloneSpec{
	//	PowerOn:  false,
	//	Template: cmd.template,
	//}
	// 7fa9e782-cba2-4061-95fc-4ebb08ec127a
	fmt.Println("Item", item.ID)
	m := vcenter.NewManager(rc)
	fr := vcenter.FilterRequest{
		Target: vcenter.Target{
			ResourcePoolID: resourceID,
			FolderID:       folderID,
		},
	}
	r, err := m.FilterLibraryItem(vw.Ctx, item.ID, fr)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	fmt.Println(11111111111, r.Networks, r.StorageGroups)
	networkKey := r.Networks[0]
	//storageKey := r.StorageGroups[0]
	deploy := vcenter.Deploy{
		DeploymentSpec: vcenter.DeploymentSpec{
			Name:               createData.VmName,
			DefaultDatastoreID: datastoreID,
			AcceptAllEULA:      true,
			NetworkMappings: []vcenter.NetworkMapping{
				{
					Key:   networkKey,
					Value: networkID,
				},
			},
			StorageMappings: []vcenter.StorageMapping{{
				Key: "",
				Value: vcenter.StorageGroupMapping{
					Type:         "DATASTORE",
					DatastoreID:  datastoreID,
					Provisioning: "thin",
				},
			}},
			StorageProvisioning: "thin",
		},
		Target: vcenter.Target{
			ResourcePoolID: resourceID,
			FolderID:       folderID,
		},
	}
	ref, err := vcenter.NewManager(rc).DeployLibraryItem(vw.Ctx, item.ID, deploy)
	if err != nil {
		fmt.Println(4444444444, err)
		panic(err)
	}
	f := find.NewFinder(vw.Client.Client)
	obj, err := f.ObjectReference(vw.Ctx, *ref)
	if err != nil {
		panic(err)
	}
	_ = obj.(*object.VirtualMachine)

	//datastores, err := finder.VirtualMachineList(vw.ctx, "*/group-v629")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to list datastore at vc %v", err)
	//	os.Exit(1)
	//}
	//fmt.Println(datastores)
}

func (vw *VmWare) CloneVM() {
	cloneData := CreateMap{
		TempName:    "xxx",
		Datacenter:  "xxx",
		Cluster:     "xxx",
		Host:        "xxx",
		Resources:   "xxx",
		Storage:     "xxx",
		VmName:      "xxx",
		SysHostName: "xxx",
		Network:     "xxx",
	}
	vmList, templateList, err := vw.GetAllVmClient()
	if err != nil {
		panic(err)
	}
	var templateNameList []string
	var vmTemplate types.ManagedObjectReference
	for _, template := range templateList {
		templateNameList = append(templateNameList, template.Name)
		if template.Name == cloneData.TempName {
			vmTemplate = template.VM
		}
	}
	if !IsContain(templateNameList, cloneData.TempName) {
		fmt.Fprintf(os.Stderr, "模版不存在，虚拟机克隆失败")
		return
	}
	dataCenterList, err := vw.GetAllDatacenter()
	if err != nil {
		panic(err)
	}
	var datacenterID, datacenterName string
	for _, datacenter := range dataCenterList {
		if datacenter.Name == cloneData.Datacenter {
			datacenterID = datacenter.Datacenter.Value
			datacenterName = datacenter.Name
		}
	}
	if datacenterName == "" {
		fmt.Fprintf(os.Stderr, "数据中心不存在，虚拟机克隆失败")
		return
	}
	hostList, err := vw.GetAllHost()
	if err != nil {
		panic(err)
	}
	var hostName string
	var hostRef types.ManagedObjectReference
	for _, host := range hostList {
		if host.Name == cloneData.Host {
			hostName = host.Name
			hostRef = host.HostSelf
		}
	}
	if hostName == "" {
		fmt.Fprintf(os.Stderr, "主机不存在，虚拟机克隆失败")
		return
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	var poolRef types.ManagedObjectReference
	for _, resource := range resourceList {
		if resource.Name == cloneData.Resources {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
			poolRef = resource.Resource
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	datastoreList, err := vw.GetAllDatastore()
	if err != nil {
		panic(err)
	}
	var datastoreID, datastoreStr string
	var datastoreRef types.ManagedObjectReference
	for _, datastore := range datastoreList {
		if datastore.Name == cloneData.Storage {
			datastoreID = datastore.Datastore.Value
			datastoreStr = datastore.Name
			datastoreRef = datastore.DatastoreSelf
		}
	}
	if datastoreStr == "" {
		fmt.Fprintf(os.Stderr, "存储中心不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("DatastoreID", datastoreID)
	networkList, err := vw.GetAllNetwork()
	if err != nil {
		panic(err)
	}
	var networkID, networkStr string
	for _, network := range networkList {
		if network["Vlan"] == cloneData.Network {
			networkStr = network["Vlan"]
			networkID = network["NetworkID"]
		}
	}

	if networkStr == "" {
		fmt.Fprintf(os.Stderr, "网络不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("NetworkID", networkID)
	clusterList, err := vw.GetAllCluster()
	if err != nil {
		panic(err)
	}
	var clusterID, clusterName string
	for _, cluster := range clusterList {
		if cluster.Name == cloneData.Cluster {
			clusterID = cluster.Cluster.Value
			clusterName = cluster.Name
		}
	}
	if clusterName == "" {
		fmt.Fprintf(os.Stderr, "集群不存在，虚拟机克隆失败")
		return
	}
	configSpecs := []types.BaseVirtualDeviceConfigSpec{}
	fmt.Println("ClusterID", clusterID)
	for _, vms := range vmList {
		if vms.Name == cloneData.VmName {
			fmt.Fprintf(os.Stderr, "虚机已存在，虚拟机克隆失败")
			return
		}
	}
	finder := find.NewFinder(vw.Client.Client)
	folders, err := finder.FolderList(vw.Ctx, "*")
	var Folder *object.Folder
	for _, folder := range folders {
		if folder.InventoryPath == "/"+cloneData.Datacenter+"/vm" {
			Folder = folder
		}
	}
	fmt.Println(Folder)
	folderList, err := vw.GetFolder()
	if err != nil {
		panic(err)
	}
	var folderRef types.ManagedObjectReference
	for _, folder := range folderList {
		if folder.Parent.Value == datacenterID && folder.Name == "vm" {
			folderRef = folder.FolderSelf
		}
	}
	fmt.Println("poolRef", poolRef)
	relocateSpec := types.VirtualMachineRelocateSpec{
		DeviceChange: configSpecs,
		Folder:       &folderRef,
		Pool:         &poolRef,
		Host:         &hostRef,
		Datastore:    &datastoreRef,
	}
	vmConf := &types.VirtualMachineConfigSpec{
		NumCPUs:  4,
		MemoryMB: 16 * 1024,
	}
	cloneSpec := &types.VirtualMachineCloneSpec{
		PowerOn:  false,
		Template: false,
		Location: relocateSpec,
		Config:   vmConf,
	}
	t := object.NewVirtualMachine(vw.Client.Client, vmTemplate)
	newFolder := object.NewFolder(vw.Client.Client, folderRef)
	fmt.Println(newFolder)
	fmt.Println(cloneData.VmName)
	fmt.Println(cloneSpec.Location)
	task, err := t.Clone(vw.Ctx, newFolder, cloneData.VmName, *cloneSpec)
	if err != nil {
		panic(err)
	}
	fmt.Println("克隆任务开始，", task.Wait(vw.Ctx))
}

func (vw *VmWare) SetIP(vm *object.VirtualMachine, ipAddr vsphere.IpAddr, os string) error {
	cam := types.CustomizationAdapterMapping{
		Adapter: types.CustomizationIPSettings{
			Ip:         &types.CustomizationFixedIp{IpAddress: ipAddr.Ip},
			SubnetMask: ipAddr.Netmask,
			Gateway:    []string{ipAddr.Gateway},
			DnsDomain:  "192.168.4.111",
		},
	}
	var customSpec types.CustomizationSpec
	if os == "Linux" {
		customSpec = types.CustomizationSpec{
			NicSettingMap: []types.CustomizationAdapterMapping{cam},
			Identity:      &types.CustomizationLinuxPrep{HostName: &types.CustomizationFixedName{Name: ipAddr.Hostname}},
		}
	} else {
		logrus.Info(fmt.Sprintf("配置windows:%s", ipAddr.Ip))
		customSpec = types.CustomizationSpec{
			NicSettingMap: []types.CustomizationAdapterMapping{cam},
			Identity: &types.CustomizationSysprep{
				UserData: types.CustomizationUserData{
					FullName:     "new",
					OrgName:      "github",
					ComputerName: &types.CustomizationFixedName{Name: "new-pc"},
				},
			},
		}
	}

	task, err := vm.Customize(vw.Ctx, customSpec)
	if err != nil {
		return err
	}
	logrus.Info("设置ip完成:", task.Name())
	return task.Wait(vw.Ctx)
}

func (vw *VmWare) FindVMByIP(ctx context.Context, c *vim25.Client, ip string) (machineName string, err error) {
	searchIndex := object.NewSearchIndex(c)
	// nil 是数据中心，你如果要指定的话，需要构造数据中心的结构体，这里就不指定了
	// ip 是你虚拟机的 ip
	// true 的意思我没搞懂
	reference, findErr := searchIndex.FindByIp(ctx, nil, ip, true)
	if findErr != nil {
		logrus.Error(findErr)
		return "", findErr
	}
	// 之所以只对 reference 进行判断而非对 err 是因为没有找到不算是 error
	// 也就是说 err 为 nil 并不代表就找到了，但是没找到 reference 一定为 nil
	if reference == nil {
		logrus.Info("vm not found")
		return "", nil
	}
	// 这类的查找的对象都是 object.Reference，你需要通过对应的方法将其转换为相应的对象
	// 比如虚拟机、文件夹、模板等~
	res := object.NewVirtualMachine(c, reference.Reference())
	name, _ := res.ObjectName(ctx)
	logrus.Info(fmt.Sprintf("根据ip查找到虚拟机名称:【%s】", name))
	return name, nil
}

func (vw *VmWare) ConnectNetwork(vm *object.VirtualMachine) error {
	devidelst, fetchErr := vm.Device(vw.Ctx)
	if fetchErr != nil {
		logrus.Error("获取虚拟机的设备列表失败，" + fetchErr.Error())
		return fetchErr
	}

	for _, device := range devidelst {
		deviceSummary := device.GetVirtualDevice().DeviceInfo.GetDescription().Summary
		fmt.Println(deviceSummary)
		if strings.Contains(deviceSummary, "Network") {
			fmt.Println("连接网卡:", devidelst.Type(device))
			connectErr := devidelst.Disconnect(device)
			if connectErr != nil {
				logrus.Error("连接网卡连接失败，" + connectErr.Error())
				return connectErr
			}
		}
	}
	return nil
}
func (vw *VmWare) PowerOn(vm *object.VirtualMachine) error {
	task, err := vm.PowerOn(vw.Ctx)
	if err != nil {
		return err
	}
	err = task.Wait(vw.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (vw *VmWare) PowerOff(vm *object.VirtualMachine) error {
	task, err := vm.PowerOff(vw.Ctx)
	if err != nil {
		return err
	}
	err = task.Wait(vw.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (vw *VmWare) MigrateVM() {
	migrateData := "测试虚机"
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		panic(err)
	}
	defer v.Destroy(vw.Ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.Ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		panic(err)
	}
	var vmTarget types.ManagedObjectReference
	for _, vm := range vms {
		if vm.Summary.Config.Name == migrateData {
			vmTarget = vm.Self
		}
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	var poolRef types.ManagedObjectReference
	for _, resource := range resourceList {
		if resource.Name == "" {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
			poolRef = resource.Resource
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机迁移失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	hostList, err := vw.GetAllHost()
	if err != nil {
		panic(err)
	}
	var hostName string
	var hostRef types.ManagedObjectReference
	for _, host := range hostList {
		if host.Name == "xxxx" {
			hostName = host.Name
			hostRef = host.HostSelf
		}
	}
	if hostName == "" {
		fmt.Fprintf(os.Stderr, "主机不存在，虚拟机迁移失败")
		return
	}
	t := object.NewVirtualMachine(vw.Client.Client, vmTarget)
	pool := object.NewResourcePool(vw.Client.Client, poolRef)
	host := object.NewHostSystem(vw.Client.Client, hostRef)
	//var priority types.VirtualMachineMovePriority
	//var state types.VirtualMachinePowerState
	task, err := t.Migrate(vw.Ctx, pool, host, "defaultPriority", "poweredOff")
	if err != nil {
		panic(err)
	}
	fmt.Println("虚拟机迁移中......")
	_ = task.Wait(vw.Ctx)
	fmt.Println("虚拟机迁移完成.....")
}

// 登录认证并获取主机信息
func VerifyHostInfo(loginInfo *contants.LoginInfo, os string) (result *contants.VerifyResult, err error) {
	if os == "Linux" {
		client, connErr := linux.NewSshClient(loginInfo)
		if connErr != nil {
			err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, loginInfo.Ip, connErr.Error()))
			logrus.Error(connErr)
			return
		}
		defer client.Logout()

		var ver contants.VerifyResult
		ver.HostName = client.GetHostName()
		ver.OSType = client.GetHostOs()
		ver.Arch = client.GetHostArch()
		result = &ver
		logrus.Info(result)
		return
	} else {
		client, connErr := windows.NewWinRmClient(loginInfo)
		if connErr != nil {
			err = errors.New(msg.ERROR_VERIFY_HOST, msg.GetMsg(msg.ERROR_VERIFY_HOST, loginInfo.Ip, connErr.Error()))
			logrus.Error(connErr)
			return
		}

		var ver contants.VerifyResult
		ver.HostName, ver.OSType, ver.Arch, err = client.GetSystemInfo()
		result = &ver
		logrus.Info(result)
		return
	}

}

func SaveHost(vmUsername string, vmPassword string, name string, os string, ip string, port uint, relicaName string, username string, vms *VmWare) (err error) {
	// 应用认证处理
	loginInfo := contants.LoginInfo{
		Ip:       ip,
		Port:     port,
		UserName: vmUsername,
		PassWord: vmPassword,
	}

	var arch = ""
	var osName = ""
	var hostName = ""
	var status = 0
	var platformId = 0

	osType := contants.LINUX
	if os == "Windows" {
		osType = contants.WINDOWS
	}

	response, verifyErr := VerifyHostInfo(&loginInfo, os)
	if verifyErr != nil {
		err = errors.New(msg.ERROR_CREATE_HOST, msg.GetMsg(msg.ERROR_CREATE_HOST, verifyErr.Error()))
		logrus.Error("================校验主机发生错误================", verifyErr)

	} else {
		status = 1
		arch = response.Arch
		osName = response.OSType
		hostName = response.HostName
	}

	platformRes, findPlatformErr := platform.PlatformMgm.FindByIp(vms.IP)
	if findPlatformErr != nil {
		logrus.Error(fmt.Sprintf("查询虚拟化平台信息报错:%s", findPlatformErr))
	} else {
		platformId = int(platformRes.Id)
	}

	if host.HostMgm.CheckIsExistByIp(ip) {
		logrus.Info(fmt.Sprintf("ip:%s已存在，更新主机信息", ip))
		obj, findErr := host.HostMgm.FindByIp(ip)
		if findErr != nil {
			logrus.Error(fmt.Sprintf("查询主机信息报错%s", findErr))
		} else {
			ctx := context.Background()
			machineName, findVmErr := vms.FindVMByIP(ctx, vms.Client.Client, ip)
			if findVmErr != nil {
				logrus.Error(findVmErr)
			} else {
				obj.Name = machineName
			}
			obj.HostName = hostName
			obj.Status = uint(status)
			obj.PlatformId = uint(platformId)
			_, updErr := host.HostMgm.UpdateHost(obj)
			if updErr != nil {
				logrus.Error("================更新主机信息发生错误================", updErr)
			}
		}
		return
	}

	// 密码入库加密
	password, err := util.AesEncrypt(loginInfo.PassWord)
	if err != nil {
		err = errors.New(msg.ERROR_UPDATE_HOST, msg.GetMsg(msg.ERROR_UPDATE_HOST, err.Error()))
		logrus.Error(err)
		return
	}

	var obj host.Host
	obj.Name = name
	obj.Desc = fmt.Sprintf("由副本：%s 挂载创建虚拟机", relicaName)
	obj.HostName = hostName
	obj.Type = osType
	obj.Status = uint(status)
	obj.Os = osName
	obj.Arch = arch
	obj.Ip = ip
	obj.Port = port
	obj.PlatformId = uint(platformId)
	obj.UserName = loginInfo.UserName
	obj.PassWord = password
	obj.CreateUser = username
	obj.AuthType = contants.PRIVATE

	_, crtErr := host.HostMgm.CreateHost(obj)
	if crtErr != nil {
		err = errors.New(msg.ERROR_CREATE_HOST, msg.GetMsg(msg.ERROR_CREATE_HOST, crtErr.Error()))
		logrus.Error(crtErr)
		return
	}
	return nil
}
