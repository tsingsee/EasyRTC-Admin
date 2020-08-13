import axios from '../http'

/*
 * 系统登录模块
 */
//获取验证码id
export const getrecordList= data => {
  return axios({
    url: '/admin/conference/history',
    method: 'post',
    data
  })
}