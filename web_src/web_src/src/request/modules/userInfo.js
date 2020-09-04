import axios from '../http'

// 用户信息请求
export const userInfo = () => {
    return axios({
        url: '/admin/passport/info',
        method: 'post'
    })
}

// 修改用户信息请求
export const editUserInfo = (data) => {
    return axios({
        url: '/admin/passport/modify',
        method: 'post',
        data
    })
}