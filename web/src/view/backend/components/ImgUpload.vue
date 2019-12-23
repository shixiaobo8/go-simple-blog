<template>
  <div class="img-upload-wrapper">
    <div class="demo-upload-list" v-for="(item, index) in uploadList" :key="index">
      <template v-if="item.status === 'finished'">
        <img :src="item.url" />
        <div class="demo-upload-list-cover">
          <Icon type="ios-trash-outline" @click.native="handleRemove(item)"></Icon>
        </div>
      </template>
      <template v-else>
        <Progress v-if="item.showProgress" :percent="item.percentage" hide-info></Progress>
      </template>
    </div>
    <Upload
      ref="upload"
      :show-upload-list="false"
      :default-file-list="defaultList"
      :on-success="handleSuccess"
      :format="['jpg','jpeg','png']"
      :max-size="2048"
      :on-format-error="handleFormatError"
      :on-exceeded-size="handleMaxSize"
      :before-upload="handleBeforeUpload"
      :multiple="multiple"
      :name="name"
      with-credentials
      type="drag"
      :action="uploadImgApi"
      style="display: inline-block;width:58px;"
    >
      <div style="width: 58px;height:58px;line-height: 58px;">
        <Icon type="ios-camera" size="20"></Icon>
      </div>
    </Upload>
  </div>
</template>

<script>
import { uploadImgApi } from '../../../api/data'
import { urlFileName } from '../../../libs/tools'

export default {
  props: {
    multiple: {
      type: Boolean,
      default() {
        return false
      }
    },
    name: {
      // 上传文件的字段名
      type: String,
      default() {
        return 'file'
      }
    },
    defaultList: {
      // 默认文件列表
      type: Array,
      default() {
        return []
      }
    },
    maxUploadNum: {
      // 限制最大上传数量
      type: Number,
      default() {
        return 1
      }
    }
  },
  data() {
    return {
      uploadImgApi,
      uploadList: []
    }
  },
  methods: {
    handleRemove(file) {
      const fileList = this.$refs.upload.fileList
      this.$refs.upload.fileList.splice(fileList.indexOf(file), 1)
      this.uploadList.splice(fileList.indexOf(file), 1)
    },
    handleSuccess(res, file) {
      let fname = urlFileName(res.data.img)
      let index = this.$refs.upload.fileList.length - 1 <= 0 ? 0 : this.$refs.upload.fileList.length - 1
      this.$refs.upload.fileList[index].name = fname
      this.$refs.upload.fileList[index].url = res.data.img

      this.uploadList.push(this.$refs.upload.fileList[index])
      // this.$refs.upload.fileList[index] = Object.assign({}, this.$refs.upload.fileList[index])
      // console.log(this.$refs.upload.fileList[index])
      // this.uploadList.splice(0, 0)
    },
    handleFormatError(file) {
      this.$Notice.warning({
        title: 'The file format is incorrect',
        desc:
          'File format of ' +
          file.name +
          ' is incorrect, please select jpg or png.'
      })
    },
    handleMaxSize(file) {
      this.$Notice.warning({
        title: 'Exceeding file size limit',
        desc: 'File  ' + file.name + ' is too large, no more than 2M.'
      })
    },
    handleBeforeUpload() {
      if (this.uploadList.length >= this.maxUploadNum) {
        this.$Notice.warning({
          title: 'only ' + this.maxUploadNum + ' images can be uploaded'
        })
        return false
      }
      return true
    },
    handleDefaultList() { // 用于兼容，defaultList修改添加属性后无法被监听到，没有找到好的解决办法
      setTimeout(() => {
        if (this.defaultList.length > 0) {
          for (let i = 0; i < this.defaultList.length; i++) {
            this.uploadList[i] = Object.assign(
              {},
              this.uploadList[i],
              this.defaultList[i]
            ) // 整个对象更新需要这样解决
          }
        } else {
          this.uploadList = []
        }

        this.uploadList = Object.assign([], this.uploadList)
      }, 500)
    }
  },
  mounted() {
    this.uploadList = this.$refs.upload.fileList
    this.handleDefaultList()
  },
  watch: {
    // 解决异步请求问题
    defaultList(val) {
      for (let i = 0; i < val.length; i++) {
        this.uploadList[i] = Object.assign(
          {},
          this.uploadList[i],
          val[i]
        ) // 整个对象更新需要这样解决
      }
      this.uploadList = Object.assign([], this.uploadList)
    }
  }
}
</script>
<style scoped>
.demo-upload-list {
  display: inline-block;
  width: 60px;
  height: 60px;
  text-align: center;
  line-height: 60px;
  border: 1px solid transparent;
  border-radius: 4px;
  overflow: hidden;
  background: #fff;
  position: relative;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.2);
  margin-right: 4px;
}
.demo-upload-list img {
  width: 100%;
  height: 100%;
}
.demo-upload-list-cover {
  display: none;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.6);
}
.demo-upload-list:hover .demo-upload-list-cover {
  display: block;
}
.demo-upload-list-cover i {
  color: #fff;
  font-size: 20px;
  cursor: pointer;
  margin: 0 2px;
}
</style>