//go:build !linux
// +build !linux

package ceph

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/sirupsen/logrus"
)

type client struct{}

var Client = client{}

func (c *client) MonCommand(command []byte) (buf []byte, err error) {
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		// 伪造测试数据
		buf = []byte("{\"fsid\":\"48b88c0f-53e3-4b1c-af86-0f1018e8cc4a\",\"health\":{\"checks\":" +
			"{\"POOL_APP_NOT_ENABLED\":{\"severity\":\"HEALTH_WARN\",\"summary\":{\"message\":\"application not " +
			"enabled on 2 pool(s)\"}}},\"status\":\"HEALTH_WARN\",\"overall_status\":\"HEALTH_WARN\"}," +
			"\"election_epoch\":40,\"quorum\":[0,1,2],\"quorum_names\":[\"ceph136\",\"ceph137\",\"ceph138\"]," +
			"\"monmap\":{\"epoch\":3,\"fsid\":\"48b88c0f-53e3-4b1c-af86-0f1018e8cc4a\",\"modified\":" +
			"\"2021-09-22 16:32:27.000021\",\"created\":\"2021-09-22 16:31:53.418042\",\"features\":{\"persistent\":" +
			"[\"kraken\",\"luminous\",\"mimic\",\"osdmap-prune\"],\"optional\":[]},\"mons\":[{\"rank\":0,\"name\":" +
			"\"ceph136\",\"addr\":\"192.168.212.136:6789/0\",\"public_addr\":\"192.168.212.136:6789/0\"},{\"rank\":1," +
			"\"name\":\"ceph137\",\"addr\":\"192.168.212.137:6789/0\",\"public_addr\":\"192.168.212.137:6789/0\"}," +
			"{\"rank\":2,\"name\":\"ceph138\",\"addr\":\"192.168.212.138:6789/0\",\"public_addr\":" +
			"\"192.168.212.138:6789/0\"}]},\"osdmap\":{\"osdmap\":{\"epoch\":513,\"num_osds\":3,\"num_up_osds\":3," +
			"\"num_in_osds\":3,\"num_remapped_pgs\":0}},\"pgmap\":{\"pgs_by_state\":[{\"state_name\":\"active+clean\"," +
			"\"count\":182}],\"num_pgs\":182,\"num_pools\":6,\"num_objects\":10946,\"data_bytes\":44128692865," +
			"\"bytes_used\":135704936448,\"bytes_avail\":1045398487040,\"bytes_total\":1181103423488},\"fsmap\":" +
			"{\"epoch\":1,\"by_rank\":[]},\"mgrmap\":{\"epoch\":20,\"active_gid\":34107,\"active_name\":\"ceph136\"," +
			"\"active_addr\":\"192.168.212.136:6804/1793\",\"available\":true,\"standbys\":[{\"gid\":34122,\"name\":" +
			"\"ceph138\",\"available_modules\":[{\"name\":\"balancer\",\"can_run\":true,\"error_string\":\"\"}," +
			"{\"name\":\"crash\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"dashboard\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"hello\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"influx\"," +
			"\"can_run\":false,\"error_string\":\"influxdb python module not found\"},{\"name\":\"iostat\"," +
			"\"can_run\":true,\"error_string\":\"\"},{\"name\":\"localpool\",\"can_run\":true,\"error_string\":\"\"}," +
			"{\"name\":\"prometheus\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"restful\",\"can_run\":" +
			"true,\"error_string\":\"\"},{\"name\":\"selftest\",\"can_run\":true,\"error_string\":\"\"},{\"name\":" +
			"\"smart\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"status\",\"can_run\":true,\"error_string\":" +
			"\"\"},{\"name\":\"telegraf\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"telemetry\"," +
			"\"can_run\":true,\"error_string\":\"\"},{\"name\":\"zabbix\",\"can_run\":true,\"error_string\":\"\"}]}," +
			"{\"gid\":34140,\"name\":\"ceph137\",\"available_modules\":[{\"name\":\"balancer\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"crash\",\"can_run\":true,\"error_string\":\"\"},{\"name\":" +
			"\"dashboard\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"hello\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"influx\",\"can_run\":false,\"error_string\":\"influxdb " +
			"python module not found\"},{\"name\":\"iostat\",\"can_run\":true,\"error_string\":\"\"},{\"name\":" +
			"\"localpool\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"prometheus\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"restful\",\"can_run\":true,\"error_string\":\"\"},{\"name\":" +
			"\"selftest\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"smart\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"status\",\"can_run\":true,\"error_string\":\"\"},{\"name\":" +
			"\"telegraf\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"telemetry\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"zabbix\",\"can_run\":true,\"error_string\":\"\"}]}],\"modules\":" +
			"[\"balancer\",\"crash\",\"iostat\",\"restful\",\"status\"],\"available_modules\":[{\"name\":\"balancer\"," +
			"\"can_run\":true,\"error_string\":\"\"},{\"name\":\"crash\",\"can_run\":true,\"error_string\":\"\"}," +
			"{\"name\":\"dashboard\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"hello\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"influx\",\"can_run\":false,\"error_string\":\"influxdb python module " +
			"not found\"},{\"name\":\"iostat\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"localpool\"," +
			"\"can_run\":true,\"error_string\":\"\"},{\"name\":\"prometheus\",\"can_run\":true,\"error_string\":\"\"}," +
			"{\"name\":\"restful\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"selftest\",\"can_run\":true," +
			"\"error_string\":\"\"},{\"name\":\"smart\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"status\"," +
			"\"can_run\":true,\"error_string\":\"\"},{\"name\":\"telegraf\",\"can_run\":true,\"error_string\":\"\"}," +
			"{\"name\":\"telemetry\",\"can_run\":true,\"error_string\":\"\"},{\"name\":\"zabbix\",\"can_run\":true," +
			"\"error_string\":\"\"}],\"services\":{}},\"servicemap\":{\"epoch\":15,\"modified\":" +
			"\"2021-10-19 21:43:42.955664\",\"services\":{\"rgw\":{\"daemons\":{\"summary\":\"\",\"ceph136\":" +
			"{\"start_epoch\":12,\"start_stamp\":\"2021-10-19 13:56:48.112669\",\"gid\":34167,\"addr\":" +
			"\"192.168.212.136:0/3198213167\",\"metadata\":{\"arch\":\"x86_64\",\"ceph_release\":\"mimic\"," +
			"\"ceph_version\":\"ceph version 13.2.10 (564bdc4ae87418a232fc901524470e1a0f76d641) mimic (stable)\"," +
			"\"ceph_version_short\":\"13.2.10\",\"cpu\":\"Intel(R) Xeon(R) CPU E5-2650 v3 @ 2.30GHz\",\"distro\":" +
			"\"centos\",\"distro_description\":\"CentOS Linux 7 (Core)\",\"distro_version\":\"7\"," +
			"\"frontend_config#0\":\"civetweb port=7480\",\"frontend_type#0\":\"civetweb\",\"hostname\":\"ceph136\"," +
			"\"kernel_description\":\"#1 SMP Thu Sep 24 19:07:15 CST 2020\",\"kernel_version\":" +
			"\"5.4.40-2009.26.el7a.x86_64\",\"mem_swap_kb\":\"8257532\",\"mem_total_kb\":\"16231892\",\"num_handles\":" +
			"\"1\",\"os\":\"Linux\",\"pid\":\"2265\",\"zone_id\":\"d957a349-27e9-4ce3-a318-49a56ec589a8\"," +
			"\"zone_name\":\"default\",\"zonegroup_id\":\"29912f02-c5b9-4d40-ba50-ad0f73deccce\",\"zonegroup_name\":" +
			"\"default\"}},\"ceph137\":{\"start_epoch\":14,\"start_stamp\":\"2021-10-19 21:39:15.519529\",\"gid\":" +
			"34146,\"addr\":\"192.168.212.137:0/3384746022\",\"metadata\":{\"arch\":\"x86_64\",\"ceph_release\":" +
			"\"mimic\",\"ceph_version\":\"ceph version 13.2.10 (564bdc4ae87418a232fc901524470e1a0f76d641) mimic " +
			"(stable)\",\"ceph_version_short\":\"13.2.10\",\"cpu\":\"Intel(R) Xeon(R) CPU E5-2650 v3 @ 2.30GHz\"," +
			"\"distro\":\"centos\",\"distro_description\":\"CentOS Linux 7 (Core)\",\"distro_version\":\"7\"," +
			"\"frontend_config#0\":\"civetweb port=7480\",\"frontend_type#0\":\"civetweb\",\"hostname\":\"ceph137\"," +
			"\"kernel_description\":\"#1 SMP Fri Apr 20 16:44:24 UTC 2018\",\"kernel_version\":" +
			"\"3.10.0-862.el7.x86_64\",\"mem_swap_kb\":\"8257532\",\"mem_total_kb\":\"16266744\",\"num_handles\":" +
			"\"1\",\"os\":\"Linux\",\"pid\":\"1313\",\"zone_id\":\"d957a349-27e9-4ce3-a318-49a56ec589a8\"," +
			"\"zone_name\":\"default\",\"zonegroup_id\":\"29912f02-c5b9-4d40-ba50-ad0f73deccce\",\"zonegroup_name\":" +
			"\"default\"}},\"ceph138\":{\"start_epoch\":15,\"start_stamp\":\"2021-10-19 21:43:42.675108\",\"gid\":" +
			"34176,\"addr\":\"192.168.212.138:0/1960341377\",\"metadata\":{\"arch\":\"x86_64\",\"ceph_release\":" +
			"\"mimic\",\"ceph_version\":\"ceph version 13.2.10 (564bdc4ae87418a232fc901524470e1a0f76d641) mimic " +
			"(stable)\",\"ceph_version_short\":\"13.2.10\",\"cpu\":\"Intel(R) Xeon(R) CPU E5-2650 v3 @ 2.30GHz\"," +
			"\"distro\":\"centos\",\"distro_description\":\"CentOS Linux 7 (Core)\",\"distro_version\":\"7\"," +
			"\"frontend_config#0\":\"civetweb port=7480\",\"frontend_type#0\":\"civetweb\",\"hostname\":\"ceph138\"," +
			"\"kernel_description\":\"#1 SMP Fri May 29 14:10:43 CST 2020\",\"kernel_version\":" +
			"\"5.4.40-2005.11.el7a.x86_64\",\"mem_swap_kb\":\"8257532\",\"mem_total_kb\":\"16231900\"," +
			"\"num_handles\":\"1\",\"os\":\"Linux\",\"pid\":\"2461\",\"zone_id\":" +
			"\"d957a349-27e9-4ce3-a318-49a56ec589a8\",\"zone_name\":\"default\",\"zonegroup_id\":" +
			"\"29912f02-c5b9-4d40-ba50-ad0f73deccce\",\"zonegroup_name\":\"default\"}}}}}}}")
		return
	}
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	return
}

