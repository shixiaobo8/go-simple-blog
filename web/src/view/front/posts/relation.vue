<template>
  <div class="posts-relation-wrapper">
    <Card :bordered="false">
      <p slot="title">
        <Icon type="ios-paper" />
        相关推荐
      </p>
      <template v-for="(item, index) in rows">
        <p class="recommend-list" :key="index" @click="readPosts(item.id)">
          <span>{{item.title}}</span>
          <span></span>
        </p>
      </template>
    </Card>
  </div>
</template>

<script>
import { getPostsRecommend } from '../../../api/data'
export default {
  name: 'Posts',
  props: {
    tags: {
      type: Array,
      default: () => []
    },
    rec: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      rows: []
    }
  },
  created() {},
  methods: {
    readPosts(id) {
      this.$router.push({ name: 'posts', params: { id } })
    },
    fetchData() {
      getPostsRecommend({ tags: this.tags })
        .then(result => {
          if (result.data.code === 200) {
            for (let i = 0; i < result.data.data.rows.length; i++) {
              this.rows.push(result.data.data.rows[i])
            }
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
    rec(val) {
      if (val === true) {
        this.fetchData()
      }
    }
  }
}
</script>

<style lang="less" scoped>
.recommend-list {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin: 2px 0;
  cursor: pointer;
  &:hover {
    text-decoration: underline;
    background: rgba(0, 154, 97, 0.08);
  }
}
</style>