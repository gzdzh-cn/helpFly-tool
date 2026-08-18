/**
 * 基础路由
 * @type { *[] }
 */

const constantRouterMap = [
  {
    path: '/',
    component: () => import('@/layouts/AppSider.vue'),
    redirect: { name: 'home', params: { id: 'home' } },
    children: [
      {
        path: '/home/:id?',
        name: 'home',
        component: () => import('@/views/home/index.vue'),
        params: { id: 'home' },
        meta: {
          title: '首页',
          icon: 'icon-shouye1',
        }
      },
      // {
      //   path: '/example',
      //   name: 'example',
      //   component: () => import('@/views/example/index.vue'),
      //   meta: {
      //     title: '示例',
      //     icon: 'icon-xitongguanli',
      //   }
      // },      
      {
        path: '/submenu',
        name: 'submenu',
        component: () => import('@/layouts/Menu.vue'),
        redirect: { name: 'consumptionList' },
        meta: {
          title: '数据',
          icon: 'icon-gform1',
        },
        children: [
          {
            path: '/submenu/consumption-list',
            name: 'consumptionList',
            component: () => import('@/views/submenu/consumption-list/index.vue'),
            meta: {
              title: '账单列表',
              icon: 'icon-xiaofeijilu',
            }
          },
          {
            path: '/submenu/form-demo',
            name: 'formDemo',
            component: () => import('@/views/submenu/form-demo/index.vue'),
            meta: {
              title: '表单演示',
              icon: 'icon-biaodan2',
            }
          },
          {
            path: '/submenu/icon-gallery',
            name: 'iconGallery',
            component: () => import('@/views/submenu/icon-gallery/index.vue'),
            meta: {
              title: '图标专栏',
              icon: 'icon-tubiao',
            }
          },
          {
            path: '/submenu/chart-demo',
            name: 'chartDemo',
            component: () => import('@/views/submenu/chart-demo/index.vue'),
            meta: {
              title: '图表分析',
              icon: 'icon-biaoge',
            }
          },
          {
            path: '/submenu/component-demo',
            name: 'componentDemo',
            component: () => import('@/views/submenu/component-demo/index.vue'),
            meta: {
              title: '通用组件',
              icon: 'icon-zujian',
            }
          },
          {
            path: '/submenu/parameter-setting',
            name: 'parameterSetting',
            component: () => import('@/views/submenu/parameter-setting/index.vue'),
            meta: {
              title: '参数设置',
              icon: 'icon-canshuguanli',
            }
          },

 
        ]
      },      
      {
        path: '/setting',
        name: 'setting',
        component: () => import('@/views/setting/index.vue'),
        meta: {
          hideInMenu:true,
        }
      },      
    ]
  },
]

export default constantRouterMap