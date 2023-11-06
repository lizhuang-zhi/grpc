import request from '@/utils/axios'

export function form(params) {
  return request({
    url: '/form',
    method: 'post',
    data: params
  })
}

export function list(params) {
  return request({
    url: '/list',
    method: 'get',
    params: params
  })
}

export function query(params) {
  return request({
    url: '/query',
    method: 'get',
    params: params
  })
}