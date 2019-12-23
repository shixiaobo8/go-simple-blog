import Main from '@/view/front/main'

export default [
  {
    path: '/',
    component: Main,
    children: [
      {
        path: '/',
        name: 'index',
        meta: {
          title: '首页',
          requireAuth: false
        },
        component: () => import('@/view/front/index')
      },
      {
        path: '/tag/:tag',
        name: 'filter-tag-posts',
        meta: {
          title: 'tag筛选',
          requireAuth: false
        },
        component: () => import('@/view/front/index')
      },
      {
        path: '/posts/:id',
        name: 'posts',
        meta: {
          title: '文章详情',
          requireAuth: false
        },
        component: () => import('@/view/front/posts/posts')
      }
    ]
  }
]