import axios from '../http'

//创建房会议室
export const addMeet = data => {
    return axios({
        url: '/admin/room/create',
        method: 'post',
        data
    })
}
//获取会议室列表
export const getMeetList = data => {
    return axios({
        url: '/admin/room/list',
        method: 'post',
        data
    })
}
//获取会议室列表
export const delMeet = data => {
    return axios({
        url: '/admin/room/delete',
        method: 'post',
        data
    })
}
//获取会议室信息
export const getMeet = data => {
    return axios({
        url: '/admin/room/info',
        method: 'post',
        data
    })
}
//修改会议室信息
export const editMeet = data => {
    return axios({
        url: '/admin/room/modify',
        method: 'post',
        data
    })
}
//获取会议室主持人oken

export const getMeetToken = data => {
    return axios({
        url: '/admin/room/token',
        method: 'post',
        data
    })
}