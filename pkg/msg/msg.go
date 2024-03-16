package msg

import (
	"fmt"
)

// 预留中英文资源文件
var OperationFlags = map[Operation]string{
	UNDEFINE_OPERAT:                    "未定义操作",
	LOGIN:                              "用户登录",
	LOGIN_SUCCESS:                      "用户(%s)登录成功!",
	LOGIN_FAILED:                       "用户(%s)登录失败[%s]!",
	LOGOUT:                             "用户注销",
	LOGOUT_SUCCESS:                     "用户(%s)注销成功!",
	LOGOUT_FAILED:                      "用户(%s)注销失败[%s]!",
	ENABLE_USER:                        "启用用户",
	ENABLE_USER_SUCCESS:                "启用用户(%s)成功!",
	ENABLE_USER_FAILED:                 "启用用户(%s)失败[%s]!",
	DISABLE_USER:                       "禁用用户",
	DISABLE_USER_SUCCESS:               "禁用用户(%s)成功!",
	DISABLE_USER_FAILED:                "禁用用户(%s)失败[%s]!",
	EDIT_USER_ROLE:                     "更改用户角色",
	EDIT_USER_ROLE_SUCCESS:             "更改用户角色为(%s)成功!",
	EDIT_USER_ROLE_FAILED:              "更改用户角色为(%s)失败[%s]!",
	CREATE_ROLE:                        "添加角色",
	CREATE_ROLE_SUCCESS:                "添加角色(%s)成功!",
	CREATE_ROLE_FAILED:                 "添加角色(%s)失败[%s]",
	EDIT_ROLE:                          "编辑角色",
	EDIT_ROLE_SUCCESS:                  "编辑角色(%s)成功!",
	EDIT_ROLE_FAILED:                   "编辑角色(%s)失败[%s]!",
	DELETE_ROLE:                        "删除角色",
	DELETE_ROLE_SUCCESS:                "删除角色(%s)成功!",
	DELETE_ROLE_FAILED:                 "删除角色(%s)失败[%s]!",
	CREATE_PLATFORM:                    "添加平台",
	CREATE_PLATFORM_SUCCESS:            "添加平台(%s)成功!",
	CREATE_PLATFORM_FAILED:             "添加平台(%s)失败[%s]!",
	EDIT_PLATFORM:                      "编辑平台",
	EDIT_PLATFORM_SUCCESS:              "编辑平台(%s)成功!",
	EDIT_PLATFORM_FAILED:               "编辑平台(%s)失败[%s]!",
	PUBLISH_PLATFORM:                   "发布平台",
	PUBLISH_PLATFORM_SUCCESS:           "发布平台(%s)成功!",
	PUBLISH_PLATFORM_FAILED:            "发布平台(%s)失败[%s]!",
	DELETE_PLATFORM:                    "删除平台",
	DELETE_PLATFORM_SUCCESS:            "删除平台(%s)成功!",
	DELETE_PLATFORM_FAILED:             "删除平台(%s)失败[%s]!",
	CREATE_HOST:                        "添加主机",
	CREATE_HOST_SUCCESS:                "添加主机(%s)成功!",
	CREATE_HOST_FAILED:                 "添加主机(%s)失败[%s]!",
	EDIT_HOST:                          "编辑主机",
	EDIT_HOST_SUCCESS:                  "编辑主机(%s)成功!",
	EDIT_HOST_FAILED:                   "编辑主机(%s)失败[%s]!",
	PUBLISH_HOST:                       "发布主机",
	PUBLISH_HOST_SUCCESS:               "发布主机(%s)成功!",
	PUBLISH_HOST_FAILED:                "发布主机(%s)失败[%s]!",
	DELETE_HOST:                        "删除主机",
	DELETE_HOST_SUCCESS:                "删除主机(%s)成功!",
	DELETE_HOST_FAILED:                 "删除主机(%s)失败[%s]!",
	UPD_HOST:                           "更新主机",
	UPD_HOST_SUCCESS:                   "更新主机(%s)成功",
	UPD_HOST_FAILED:                    "更新主机(%s)失败",
	CREATE_POOL:                        "创建存储池",
	CREATE_POOL_SUCCESS:                "创建存储(%s)成功!",
	CREATE_POOL_FAILED:                 "创建存储(%s)失败[%s]!",
	DELETE_POOL:                        "删除存储池",
	DELETE_POOL_SUCCESS:                "删除存储(%s)成功!",
	DELETE_POOL_FAILED:                 "删除存储(%s)失败[%s]!",
	CREATE_IMAGE_BY_REPLICA:            "由副本生成镜像",
	CREATE_IMAGE_BY_REPLICA_SUCCESS:    "由副本(%s)生成镜像(%s)成功!",
	CREATE_IMAGE_BY_REPLICA_FAILED:     "由副本(%s)生成镜像(%s)失败[%s]!",
	CREATE_IMAGE_BY_SNAPSHOT:           "由快照生成镜像",
	CREATE_IMAGE_BY_SNAPSHOT_SUCCESS:   "由快照(%s)生成镜像(%s)成功!",
	CREATE_IMAGE_BY_SNAPSHOT_FAILED:    "由快照(%s)生成镜像(%s)失败[%s]!",
	PUBLISH_IMAGE:                      "发布镜像",
	PUBLISH_IMAGE_SUCCESS:              "发布镜像(%s)成功!",
	PUBLISH_IMAGE_FAILED:               "发布镜像(%s)失败[%s]!",
	DELETE_IMAGE:                       "销毁镜像",
	DELETE_IMAGE_SUCCESS:               "销毁镜像(%s)成功!",
	DELETE_IMAGE_FAILED:                "销毁镜像(%s)失败[%s]!",
	CREATE_REPLICA:                     "创建新副本",
	CREATE_REPLICA_SUCCESS:             "创建新副本(%s)成功!",
	CREATE_REPLICA_FAILED:              "创建新副本(%s)失败[%s]!",
	CREATE_REPLICA_BY_IMAGE:            "由镜像生成副本",
	CREATE_REPLICA_BY_IMAGE_BEGIN:      "开始镜像(%s)生成副本(%s)",
	CREATE_REPLICA_BY_IMAGE_SUCCESS:    "由镜像(%s)生成副本(%s)成功!",
	CREATE_REPLICA_BY_IMAGE_FAILED:     "由镜像(%s)生成副本(%s)失败[%s]!",
	CREATE_REPLICA_BY_SNAPSHOT:         "由快照生成副本",
	CREATE_REPLICA_BY_SNAPSHOT_SUCCESS: "由快照(%s)生成副本(%s)成功!",
	CREATE_REPLICA_BY_SNAPSHOT_FAILED:  "由快照(%s)生成副本(%s)失败[%s]!",
	DELETE_REPLICA:                     "销毁副本",
	DELETE_REPLICA_SUCCESS:             "销毁副本(%s)成功!",
	DELETE_REPLICA_FAILED:              "销毁副本(%s)失败[%s]!",
	MOUNT_REPLICA:                      "挂载副本",
	MOUNT_REPLICA_SUCCESS:              "挂载副本(%s)成功!",
	MOUNT_REPLICA_FAILED:               "挂载副本(%s)失败[%s]!",
	UNMOUNT_REPLICA:                    "卸载副本",
	UNMOUNT_REPLICA_SUCCESS:            "卸载副本(%s)成功!",
	UNMOUNT_REPLICA_FAILED:             "卸载副本(%s)失败[%s]!",
	CREATE_SNAPSHOT:                    "创建副本快照",
	CREATE_SNAPSHOT_SUCCESS:            "创建副本(%s)的快照(%s)成功!",
	CREATE_SNAPSHOT_FAILED:             "创建副本(%s)的快照(%s)失败[%s]!",
	ROLLBACK_SNAPSHOT:                  "回滚快照",
	ROLLBACK_SNAPSHOT_SUCCESS:          "回滚快照(%s)成功!",
	ROLLBACK_SNAPSHOT_FAILED:           "回滚快照(%s)失败[%s]!",
	DELETE_SNAPSHOT:                    "删除快照",
	DELETE_SNAPSHOT_SUCCESS:            "删除快照(%s)成功!",
	DELETE_SNAPSHOT_FAILED:             "删除快照(%s)失败[%s]!",
	DELETE_SNAPSHOT_BY_REPLICA:         "删除副本所有快照",
	DELETE_SNAPSHOT_BY_REPLICA_SUCCESS: "删除副本(%s)快照(%s)成功!",
	DELETE_SNAPSHOT_BY_REPLICA_FAILED:  "删除副本(%s)快照(%s)失败[%s]!",
	UPLOAD_SCRIPT_FILE:                 "上传脚本",
	UPLOAD_SCRIPT_FILE_SUCCESS:         "上传脚本(%s)成功!",
	UPLOAD_SCRIPT_FILE_FAILED:          "上传脚本(%s)失败[%s]!",
	DOWNLOAD_SCRIPT_FILE:               "下载脚本",
	DOWNLOAD_SCRIPT_FILE_SUCCESS:       "下载脚本(%s)成功!",
	DOWNLOAD_SCRIPT_FILE_FAILED:        "下载脚本(%s)失败[%s]!",
	DELETE_SCRIPT_FILE:                 "删除脚本",
	DELETE_SCRIPT_FILE_SUCCESS:         "删除脚本(%s)成功!",
	DELETE_SCRIPT_FILE_FAILED:          "删除脚本(%s)失败[%s]!",
	DEPLOY_CLIENT:                      "部署客户端",
	DEPLOY_CLIENT_FAILED:               "%s部署客户端失败:%s",
	DEPLOY_CLIENT_START:                "开始部署客户端:%s",
	DEPLOY_CLIENT_SUCCESS:              "部署客户端成功",
}

