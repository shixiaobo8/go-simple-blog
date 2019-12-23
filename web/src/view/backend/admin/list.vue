<template>
  <div class="admin-user-wrapper">
    <Table border :columns="columns" :data="rows">
      <template slot-scope="{ row }" slot="avatar">
        <Avatar :src="row.avatar" />
      </template>
      <template slot-scope="{ row, index }" slot="action">
        <Button type="primary" size="small" style="margin-right: 5px" @click="edit(index)">编辑</Button>
        <Button type="error" size="small" @click="remove(index)">删除</Button>
      </template>
    </Table>

    <Modal
      v-model="showModel"
      @on-ok="handleSubmit"
      @on-cancel="handleCancel"
      :loading="true"
      title="编辑框"
      width="50%"
    >
      <Form
        :model="form"
        ref="adminUserForm"
        :label-width="80"
        :rules="rules"
        style="max-width:80%;margin:0 auto;"
      >
        <FormItem label="用户昵称" prop="nickname">
          <Input v-model="form.nickname" placeholder="请输入..."></Input>
        </FormItem>
        <FormItem label="用户头像" prop="avatar">
          <cropper
            :src="form.avatar"
            crop-button-text="上传头像"
            @on-crop="handleCroped"
            style="height:300px;"
          ></cropper>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<script>
import { list, updateUser } from '../../../api/backend/adminUser'
import { objInit } from '../../../libs/tools'
import { formatDate } from '../../../libs/date'
import Cropper from '@/components/cropper'
import { uploadImg } from '@/api/data'

export default {
  name: 'admin-user-list',
  components: {
    Cropper
  },
  data() {
    return {
      columns: [],
      rows: [],
      showModel: false,
      form: {
        id: 0,
        nickname: '',
        avatar: ''
      }
    }
  },
  mounted() {
    list()
      .then(result => {
        let _data = result.data.data.rows
        for (let i = 0; i < _data.length; i++) {
          _data[i]['last_login_time'] = _data[i]['last_login_time']
            ? formatDate(new Date(_data[i]['last_login_time']))
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
        this.$Message['error']({
          background: true,
          content: err
        })
      })
  },
  methods: {
    edit(index) {
      this.showModel = true
      this.form.id = this.rows[index].id
      this.form.nickname = this.rows[index].nickname
      this.form.avatar = this.rows[index].avatar
    },
    handleCancel() {
      this.showModel = false
      objInit(this.form)
    },
    handleSubmit() {
      this.$refs.adminUserForm.validate(valid => {
        if (valid) {
          updateUser(this.form).then(r => {
            if (r.data.code === 200) {
              this.$Message.success('修改成功')
            }
            this.showModel = false
          })
        } else {
          this.showModel = false
        }
      })
    },
    remove(index) {
      this.rows.splice(index, 1)
    },
    handleCroped(blob) {
      // 上传头像
      const formData = new FormData()
      formData.append('file', blob)
      uploadImg(formData).then(r => {
        this.form.avatar = r.data.data.img
        this.$Message.success('上传成功')
      })
    }
  },
  computed: {
    rules() {
      // 计算属性
      return {
        nickname: [
          {
            required: true,
            type: 'string',
            message: '昵称不能为空',
            trigger: 'blur'
          },
          {
            type: 'string',
            min: 3,
            message: '昵称字符长度不得小于3',
            trigger: 'blur'
          },
          { max: 15, message: '昵称字符长度不得多于15', trigger: 'blur' }
        ]
      }
    }
  }
}
</script>
