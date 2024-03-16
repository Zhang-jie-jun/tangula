import {
  login,
  logout,
  getUserInfo,
  getUserPermission,
  getRoleName
} from '@/api/user'
import { setToken, getToken } from '@/libs/util'
import {Message} from "iview"
export default {
  state: {
    userName: '',
    userId: '',
    userBusi: '',
    avatorImgPath: '',
    token: getToken(),
    access: '',
    permits: [],
    permission: {},
    hasGetInfo: false,
    unreadCount: 0,
    messageUnreadList: [],
    messageReadedList: [],
    messageTrashList: [],
    messageContentStore: {}
  },
  mutations: {
    setAvator (state, avatorPath) {
      state.avatorImgPath = avatorPath
    },
    setUserId (state, id) {
      state.userId = id
    },
    setUserName (state, name) {
      state.userName = name
    },
    setUserBusi (state, name) {
      state.userBusi = name
    },
    setAccess (state, access) {
      console.log('setAccess:', access)
      state.access = access
    },
    setPermission (state, perm) {
      state.permission = perm
    },
    setToken (state, token) {
      console.log('setToken:', token)
      state.token = token
      setToken(token)
    },
    setHasGetInfo (state, status) {
      state.hasGetInfo = status
    },
    setMessageCount (state, count) {
      state.unreadCount = count
    },
    setMessageUnreadList (state, list) {
      state.messageUnreadList = list
    },
    setMessageReadedList (state, list) {
      state.messageReadedList = list
    },
    setMessageTrashList (state, list) {
      state.messageTrashList = list
    },
    updateMessageContentStore (state, { msg_id, content }) {
      state.messageContentStore[msg_id] = content
    },
    moveMsg (state, { from, to, msg_id }) {
      const index = state[from].findIndex(_ => _.msg_id === msg_id)
      const msgItem = state[from].splice(index, 1)[0]
      msgItem.loading = false
      state[to].unshift(msgItem)
    },
    setPermits (state, permits) {
      state.permits = permits
    }
  },
  getters: {
    messageUnreadCount: state => state.messageUnreadList.length,
    messageReadedCount: state => state.messageReadedList.length,
    messageTrashCount: state => state.messageTrashList.length
  },
  actions: {
    // 登录
    handleLogin ({ commit }, { userName, password }) {
      userName = userName.trim()
      return new Promise((resolve, reject) => {
        login({
          userName,
          password
        }).then(res => {
          const resData = res.data
          if (resData.code === 200) {
            commit('setToken', resData.response.token)
            resolve()
          } else {
            reject(resData.message)
          }
        }).catch(err => {
          reject(err)
        })
      })
    },
    // 退出登录
    handleLogOut ({ state, commit }) {
      return new Promise((resolve, reject) => {
        logout(state.token).then(() => {
          commit('setToken', '')
          commit('setAccess', [])
          resolve()
        }).catch(err => {
          reject(err)
        })
        // 如果你的退出登录无需请求接口，则可以直接使用下面三行代码而无需使用logout调用接口
        // commit('setToken', '')
        // commit('setAccess', [])
        // resolve()
      })
    },
    // 获取用户相关信息
    getUserInfo ({ state, commit }) {
      return new Promise((resolve, reject) => {
        try {
          getUserInfo().then(res => {
            console.log('getUserInfo:', res)
            const data = res.data
            if (data.response.account) {
              if (data.response.status === -1) {
                Message.error('用户已被禁用！')
                commit('setToken', '')
                commit('setAccess', [])
                // eslint-disable-next-line prefer-promise-reject-errors
                reject('用户已被禁用')
              }
              commit('setUserName', data.response.account)
              commit('setAccess', getRoleName(data.response.role_id))
            }
            // commit('setAvator', data.avator)
            // commit('setUserBusi', data.busi)
            // commit('setUserId', data.user_id)
            // commit('setAccess', data.access)
            // commit('setHasGetInfo', true)
            resolve(data)
          }).catch(err => {
            reject(err)
          })
        } catch (error) {
          reject(error)
        }
      })
    },
    // 获取用户权限
    getUserPermits ({ commit }, userName) {
      userName = userName.trim()
      return new Promise((resolve, reject) => {
        getUserPermission({
          userName
        }).then(res => {
          const data = res.data.data
          commit('setPermits', data)
          resolve()
        }).catch(err => {
          reject(err)
        })
      })
    }
  }
}
