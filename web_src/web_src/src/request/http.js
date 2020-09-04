import axios from 'axios'
import qs from 'qs'
import config from './config'
import router from '@/router'
import store from '@/store'
import { Loading, Message } from 'element-ui'

export default function $axios(options) {
  let loading = null
  let that = this
  return new Promise((resolve, reject) => {
    const service = axios.create({
      baseURL: config.baseURL,
      headers: {
        'Content-Type': 'application/json;charset=UTF-8',
      }
    })

    // request拦截器
    service.interceptors.request.use(
      config => {
        const token = store.state.token
        // 1.loading动画
        loading = Loading.service({
          text: '正在加载中......',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        })
        // 2. 带上token
        // if (token) {
        //   config.headers.token = token
        // } else {
        //   // 重定向到登录页面
        //   router.push('/login')
        // }
        // 3. 根据请求方法，序列化传来的参数，根据后端需求是否序列化
        // if (config.method === 'post') {
        //   config.data = qs.stringify(config.data)
        // }
        return config
      },
      error => {
        console.log(error, '出错l');
        return Promise.reject(error)
      }
    )

    // response拦截器
    service.interceptors.response.use(
      response => {
         // 若未登录请求返回203就强行跳转到登录页
        if(response.status==203){
          router.push('/login')
        }
        return new Promise((resolve, reject) => {
          // 请求成功后关闭加载框
          if (loading) {
            loading.close()
          }
          const res = response.data
          if (response.status === 200) {
            resolve(res)
          } else {
            reject(res)
          }
        })
      },
      error => {
        console.log('Response Error' + error)
        // 请求成功后关闭加载框
        if (loading) {
          loading.close()
        }
        // 断网处理或者请求超时
        if (!error.response) {
          // 请求超时
          if (error.message.includes('timeout')) {
            console.log('超时了')
          } else {
            // 断网，可以展示断网组件
            console.log('断网了')
            Message.error('断网了')
          }
          return Promise.reject(error)
        }
        const status = error.response.status
        switch (status) {
          case 500:
            break
          case 404:
            break
          case 401:
            store.dispatch('delUser')
            setTimeout(() => {
              router.replace({
                path: '/login',
                query: {
                  redirect: router.currentRoute.fullPath
                }
              })
            }, 1000)
            break
          case 400:
            Message({
              message: error.response.data.error,
              type: 'error',
              duration: 5 * 1000
            })
            break
          default:
        }
        return Promise.reject(error)
      }
    )

    // 请求处理
    service(options).then(res => {
      resolve(res)
      return false
    }).catch(error => {
      reject(error)
    })
  })
}
