import axios from 'axios'
import { StorageGetItem } from "@/utils/storage";

// 创建axios实例
let service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // api的base_url
  validateStatus: function (status) {
    return (status >= 200 && status < 300) || status === 401;
  },
  timeout: 30000, // 请求超时时间30s
})

// request拦截器
service.interceptors.request.use(config => {
  const token = StorageGetItem('token') || ""
  if (token != "") {
    config.headers.Authorization = encodeURIComponent(`${token}`)
  }
  return config
}, error => {
  Promise.reject(error)
})

// respone拦截器
service.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    var message = error.message
    if (error.response) {
      message = error.response.data
    }
    return Promise.reject(message)
  }
)

export default service
