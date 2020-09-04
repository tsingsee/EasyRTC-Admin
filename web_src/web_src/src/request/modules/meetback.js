import axios from '../http'

// 获取会议录像列表
export const getVideoList = (data) => {
    return axios({
        url: '/admin/record/list',
        method: 'post',
        data
    })
}
// 录像删除
export const deleVideoL = (data) => {
    return axios({
        url: '/admin/record/delete',
        method: 'post',
        data
    })
}

