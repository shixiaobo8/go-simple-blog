<template>
  <div class="posts-list-wrapper">
    <Row class="search-wrapper">
      <Col span="4" style="padding-right:5px">
        <Select v-model="searchType" size="large">
          <Option
            v-for="(option, index) in searchValues"
            :value="option.value"
            :key="index"
          >{{option.label}}</Option>
        </Select>
      </Col>
      <Col span="6" style="padding-right:5px">
        <Input v-if="searchType === 1" v-model="searchModel" size="large" placeholder="请输入标题" />
        <Select
          v-else-if="searchType === 2"
          v-model="searchModel"
          size="large"
          multiple
          filterable
          remote
          :remote-method="remoteFilter"
          :loading="loading"
        >
          <Option
            v-for="(option, index) in options"
            :value="option.value"
            :key="index"
          >{{option.label}}</Option>
        </Select>
      </Col>
      <Col span="1">
        <Button type="primary" icon="ios-search" size="large" @click="handleRequestList(1)">搜索</Button>
      </Col>
    </Row>

    <div class="table-wrapper">
      <Table border :columns="columns" :data="rows">
        <template slot-scope="{ row }" slot="avatar">
          <Avatar :src="row.avatar" />
        </template>
        <template slot-scope="{ row, index }" slot="action">
          <Button type="primary" size="small" style="margin-right: 5px" @click="edit(index)">编辑</Button>
          <Button type="error" size="small" @click="remove(index, row.id)">删除</Button>
        </template>
      </Table>
      <Spin size="large" fix v-if="spinShow"></Spin>
    </div>

    <div class="page-wrapper">
      <Page
        ref="page"
        :total="total"
        :current.sync="current"
        show-elevator
        show-sizer
        @on-change="handleRequestList"
      />
    </div>
  </div>
</template>

<script>
import { list, del, filter } from '../../../api/backend/posts'
import { formatDate } from '../../../libs/date'

export default {
  name: 'admin-posts-list',
  components: {},
  data() {
    return {
      columns: [],
      rows: [],
      total: 0,
      spinShow: true,
      current: 1,
      searchType: 1,
      searchValues: [
        { label: '标题', value: 1, key: 'title', type: String },
        { label: '标签', value: 2, key: 'labels', type: Array }
      ],
      loading: false,
      options: [],
      searchModel: ''
    }
  },
  methods: {
    edit(index) {
      // 跳转编辑页面
      this.$router.push({
        name: 'admin-posts-write',
        query: { id: this.rows[index].id }
      })
    },
    remove(index, id) {
      del(id)
        .then(result => {
          if (result.data.code === 200) {
            this.$Message['success']({
              background: true,
              content: '删除成功'
            })
            this.rows.splice(index, 1)
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
    handleRequestList(p) {
      let d = {
        page: isNaN(p) ? 1 : p,
        page_size: this.$refs.page.currentPageSize,
        search: this.searchParams()
      }

      if (this.current !== p) {
        this.current = p
      }

      list(d)
        .then(result => {
          this.spinShow = false
          let _data = result.data.data.rows
          this.total = result.data.data.total

          for (let i = 0; i < _data.length; i++) {
            _data[i]['created_at'] = _data[i]['created_at']
              ? formatDate(new Date(_data[i]['created_at']))
              : ''

            _data[i]['update_at'] = _data[i]['update_at']
              ? formatDate(new Date(_data[i]['update_at']))
              : ''
          }

          let _columns = result.data.data.columns
          for (let i = 0; i < _columns.length; i++) {
            if (_columns[i]['key'] === 'avatar') {
              _columns[i]['slot'] = _columns[i]['key']
            }
          }
          this.$set(this, 'columns', _columns)
          this.columns.push({
            title: '操作',
            fixed: 'right',
            slot: 'action'
          })
          this.$set(this, 'rows', _data)
        })
        .catch(err => {
          this.spinShow = false
          this.$Message['error']({
            background: true,
            content: err
          })
        })
    },
    remoteFilter(query) {
      if (query !== '') {
        this.loading = true
        let data = {
          type: this.searchType,
          query: query
        }

        filter(data)
          .then(result => {
            this.loading = false
            if (result.data.code === 200) {
              this.options.splice(0, this.options.length)
              for (let i = 0; i < result.data.data.tags.length; i++) {
                this.options.push(result.data.data.tags[i])
              }
            } else {
              this.$Message['error']({
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
    searchParams() {
      let search = {}
      for (let i = 0; i < this.searchValues.length; i++) {
        if (this.searchValues[i].value === this.searchType) {
          search[this.searchValues[i].key] = this.searchModel
        } else {
          let Type = this.searchValues[i].type
          search[this.searchValues[i].key] = new Type()
        }
      }
      return search
    }
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.handleRequestList(vm.current)
    })
  }
}
</script>

<style lang="less" scoped>
.page-wrapper {
  margin: 15px 20px 0px 0px;
  float: right;
}
.search-wrapper {
  margin: 0 0 20px 0;
}
.table-wrapper {
  position: relative;
}
</style>
