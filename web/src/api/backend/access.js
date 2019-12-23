import axios, { baseUrl } from '@/libs/api.request'
import md5 from 'js-md5'

// 获取验证码
export const captCha = () => {
  return new Promise((resolve, reject) => {
    axios.request({
      url: '/admin/access/captcha',
      method: 'post'
    }).then(val => {
      val.data.data.cpt = baseUrl + '/admin/access/captcha' + '/' + val.data.data.cpt
      resolve(val)
    }).catch(err => {
      reject(err)
    })
  })
}

// 登录
export const login = ({ userName, password, cptId, cptCode }) => {
  password = md5(password)
  const data = {
    userName,
    password,
    cptId,
    cptCode
  }
  return axios.request({
    url: '/admin/access/login',
    data,
    method: 'post'
  })
}

// 获取用户信息
export const getUserInfo = (id) => {
  return axios.request({
    url: '/admin/adminUser/userInfo',
    params: {
      id
    },
    method: 'get'
  })
}

// 退出
export const logout = (token) => {
  return axios.request({
    url: '/admin/access/logout',
    method: 'post'
  })
}
