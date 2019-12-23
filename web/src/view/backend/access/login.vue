<style lang="less">
  @import './login.less';
</style>

<template>
  <div class="login">
    <div class="login-con">
      <Card icon="log-in" title="欢迎登录" :bordered="false">
        <div class="form-con">
          <login-form @on-success-valid="handleSubmit" ref="loginForm"></login-form>
          <p class="login-tip">请完整填写信息</p>
        </div>
      </Card>
    </div>
  </div>
</template>

<script>
import LoginForm from './login-form'
import { mapActions } from 'vuex'
export default {
  components: {
    LoginForm
  },
  methods: {
    ...mapActions([
      'handleLogin',
      'getUserInfo'
    ]),
    handleSubmit ({ userName, password, cptId, cptCode }) {
      this.handleLogin({ userName, password, cptId, cptCode }).then(res => {
        if (res.code !== 200) {
          this.$refs.loginForm.handleGetCaptcha(true)
          this.$Message['error']({
            background: true,
            content: res.msg
          })
          return
        }

        this.getUserInfo(res.data.id).then(res => {
          this.$router.push({
            name: this.$config.adminHomeName
          })
        })
      })
    }
  }
}
</script>

<style>

</style>
