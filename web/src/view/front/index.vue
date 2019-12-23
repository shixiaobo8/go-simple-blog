<template>
  <div class="wrapper clearfix">
  <Row type="flex" justify="center">
  <Col :xs="0" :sm="0" :md="0" :lg="{span: 4}">
    <Card>
      <div style="text-align:center">
        <img src="../../../src/assets/images/logo.jpg" style="max-width:100%" />
        <h3>A high quality UI Toolkit based on Vue.js</h3>
      </div>
    </Card>
  </Col>
  <Col :xs="24" :sm="24" :md="{span: 16}" :lg="{span: 16}">
    <PostsList :tag="tag"></PostsList>
  </Col>
  </Row>
  </div>
</template>
<script>
import PostsList from './posts/list'
export default {
  name: 'index',
  components: {
    PostsList
  },
  data() {
    return {

    }
  },
  computed: {
    tag() {
      return this.$route.params.tag ? this.$route.params.tag : ''
    }
  }
}
</script>

<style lang="less">
#app {
  overflow: auto;
}
.clearfix:after {
  clear: both;
  visibility: hidden;
  font-size: 0;
  height: 0;
  content: "";
  display: table;
}
::-webkit-scrollbar {
  width: .5rem;
  height: .5rem;
  background: hsla(0, 0%, 100%, 0.6);
}

::-webkit-scrollbar-track {
  border-radius: 0;
}

::-webkit-scrollbar-thumb {
  border-radius: 0;
  background-color: rgba(95,95,95,.4);
  transition: all .2s;
  border-radius: .5rem;

  &:hover {
    background-color: rgba(95,95,95, .7);
  }
}
</style>