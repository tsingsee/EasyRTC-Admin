import axios from '../http'

/*
 * 系统登录模块
 */
//获取验证码id
export const getCaptchaId = data => {
  return axios({
    url: '/admin/captcha-id',
    method: 'post',
    data
  })
}
// 注册
export const sigin = data => {
  return axios({
    url: '/admin/passport/signup',
    method: 'post',
    data
  })
}

// 登录
export const login = data => {
  return axios({
    url: '/admin/passport/login',
    method: 'post',
    data
  })
}

// 登出
export const logout = () => {
  return axios({
    url: '/admin/passport/logout',
    method: 'post'
  })
}

