import axios from '@/libs/api.request'
import { limitOffset, limitOffset50 } from '@/libs/tools.js'


//获取副本镜像类型
export const getTypes = () => {
  return axios.request({
    url: '/store_pool/replica/support',
    method: 'get'
  })
}

// 存储池列表
export const getAllPools = () => {
  return axios.request({
    url: '/resource/store_pool',
    method: 'get'
  })
}

/**
 * 创建副本
 */
export const crtBackup = (params) => {
  return axios.request({
    url: '/store_pool/replica',
    method: 'post',
    data: params
  })
}

/**
 * 编辑副本
 */
export const editBackup = (params) => {
  let hostId = params.id
  let url = '/store_pool/replica/' + hostId
  return axios.request({
    url: url,
    method: 'put',
    data: params
  })
}

// 副本列表
export const getBackupList = (page,size, params) => {
  var lo = limitOffset(page,size)
  return axios.request({
    url: '/store_pool/replica',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

/**
 * 查看副本详情
 */
export const backupDetails = (params) => {
  let url = '/store_pool/replica/' + params
  return axios.request({
    url: url,
    method: 'get'
  })
}
/**
 * 删除副本
 */
export const deleteBackup = (params) => {
  let id = params.id
  let url = '/store_pool/replica/' + id
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}

// 获取数据源
export const getResource = (params) => {
  let id = params.id
  let url = '/resource/platform/' + id + '/vmware/datasources'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

// 根据位置路径获取主机列表
export const getVMwareHostsByPath = (params) => {
  let id = params.id
  let url = '/resource/platform/' + id + '/vmware/hostsByPath'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

// 获取CAS平台主机池列表
export const getCasHostPool = (params) => {
  let id = params.id
  let url = '/resource/platform/' + id+ '/cas/hostPool'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

// 获取CAS平台主机列表
export const getCasHosts = (params) => {
  let id = params.id
  let url = '/resource/platform/' + id + '/cas/hosts'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

// 获取CAS平台主机详情
export const getCasHostInfo = (params) => {
  let id = params.id
  let url = '/resource/platform/' + id + '/cas/hostInfo'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

// 获取CAS平台虚拟机列表
export const getCasVmList = (params) => {
  let url = '/resource/platform/cas/vms'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

// 获取CAS平台存储列表
export const getCasStorageList = (params) => {
  let url = '/resource/platform/cas/storages'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

// 获取fusioncompute平台主机列表
export const getFcHosts = (params) => {
  let id = params.id
  let url = '/resource/platform/' + id + '/fusioncompute/hosts'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}



/**
 * 获取支持的副本挂载类型
 */
export const getMountType = () => {
  let url = '/store_pool/replica/mount_type'
  return axios.request({
    url: url,
    method: 'get'
  })
}

/**
 * 挂载副本
 */
export const mountBackup = (params) => {
  let url = '/store_pool/replica/mount'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}
/**
 * 批量挂载副本
 */
export const mountBackupBatch = (params) => {
  let url = '/store_pool/replica/batch/mount'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}
/**
 * 卸载副本
 */
export const unmountBackup = (params) => {
  let id = params.id
  let url = '/store_pool/replica/' + id +'/unmount'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 批量卸载副本
 */
export const unMountBackupBatch = (params) => {
  let url = '/store_pool/replica/batch/unmount'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 批量删除副本
 */
export const deleteBackupBatch = (params) => {
  let url = '/store_pool/replica/batch'
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}

// 副本挂载记录列表
export const getInstanceList = (page, params) => {
  let id = params.id
  let url = '/store_pool/replica/' + id +'/mount/instance'
  var lo = limitOffset(page)
  return axios.request({
    url: url,
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

/**
 * 副本挂载记录详情
 */
export const getInstanceDetails = (params) => {
  let id = params.id
  let url = '/store_pool/replica/mount/instance/'+ id +'/logs'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

/**
 * 生成镜像
 */
export const crtMirror = (params) => {
  let id = params.id
  let url = '/store_pool/replica/' + id + '/image'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 镜像列表
 */
export const getMirrorList = (page,size, params) => {
  var lo = limitOffset(page,size)
  return axios.request({
    url: '/store_pool/image',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

/**
 * 查看镜像详情
 */
export const mirrorDetails = (params) => {
  let url = '/store_pool/image/' + params
  return axios.request({
    url: url,
    method: 'get'
  })
}

/**
 * 发布镜像
 */
export const publishMirror = (params) => {
  let hostId = params.id
  let url = '/store_pool/image/' + hostId + '/publish'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 编辑镜像
 */
export const editImg = (params) => {
  let imgId = params.id
  let url = '/store_pool/image/' + imgId
  return axios.request({
    url: url,
    method: 'put',
    data: params
  })
}

/**
 * 生成副本
 */
export const newBackup = (params) => {
  let id = params.id
  let url = '/store_pool/image/' + id + '/replica'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 删除镜像
 */
export const deleteMirror = (params) => {
  let id = params.id
  let url = '/store_pool/image/' + id
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}

/**
 * 快照列表
 */
export const getSnapshotList = (page, params) => {
  var lo = limitOffset(page)
  let id = params.id
  let url = '/store_pool/replica/'+id+'/snapshot'
  return axios.request({
    url: url,
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

/**
 * 创建快照
 */
export const crtSnapshot = (params) => {
  let id = params.id
  let url = '/store_pool/replica/' + id + '/snapshot'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 快照生成副本
 */
export const snapshotToBackup = (params) => {
  let id = params.id
  let url = '/store_pool/snapshot/' + id + '/replica'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 快照生成镜像
 */
export const snapshotToMirror = (params) => {
  let id = params.id
  let url = '/store_pool/snapshot/' + id + '/image'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 回滚快照
 */
export const rollback = (params) => {
  let id = params.id
  let url = '/store_pool/snapshot/' + id + '/rollback'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 删除快照
 */
export const deleteSnapshot = (params) => {
  let id = params.id
  let url = '/store_pool/snapshot/' + id
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}

/**
 * 操作记录列表
 */
export const getRecordList = (page, params) => {
  var lo = limitOffset(page)
  return axios.request({
    url: '/log/record',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}


/**
 * 操作记录详情
 */
export const getRecordDetails = (params) => {
  let id = params.id
  let url = '/log/record/' + id + '/detail'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}


/**
 * 脚本列表
 */
export const getScriptList = (page, params) => {
  var lo = limitOffset(page)
  return axios.request({
    url: '/script/browse',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}
/**
 * 脚本列表(多个)
 */
export const getAllScriptList = () => {
  return axios.request({
    url: '/script/browse',
    method: 'get',
    params: Object.assign({ count: 99, index: 0 })
  })
}


/**
 * 删除脚本
 */
export const deleteScript = (params) => {
  let hostId = params.id
  let url = '/script/' + hostId+'/delete'
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}


/**
 * 脚本详情
 */
export const getScriptDetails = (params) => {
  let id = params.id
  let url = '/script/'+ id +'/content'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}
/**
 * cas配置文件详情
 */
export const getCasFileDetails = (params) => {
  let url = '/script/replica/getFile'
  const p = Object.assign(params)
  return axios.request({
    url: url,
    method: 'get',
    params: p
  })
}

/**
 * 脚本下载
 */
export const downloadScript = (params) => {
  let id = params.id
  let url = '/script/'+ id +'/download'
  return axios.request({
    url: url,
    method: 'get'
  })
}

/**
 * 执行casjenkins任务
 */
export const doCasJenkins = (params) => {
  let url = '/store_pool/replica/doCompability'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}


export const getImageStatus = (val) => {
  const statusNum = Number(val)
  const statusMap = {
    1024: '空闲',
    2048: '挂载中',
    4096: '已挂载',
    8192: '卸载中',
    16384: '克隆中'
  }
  if (statusMap.hasOwnProperty(statusNum)) {
    return statusMap[statusNum]
  } else {
    return '未挂载'
  }
}

export const choiceImageStateColor = (state) => {
  let color = 'default'
  switch (state) {
    case 1024:
      color = 'default'
      break
    case 2048:
      color = 'primary'
      break
    case 4096:
      color = 'success'
      break
    case 8192:
      color = 'primary'
      break
    case 16384:
      color = 'warning'
      break
    default:
      color = 'warning'
  }
  return color
}

export const choiceTypeColor = (state) => {
  let color
  switch (state) {
    case 1003:
      color = '#6666CC'
      break
    case 1004:
      color = '#0099CC'
      break
    case 1005:
      color = '#CC6600'
      break
    case 1006:
      color = '#F433FF'
      break
    default:
      color = ''
  }
  return color
}

export const toGB = (limit) =>{
  var size = "";
  if(limit < 0.1 * 1024){                            //小于0.1KB，则转化成B
    size = limit.toFixed(2) + "B"
  }else if(limit < 0.1 * 1024 * 1024){            //小于0.1MB，则转化成KB
    size = (limit/1024).toFixed(2) + "KB"
  }else if(limit < 0.1 * 1024 * 1024 * 1024){        //小于0.1GB，则转化成MB
    size = (limit/(1024 * 1024)).toFixed(2) + "MB"
  }else{                                            //其他转化成GB
    size = (limit/(1024 * 1024 * 1024)).toFixed(2)
  }
  return size;
}


export const getRecordStatus = (val) => {
  const statusNum = Number(val)
  const statusMap = {
    1: '成功',
    2: '失败',
    3: '执行中',
    4: '异常'
  }
  if (statusMap.hasOwnProperty(statusNum)) {
    return statusMap[statusNum]
  } else {
    return '异常'
  }
}

export const choiceRecordStateColor = (state) => {
  let color = 'warning'
  switch (state) {
    case 1:
      color = 'success'
      break
    case 2:
      color = 'error'
      break
    case 3:
      color = 'warning'
      break
    case 4:
      color = 'error'
      break
    default:
      color = 'error'
  }
  return color
}
