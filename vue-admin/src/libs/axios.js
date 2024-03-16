import axios from 'axios'
import { getToken, setToken } from '@/libs/util'
import { refreshToken } from '@/api/user'
import {Message} from "iview"

class HttpRequest {
  constructor (baseUrl = baseURL) {
    this.baseUrl = baseUrl
    this.queue = {}
  }
  getInsideConfig () {
    const config = {
      baseURL: this.baseUrl,
      headers: {
        //
      }
    }
    return config
  }
  destroy (url) {
    delete this.queue[url]
    if (!Object.keys(this.queue).length) {
      // Spin.hide()
    }
  }
  request (options) {
    const instance = axios.create()
    options = Object.assign(this.getInsideConfig(), options)
    // 请求判断
    instance.interceptors.request.use(config => {
      let token = getToken()
      config.headers.Authorization = 'Tangula ' + token
      return config
    }, error => {
      Message.error('请求head有误')
      return Promise.reject(error)
    })
    // 返回状态判断
    instance.interceptors.response.use(res => {
      // 刷新token
      if (options.url !== '/login' && options.url !== '/logout' && options.url !== '/auth/token') {
        refreshToken().then(res => {
          let newToken = res.data.response.token
          setToken(newToken)
        }).catch(err => {
          Message.error('token已失效:' + err)
          setToken('')
          location.href = '/login'
        })
      }
      return res
    }, error => {
      console.log('接口返回报错', error)
      if (options.url !== '/login' && options.url !== '/logout' && error.response.status === 401) {
        // 401 说明 token 验证失败
        // 可以直接跳转到登录页面，重新登录获取 token
        Message.error('token已失效')
        setToken('')
        location.href = '/login'
      } else if (error.response.status === 500) {
        // 服务器错误
        Message.error('服务端错误')
        return Promise.reject(error.response.data)
      }else if (error.response.status === 404) {
        location.href = '/404'
        return Promise.reject(error.response.data)
      }
      // 返回 response 里的错误信息
      return Promise.reject(error.response.data)
    })
    return instance(options)
  }
}
export default HttpRequest
