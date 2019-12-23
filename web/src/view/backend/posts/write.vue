<template>
  <div>
    <Row>
      <Col :md="20" :lg="20" :xs="24">
        <Form :model="form" :label-width="120" :rules="rules">
          <FormItem label="标题" prop="title">
            <Input v-model="form.title" placeholder="输入标题" />
          </FormItem>
          <FormItem label="标签" prop="tags">
            <Input v-model="form.tags" placeholder="输入标签，空格分开" />
          </FormItem>
          <FormItem label="封面" props="cover">
            <ImgUpload
              ref="imgupload"
              :default-list="defaultList"
              :key="$route.query.id === undefined ? 0 : $route.query.id"
            />
          </FormItem>
          <FormItem label="内容" props="content">
            <Editor
              ref="editor"
              :content="form.content"
              :key="$route.query.id === undefined ? 0 : $route.query.id"
            />
          </FormItem>
          <FormItem label>
            <Button @click="handleSubmit" type="primary" style="float:right;">发表文章</Button>
          </FormItem>
        </Form>
      </Col>
    </Row>
  </div>
</template>

<script>
import Editor from '_c/editor'
import { update as updatePosts } from '../../../api/backend/posts'

import ImgUpload from '../components/ImgUpload'
import { urlFileName } from '../../../libs/tools'
import { mapMutations, mapActions, mapGetters } from 'vuex'

// mapActions辅助函数可以将methods映射为store.dispath调用（需要在根结点注入store）
/* 比如  ...mapActions([
      'increment', // 将 `this.increment()` 映射为 `this.$store.dispatch('increment')`
    ]),
*/

export default {
  name: 'admin-posts-write',
  components: {
    Editor,
    ImgUpload
  },
  data() {
    return {
      form: {
        id: 0,
        title: '',
        content: '',
        cover: '',
        tags: ''
      }
    }
  },
  methods: {
    ...mapMutations(['closeTag']),
    ...mapActions('posts', ['fetchData', 'actionRemoveForm']),
    handleFetchData() {
      this.fetchData(this.id)
        .then(r => {
          this.form = Object.assign({}, this.getForm(this.id))
          this.$nextTick(() => {
            this.$refs.imgupload.handleDefaultList()
          })
        })
        .catch(e => {
          this.$Message['error']({
            background: true,
            content: '数据请求失败'
          })
        })
    },
    handleSubmit() {
      this.writedContent()
      this.uploadedCover()
      updatePosts(this.form)
        .then(r => {
          if (r.data.code === 200) {
            this.$Message['success']({
              background: true,
              content: r.data.msg
            })
            this.closeTag(this.$route)
            this.$router.push({ name: 'admin-posts-list' })
          } else {
            this.$Message['error']({
              background: true,
              content: r.data.msg
            })
          }
        })
        .catch(e => {
          this.$Message['error']({
            background: true,
            content: '服务器连接失败'
          })
        })
    },
    writedContent() {
      this.form.content = this.$refs.editor.text
    },
    uploadedCover() {
      if (this.$refs.imgupload.uploadList.length > 0) {
        this.form.cover = this.$refs.imgupload.uploadList[0].url
      } else {
        this.form.cover = ''
      }
    },
    metaTitle(route) {
      if (route.query.id) {
        route.meta.title = '编辑文章'
      } else {
        route.meta.title = '发表文章'
      }
      this.$store.dispatch('tagTitle', route) // 同步修改tag显示标题
    }
  },
  computed: {
    ...mapGetters(['tagListExists']),
    ...mapGetters('posts', ['getForm']),
    defaultList() {
      return this.form.cover !== '' && this.form.cover
        ? [
          {
            name: urlFileName(this.form.cover),
            url: this.form.cover
          }
        ]
        : []
    },
    id() {
      let id = parseInt(this.$route.query.id)
      return isNaN(id) ? 0 : id
    },
    rules() {
      return {
        title: {
          required: true,
          type: 'string',
          message: '标题不能为空',
          trigger: 'blur'
        }
      }
    }
  },
  created() {
    this.handleFetchData()
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.metaTitle(to)
      if (to.query.id === undefined && !vm.tagListExists(vm.$route)) {
        this.actionRemoveForm(0)
      }
    })
  },
  beforeRouteUpdate(to, from, next) {
    // https://router.vuejs.org/zh/guide/advanced/navigation-guards.html#组件内的守卫
    this.metaTitle(to)
    next()
  },
  beforeRouteLeave(to, from, next) {
    next()
  },
  watch: {
    // 解决页面不刷新, 虽然有beforeRouteUpdate，但是它仅限于route的fullpath一样的情况，现在编辑和发表组件复用可以实现刷新
    // 如果打开一个发表页，然后从列表页打开编辑页，此时由于组件复用，并不会重新请求数据，需要
    $route(to, from) {
      if (to.fullPath !== from.fullPath) {
        this.handleFetchData()
      }

      // 用户如果已经修改表单，然后切换了路由，再重新切换回来需要保留表单信息
      if (
        (from.query.id === 0 || from.query.id === undefined) &&
        (to.query.id !== 0 || to.query.id !== undefined)
      ) {
        // 用户编辑之后需要保留上一次的编辑信息
        this.writedContent()
        this.uploadedCover()
        this.$store.dispatch('posts/actionSetForm', { id: 0, form: this.form })
      }
    }
  }
}
</script>
