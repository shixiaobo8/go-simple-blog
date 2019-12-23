import axios from '@/libs/api.request'

// 获取文章列表
export const list = (data) => {
  return axios.request({
    url: '/admin/posts/list',
    data,
    method: 'post'
  })
}

// 获取单篇文章信息
export const info = (id) => {
  return axios.request({
    url: '/admin/posts/info',
    params: {
      id
    },
    method: 'get'
  })
}

// 发表更新文章
export const update = (data) => {
  return axios.request({
    url: '/admin/posts/update',
    data: data,
    method: 'post'
  })
}

// 删除
export const del = (id) => {
  return axios.request({
    url: '/admin/posts/delete',
    data: {
      id
    },
    method: 'post'
  })
}

// 筛选搜索项
export const filter = ({ query }) => {
  return axios.request({
    url: '/admin/posts/filterTag',
    data: { query },
    method: 'post'
  })
}