var MassageFlags = map[Massage]string{
	FAILURE:       "failure",
	SUCCESS:       "success",
	UNDEFINED:     "未知错误:%s",
	URL_NOT_FOUND: "请求资源不存在!",

	ERROR:                "程序发生错误:%s",
	ERROR_INVALID_PARAMS: "请求参数错误:%s",
	ERROR_GET_USER_INFO:  "获取用户信息发生错误:%s",
	ERROR_USER_IS_EXIST:  "该用户已存在",
	ERROR_USER_NOT_EXIST: "用户不存在",
	ERROR_SET_USER_ROLE:  "设置用户角色发生错误:%s",
	ERROR_ENABLE_USER:    "启用用户发生错误:%s",
	ERROR_DISABLE_USER:   "禁用用户发生错误:%s",
	ERROR_GET_ROLE_INFO:  "获取角色信息发生错误:%s",
	ERROR_ROLE_IS_EXIST:  "该角色已存在",

	AUTH_PERMISSION_DENIED:   "操作权限不足",
	AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	AUTH_TOKEN:               "Token生成失败",
	AUTH_ERROR:               "Token错误",

	ERROR_STORE_POOL_NAME_IS_EXIST:        "存储池名称已存在",
	ERROR_CREATE_STORE_POOL:               "创建存储池发生错误:%s",
	ERROR_GET_STORE_POOL_INFO:             "获取存储池信息发生错误:%s",
	ERROR_DELETE_STORE_POOL:               "删除存储池发生错误:%s",
	ERROR_GET_IMAGE_INFO:                  "获取镜像信息发生错误:%s",
	ERROR_DELETE_IMAGE:                    "删除镜像发生错误:%s",
	ERROR_PUBLISH_IMAGE:                   "发布镜像发生错误:%s",
	ERROR_REPLICA_NAME_IS_EXIST:           "副本名称已存在!",
	ERROR_IMAGE_TO_REPLICA:                "由镜像生成副本发生错误:%s",
	ERROR_CREATE_REPLICA:                  "创建副本发生错误:%s",
	ERROR_STORE_INSUFFICIENT_SPACE:        "存储空间不足，可用存储大小为:%ld",
	ERROR_GET_REPLICA_INFO:                "获取副本信息发生错误:%s",
	ERROR_DELETE_REPLICA:                  "删除副本发生错误:%s",
	ERROR_IMAGE_NAME_IS_EXIST:             "镜像名称已存在!",
	ERROR_REPLICA_TO_IMAGE:                "由副本生成镜像发生错误:%s",
	ERROR_SNAPSHOT_NAME_IS_EXIST:          "快照名称已存在!",
	ERROR_CREATE_SNAPSHOT:                 "创建副本快照发生错误:%s",
	ERROR_GET_SNAPSHOT_INFO:               "获取副本快照信息发生错误:%s",
	ERROR_DELETE_SNAPSHOT:                 "删除副本快照发生错误:%s",
	ERROR_ROLL_BACK_SNAPSHOT:              "回滚副本快照发生错误:%s",
	ERROR_REPLICAT_STATUS:                 "副本未处于空闲状态，不允许修改大小和类型",
	ERROR_SNAPSHOT_TO_IMAGE:               "由副本快照生成镜像发生错误:%s",
	ERROR_SNAPSHOT_TO_REPLICA:             "由副本快照生成副本发生错误:%s",
	ERROR_EXIST_MOUNT_STATUS_REPLICA:      "副本(%s)处于挂载状态，不允许执行该操作!",
	ERROR_EDIT_REPLICA:                    "编辑副本发生错误:%s",
	ERROR_EDIT_IMAGE:                      "编辑镜像发生错误:%s",
	ERROR_REPLICA_CAPACITY_EXPANSION_ONLY: "副本只允许扩容，不允许减容!",

	ERROR_PLATFORM_IS_EXIST:                   "该虚拟化(云)平台已存在",
	ERROR_PLATFORM_NOT_EXIST:                  "虚拟化(云)平台不存在",
	ERROR_CREATE_PLATFORM:                     "创建虚拟化(云)平台发生错误:%s",
	ERROR_VERIFY_PLATFORM:                     "虚拟化(云)平台(%s)认证失败:%s",
	ERROR_UPDATE_PLATFORM:                     "编辑虚拟化(云)平台发生错误:%s",
	ERROR_DELETE_PLATFORM:                     "删除虚拟化(云)平台发生错误:%s",
	ERROR_PUBLISH_PLATFORM:                    "发布虚拟化(云)平台发生错误:%s",
	ERROR_PLATFORM_NAME_IS_EXIST:              "虚拟化(云)平台名称已存在!",
	ERROR_PLATFORM_IP_IS_EXIST:                "虚拟化(云)平台IP已存在!",
	ERROR_GET_PLATFORM_INFO:                   "获取虚拟化(云)平台信息发生错误:%s",
	ERROR_PLATFORM_TYPE_NO_MATCH:              "平台类型不匹配!",
	ERROR_GET_VMWARE_DATASOURCE_FAILED:        "获取VMware平台数据源发生错误:%s",
	ERROR_SHOW_TYPE_NON_SUPPORT:               "浏览方式不支持!",
	ERROR_PLATFORM_EXIST_MOUNT_STATUS_REPLICA: "平台上存在处于挂载状态的副本，不允许执行该操作!",
	ERROR_GET_VMWARE_HOSTS_FAILED:             "获取VMware主机失败:%s",

	ERROR_GET_HOST_INFO:         "获取主机信息发生错误:%s",
	ERROR_CREATE_HOST:           "创建主机发生错误:%s",
	ERROR_VERIFY_HOST:           "主机(%s)认证失败:%s",
	ERROR_UPDATE_HOST:           "编辑主机发生错误:%s",
	ERROR_DELETE_HOST:           "删除主机发生错误:%s",
	ERROR_PUBLISH_HOST:          "发布主机发生错误:%s",
	ERROR_HOST_NAME_IS_EXIST:    "主机名称已存在!",
	ERROR_HOST_IP_IS_EXIST:      "主机IP已存在!",
	ERROR_UPD_HOST_INFO:         "更新主机信息失败",
	ERROR_HOST_DEPLOYING:        "%s正在部署中,请勿重复部署",
	ERROR_HOST_START_JENKINSJOB: "请求执行部署任务失败",
	ERROR_GET_DEPLOYAPP:         "%s 部署记录不存在",
	ERROR_UPD_DEPLOYAPP:         "%s 更新部署记录失败",

	ERROR_GET_LOG_RECORD_INFO:    "获取操作记录发生错误:%s",
	ERROR_GET_LOG_DETAIL_INFO:    "获取日志详情发生错误:%s",
	ERROR_UPDATE_LOG_RECORD_INFO: "更新操作记录状态发生错误:%s",
	ERROR_CREATE_LOG_RECORD_INFO: "创建操作记录发生错误:%s",
	ERROR_CREATE_LOG_DETAIL_INFO: "创建日志详情发生错误:%s",

	ERROR_CEPH_NOT_CONNECT:               "操作ceph异常，未连接ceph集群",
	ERROR_CEPH_POOL_NOT_OPEN:             "操作ceph异常，打开存储池发生错误:%s",
	ERROR_CEPH_IMAGE_NOT_OPEN:            "操作ceph异常，打开镜像发生错误:%s",
	ERROR_CEPH_GET_CLUSTER_STATUS:        "操作ceph异常，获取ceph集群信息发生错误:%s",
	ERROR_CEPH_CREATE_STORE_POOL:         "操作ceph异常，创建存储池发生错误:%s",
	ERROR_CEPH_DELETE_STORE_POOL:         "操作ceph异常，删除存储池发生错误:%s",
	ERROR_CEPH_CREATE_IMAGE:              "操作ceph异常，创建镜像发生错误:%s",
	ERROR_CEPH_DELETE_IMAGE:              "操作ceph异常，删除镜像发生错误:%s",
	ERROR_CEPH_CLONE_IMAGE:               "操作ceph异常，克隆镜像发生错误:%s",
	ERROR_CEPH_RESIZE_IMAGE:              "操作ceph异常，重设镜像大小发生错误:%s",
	ERROR_CEPH_RENAME_IMAGE:              "操作ceph异常，修改镜像名称发生错误:%s",
	ERROR_CEPH_FLUSH_IMAGE:               "操作ceph异常，刷新镜像缓存发生错误:%s",
	ERROR_CEPH_CREATE_IMAGE_SNAPSHOT:     "操作ceph异常，创建镜像快照发生错误:%s",
	ERROR_CEPH_DELETE_IMAGE_SNAPSHOT:     "操作ceph异常，删除镜像快照发生错误:%s",
	ERROR_CEPH_PROTECT_IMAGE_SNAPSHOT:    "操作ceph异常，保护镜像快照发生错误:%s",
	ERROR_CEPH_UN_PROTECT_IMAGE_SNAPSHOT: "操作ceph异常，解除镜像快照保护发生错误:%s",
	ERROR_CEPH_GET_IMAGE_FEATURES:        "操作ceph异常，获取镜像特性发生错误:%s",
	ERROR_CEPH_COPY_IMAGE:                "操作ceph异常，复制镜像发生错误:%s",
	ERROR_CEPH_ROLLBACK_IMAGE_SNAPSHOT:   "操作ceph异常，回滚镜像快照发生错误:%s",
	ERROR_CEPH_FLATTEN_IMAGE:             "操作ceph异常，分离镜像快照发生错误:%s",
	ERROR_CEPH_NON_SUPPORT_WINDOWS:       "操作ceph异常，不支持Windows平台！",
	ERROR_CEPH_MAP_IMAGE:                 "操作ceph异常，创建镜像映射发生错误:%s",
	ERROR_CEPH_GET_MAP_INFO:              "操作ceph异常，获取镜像映射信息发生错误:%s",
	ERROR_CEPH_UN_MAP_IMAGE:              "操作ceph异常，取消镜像映射发生错误:%s",

	ERROR_RUN_SYSTEM_COMMAND:           "调用系统命令发生错误:%s",
	ERROR_EXPECTED_RESULT_IS_INCORRECT: "预期结果不匹配!",
	ERROR_DEV_NOT_EXISTS:               "设备(%s)不存在!",
	ERROR_FILE_NOT_EXISTS:              "文件(%s)不存在!",
	ERROR_DIR_NOT_EXISTS:               "目录(%s)不存在!",
	ERROR_CREATE_DIR:                   "创建目录(%s)发生错误:%s",
	ERROR_MOUNT_POINT_OCCUPY:           "挂载目录(%s)处于占用中!",
	ERROR_MOUNT_FILED:                  "挂载文件系统发生错误:%s",
	ERROR_UNMOUNT_FILED:                "卸载文件系统发生错误:%s",
	ERROR_DIR_IS_MOUNT_POINT:           "目录(%s)是一个挂载点!",
	ERROR_DIR_NOT_MOUNT_POINT:          "目录(%s)不是一个挂载点!",
	ERROR_GET_SERVER_STATUS:            "获取服务(%s)的状态发生错误:%s",
	ERROR_START_SERVER:                 "启动服务(%s)发生错误:%s",
	ERROR_RESTART_SERVER:               "重启服务(%s)发生错误:%s",
	ERROR_STOP_SERVER:                  "停止服务(%s)发生错误:%s",
	ERROR_RELOAD_SERVER:                "重新加载服务(%s)发生错误:%s",
	ERROR_ENABLE_SERVER:                "设置服务(%s)开机自动启动发生错误:%s",
	ERROR_DISABLE_SERVER:               "禁止服务(%s)开机自动启动发生错误:%s",
	ERROR_ADD_NFS_SHARE:                "添加NFS共享发生错误:%s",
	ERROR_REMOVE_NFS_SHARE:             "取消NFS共享发生错误:%s",
	ERROR_SERVER_NOT_START:             "服务(%s)未启动!",
	ERROR_MOUNT_POINT_NOT_INCORRECT:    "检测到设备(%s)挂载到非预期挂载点(%s)上!",
	ERROR_GET_LOCAL_ADDR:               "获取本机IP发生错误:%s",
	ERROR_RUN_SHELL_SCRIPT:             "运行shell脚本发生错误:%s",
	ERROR_NON_SUPPORT_NOT_SHELL:        "不支持执行shell以外的脚本",

	ERROR_LOGIN_REMOTE_HOST:             "登录远程主机(%s)发生错误:%s",
	ERROR_RUN_REMOTE_COMMAND:            "在远程主机(%s)上执行命令(%s)发生错误:%s",
	ERROR_READ_STDOUT_FAILED:            "读取结果失败!",
	ERROR_REMOTE_HOST_MOUNT_FAILED:      "在远程主机(%s)上挂载NFS发生错误:%s",
	ERROR_REMOTE_HOST_UN_MOUNT_FAILED:   "在远程主机(%s)上卸载NFS发生错误:%s",
	ERROR_GET_REMOTE_HOST_INFO:          "获取远程主机(%s)系统信息发生错误:%s",
	ERROR_REMOTE_HOST_NFS_NOT_START:     "检测到远程主机(%s)NFS客户端未安装或已停止",
	ERROR_GET_REMOTE_HOST_DRIVE:         "获取远程主机(%s)可用盘符发生错误:%s",
	ERROR_REMOTE_HOST_NO_DRIVE:          "检测到远程主机(%s)没有可用盘符",
	ERROR_GET_REMOTE_HOST_DRIVE_BY_PATH: "通过共享路径获取远程主机(%s)挂载盘符发生错误:%s",
	ERROR_REMOTE_HOST_CREATE_SHORTCUT:   "在远程主机(%s)上创建桌面快捷方式发生错误",
	ERROR_REMOTE_HOST_DELETE_SHORTCUT:   "在远程主机(%s)上删除桌面快捷方式发生错误",
	ERROR_REMOTE_HOST_GET_DESKTOP_PATH:  "获取远程主机(%s)桌面路径发生错误:%s",
	ERROR_NOT_FOUND_WITH_DRIVER:         "未找到匹配的盘符!",

	ERROR_LOGIN_VSPHERE:                  "登录VMware平台(%s)发生错误:%s",
	ERROR_CREATE_VSPHERE_CONTAINER_VIEW:  "创建vSphere容器视图发生错误:%s",
	ERROR_GET_DATACENTER_FAILED:          "获取VMware平台(%s)数据中心发生错误:%s",
	ERROR_GET_FOLDER_FAILED:              "获取VMware平台(%s)文件夹发生错误:%s",
	ERROR_FIND_BY_INVENTORY_PATH:         "根据路径(%s)查找对象发生错误:%s",
	ERROR_OBJECT_IS_NULL:                 "查询对象为空!",
	ERROR_FIND_BY_INVENTORY_UUID:         "根据UUID(%s)查找对象发生错误:%s",
	ERROR_GET_INVENTORY_OBJECT:           "查询对象属性发生错误:%s",
	ERROR_GET_ALL_HOST_OBJECT:            "获取所有主机对象发生错误:%s",
	ERROR_GET_HOST_OBJECT_BY_NAME:        "根据主机名称获取主机对象发生错误:%s",
	ERROR_GET_HOST_MANAGE_OBJECT_BY_NAME: "根据主机名称获取主机存储管理对象发生错误:%s",
	ERROR_GET_SUB_OBJECT_UN_DEFINE:       "根据路径获取子对象发生错误，未定义对象类型!",
	ERROR_CONFIG_VM_FAILED:               "自定义配置虚拟机发生错误:%s",
	ERROR_POWER_ON_VM_FAILED:             "打开虚拟机电源发生错误:%s",
	ERROR_POWER_OFF_VM_FAILED:            "关闭虚拟机电源发生错误:%s",
	ERROR_CREATE_VM_SNAPSHOT_FAILED:      "创建虚拟机快照发生错误:%s",
	ERROR_RECOVER_VM_SNAPSHOT_FAILED:     "恢复虚拟机快照发生错误:%s",
	ERROR_CREATE_NFS_DATASTORE_FAILED:    "创建NFS存储发生错误:%s",
	ERROR_REGISTER_VM_FAILED:             "注册虚拟机发生错误:%s",
	ERROR_UN_REGISTER_VM_FAILED:          "取消注册虚拟机发生错误:%s",
	ERROR_GET_VM_BASE_INFO:               "获取虚拟机基本信息发生错误:%s",
	ERROR_VMX_FILE_NOT_EXIST:             "未找到虚拟机vmx文件!",
	ERROR_GET_VM_UUID_BY_VMX:             "从vmx文件中解析虚拟机uuid失败!",
	ERROR_GET_VM_NAME_BY_VMX:             "从vmx文件中解析虚拟机名称失败!",
	ERROR_ANALYTIC_MOUNT_PARAMETER:       "解析挂载参数发生错误:%s",
	ERROR_CHEACK_NFS_STORE:               "检查NFS存储是否存在发生错误:%s",
	ERROR_REMOVE_NFS_DATASTORE_FAILED:    "卸载NFS存储发生错误:%s",
	ERROR_GET_HOST_OBJECT_BY_PATH:        "根据主机路径获取主机对象发生错误:%s",

	ERROR_NONSUPPORT_APP_TYPE:             "不支持应用类型",
	ERROR_NONSUPPORT_FUNCTION:             "暂不支持该功能",
	ERROR_QUERY_RESOURCE_AUTH:             "查询资源失败:%s",
	ERROR_EDIT_RESOURCE_AUTH:              "编辑资源失败:%s",
	ERROR_DELETE_RESOURCE_AUTH:            "删除资源失败:%s",
	ERROR_PUBLISH_RESOURCE_AUTH:           "发布资源失败:%s",
	ERROR_PARAM_IS_NULL:                   "参数不允许为空!",
	INFO_START_MOUNT:                      "开始挂载副本!",
	INFO_START_RBD_MAP:                    "开始创建rbd映射!",
	INFO_RBD_MAP_SUCEESS:                  "创建rbd映射成功!",
	INFO_RBD_MAP_FAILED:                   "创建rbd映射失败!",
	INFO_START_FORMAT_REPLICA:             "开始格式化文件系统!",
	INFO_FORMAT_REPLICA_SUCCESS:           "格式化文件系统成功!",
	INFO_FORMAT_REPLICA_FAILED:            "格式化文件系统失败!",
	INFO_START_MOUNT_FILESYSTEM:           "开始挂载文件系统!",
	INFO_MOUNT_FILESYSTEM_SUCCESS:         "挂载文件系统成功!",
	INFO_MOUNT_FILESYSTEM_FAILED:          "挂载文件系统失败!",
	INFO_ADD_NFS_SHARE_SUCCESS:            "添加NFS共享成功!",
	INFO_ADD_NFS_SHARE_FAILED:             "添加NFS共享失败!",
	INFO_MOUNT_POINT_NAME:                 "挂载位置为:%s",
	ERROR_REPLICA_STATUS_NOT_MOUNT:        "副本未处于空闲状态，不允许挂载!",
	ERROR_REPLICA_STATUS_IS_MOUNT:         "副本未处于挂载状态，不允许卸载!",
	ERROR_GET_TARGET_PLATFORM_FAILED:      "获取挂载目标平台信息失败!",
	INFO_START_UNMOUNT:                    "开始卸载副本!",
	ERROR_GET_MAP_DEV_FAILED:              "获取rbd映射设备失败!",
	INFO_REMOVE_NFS_SHARE_SUCCESS:         "解除NFS共享成功!",
	INFO_REMOVE_NFS_SHARE_FAILED:          "解除NFS共享失败!",
	INFO_START_UNMOUNT_FILESYSTEM:         "开始卸载文件系统!",
	INFO_UNMOUNT_FILESYSTEM_SUCCESS:       "卸载文件系统成功!",
	INFO_UNMOUNT_FILESYSTEM_FAILED:        "卸载文件系统失败!",
	INFO_START_RBD_UN_MAP:                 "开始解除rbd映射!",
	INFO_RBD_UN_MAP_SUCEESS:               "解除rbd映射成功!",
	INFO_RBD_UN_MAP_FAILED:                "解除rbd映射失败!",
	ERROR_AES_DECRYPT:                     "AES解密失败!",
	ERROR_APP_INVALID_PARAMS:              "解析应用参数(%s)错误",
	ERROR_APP_PROCESS_ERROR:               "应用处理出错:%s",
	INFO_NFS_SHARE_PATH:                   "NFS共享路径为:%s",
	INFO_START_CREATE_NFS_STORE:           "开始创建NFS存储(%s)!",
	INFO_CREATE_NFS_STORE_SUCCESS:         "创建NFS存储(%s)成功!",
	INFO_CREATE_NFS_STORE_FAILED:          "创建NFS存储(%s)失败!",
	INFO_NFS_STORE_IS_EXIST:               "检测到NFS存储(%s)已存在!",
	INFO_START_REGISTER_VM_BY_VMWARE:      "开始注册虚拟机(%s)!",
	INFO_REGISTER_VM_FAILED_BY_VMWARE:     "注册虚拟机(%s)失败!",
	INFO_REGISTER_VM_SUCCESS_BY_VMWARE:    "注册虚拟机(%s)失败!",
	INFO_START_CONFIG_VM_BY_VMWARE:        "开始重新配置虚拟机(%s)!",
	INFO_CONFIG_VM_FAILED_BY_VMWARE:       "重新配置虚拟机(%s)失败!",
	INFO_CONFIG_VM_SUCCESS_BY_VMWARE:      "重新配置虚拟机(%s)成功!",
	INFO_START_POWER_ON_VM_BY_VMWARE:      "开始打开虚拟机(%s)电源!",
	INFO_POWER_ON_VM_FAILED_BY_VMWARE:     "打开虚拟机(%s)电源失败!",
	INFO_POWER_ON_VM_SUCCESS_BY_VMWARE:    "打开虚拟机(%s)电源成功!",
	INFO_START_UN_REGISTER_VM_BY_VMWARE:   "开始取消虚拟机(%s)注册!",
	INFO_UN_REGISTER_VM_BY_VMWARE_FAILED:  "取消虚拟机(%s)注册失败!",
	INFO_UN_REGISTER_VM_BY_VMWARE_SUCCESS: "取消虚拟机(%s)注册成功!",
	INFO_START_DELETE_NFS_STORE:           "开始删除NFS存储(%s)!",
	INFO_DELETE_NFS_STORE_FAILED:          "删除NFS存储(%s)失败!",
	INFO_DELETE_NFS_STORE_SUCCESS:         "删除NFS存储(%s)成功!",
	INFO_START_CREATE_SHORTCUT_STORE:      "开始创建桌面快捷方式(%s)!",
	INFO_CREATE_SHORTCUT_FALIED:           "创建桌面快捷方式(%s)失败!",
	INFO_CREATE_SHORTCUT_SUCCESS:          "创建桌面快捷方式(%s)成功!",
	INFO_START_DELETE_SHORTCUT_STORE:      "开始删除桌面快捷方式(%s)!",
	INFO_DELETE_SHORTCUT_FAILED:           "删除桌面快捷方式(%s)失败!",
	INFO_DELETE_SHORTCUT_SUCCESS:          "删除桌面快捷方式(%s)成功!",
	ERROR_REPLICA_TYPE_MISMATCHING:        "副本类型与挂载目标平台类型不匹配!",
	ERROR_CREATE_FILE_CACHE:               "创建缓存目录发生错误:%s",
	ERROR_BROWSE_FILE:                     "浏览文件列表发生错误:%s",
	ERROR_UPLOAD_FILE:                     "上传文件发生错误:%s",
	ERROR_DOWNLOAD_FILE:                   "下载文件发生错误:%s",
	INFO_START_RUNNING_SHELL_SCRIPT:       "开始运行shell脚本(%s)",
	INFO_RUNNING_SHELL_SCRIPT_FAILED:      "运行shell脚本(%s)失败!",
	INFO_RUNNING_SHELL_SCRIPT_SUCCESS:     "运行shell脚本(%s)成功!",
	ERROR_PROGRAM_CRASHES:                 "捕获到程序运行发生崩溃!",
	ERROR_UPLOAD_FILE_SZIE_IS_TOO_BIG:     "文件大小超出限制",
	ERROR_UPLOAD_FILE_IS_NOT_SHELL:        "文件类型不是shell脚本文件",
	ERROR_GET_SCRIPT_FILE_CONTENT:         "获取脚本文件内容发生错误:%s",
	ERROR_DELETE_SCRIPT_FILE:              "删除脚本发生错误:%s",
	ERROR_GET_SCRIPT_FILE_INFO:            "获取脚本信息发生错误:%s",
	INFO_MOUNT_THE_WAY:                    "本次执行的挂载方式为:%s",
	ERROR_GET_SERVER_INFO:                 "获取服务器状态信息发生错误:%s",
	INFO_GROWFS_REPLICA_SUCCESS:           "同步文件系统成功",
	INFO_GROWFS_REPLICA_FAILED:            "同步文件系统失败",
	ERROR_DO_JENKINS_JOB:                  "执行jenk任务(%s)失败:%s",
	ERROR_VMNAME_EXSIT:                    "虚拟机名称已存在，请勿重复",
	ERROR_IP_EXSIT:                        "ip地址冲突,请更换一个",
	ERROR_VMWARE_TYPE:                     "ESXi类型平台不支持创建虚拟机",
	INFO_START_VM:                         "开始注册虚拟机",
	INFO_START_POWERON:                    "开始打开虚拟机电源",
	ERROR_SAVEHOST:                        "保存虚拟机信息报错:%s",
	INFO_SAVEHOST:                         "保存虚拟机信息成功",
	INFO_START_POWEROFF:                   "开始关闭虚拟机电源%s",
	ERROR_SET_IPADDR:                      "配置ip地址错误:%s",
	//CAS
	ERROR_LOGIN_CAS:                 "登录CAS平台(%s)错误:%s",
	ERRIR_CREATE_POOL:               "创建CAS存储池(%s)错误:%s",
	ERRIR_START_POOL:                "启动CAS存储池(%s)错误:%s",
	ERRIR_CREATE_CAS:                "创建CAS虚拟机(%s)错误:%s",
	ERRIR_QUERY_CAS:                 "查询CAS虚拟机(%s)错误:%s",
	ERRIR_DELETE_CAS:                "删除CAS虚拟机(%s)错误:%s",
	ERRIR_DELETE_POOL:               "删除CAS存储池(%s)错误:%s",
	ERROR_GET_CAS_DATASOURCE_FAILED: "获取CAS平台数据源发生错误:%s",
	INFO_CREATE_STORAGE:             "开始添加CAS平台存储",
	INFO_CREATE_VM:                  "开始创建CAS平台虚拟机",
	INFO_CREATE_VM_SUCCESS:          "创建CAS平台虚拟机【%s】成功!",

	//FUSIONCOMPUTE
	ERROR_REQUEST_FUSIONCOMPUTE:               "请求FUSIONCOMPUTE平台API(%s)错误:%s",
	ERROR_LOGIN_FUSIONCOMPUTE:                 "登录FUSIONCOMPUTE平台(%s)错误:%s",
	ERROR_GET_FUSIONCOMPUTE_INFO_FAILED:       "获取FUSIONCOMPUTE平台信息错误:%s",
	ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED: "获取FUSIONCOMPUTE平台(%s)资源失败:%s",

	ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED: "FUSIONCOMPUTE平台(%s)添加存储资源失败:%s",
	ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED:    "FUSIONCOMPUTE平台(%s)关联主机失败:%s",
	ERROR_FUSIONCOMPUTE_MOUNT_FAILED:                "挂载FUSIONCOMPUTE平台(%s)失败:%s",
	ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED:              "卸载FUSIONCOMPUTE平台(%s)失败:%s",

	ERROR_FUSIONCOMPUTE_REFRESH:       "FUSIONCOMPUTE平台(%s)刷新存储设备失败:%s",
	ERROR_FUSIONCOMPUTE_ADD_DATASTORE: "FUSIONCOMPUTE平台(%s)添加数据存储失败:%s",
	ERROR_FUSIONCOMPUTE_DEL_DATASTORE: "FUSIONCOMPUTE平台(%s)删除数据存储失败:%s",
	ERROR_FUSIONCOMPUTE_GET_DATASTORE: "FUSIONCOMPUTE平台(%s)查询数据存储失败:%s",
	ERROR_FUSIONCOMPUTE_CRT_VM:        "FUSIONCOMPUTE平台(%s)创建虚拟机失败:%s",
	ERROR_FUSIONCOMPUTE_GET_VM_INFO:   "FUSIONCOMPUTE平台(%s)获取虚拟机信息失败:%s",

	INFO_GET_FUSIONCOMPUTE_DATASOURCE:              "获取FUSIONCOMPUTE平台(%s)资源:%s",
	INFO_START_ADD_FUSIONCOMPUTE_STORAGERE_SOURCE:  "开始在FUSIONCOMPUTE平台(%s)添加存储资源:%s",
	INFO_SUCESS_ADD_FUSIONCOMPUTE_STORAGERE_SOURCE: "成功在FUSIONCOMPUTE平台(%s)添加存储资源:%s",
	INFO_START_ADD_FUSIONCOMPUTE_RESOURCE_HOST:     "开始在FUSIONCOMPUTE平台(%s)关联主机:%s",
	INFO_SUCESS_ADD_FUSIONCOMPUTE_RESOURCE_HOST:    "成功在FUSIONCOMPUTE平台(%s)关联主机:%s",
	INFO_START_FUSIONCOMPUTE_REFRESH:               "开始在FUSIONCOMPUTE平台(%s)扫描存储设备",
	INFO_SUCESS_FUSIONCOMPUTE_REFRESH:              "成功在FUSIONCOMPUTE平台(%s)扫描存储设备:%s",
	INFO_SUCESS_FUSIONCOMPUTE_ADD_DATASTORE:        "成功在FUSIONCOMPUTE平台(%s)添加数据存储:%s",
	INFO_SUCESS_FUSIONCOMPUTE_DEL_DATASTORE:        "成功在FUSIONCOMPUTE平台(%s)删除数据存储:%s",
	INFO_START_FUSIONCOMPUTE_CRT_VM:                "开始在FUSIONCOMPUTE平台(%s)创建虚拟机:%s",
	INFO_SUCESS_FUSIONCOMPUTE_CRT_VM:               "FUSIONCOMPUTE平台(%s)创建虚拟机:%s成功",
	INFO_SUCESS_FUSIONCOMPUTE_MV_DISK:              "FUSIONCOMPUTE平台(%s)虚拟机:%s替换磁盘文件成功",
}

func GetOperation(code Operation, args ...interface{}) string {
	msg, ok := OperationFlags[code]
	if ok {
		return fmt.Sprintf(msg, args...)
	}
	return OperationFlags[UNDEFINE_OPERAT]
}

func GetMsg(code Massage, args ...interface{}) string {
	msg, ok := MassageFlags[code]
	if ok {
		return fmt.Sprintf(msg, args...)
	}
	return MassageFlags[FAILURE]
}