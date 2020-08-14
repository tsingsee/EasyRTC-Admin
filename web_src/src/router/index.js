import Vue from 'vue'
import VueRouter from 'vue-router'
import AdminTE from '../views/common/AdminTE.vue'

Vue.use(VueRouter)
const Login = () => import('../views/login/index.vue') //登录页

const Sigin = () => import('../views/login/sigin.vue') //注册页

const MeetIndex = () => import('../views/Meet/index.vue') //会议室列表

const EditMeet=()=>import('../views/Meet/editMeet.vue') //会议室编辑


const UserInfo=()=>import('../views/UserInfo/index.vue') //用户信息模块

const MeetBlack=()=>import('../views/MeetBlack/index.vue') //会议回看列表

const MeetProc=()=>import('../views/MeetProc/index.vue') //会议记录列表

const routes = [
  {
    path: '*',
    redirect: '/login',
  },

  {
    path: '/login',
    component: Login,
  },
  {
    path: '/signin',
    component: Sigin,
  },
  {
    path: '/',
    component: AdminTE,
    name:AdminTE,
    children: [
      {
        path: '',
        redirect: '/MeetIndex'
      },
      {
        path: '/MeetIndex',
        component: MeetIndex,
        props: true
      },
      {
        path: '/EditMeet',
        component: EditMeet,
        props: true
      },
      {
        path: '/UserInfo',
        component: UserInfo,
        props: true
      },
      {
        path: '/MeetBlack',
        component: MeetBlack,
        props: true
      },
      {
        path: '/MeetProc',
        component: MeetProc,
        props: true
      },
    ]
  },

]

const router = new VueRouter({
  routes
})

export default router
