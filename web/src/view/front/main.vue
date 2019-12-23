<template>
  <div id="main-wrapper">
    <Header></Header>
    <Row type="flex" justify="center" class="centre-wrapper">
      <Col :xs="24" :sm="24" :md="22" :lg="22">
        <router-view></router-view>
      </Col>
    </Row>
  </div>
</template>

<script>
import Header from './header'
export default {
  name: 'Main',
  components: {
    Header
  },
  data() {
    return {}
  }
}
</script>

<style lang="less">
body {
  background-color: #f9f9f9 !important;
}
.centre-wrapper {
  padding-top: 10px;
}
</style>