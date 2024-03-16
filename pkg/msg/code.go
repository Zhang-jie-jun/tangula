package msg

// 操作定义(2000 --- 3000)
type Operation int

const (
	UNDEFINE_OPERAT Operation = -1
	LOGIN           Operation = 2000 + iota
	LOGIN_SUCCESS
	LOGIN_FAILED
	LOGOUT
	LOGOUT_SUCCESS
	LOGOUT_FAILED
	ENABLE_USER
	ENABLE_USER_SUCCESS
	ENABLE_USER_FAILED
	DISABLE_USER
	DISABLE_USER_SUCCESS
	DISABLE_USER_FAILED
	EDIT_USER_ROLE
	EDIT_USER_ROLE_SUCCESS
	EDIT_USER_ROLE_FAILED
	CREATE_ROLE
	CREATE_ROLE_SUCCESS
	CREATE_ROLE_FAILED
	EDIT_ROLE
	EDIT_ROLE_SUCCESS
	EDIT_ROLE_FAILED
	DELETE_ROLE
	DELETE_ROLE_SUCCESS
	DELETE_ROLE_FAILED
	CREATE_PLATFORM
	CREATE_PLATFORM_SUCCESS
	CREATE_PLATFORM_FAILED
	EDIT_PLATFORM
	EDIT_PLATFORM_SUCCESS
	EDIT_PLATFORM_FAILED
	PUBLISH_PLATFORM
	PUBLISH_PLATFORM_SUCCESS
	PUBLISH_PLATFORM_FAILED
	DELETE_PLATFORM
	DELETE_PLATFORM_SUCCESS
	DELETE_PLATFORM_FAILED
	CREATE_HOST
	CREATE_HOST_SUCCESS
	CREATE_HOST_FAILED
	EDIT_HOST
	EDIT_HOST_SUCCESS
	EDIT_HOST_FAILED
	PUBLISH_HOST
	PUBLISH_HOST_SUCCESS
	PUBLISH_HOST_FAILED
	UPD_HOST
	UPD_HOST_SUCCESS
	UPD_HOST_FAILED
	DELETE_HOST
	DELETE_HOST_SUCCESS
	DELETE_HOST_FAILED
	CREATE_POOL
	CREATE_POOL_SUCCESS
	CREATE_POOL_FAILED
	DELETE_POOL
	DELETE_POOL_SUCCESS
	DELETE_POOL_FAILED
	CREATE_IMAGE_BY_REPLICA
	CREATE_IMAGE_BY_REPLICA_SUCCESS
	CREATE_IMAGE_BY_REPLICA_FAILED
	CREATE_IMAGE_BY_SNAPSHOT
	CREATE_IMAGE_BY_SNAPSHOT_SUCCESS
	CREATE_IMAGE_BY_SNAPSHOT_FAILED
	PUBLISH_IMAGE
	PUBLISH_IMAGE_SUCCESS
	PUBLISH_IMAGE_FAILED
	DELETE_IMAGE
	DELETE_IMAGE_SUCCESS
	DELETE_IMAGE_FAILED
	CREATE_REPLICA
	CREATE_REPLICA_SUCCESS
	CREATE_REPLICA_FAILED
	DELETE_REPLICA
	DELETE_REPLICA_SUCCESS
	DELETE_REPLICA_FAILED
	MOUNT_REPLICA
	MOUNT_REPLICA_SUCCESS
	MOUNT_REPLICA_FAILED
	UNMOUNT_REPLICA
	UNMOUNT_REPLICA_SUCCESS
	UNMOUNT_REPLICA_FAILED
	CREATE_SNAPSHOT
	CREATE_SNAPSHOT_SUCCESS
	CREATE_SNAPSHOT_FAILED
	ROLLBACK_SNAPSHOT
	ROLLBACK_SNAPSHOT_SUCCESS
	ROLLBACK_SNAPSHOT_FAILED
	DELETE_SNAPSHOT
	DELETE_SNAPSHOT_SUCCESS
	DELETE_SNAPSHOT_FAILED
	CREATE_REPLICA_BY_IMAGE
	CREATE_REPLICA_BY_IMAGE_BEGIN
	CREATE_REPLICA_BY_IMAGE_SUCCESS
	CREATE_REPLICA_BY_IMAGE_FAILED
	CREATE_REPLICA_BY_SNAPSHOT
	CREATE_REPLICA_BY_SNAPSHOT_SUCCESS
	CREATE_REPLICA_BY_SNAPSHOT_FAILED
	DELETE_SNAPSHOT_BY_REPLICA
	DELETE_SNAPSHOT_BY_REPLICA_SUCCESS
	DELETE_SNAPSHOT_BY_REPLICA_FAILED
	UPLOAD_SCRIPT_FILE
	UPLOAD_SCRIPT_FILE_SUCCESS
	UPLOAD_SCRIPT_FILE_FAILED
	DOWNLOAD_SCRIPT_FILE
	DOWNLOAD_SCRIPT_FILE_SUCCESS
	DOWNLOAD_SCRIPT_FILE_FAILED
	DELETE_SCRIPT_FILE
	DELETE_SCRIPT_FILE_SUCCESS
	DELETE_SCRIPT_FILE_FAILED
	DEPLOY_CLIENT
	DEPLOY_CLIENT_FAILED
	DEPLOY_CLIENT_START
	DEPLOY_CLIENT_SUCCESS
)

