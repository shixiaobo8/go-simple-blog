import {
  captCha,
  login,
  logout,
  getUserInfo
} from '@/api/backend/access'
import { setToken, getToken } from '@/libs/util'

export default {
  state: {
    userName: '',
    userId: '',
    avatarImgPath: '',
    token: getToken(),
    access: '',
    hasGetInfo: false,
    unreadCount: 0,
    messageUnreadList: [],
    messageReadedList: [],
    messageTrashList: [],
    messageContentStore: {},
    cpt: {
      cpt: '',
      id: ''
    } // 验证码
  },
  mutations: {
    setCpt (state, cpt) {
      state.cpt = Object.assign(state.cpt, cpt)
    },
    setAvatar (state, avatarPath) {
      state.avatarImgPath = avatarPath
    },
    setUserId (state, id) {
      state.userId = id
    },
    setUserName (state, name) {
      state.userName = name
    },
    setAccess (state, access) {
      state.access = access
    },
    setToken (state, token) {
      state.token = token
      setToken(token)
    },
    setHasGetInfo (state, status) {
      state.hasGetInfo = status
    }
  },
  getters: {
  },
  actions: {
    getCaptcha ({ commit }) {
      return new Promise((resolve, reject) => {
        captCha().then(res => {
          const r = res.data
          commit('setCpt', r.data)
          resolve(r.data)
        }).catch(err => {
          reject(err)
        })
      })
    },
    // 登录
    handleLogin ({ commit }, { userName, password, cptId, cptCode }) {
      userName = userName.trim()
      cptCode = cptCode.trim()

      return new Promise((resolve, reject) => {
        login({
          userName,
          password,
          cptId,
          cptCode
        }).then(res => {
          const data = res.data
          commit('setToken', data.data.token)
          resolve(data)
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
      })
    },
    // 获取用户相关信息
    getUserInfo ({ state, commit }, id) {
      return new Promise((resolve, reject) => {
        try {
          getUserInfo(id).then(res => {
            const rd = res.data
            commit('setAvatar', rd.data.avatar)
            commit('setUserName', rd.data.nickname)
            commit('setUserId', rd.data.id)
            commit('setHasGetInfo', true)
            resolve(rd)
          }).catch(err => {
            reject(err)
          })
        } catch (error) {
          reject(error)
        }
      })
    }
  }
}