// 获取集群信息
func (c *client) GetClusterInfo() (info *ClusterInfo, err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		// 伪造测试数据
		TB := uint64(1024 * 1024 * 1024 * 1024)
		info := ClusterInfo{TotalSize: TB, UsedSize: 0, AvailSize: TB, TotalObject: 0}
		return &info, nil
	}
	return
}

// 创建存储池
func (c *client) CreatePool(poolName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 删除存储池
func (c *client) DeletePool(poolName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 创建image
func (c *client) CreateImage(poolName, imageName string, size uint64) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 删除image
func (c *client) DeleteImage(poolName, imageName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 复制image
func (c *client) CopyImage(poolName, imageName, destName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 从快照克隆image[克隆前需要检查快照是否受保护]
func (c *client) CloneImageBySnapshot(poolName, imageName, snapName, destName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 分离image快照依赖
func (c *client) FlattenImage(poolName, imageName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 重置image大小
func (c *client) ReSizeImage(poolName, imageName string, size uint64) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 重命名镜像
func (c *client) ReNameImage(poolName, imageName, destImageName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 刷新缓存数据到镜像
func (c *client) FlushImage(poolName, imageName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	return
}

// 创建image快照
func (c *client) CreateSnapshot(poolName, imageName, snapName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 删除image快照
func (c *client) DeleteSnapshot(poolName, imageName, snapName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 回滚image快照
func (c *client) RollbackSnapshot(poolName, imageName, snapName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 保护快照
func (c *client) ProtectSnapShot(poolName, imageName, snapName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

// 解除快照保护
func (c *client) UnProtectSnapShot(poolName, imageName, snapName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

func (c *client) MapRBDImage(poolName, imageName string) (devPath string, err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		// 伪造测试数据
		devPath = "/dev/rbd0"
		return devPath, nil
	}
	return
}

func (c *client) UnMapRBDImage(poolName, imageName string) (err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		return nil
	}
	return
}

func (c *client) ShowMapRBDImage() (mapInfos []MapInfo, err error) {
	err = errors.New(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS, msg.GetMsg(msg.ERROR_CEPH_NON_SUPPORT_WINDOWS))
	logrus.Error(err)
	// Windows下调试模式不抛错
	if contants.AppCfg.System.RunMode == contants.DEBUG {
		// 伪造测试数据
		mapInfos = append(mapInfos, MapInfo{PoolName: "test_pool", ImageName: "test_image", DevPath: "/dev/rbd0"})
		return mapInfos, nil
	}
	return
}
