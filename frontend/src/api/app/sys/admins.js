import request from '@/utils/request'

export function requestList(query) {
  return request({
    url: '/user/list',
    method: 'get',
    params: query
  })
}

export function requestDetail(id) {
  return request({
    url: '/user/detail',
    method: 'get',
    params: { id }
  })
}

export function requestUpdate(data) {
  return request({
    url: '/user/update',
    method: 'post',
    data
  })
}

export function requestCreate(data) {
  return request({
    url: '/user/create',
    method: 'post',
    data
  })
}

export function requestDelete(data) {
  return request({
    url: '/user/delete',
    method: 'post',
    data
  })
}

export function requestAdminsRoleIDList(adminsid) {
  return request({
    url: '/user/adminsroleidlist',
    method: 'get',
    params: { adminsid }
  })
}

export function requestSetRole(adminsid, data) {
  return request({
    url: '/user/setrole',
    method: 'post',
    params: { adminsid },
    data
  })
}

