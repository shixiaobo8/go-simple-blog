import Vue from 'vue'
import Router from 'vue-router'
import routes from './routers'
import store from '@/store'
import iView from 'iview'
import { setToken, getToken, canTurnTo, setTitle } from '@/libs/util'
import config from '@/config'
const { adminHomeName } = config

// 解决NavigationDuplicated
const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}

// 使用vue-router
Vue.use(Router)
const router = new Router({
  routes,
  mode: 'history'
})

// 后台入口路由名
const ADMIN_LOGIN_PAGE_NAME = 'admin-access-login'

// 根据权限进行跳转
const turnTo = (to, access, next) => {
  if (canTurnTo(to.name, access, routes)) next() // 有权限，可访问
  else next({ replace: true, name: 'error_401' }) // 无权限，重定向到401页面
}

/*
  导航守卫-全局
  参数或查询的改变并不会触发进入/离开的导航守卫。你可以通过观察 $route 对象来应对这些变化
  或使用 beforeRouteUpdate 的组件内守卫
  @see https://router.vuejs.org/zh/guide/advanced/navigation-guards.html#%E5%85%A8%E5%B1%80%E5%89%8D%E7%BD%AE%E5%AE%88%E5%8D%AB
*/
router.beforeEach((to, from, next) => {
  iView.LoadingBar.start()
  const token = getToken()

  if (to.meta.requireAuth) { // 表示后台需要验证权限
    if (!token && to.name !== ADMIN_LOGIN_PAGE_NAME) {
      // 未登录且要跳转的页面不是登录页
      next({
        name: ADMIN_LOGIN_PAGE_NAME // 跳转到登录页
      })
    } else if (!token && to.name === ADMIN_LOGIN_PAGE_NAME) {
      // 未登陆且要跳转的页面是登录页
      next() // 跳转
    } else if (token && to.name === ADMIN_LOGIN_PAGE_NAME) {
      // 已登录且要跳转的页面是登录页
      next({
        name: adminHomeName
      })
    } else {
      if (store.state.user.hasGetInfo) {
        turnTo(to, store.state.user.access, next)
      } else {
        store.dispatch('getUserInfo').then(user => {
          // 拉取用户信息，通过用户权限和跳转的页面的name来判断是否有权限访问;access必须是一个数组，如：['super_admin'] ['super_admin', 'admin']
          turnTo(to, user.access, next)
        }).catch(() => {
          setToken('')
          next({
            name: 'login'
          })
        })
      }
    }
  } else {
    next() // 跳转
  }
})

router.afterEach(to => {
  setTitle(to, router.app)
  iView.LoadingBar.finish()
  window.scrollTo(0, 0)
})

export default router