// 消息定义
type Massage int

const (
	// 接口状态
	UNDEFINED     Massage = -1
	SUCCESS       Massage = 200
	FAILURE       Massage = 400
	URL_NOT_FOUND Massage = 404
	// 用户管理
	ERROR                Massage = 10001
	ERROR_INVALID_PARAMS Massage = 10002
	ERROR_GET_USER_INFO  Massage = 10003
	ERROR_USER_IS_EXIST  Massage = 10004
	ERROR_USER_NOT_EXIST Massage = 10005
	ERROR_SET_USER_ROLE  Massage = 10006
	ERROR_ENABLE_USER    Massage = 10007
	ERROR_DISABLE_USER   Massage = 10008
	ERROR_GET_ROLE_INFO  Massage = 10009
	ERROR_ROLE_IS_EXIST  Massage = 10010
	// 权限认证
	AUTH_PERMISSION_DENIED   Massage = 20001
	AUTH_CHECK_TOKEN_FAIL    Massage = 20002
	AUTH_CHECK_TOKEN_TIMEOUT Massage = 20003
	AUTH_TOKEN               Massage = 20004
	AUTH_ERROR               Massage = 20005
	// Ceph管理
	ERROR_STORE_POOL_NAME_IS_EXIST        Massage = 30004
	ERROR_CREATE_STORE_POOL               Massage = 30005
	ERROR_GET_STORE_POOL_INFO             Massage = 30006
	ERROR_DELETE_STORE_POOL               Massage = 30007
	ERROR_GET_IMAGE_INFO                  Massage = 30008
	ERROR_DELETE_IMAGE                    Massage = 30009
	ERROR_PUBLISH_IMAGE                   Massage = 30010
	ERROR_REPLICA_NAME_IS_EXIST           Massage = 30011
	ERROR_IMAGE_TO_REPLICA                Massage = 30012
	ERROR_CREATE_REPLICA                  Massage = 30013
	ERROR_STORE_INSUFFICIENT_SPACE        Massage = 30014
	ERROR_GET_REPLICA_INFO                Massage = 30015
	ERROR_DELETE_REPLICA                  Massage = 30016
	ERROR_IMAGE_NAME_IS_EXIST             Massage = 30017
	ERROR_REPLICA_TO_IMAGE                Massage = 30018
	ERROR_SNAPSHOT_NAME_IS_EXIST          Massage = 30019
	ERROR_CREATE_SNAPSHOT                 Massage = 30020
	ERROR_GET_SNAPSHOT_INFO               Massage = 30021
	ERROR_DELETE_SNAPSHOT                 Massage = 30022
	ERROR_ROLL_BACK_SNAPSHOT              Massage = 30023
	ERROR_REPLICAT_STATUS                 Massage = 30024
	ERROR_SNAPSHOT_TO_IMAGE               Massage = 30025
	ERROR_SNAPSHOT_TO_REPLICA             Massage = 30026
	ERROR_EXIST_MOUNT_STATUS_REPLICA      Massage = 30027
	ERROR_EDIT_REPLICA                    Massage = 30028
	ERROR_REPLICA_CAPACITY_EXPANSION_ONLY Massage = 30029
	ERROR_EDIT_IMAGE                      Massage = 30030
	// 平台管理
	ERROR_PLATFORM_IS_EXIST                   Massage = 40001
	ERROR_PLATFORM_NOT_EXIST                  Massage = 40002
	ERROR_CREATE_PLATFORM                     Massage = 40003
	ERROR_VERIFY_PLATFORM                     Massage = 40004
	ERROR_UPDATE_PLATFORM                     Massage = 40005
	ERROR_DELETE_PLATFORM                     Massage = 40006
	ERROR_PUBLISH_PLATFORM                    Massage = 40007
	ERROR_PLATFORM_NAME_IS_EXIST              Massage = 40008
	ERROR_PLATFORM_IP_IS_EXIST                Massage = 40009
	ERROR_GET_PLATFORM_INFO                   Massage = 40010
	ERROR_PLATFORM_TYPE_NO_MATCH              Massage = 40011
	ERROR_GET_VMWARE_DATASOURCE_FAILED        Massage = 40012
	ERROR_SHOW_TYPE_NON_SUPPORT               Massage = 40013
	ERROR_PLATFORM_EXIST_MOUNT_STATUS_REPLICA Massage = 40014
	ERROR_GET_VMWARE_HOSTS_FAILED             Massage = 40015
	// 主机管理
	ERROR_GET_HOST_INFO         Massage = 50001
	ERROR_CREATE_HOST           Massage = 50002
	ERROR_VERIFY_HOST           Massage = 50003
	ERROR_UPDATE_HOST           Massage = 50004
	ERROR_DELETE_HOST           Massage = 50005
	ERROR_PUBLISH_HOST          Massage = 50006
	ERROR_HOST_NAME_IS_EXIST    Massage = 50007
	ERROR_HOST_IP_IS_EXIST      Massage = 50008
	ERROR_UPD_HOST_INFO         Massage = 50009
	ERROR_HOST_DEPLOYING        Massage = 50010
	ERROR_HOST_START_JENKINSJOB Massage = 50011
	//客户端部署相关
	ERROR_GET_DEPLOYAPP Massage = 50012
	ERROR_UPD_DEPLOYAPP Massage = 50013

	// 日志管理
	ERROR_GET_LOG_RECORD_INFO    Massage = 70001
	ERROR_GET_LOG_DETAIL_INFO    Massage = 70002
	ERROR_UPDATE_LOG_RECORD_INFO Massage = 70003
	ERROR_CREATE_LOG_RECORD_INFO Massage = 70004
	ERROR_CREATE_LOG_DETAIL_INFO Massage = 70005
	// Ceph操作
	ERROR_CEPH_NOT_CONNECT               Massage = 80001
	ERROR_CEPH_POOL_NOT_OPEN             Massage = 80002
	ERROR_CEPH_IMAGE_NOT_OPEN            Massage = 80003
	ERROR_CEPH_GET_CLUSTER_STATUS        Massage = 80004
	ERROR_CEPH_CREATE_STORE_POOL         Massage = 80005
	ERROR_CEPH_DELETE_STORE_POOL         Massage = 80006
	ERROR_CEPH_CREATE_IMAGE              Massage = 80007
	ERROR_CEPH_DELETE_IMAGE              Massage = 80008
	ERROR_CEPH_CLONE_IMAGE               Massage = 80009
	ERROR_CEPH_RESIZE_IMAGE              Massage = 80010
	ERROR_CEPH_RENAME_IMAGE              Massage = 80011
	ERROR_CEPH_FLUSH_IMAGE               Massage = 80012
	ERROR_CEPH_CREATE_IMAGE_SNAPSHOT     Massage = 80013
	ERROR_CEPH_DELETE_IMAGE_SNAPSHOT     Massage = 80014
	ERROR_CEPH_PROTECT_IMAGE_SNAPSHOT    Massage = 80015
	ERROR_CEPH_UN_PROTECT_IMAGE_SNAPSHOT Massage = 80016
	ERROR_CEPH_GET_IMAGE_FEATURES        Massage = 80017
	ERROR_CEPH_COPY_IMAGE                Massage = 80018
	ERROR_CEPH_ROLLBACK_IMAGE_SNAPSHOT   Massage = 80019
	ERROR_CEPH_FLATTEN_IMAGE             Massage = 80020
	ERROR_CEPH_NON_SUPPORT_WINDOWS       Massage = 80021
	ERROR_CEPH_MAP_IMAGE                 Massage = 80022
	ERROR_CEPH_GET_MAP_INFO              Massage = 80023
	ERROR_CEPH_UN_MAP_IMAGE              Massage = 80024
	// 系统调用
	ERROR_RUN_SYSTEM_COMMAND           Massage = 90001
	ERROR_EXPECTED_RESULT_IS_INCORRECT Massage = 90002
	ERROR_DEV_NOT_EXISTS               Massage = 90003
	ERROR_FILE_NOT_EXISTS              Massage = 90004
	ERROR_DIR_NOT_EXISTS               Massage = 90005
	ERROR_CREATE_DIR                   Massage = 90006
	ERROR_MOUNT_POINT_OCCUPY           Massage = 90007
	ERROR_MOUNT_FILED                  Massage = 90008
	ERROR_UNMOUNT_FILED                Massage = 90009
	ERROR_DIR_IS_MOUNT_POINT           Massage = 90010
	ERROR_DIR_NOT_MOUNT_POINT          Massage = 90011
	ERROR_GET_SERVER_STATUS            Massage = 90012
	ERROR_START_SERVER                 Massage = 90013
	ERROR_RESTART_SERVER               Massage = 90014
	ERROR_STOP_SERVER                  Massage = 90015
	ERROR_RELOAD_SERVER                Massage = 90016
	ERROR_ENABLE_SERVER                Massage = 90017
	ERROR_DISABLE_SERVER               Massage = 90018
	ERROR_ADD_NFS_SHARE                Massage = 90019
	ERROR_REMOVE_NFS_SHARE             Massage = 90020
	ERROR_SERVER_NOT_START             Massage = 90021
	ERROR_MOUNT_POINT_NOT_INCORRECT    Massage = 90022
	ERROR_GET_LOCAL_ADDR               Massage = 90023
	ERROR_RUN_SHELL_SCRIPT             Massage = 90024
	ERROR_NON_SUPPORT_NOT_SHELL        Massage = 90025
	// 远程主机调用
	ERROR_LOGIN_REMOTE_HOST             Massage = 100001
	ERROR_RUN_REMOTE_COMMAND            Massage = 100002
	ERROR_READ_STDOUT_FAILED            Massage = 100003
	ERROR_REMOTE_HOST_MOUNT_FAILED      Massage = 100004
	ERROR_REMOTE_HOST_UN_MOUNT_FAILED   Massage = 100005
	ERROR_GET_REMOTE_HOST_INFO          Massage = 100006
	ERROR_REMOTE_HOST_NFS_NOT_START     Massage = 100007
	ERROR_GET_REMOTE_HOST_DRIVE         Massage = 100008
	ERROR_REMOTE_HOST_NO_DRIVE          Massage = 100009
	ERROR_GET_REMOTE_HOST_DRIVE_BY_PATH Massage = 100010
	ERROR_REMOTE_HOST_CREATE_SHORTCUT   Massage = 100011
	ERROR_REMOTE_HOST_DELETE_SHORTCUT   Massage = 100012
	ERROR_REMOTE_HOST_GET_DESKTOP_PATH  Massage = 100013
	ERROR_NOT_FOUND_WITH_DRIVER         Massage = 100014
	// VMware
	ERROR_LOGIN_VSPHERE                  Massage = 200001
	ERROR_CREATE_VSPHERE_CONTAINER_VIEW  Massage = 200002
	ERROR_GET_DATACENTER_FAILED          Massage = 200003
	ERROR_GET_FOLDER_FAILED              Massage = 200004
	ERROR_FIND_BY_INVENTORY_PATH         Massage = 200005
	ERROR_OBJECT_IS_NULL                 Massage = 200006
	ERROR_FIND_BY_INVENTORY_UUID         Massage = 200007
	ERROR_GET_INVENTORY_OBJECT           Massage = 200008
	ERROR_GET_ALL_HOST_OBJECT            Massage = 200009
	ERROR_GET_HOST_OBJECT_BY_NAME        Massage = 200010
	ERROR_GET_HOST_MANAGE_OBJECT_BY_NAME Massage = 200011
	ERROR_GET_SUB_OBJECT_UN_DEFINE       Massage = 200012
	ERROR_CONFIG_VM_FAILED               Massage = 200013
	ERROR_POWER_ON_VM_FAILED             Massage = 200014
	ERROR_POWER_OFF_VM_FAILED            Massage = 200015
	ERROR_CREATE_VM_SNAPSHOT_FAILED      Massage = 200016
	ERROR_RECOVER_VM_SNAPSHOT_FAILED     Massage = 200017
	ERROR_CREATE_NFS_DATASTORE_FAILED    Massage = 200018
	ERROR_REGISTER_VM_FAILED             Massage = 200019
	ERROR_UN_REGISTER_VM_FAILED          Massage = 200020
	ERROR_GET_VM_BASE_INFO               Massage = 200021
	ERROR_VMX_FILE_NOT_EXIST             Massage = 200022
	ERROR_GET_VM_UUID_BY_VMX             Massage = 200023
	ERROR_GET_VM_NAME_BY_VMX             Massage = 200024
	ERROR_ANALYTIC_MOUNT_PARAMETER       Massage = 200025
	ERROR_CHEACK_NFS_STORE               Massage = 200026
	ERROR_REMOVE_NFS_DATASTORE_FAILED    Massage = 200027
	ERROR_GET_HOST_OBJECT_BY_PATH        Massage = 200028
	INFO_START_VM                        Massage = 200029
	INFO_START_POWERON                   Massage = 200030
	ERROR_SAVEHOST                       Massage = 200031
	INFO_SAVEHOST                        Massage = 200032
	INFO_START_POWEROFF                  Massage = 200033
	ERROR_SET_IPADDR                     Massage = 200034
	// cas
	ERROR_LOGIN_CAS                 Massage = 300001
	ERRIR_CREATE_POOL               Massage = 300002
	ERRIR_START_POOL                Massage = 300003
	ERRIR_CREATE_CAS                Massage = 300004
	ERRIR_QUERY_CAS                 Massage = 300005
	ERRIR_DELETE_CAS                Massage = 300006
	ERRIR_DELETE_POOL               Massage = 300007
	ERROR_GET_CAS_DATASOURCE_FAILED Massage = 300008
	INFO_CREATE_STORAGE             Massage = 300009
	INFO_CREATE_VM                  Massage = 300010
	INFO_CREATE_VM_SUCCESS          Massage = 300011

	// fusioncompute
	ERROR_REQUEST_FUSIONCOMPUTE               Massage = 400000
	ERROR_LOGIN_FUSIONCOMPUTE                 Massage = 400001
	ERROR_GET_FUSIONCOMPUTE_INFO_FAILED       Massage = 400002
	ERROR_GET_FUSIONCOMPUTE_DATASOURCE_FAILED Massage = 400003
	ERROR_FUSIONCOMPUTE_MOUNT_FAILED          Massage = 400004
	ERROR_FUSIONCOMPUTE_UNMOUNT_FAILED        Massage = 400014

	ERROR_FUSIONCOMPUTE_ADD_STORAGERE_SOURCE_FAILED Massage = 400005
	ERROR_FUSIONCOMPUTE_ADD_RESOURCE_HOST_FAILED    Massage = 400006
	ERROR_FUSIONCOMPUTE_REFRESH                     Massage = 400007
	ERROR_FUSIONCOMPUTE_ADD_DATASTORE               Massage = 400008
	ERROR_FUSIONCOMPUTE_DEL_DATASTORE               Massage = 400009
	ERROR_FUSIONCOMPUTE_GET_DATASTORE               Massage = 400010
	ERROR_FUSIONCOMPUTE_CRT_VM                      Massage = 400011
	ERROR_FUSIONCOMPUTE_GET_VM_INFO                 Massage = 400012

	INFO_GET_FUSIONCOMPUTE_DATASOURCE              Massage = 410000
	INFO_START_ADD_FUSIONCOMPUTE_STORAGERE_SOURCE  Massage = 410001
	INFO_SUCESS_ADD_FUSIONCOMPUTE_STORAGERE_SOURCE Massage = 410002
	INFO_START_ADD_FUSIONCOMPUTE_RESOURCE_HOST     Massage = 410003
	INFO_SUCESS_ADD_FUSIONCOMPUTE_RESOURCE_HOST    Massage = 410004
	INFO_START_FUSIONCOMPUTE_REFRESH               Massage = 410005
	INFO_SUCESS_FUSIONCOMPUTE_REFRESH              Massage = 410006
	INFO_SUCESS_FUSIONCOMPUTE_ADD_DATASTORE        Massage = 410007
	INFO_SUCESS_FUSIONCOMPUTE_DEL_DATASTORE        Massage = 410008
	INFO_START_FUSIONCOMPUTE_CRT_VM                Massage = 410009
	INFO_SUCESS_FUSIONCOMPUTE_CRT_VM               Massage = 410010
	INFO_SUCESS_FUSIONCOMPUTE_MV_DISK              Massage = 410011
	// 其他消息
	ERROR_NONSUPPORT_APP_TYPE             Massage = 990001
	ERROR_NONSUPPORT_FUNCTION             Massage = 990002
	ERROR_QUERY_RESOURCE_AUTH             Massage = 990003
	ERROR_EDIT_RESOURCE_AUTH              Massage = 990004
	ERROR_DELETE_RESOURCE_AUTH            Massage = 990005
	ERROR_PUBLISH_RESOURCE_AUTH           Massage = 990006
	ERROR_PARAM_IS_NULL                   Massage = 990007
	INFO_START_MOUNT                      Massage = 990008
	INFO_START_RBD_MAP                    Massage = 990009
	INFO_RBD_MAP_SUCEESS                  Massage = 990010
	INFO_RBD_MAP_FAILED                   Massage = 990011
	INFO_START_FORMAT_REPLICA             Massage = 990012
	INFO_FORMAT_REPLICA_SUCCESS           Massage = 990013
	INFO_FORMAT_REPLICA_FAILED            Massage = 990014
	INFO_START_MOUNT_FILESYSTEM           Massage = 990015
	INFO_MOUNT_FILESYSTEM_SUCCESS         Massage = 990016
	INFO_MOUNT_FILESYSTEM_FAILED          Massage = 990017
	INFO_ADD_NFS_SHARE_SUCCESS            Massage = 990018
	INFO_ADD_NFS_SHARE_FAILED             Massage = 990019
	INFO_MOUNT_POINT_NAME                 Massage = 990020
	ERROR_REPLICA_STATUS_NOT_MOUNT        Massage = 990021
	ERROR_REPLICA_STATUS_IS_MOUNT         Massage = 990022
	ERROR_GET_TARGET_PLATFORM_FAILED      Massage = 990023
	INFO_START_UNMOUNT                    Massage = 990024
	ERROR_GET_MAP_DEV_FAILED              Massage = 990025
	INFO_REMOVE_NFS_SHARE_SUCCESS         Massage = 990026
	INFO_REMOVE_NFS_SHARE_FAILED          Massage = 990027
	INFO_START_UNMOUNT_FILESYSTEM         Massage = 990028
	INFO_UNMOUNT_FILESYSTEM_SUCCESS       Massage = 990029
	INFO_UNMOUNT_FILESYSTEM_FAILED        Massage = 990030
	INFO_START_RBD_UN_MAP                 Massage = 990031
	INFO_RBD_UN_MAP_SUCEESS               Massage = 990032
	INFO_RBD_UN_MAP_FAILED                Massage = 990033
	ERROR_AES_DECRYPT                     Massage = 990034
	ERROR_APP_INVALID_PARAMS              Massage = 990035
	ERROR_APP_PROCESS_ERROR               Massage = 990036
	INFO_NFS_SHARE_PATH                   Massage = 990037
	INFO_START_CREATE_NFS_STORE           Massage = 990038
	INFO_CREATE_NFS_STORE_SUCCESS         Massage = 990039
	INFO_CREATE_NFS_STORE_FAILED          Massage = 990040
	INFO_NFS_STORE_IS_EXIST               Massage = 990041
	INFO_START_REGISTER_VM_BY_VMWARE      Massage = 990042
	INFO_REGISTER_VM_FAILED_BY_VMWARE     Massage = 990043
	INFO_REGISTER_VM_SUCCESS_BY_VMWARE    Massage = 990044
	INFO_START_CONFIG_VM_BY_VMWARE        Massage = 990045
	INFO_CONFIG_VM_FAILED_BY_VMWARE       Massage = 990046
	INFO_CONFIG_VM_SUCCESS_BY_VMWARE      Massage = 990047
	INFO_START_POWER_ON_VM_BY_VMWARE      Massage = 990048
	INFO_POWER_ON_VM_FAILED_BY_VMWARE     Massage = 990049
	INFO_POWER_ON_VM_SUCCESS_BY_VMWARE    Massage = 990050
	INFO_START_UN_REGISTER_VM_BY_VMWARE   Massage = 990051
	INFO_UN_REGISTER_VM_BY_VMWARE_FAILED  Massage = 990052
	INFO_UN_REGISTER_VM_BY_VMWARE_SUCCESS Massage = 990053
	INFO_START_DELETE_NFS_STORE           Massage = 990054
	INFO_DELETE_NFS_STORE_FAILED          Massage = 990055
	INFO_DELETE_NFS_STORE_SUCCESS         Massage = 990056
	INFO_START_CREATE_SHORTCUT_STORE      Massage = 990057
	INFO_CREATE_SHORTCUT_FALIED           Massage = 990058
	INFO_CREATE_SHORTCUT_SUCCESS          Massage = 990059
	INFO_START_DELETE_SHORTCUT_STORE      Massage = 990060
	INFO_DELETE_SHORTCUT_FAILED           Massage = 990061
	INFO_DELETE_SHORTCUT_SUCCESS          Massage = 990062
	ERROR_REPLICA_TYPE_MISMATCHING        Massage = 990063
	ERROR_CREATE_FILE_CACHE               Massage = 990064
	ERROR_BROWSE_FILE                     Massage = 990065
	ERROR_UPLOAD_FILE                     Massage = 990066
	ERROR_DOWNLOAD_FILE                   Massage = 990067
	INFO_START_RUNNING_SHELL_SCRIPT       Massage = 990068
	INFO_RUNNING_SHELL_SCRIPT_FAILED      Massage = 990069
	INFO_RUNNING_SHELL_SCRIPT_SUCCESS     Massage = 990070
	ERROR_PROGRAM_CRASHES                 Massage = 990071
	ERROR_UPLOAD_FILE_SZIE_IS_TOO_BIG     Massage = 990072
	ERROR_UPLOAD_FILE_IS_NOT_SHELL        Massage = 990073
	ERROR_GET_SCRIPT_FILE_CONTENT         Massage = 990074
	ERROR_DELETE_SCRIPT_FILE              Massage = 990075
	ERROR_GET_SCRIPT_FILE_INFO            Massage = 990076
	INFO_MOUNT_THE_WAY                    Massage = 990077
	ERROR_GET_SERVER_INFO                 Massage = 990078
	INFO_GROWFS_REPLICA_SUCCESS           Massage = 990079
	INFO_GROWFS_REPLICA_FAILED            Massage = 990080
	ERROR_DO_JENKINS_JOB                  Massage = 990081
	ERROR_IP_EXSIT                        Massage = 990082
	ERROR_VMWARE_TYPE                     Massage = 990083
	ERROR_VMNAME_EXSIT                    Massage = 990084
)