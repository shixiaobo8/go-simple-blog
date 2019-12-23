import axios from '@/libs/api.request'

// 获取后台用户列表
export const list = () => {
  return axios.request({
    url: '/admin/adminUser/list',
    method: 'post'
  })
}

// 更新后台用户信息
export const updateUser = info => {
  return axios.request({
    url: '/admin/adminUser/updateUser',
    data: info,
    method: 'post'
  })
}
