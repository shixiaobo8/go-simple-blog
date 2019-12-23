import axios, { baseUrl } from '@/libs/api.request'

export const uploadImgApi = baseUrl + '/admin/upload/image'

// 图片上传
export const uploadImg = formData => {
  return axios.request({
    url: '/admin/upload/image',
    data: formData,
    method: 'post'
  })
}

// 文章列表
export const getPostsList = data => {
  return axios.request({
    url: '/api/posts/list',
    data,
    method: 'post'
  })
}

// 文章内容
export const getPostsInfo = ({ id }) => {
  return axios.request({
    url: '/api/posts/info',
    params: {
      id
    },
    method: 'get'
  })
}

// 推荐
export const getPostsRecommend = ({ tags }) => {
  return axios.request({
    url: '/api/posts/recommend',
    data: {
      tags
    },
    method: 'post'
  })
}