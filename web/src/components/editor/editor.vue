<template>
  <div class="editor-wrapper">
    <input
      type="file"
      name="file"
      id="editor-upload-image"
      style="display:none"
      @click="handleClick"
      @change="handleUpload"
    />
    <quill-editor
      v-model="text"
      :options="editorOption"
      ref="QuillEditor">
    </quill-editor>
  </div>
</template>

<script>
import Vue from 'vue'
import { uploadImg } from '@/api/data'
import VueQuillEditor, { Quill } from 'vue-quill-editor'
import ImageResize from 'quill-image-resize-module'
import 'quill/dist/quill.core.css'
import 'quill/dist/quill.snow.css'

Quill.register('modules/imageResize', ImageResize)
Vue.use(VueQuillEditor)
// 工具栏配置
const toolbarOptions = [
  ['bold', 'italic', 'underline', 'strike'], // toggled buttons
  ['blockquote', 'code-block'],

  [{ header: 1 }, { header: 2 }], // custom button values
  [{ list: 'ordered' }, { list: 'bullet' }],
  [{ script: 'sub' }, { script: 'super' }], // superscript/subscript
  [{ indent: '-1' }, { indent: '+1' }], // outdent/indent
  [{ direction: 'rtl' }], // text direction

  [{ size: ['small', false, 'large', 'huge'] }], // custom dropdown
  [{ header: [1, 2, 3, 4, 5, 6, false] }],

  [{ color: [] }, { background: [] }], // dropdown with defaults from theme
  [{ font: [] }],
  [{ align: [] }],
  ['link', 'image', 'video'],
  ['clean'] // remove formatting button
]

export default {
  name: 'editor',
  components: {
    VueQuillEditor
  },
  props: {
    content: {
      type: String,
      default() {
        return ''
      }
    }
  },
  data() {
    return {
      text: '',
      editorOption: {
        modules: {
          imageResize: {},
          toolbar: {
            container: toolbarOptions, // 工具栏
            handlers: {
              image: function(value) {
                if (value) {
                  // 调用iview图片上传
                  document.querySelector('#editor-upload-image').click()
                } else {
                  this.quill.format('image', false)
                }
              }
            }
          }
        }
      }
    }
  },
  methods: {
    handleSuccess(res, file) {
      let quill = this.$refs.QuillEditor.quill
      if (res.data.code === 200) {
        // 获取光标所在位置
        let length = quill.getSelection().index
        // 插入图片，res为服务器返回的图片链接地址
        quill.insertEmbed(length, 'image', res.data.data.img)
        // 调整光标到最后
        quill.setSelection(length + 1)
      } else {
        this.$Message['error']({
          background: true,
          content: '图片上传失败'
        })
      }
    },
    handleClick() {
      // to do
    },
    handleBeforeUpload(file) {
      if (file instanceof File === false) {
        return false
      }
      const allow = ['jpg', 'jpeg', 'gif', 'png']
      const arr = file.name.split('.')
      if (allow.includes(arr[arr.length - 1]) === false) {
        this.$Message['error']({
          background: true,
          content: '图片格式不支持'
        })
        return false
      }

      return true
    },
    handleUpload() {
      const file = document.getElementById('editor-upload-image').files[0]
      if (!this.handleBeforeUpload(file)) {
        return false
      }
      const formData = new FormData()
      formData.append('file', file)
      uploadImg(formData)
        .then(r => {
          this.handleSuccess(r, file)
        })
        .catch(e => {})
    }
  },
  beforeMount: function() { // 这里需要防止父组件ajax异步未请求完成，需要配合watch
    this.text = this.content
  },
  watch: {// 使用watch监听变化，及时更改子组件数据
    content: {
      handler(val) {
        this.text = val
      },
      immediate: true
    }
  }
}
</script>

<style lang="less">
.ql-editor {
  min-height: 10rem;
  background-color: white;
}
</style>
