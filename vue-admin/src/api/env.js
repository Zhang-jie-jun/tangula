import axios from '@/libs/api.request'
import { limitOffset, limitOffset50 } from '@/libs/tools.js'

// 查询主机列表
export const getHostList = (page, params) => {
  var lo = limitOffset(page)
  return axios.request({
    url: '/resource/host',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

export const getHostListAll = (params) => {

  return axios.request({
    url: '/resource/host',
    method: 'get',
    params: Object.assign(params, { count: 50, index: 0 })
  })
}
// 获取主机资源类型
export const getHostTypes = () => {
  return axios.request({
    url: '/resource/host/support',
    method: 'get'
  })
}

// 查看主机详情
export const getHostDetails = (params) => {
  let url = '/resource/host/' + params
  return axios.request({
    url: url,
    method: 'get'
  })
}

/**
 * 创建主机
 */
export const crtHost = (params) => {
  return axios.request({
    url: '/resource/host',
    method: 'post',
    data: params
  })
}

/**
 * 编辑主机
 */
export const editHost = (params) => {
  let hostId = params.id
  let url = '/resource/host/' + hostId
  return axios.request({
    url: url,
    method: 'put',
    data: params
  })
}

/**
 * 更新主机
 */
export const updateHost = (params) => {
  let hostId = params.id
  let url = '/resource/host/' + hostId + '/update'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 发布主机
 */
export const publishHost = (params) => {
  let hostId = params.id
  let url = '/resource/host/' + hostId + '/publish'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 删除主机
 */
export const deletehHost = (params) => {
  let hostId = params.id
  let url = '/resource/host/' + hostId
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}
/**
 * 部署客户端
 */
export const deployClient = (params) => {
  let url = '/resource/host/deployClient'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 查看部署记录
 */
export const getDeployRecord = (page, params) => {
  let hostId = params.id
  let url = '/resource/host/deployRecord/' + hostId
  var lo = limitOffset(page)
  return axios.request({
    url: url,
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

// 获取平台资源类型
export const getPlatformTypes = () => {
  return axios.request({
    url: '/resource/platform/support',
    method: 'get'
  })
}

// 查询平台列表
export const getPlatList = (page, params) => {
  var lo = limitOffset(page)
  return axios.request({
    url: '/resource/platform',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

export const getPlatListAll = (params) => {
  return axios.request({
    url: '/resource/platform',
    method: 'get',
    params: Object.assign(params, { count: 50, index: 0 })
  })
}

/**
 * 创建平台
 */
export const crtPlat = (params) => {
  return axios.request({
    url: '/resource/platform',
    method: 'post',
    data: params
  })
}

//查看平台详情
export const getPlatDetails = (params) => {
  let url = '/resource/platform/' + params
  return axios.request({
    url: url,
    method: 'get'
  })
}


/**
 * 编辑平台
 */
export const editPlat = (params) => {
  let hostId = params.id
  let url = '/resource/platform/' + hostId
  return axios.request({
    url: url,
    method: 'put',
    data: params
  })
}

/**
 * 发布平台
 */
export const publishPlat = (params) => {
  let hostId = params.id
  let url = '/resource/platform/' + hostId + '/publish'
  return axios.request({
    url: url,
    method: 'post',
    data: params
  })
}

/**
 * 删除平台
 */
export const deletePlat = (params) => {
  let hostId = params.id
  let url = '/resource/platform/' + hostId
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}

/**
 * 创建存储池
 */
export const crtPool = (params) => {
  return axios.request({
    url: '/resource/store_pool',
    method: 'post',
    data: params
  })
}

// 存储池列表
export const getPoolList = (page, params) => {
  var lo = limitOffset(page)
  return axios.request({
    url: '/resource/store_pool',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

/**
 * 存储池详情
 */
export const getPoolDetails = (params) => {
  let url = '/resource/store_pool/' + params
  return axios.request({
    url: url,
    method: 'get'
  })
}

/**
 * 删除存储池
 */
export const deletePool = (params) => {
  let hostId = params.id
  let url = '/resource/store_pool/' + hostId
  return axios.request({
    url: url,
    method: 'delete',
    data: params
  })
}

/**
 * ceph信息
 */
export const getCephInfo = () => {
  let url = '/server/ceph'
  return axios.request({
    url: url,
    method: 'get'
  })
}

/**
 * 服务器信息
 */
export const getSystemInfo = () => {
  let url = '/server/system'
  return axios.request({
    url: url,
    method: 'get'
  })
}

/**
 * 平台信息
 */
export const getDashboardInfo = () => {
  let url = '/server/dashboard'
  return axios.request({
    url: url,
    method: 'get'
  })
}
export const getAllJson = (jsons, name, sign) => {
  const result = []
  if (name === '' || name === undefined) {
    name = 'json'
  }
  for (let key in jsons) {
    var k = name + sign + key
    if (!(jsons[key] instanceof Object)) {
      console.log(k + ' = ' + jsons[key]) // 如果不是Object则打印键值
    } else {
      getAllJson(jsons[key], k, sign) // 如果是Object则递归
    }
  }
  return result
}

export const jsonToTree = (item) => {
  const result = []
  for (let key in item) {
    // 读取 map 的键值映射
    let children = item[key]
    // 如果有子节点，递归
    if (item[key] instanceof Object) {
      let childrenTmp = jsonToTree(children)
      result.push({
        name: key,
        children: childrenTmp
      })
    } else {
      result.push({
        name: key,
        value: item[key]
      })
    }
  }
  return result
}

export const getDeployStatus = (val) => {
  const statusNum = Number(val)
  const statusMap = {
    1: '未部署',
    2: '部署中',
    3: '成功',
    4: '失败'
  }
  if (statusMap.hasOwnProperty(statusNum)) {
    return statusMap[statusNum]
  } else {
    return '未部署'
  }
}

export const choiceDeployStateColor = (state) => {
  let color = 'default'
  switch (state) {
    case 1:
      color = 'default'
      break
    case 2:
      color = 'warning'
      break
    case 3:
      color = 'success'
      break
    case 4:
      color = 'error'
      break
    default:
      color = 'default'
  }
  return color
}
