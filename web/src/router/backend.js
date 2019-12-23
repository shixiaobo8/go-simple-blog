import Main from '@/components/main'
// import parentView from '@/components/parent-view'
/**
 * iview-admin中meta除了原生参数外可配置的参数:
 * meta: {
 *  title: { String|Number|Function }
 *         显示在侧边栏、面包屑和标签栏的文字
 *         使用'{{ 多语言字段 }}'形式结合多语言使用，例子看多语言的路由配置;
 *         可以传入一个回调函数，参数是当前路由对象，例子看动态路由和带参路由
 *  hideInBread: (false) 设为true后此级路由将不会出现在面包屑中，示例看QQ群路由配置
 *  hideInMenu: (false) 设为true后在左侧菜单不会显示该页面选项
 *  notCache: (false) 设为true后页面在切换标签后不会缓存，如果需要缓存，无需设置这个字段，而且需要设置页面组件name属性和路由配置的name一致
 *  access: (null) 可访问该页面的权限数组，当前路由设置的权限会影响子路由
 *  icon: (-) 该页面在左侧菜单、面包屑和标签导航处显示的图标，如果是自定义图标，需要在图标名称前加下划线'_'
 *  beforeCloseName: (-) 设置该字段，则在关闭当前tab页时会去'@/router/before-close.js'里寻找该字段名对应的方法，作为关闭前的钩子函数
 * }
 */

export default [{
  path: '/admin/access/login',
  name: 'admin-access-login',
  meta: {
    title: 'Login - 登录',
    hideInMenu: true
  },
  component: () => import('@/view/backend/access/login.vue')
},
{
  path: '/admin/home',
  name: 'admin-home',
  redirect: '/admin/home/index',
  component: Main,
  meta: {
    hideInMenu: true,
    notCache: true
  },
  children: [{
    path: '/admin/home/index',
    name: 'admin-home-index',
    meta: {
      hideInMenu: true,
      title: '首页',
      notCache: true,
      icon: 'md-home',
      requireAuth: true
    },
    component: () => import('@/view/backend/home/home')
  }]
},
{
  path: '/admin/adminUser',
  name: 'admin-user',
  component: Main,
  meta: {
    hideInBread: true,
    requireAuth: true
  },
  children: [{
    path: '/admin/adminUser/list',
    name: 'admin-adminUser-list',
    meta: {
      title: '后台用户',
      icon: 'md-people',
      requireAuth: true
    },
    component: () => import('@/view/backend/admin/list')
  }]
},
{
  path: '/admin/posts',
  name: 'admin-posts',
  component: Main,
  meta: {
    title: '文章管理',
    icon: 'ios-bookmarks',
    hideInBread: true,
    require_Auth: true
  },
  children: [{
    path: 'list',
    name: 'admin-posts-list',
    meta: {
      title: '文章列表',
      hideInBread: true,
      requireAuth: true
    },
    component: () => import('@/view/backend/posts/list')
  },
  {
    path: 'write',
    name: 'admin-posts-write',
    meta: {
      title: '发表文章',
      hideInBread: true,
      requireAuth: true
    },
    component: () => import('@/view/backend/posts/write')
  }]
},
{
  path: '/argu',
  name: 'argu',
  meta: {
    hideInMenu: true
  },
  component: Main,
  children: [{
    path: 'params/:id',
    name: 'params',
    meta: {
      icon: 'md-flower',
      title: route => `{{ params }}-${route.params.id}`,
      notCache: true,
      beforeCloseName: 'before_close_normal'
    },
    component: () => import('@/view/argu-page/params.vue')
  },
  {
    path: 'query',
    name: 'query',
    meta: {
      icon: 'md-flower',
      title: route => `{{ query }}-${route.query.id}`,
      notCache: true
    },
    component: () => import('@/view/argu-page/query.vue')
  }]
},
{
  path: '/401',
  name: 'error_401',
  meta: {
    hideInMenu: true
  },
  component: () => import('@/view/error-page/401.vue')
},
{
  path: '/500',
  name: 'error_500',
  meta: {
    hideInMenu: true
  },
  component: () => import('@/view/error-page/500.vue')
},
{
  path: '*',
  name: 'error_404',
  meta: {
    hideInMenu: true
  },
  component: () => import('@/view/error-page/404.vue')
}]
