import Main from '@/components/main'
import {f} from "vuedraggable";
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

export default [
  {
    path: '/',
    name: '_home',
    redirect: '/home',
    component: Main,
    meta: {
      // hideInMenu: true,
      // notCache: true
    },
    children: [
      {
        path: '/home',
        name: 'home',
        meta: {
          // hideInMenu: true,
          title: '首页',
          // notCache: true,
          icon: 'md-home'
        },
        component: () => import('@/view/single-page/home')
      }
    ]
  },
  {
    path: '/server',
    name: 'server',
    component: Main, // 继承main页面
    meta: {
      icon: 'ios-cloudy-outline',
      title: '服务状态',
      hideInMenu: false
    },
    children: [
      {
        path: '/ceph',
        name: 'ceph',
        meta: {
          icon: 'ios-cloudy-outline',
          title: 'CEPH信息'
        },
        component: () => import('@/view/server/ceph.vue')
      }
    ]
  },
  {
    path: '/env',
    name: 'env',
    component: Main, // 继承main页面
    meta: {
      icon: 'ios-desktop-outline',
      title: '环境管理',
      hideInMenu: false
    },
    children: [
      {
        path: '/machineList',
        name: 'machineList',
        meta: {
          icon: 'md-hammer',
          title: '主机列表'
        },
        component: () => import('@/view/env/machine.vue')
      },
      {
        path: '/platformList',
        name: 'platformList',
        meta: {
          icon: 'md-hammer',
          title: '平台列表'
        },
        component: () => import('@/view/env/platform.vue')
      },

      {
        path: '/poolList',
        name: 'poolList',
        meta: {
          icon: 'md-hammer',
          title: '存储池列表'
        },
        component: () => import('@/view/env/pool.vue')
      }
    ]
  },
  {
    path: '/data',
    name: 'backup',
    component: Main, // 继承main页面
    meta: {
      icon: 'md-filing',
      title: '数据管理',
      hideInMenu: false
    },
    children: [
      {
        path: '/replica',
        name: 'replica',
        meta: {
          icon: 'md-hammer',
          title: '副本列表'
        },
        component: () => import('@/view/data/replica.vue')
      },
      {
        path: '/mirror',
        name: 'mirror',
        meta: {
          icon: 'md-hammer',
          title: '镜像列表'
        },
        component: () => import('@/view/data/mirror.vue')
      },
      {
        path: '/script',
        name: 'script',
        meta: {
          icon: 'md-hammer',
          title: '脚本列表'
        },
        component: () => import('@/view/data/script.vue')
      }
    ]
  },
  {
    path: '/quality',
    name: 'quality',
    component: Main, // 继承main页面
    meta: {
      icon: 'md-desktop',
      title: '质量管理',
      hideInMenu: true
    },
    children: [
      {
        path: 'metaCompare',
        name: 'metaCompare',
        meta: {
          icon: 'md-hammer',
          title: '元数据校验',
          hideInMenu: false
        },
        component: () => import('@/view/quality/metaCompare.vue')
      }
    ]

  },
  {
    path: '/access',
    name: 'Access',
    component: Main, // 继承main页面
    meta: {
      icon: 'md-person-add',
      title: '用户管理',
      hideInMenu: false,
      access: ['SUPER_ADMIN', 'ADMIN']
    },
    children: [
      {
        path: '/userList',
        name: 'userList',
        meta: {
          icon: 'md-hammer',
          title: '用户列表'
        },
        component: () => import('@/view/permission/users.vue')
      },
      {
        path: '/roleList',
        name: 'roleList',
        meta: {
          icon: 'md-hammer',
          title: '角色列表'
        },
        component: () => import('@/view/permission/roles.vue')
      },
      {
        path: '/record',
          name: 'record',
        meta: {
        icon: 'md-hammer',
          title: '操作记录'
      },
        component: () => import('@/view/permission/record.vue')
      }
    ]
  },
  {
    path: '/login',
    name: 'login',
    meta: {
      title: 'Login - 登录',
      hideInMenu: true
    },
    component: () => import('@/view/login/login.vue')
  },
  {
    path: '/message',
    name: 'message',
    component: Main,
    meta: {
      hideInBread: true,
      hideInMenu: true
    },
    children: [
      {
        path: 'message_page',
        name: 'message_page',
        meta: {
          icon: 'md-notifications',
          title: '消息中心'
        },
        component: () => import('@/view/single-page/message/index.vue')
      }
    ]
  },
  {
    path: '/error_logger',
    name: 'error_logger',
    meta: {
      hideInBread: true,
      hideInMenu: true
    },
    component: Main,
    children: [
      {
        path: 'error_logger_page',
        name: 'error_logger_page',
        meta: {
          icon: 'ios-bug',
          title: '错误收集'
        },
        component: () => import('@/view/single-page/error-logger.vue')
      }
    ]
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
    path: '/404',
    name: 'error_404',
    meta: {
      hideInMenu: true
    },
    component: () => import('@/view/error-page/404.vue')
  }

]
