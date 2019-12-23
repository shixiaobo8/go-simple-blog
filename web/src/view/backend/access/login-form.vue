<template>
  <Form ref="loginForm" :model="form" :rules="rules" @keydown.enter.native="handleSubmit">
    <FormItem prop="userName">
      <Input v-model="form.userName" placeholder="请输入用户名">
        <span slot="prepend">
          <Icon :size="16" type="ios-person"></Icon>
        </span>
      </Input>
    </FormItem>
    <FormItem prop="password">
      <Input type="password" v-model="form.password" placeholder="请输入密码">
        <span slot="prepend">
          <Icon :size="14" type="md-lock"></Icon>
        </span>
      </Input>
    </FormItem>
    <FormItem prop="cptCode">
    <div class="cpt-img">
        <img :src="cpt.cpt" alt="验证码" @click="handleGetCaptcha">
    </div>
    <div class="cpt-code">
      <Input type="text" v-model="form.cptCode" placeholder="请输入验证码">
      <span slot="prepend">
        <Icon :size="14" type="ios-cafe-outline"></Icon>
      </span>
    </Input>
    </div>
    </FormItem>
    <FormItem>
      <Button @click="handleSubmit" type="primary" long>登录</Button>
    </FormItem>
  </Form>
</template>
<script>
import { mapActions, mapState } from 'vuex'
export default {
  name: 'LoginForm',
  props: {
    // 父组件可传递
    userNameRules: {
      type: Array,
      default: () => {
        return [
          { required: true, message: '账号不能为空', trigger: 'blur' },
          { type: 'string', min: 3, message: '账号字符长度不得小于3', trigger: 'blur' },
          { max: 20, message: '账号字符长度不得多于20', trigger: 'blur' }
        ]
      }
    },
    passwordRules: {
      type: Array,
      default: () => {
        return [
          { required: true, message: '密码不能为空', trigger: 'blur' }
        ]
      }
    }
  },
  data () {
    return {
      form: {
        userName: '',
        password: '',
        cptId: '',
        cptCode: ''
      },
      flag: true // 阻止频繁刷新验证码
    }
  },
  computed: {
    ...mapState({
      cpt: state => state.user.cpt
    }),
    rules () { // 计算属性
      return {
        userName: this.userNameRules,
        password: this.passwordRules,
        cptCode: { required: true, type: 'string', message: '验证码不能为空', trigger: 'blur' }
      }
    }
  },
  mounted: function () {
    this.getCaptcha() // 自动映射为this.$store.dispatch('getCaptcha')
  },
  methods: {
    ...mapActions([
      'getCaptcha' // 获取验证码
    ]),
    handleGetCaptcha (force) {
      if (force === true) {
        this.getCaptcha()
        this.flag = true
        return
      }
      const interval = 10 * 1000
      if (this.flag === false) {
        return
      }
      this.flag = false
      this.getCaptcha()
      setTimeout(() => {
        this.flag = true
      }, interval)
    },
    handleSubmit () {
      // @see https://www.iviewui.com/components/form
      // 对整个表单进行校验，参数为检验完的回调，会返回一个 Boolean 表示成功与失败，支持 Promise
      this.$refs.loginForm.validate((valid) => {
        if (valid) {
          this.$emit('on-success-valid', {
            userName: this.form.userName,
            password: this.form.password,
            cptId: this.cpt.id,
            cptCode: this.form.cptCode
          })
        }
      })
    }
  }
}
</script>

<style lang="less" scoped>
  .cpt-img {
    width: 35%;
    display: inline-block;
    cursor: pointer;
  }
  .cpt-code {
    max-width: 65%;
    display: inline-block;
  }
</style>>
