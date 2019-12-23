import axios from 'axios'
import { getToken, setToken } from '@/libs/util'
// import store from '@/store'
import { Message } from 'iview'

axios.defaults.withCredentials = true

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
  interceptors (instance, url) {
    // 请求拦截
    instance.interceptors.request.use(config => {
      this.queue[url] = true
      return config
    }, error => {
      return Promise.reject(error)
    })
    // 响应拦截
    instance.interceptors.response.use(res => {
      this.destroy(url)
      const { data, status } = res
      return { data, status }
    }, error => {
      this.destroy(url)
      let errorInfo = error.response
      if (!errorInfo) {
        const { request: { statusText, status }, config } = JSON.parse(JSON.stringify(error))
        errorInfo = {
          statusText,
          status,
          request: { responseURL: config.url }
        }
      }

      // 拦截重定向
      if (errorInfo.status === 302 && errorInfo.data.hasOwnProperty('data') && errorInfo.data.data.hasOwnProperty('redirect_to')) {
        window.location = errorInfo.data.data.redirect_to
        return new Promise(() => {}) // 中断promise继续向下执行
      }

      if (errorInfo.status === 401) {
        if (errorInfo.data.data.hasOwnProperty('token') && errorInfo.data.data.token) {
          Message.warning('refresh token')
          // 重新设置token
          setToken(errorInfo.data.data.token)
          errorInfo.config.headers.Authorization = getToken()
          return axios(errorInfo.config)
        } else {
          Message.error('长时间没有操作，系统即将自动退出，请重新登录')
          setTimeout(() => {
            window.location = '/'
          }, 2000)
        }
        return new Promise(() => {})
      }

      return Promise.reject(error)
    })
  }
  request (options) {
    if (!options.hasOwnProperty('headers')) {
      options.headers = {}
    }
    const token = getToken()
    options.headers.Authorization = token === false ? '' : token
    const instance = axios.create()
    options = Object.assign(this.getInsideConfig(), options)
    this.interceptors(instance, options.url)
    return instance(options)
  }
}
export default HttpRequest
