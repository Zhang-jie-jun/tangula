import axios from '@/libs/api.request'
import { limitOffset } from '@/libs/tools.js'

export const login = ({ userName, password }) => {
  const data = {
    username: userName,
    password: password
  }
  return axios.request({
    url: '/login',
    data,
    method: 'post'
  })
}

export const getUserInfo = () => {
  return axios.request({
    url: '/auth/user',
    method: 'get'
  })
}

export const refreshToken = () => {
  return axios.request({
    url: '/auth/token',
    method: 'get'
  })
}

export const getUserList = (page, params) => {
  var lo = limitOffset(page)
  return axios.request({
    url: '/auth/users',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

export const getRoleList = (page, params) => {
  var lo = limitOffset(page)
  return axios.request({
    url: '/auth/role',
    method: 'get',
    params: Object.assign(params, { count: lo[0], index: lo[1] })
  })
}

export const disableUser = (params) => {
  let userId = params.id
  let url = '/auth/user/' + userId + '/disable'
  return axios.request({
    url: url,
    method: 'put'
  })
}

export const enableUser = (params) => {
  let userId = params.id
  let url = '/auth/user/' + userId + '/enable'
  return axios.request({
    url: url,
    method: 'put'
  })
}

export const setRole = (params) => {
  let userId = params.id
  let url = '/auth/user/' + userId + '/set_role'
  return axios.request({
    url: url,
    method: 'put',
    data: params
  })
}

export const logout = (token) => {
  return axios.request({
    url: '/logout',
    method: 'post'
  })
}

export const getUserPermission = userName => {
  return axios.request({
    url: 'auth/permission',
    method: 'post',
    data: userName

  })
}

export const getRoleName = (val) => {
  const statusNum = Number(val)
  const statusMap = {
    1: '普通用户',
    2: 'ADMIN',
    0: 'SUPER_ADMIN'
  }
  if (statusMap.hasOwnProperty(statusNum)) {
    return statusMap[statusNum]
  } else {
    return '普通用户'
  }
}
