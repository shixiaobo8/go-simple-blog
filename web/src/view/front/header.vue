<template>
  <Row>
    <div class="header">
      <Row class="navbar">
        <Col :md="20" :xs="24">
          <Row type="flex" justify="end" :xs="24" class="menu">
            <template v-for="(item, index) in menu">
              <Col
                :md="3"
                :xs="4"
                class="menu-item"
                :key="index"
                @click.native="handleMenuClick(item)"
              >{{item.title}}</Col>
            </template>
          </Row>
        </Col>
      </Row>
    </div>
  </Row>
</template>

<script>
export default {
  name: 'Header',
  data() {
    return {
      menu: [
        {
          title: '首页',
          to: { name: 'index' }
        },
        {
          title: 'Go',
          to: { name: 'filter-tag-posts', params: { tag: 'golang' } }
        },
        {
          title: 'PHP',
          to: { name: 'filter-tag-posts', params: { tag: 'php' } }
        },
        {
          title: 'Python',
          to: { name: 'filter-tag-posts', params: { tag: 'python' } }
        }
      ]
    }
  },
  methods: {
    handleMenuClick(item) {
      this.$router.push(item.to)
    }
  }
}
</script>

<style lang="less" scoped>
.header {
  position: relative;
  height: 4rem;
  z-index: 1;
}
.navbar {
  width: 100%;
  position: fixed;
  top: 0;
  background-color: #212121;
  display: flex;
  justify-content: space-between;
  height: 112px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  height: 4rem;
}
.menu {
  .menu-item {
    padding-top: 1rem;
    text-align: center;
    font-size: 0.9rem;
    color: #9e9e9e;
    cursor: pointer;
    &:hover {
      color: #fff;
    }
  }
}
</style>
