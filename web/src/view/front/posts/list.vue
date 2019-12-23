<template>
  <div class="posts-list-wrapper">
    <List item-layout="vertical">
      <ListItem
        v-for="(item, index) in list"
        :key="index"
        @mouseover.native="itemMouseOver(index)"
        @mouseleave.native="itemMouseLeave(index)"
        :class="{'item-active':index === active}"
        class="posts-item"
      >
        <ListItemMeta
          :title="item.title"
          :description="item.description"
          style="cursor: pointer;"
          @click.native="readPosts(item.id)"
        />
        {{ item.content }}
        <template slot="action">
          <li>
            <Icon type="ios-thumbs-up-outline" />234
          </li>
          <li>
            <Icon type="ios-chatbubbles-outline" />345
          </li>
          <li>
            <Icon type="ios-chatbubbles-outline" />
            {{ formatTime(item.created_at) }}
          </li>
        </template>
        <template slot="extra">
          <img :src="item.cover" style="width: 12rem;height: 12rem;" />
        </template>
      </ListItem>
    </List>

    <Row>
      <Col class="spin-col" span="24" v-if="spinShow">
        <Spin fix>
          <Icon type="ios-loading" size="36" class="spin-icon-load"></Icon>
          <div>正在载入...</div>
        </Spin>
      </Col>
      <div class="scroll-tips" v-if="!spinShow">
        <p>
          <Icon type="md-sad" />没有更多啦!
        </p>
      </div>
    </Row>
  </div>
</template>

<script>
import { getPostsList } from '../../../api/data'
import { formatDate } from '../../../libs/date'

export default {
  name: 'PostsList',
  props: {
    tag: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      page: 1,
      page_size: 5,
      searchType: 0,
      searchValues: [
        { label: '标题', value: 1, key: 'title', type: String },
        { label: '标签', value: 2, key: 'tag', type: String }
      ],
      list: [],
      timer: null,
      spinShow: false,
      active: -1
    }
  },
  created() {
    this.fetchData()
    this.handleReachBottom()
  },
  methods: {
    formatTime(t) {
      return formatDate(new Date(t))
    },
    fetchData(p) {
      let d = {
        page: isNaN(p) ? 1 : p,
        page_size: this.page_size,
        search: this.searchParams()
      }

      getPostsList(d)
        .then(result => {
          if (result.data.code === 200) {
            if (result.data.data.rows.length <= this.page_size) {
              this.spinShow = false
              this.removeScrollHandler()
            }
            for (let i = 0; i < result.data.data.rows.length; i++) {
              this.list.push(result.data.data.rows[i])
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
    },
    searchParams() {
      let search = {}
      search.tag = this.tag
      return search
    },
    scrollHandler() {
      clearTimeout(this.timer)
      this.timer = setTimeout(() => {
        // 变量scrollTop是滚动条滚动时，距离顶部的距离
        const scrollTop =
          document.documentElement.scrollTop || document.body.scrollTop
        // 变量windowHeight是可视区的高度
        const windowHeight =
          document.documentElement.clientHeight || document.body.clientHeight
        // 变量scrollHeight是滚动条的总高度
        const scrollHeight =
          document.documentElement.scrollHeight || document.body.scrollHeight

        if (scrollTop + windowHeight >= scrollHeight) {
          this.spinShow = true
          this.fetchData(this.page)
          this.page++
        }
      }, 1000)
    },
    removeScrollHandler() {
      if (this.timer !== null) {
        clearTimeout(this.timer)
      }
      window.removeEventListener('scroll', this.scrollHandler, true)
    },
    handleReachBottom() {
      window.addEventListener('scroll', this.scrollHandler, true)
    },
    readPosts(id) {
      this.$router.push({ name: 'posts', params: { id } })
    },
    itemMouseOver(index) {
      this.active = index
    },
    itemMouseLeave(index) {
      this.active = -1
    }
  },
  beforeDestroy() {
    this.removeScrollHandler()
  },
  watch: {
    tag(val) {
      this.list = []
      this.page = 1
      this.fetchData(1)
    }
  }
}
</script>

<style lang="less" scoped>
.spin-col {
  height: 160px;
  position: relative;
  border: none;
}
.spin-icon-load {
  animation: ani-demo-spin 1s linear infinite;
}
@keyframes ani-demo-spin {
  from {
    transform: rotate(0deg);
  }
  50% {
    transform: rotate(180deg);
  }
  to {
    transform: rotate(360deg);
  }
}
.scroll-tips {
  text-align: center;
  font-size: 14px;
}
.item-active {
  background-color: #ccc;
}
.posts-item {
  padding: 10px;
  background-color: white;
  transition: 0.4s all;
  border-bottom: none;
  margin: 5px 0;
  border-radius: 2px;
  box-sizing: border-box;
  box-shadow: 0px 1px 3px rgba(0, 0, 0, 0.2);
  &:hover {
    box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.3);
  }
}
</style>