<template>
  <div class="posts-read-wrapper">
    <Row :gutter="20" type="flex" justify="center">
    <Col :xs="0" :sm="0" :md="0" :lg="{span: 6}">
      <Relation :tags="posts.tags" :rec="rec"></Relation>
    </Col>
    <Col :xs="24" :sm="24" :md="{span: 16}" :lg="{span: 16}">
      <Card shadow>
          <p slot="title">{{ posts.title }}</p>
          <div v-html="posts.content" class="posts-content"></div>
      </Card>
    </Col>
    </Row>
  </div>
</template>

<script>
import { getPostsInfo } from '../../../api/data'
import Relation from '../../front/posts/relation'

export default {
  name: 'Posts',
  components: {
    Relation
  },
  data() {
    return {
      posts: {
        tags: []
      },
      rec: false // 控制在当前组件异步请求成功后，realtion组件再进行请求
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      let d = {
        id: isNaN(this.$route.params.id) ? 0 : this.$route.params.id
      }
      if (d.id === 0) {
        return
      }

      getPostsInfo(d)
        .then(result => {
          if (result.data.code === 200) {
            this.posts = result.data.data
            this.rec = true
          } else {
            this.$Message['warning']({
              background: true,
              content: result.data.msg
            })
          }
        })
        .catch(err => {
          this.$Message['error']({
            background: true,
            content: err
          })
        })
    }
  },
  watch: {
    $route(to, from) {
      if (to.fullPath !== from.fullPath) {
        this.fetchData()
      }
    }
  }
}
</script>

<style lang="less">
  .posts-content {
    line-height: 1.5rem;
    img {
      max-width: 100%;
    }
    p {
      text-indent: 2em;
    }
  }
</style>