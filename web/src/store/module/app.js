import {
  getBreadCrumbList,
  setTagNavListInLocalstorage,
  getMenuByRouter,
  getTagNavListFromLocalstorage,
  getHomeRoute,
  getNextRoute,
  routeHasExist,
  routeEqual,
  getRouteTitleHandled,
  localSave,
  localRead
} from '@/libs/util'
import router from '@/router'
import backendRouters from '@/router/backend'
import config from '@/config'
const { adminHomeName } = config

const closePage = (state, route) => {
  const nextRoute = getNextRoute(state.tagNavList, route)
  state.tagNavList = state.tagNavList.filter(item => {
    return !routeEqual(item, route)
  })
  router.push(nextRoute)
}

export default {
  state: {
    breadCrumbList: [],
    tagNavList: [],
    homeRoute: {},
    local: localRead('local')
  },
  getters: {
    // 只对后端router进行处理
    menuList: (state, getters, rootState) => getMenuByRouter(backendRouters, rootState.user.access),
    tagListExists: (state) => (route) => { // 检查是否存在某个tag
      let tag = state.tagNavList.filter(item => routeEqual(item, route))
      return !!tag[0]
    }
  },
  mutations: {
    setBreadCrumb (state, route) {
      state.breadCrumbList = getBreadCrumbList(route, state.homeRoute)
    },
    setHomeRoute (state, routes) {
      state.homeRoute = getHomeRoute(routes, adminHomeName)
    },
    setTagNavList (state, list) {
      let tagList = []
      if (list) {
        tagList = [...list]
      } else tagList = getTagNavListFromLocalstorage() || []
      if (tagList[0] && tagList[0].name !== adminHomeName) tagList.shift()
      let homeTagIndex = tagList.findIndex(item => item.name === adminHomeName)
      if (homeTagIndex > 0) {
        let homeTag = tagList.splice(homeTagIndex, 1)[0]
        tagList.unshift(homeTag)
      }
      state.tagNavList = tagList
      setTagNavListInLocalstorage([...tagList])
    },
    closeTag (state, route) {
      let tag = state.tagNavList.filter(item => routeEqual(item, route))
      route = tag[0] ? tag[0] : null
      if (!route) return
      closePage(state, route)
    },
    addTag (state, { route, type = 'unshift' }) {
      let router = getRouteTitleHandled(route)
      if (!routeHasExist(state.tagNavList, router)) {
        if (type === 'push') state.tagNavList.push(router)
        else {
          if (router.name === adminHomeName) state.tagNavList.unshift(router)
          else state.tagNavList.splice(1, 0, router)
        }
        setTagNavListInLocalstorage([...state.tagNavList])
      }
    },
    setLocal (state, lang) {
      localSave('local', lang)
      state.local = lang
    },
    setTagTitle (state, route) { // 在组件中对route.meta.title信息修改后需要同步到tag标题
      if (route && route.meta && route.meta.title) {
        let title = route.meta.title
        let tag = state.tagNavList.filter(item => routeEqual(item, route))
        route = tag[0] ? tag[0] : null
        if (route !== null) {
          route.meta.title = title
        }
      }
    }
  },
  actions: {
    tagTitle({ commit }, route) {
      commit('setTagTitle', route)
    }
  }
}
