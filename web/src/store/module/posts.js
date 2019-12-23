import { info } from '../../api/backend/posts'
import { deepCopy } from '../../libs/tools'

export default {
  namespaced: true, // 开启namespace
  state: {
    form: {}
  },
  getters: {
    getForm(state) {
      return (id) => {
        return state.form[id]
      }
    }
  },
  mutations: {
    setForm(state, payload) {
      state.form[payload.id] = deepCopy(payload.form)
    },
    removeForm(state, payload) {
      delete state.form[payload.id]
    }
  },
  actions: { // action提交mutation，action可以包含异步操作，而mutation必须是同步
    fetchData({ commit, state }, id) { // action函数接受一个与store实例具有相同方法和属性的context对象 因此你可以调用 context.commit 提交一个 mutation
      /**
       * 对于模块内部的 action，局部状态通过 context.state 暴露出来，根节点状态则为 context.rootStat
       * 所以context并不是store实例，而是store下期中一个模块实例
       */
      let form = {}
      if (id) {
        return new Promise((resolve, reject) => {
          info(id).then(r => {
            if (r.data.code === 200) {
              form.id = r.data.data.id
              form.title = r.data.data.title
              form.content = r.data.data.content
              form.cover = r.data.data.cover
              form.tags = r.data.data.tags ? r.data.data.tags.join(' ') : ''

              commit('setForm', { id, form })
            }
            resolve(r.data)
          }).catch(e => {
            reject(e)
          })
        })
      } else {
        if (!state.form.hasOwnProperty(id)) {
          let form = {
            id: id,
            title: '',
            content: '',
            cover: '',
            tags: ''
          }
          commit('setForm', { id, form })
        }
        return Promise.resolve()
      }
    },
    actionSetForm({ commit }, { id, form }) {
      commit('setForm', { id, form })
    },
    actionRemoveForm({ commit }, { id }) {
      commit('removeForm', { id })
    }
  }
}
