// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store'
import iView from 'view-design'
import 'view-design/dist/styles/iview.css'
import i18n from '@/locale'
import config from '@/config'
import importDirective from '@/directive'
import installPlugin from '@/plugin'
import './index.less'
import '@/assets/icons/iconfont.css'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import TreeTable from 'tree-table-vue'
import preventReClick from './utils/preventRepeatClick.js'

Vue.use(ElementUI);
// 实际打包时应该不引入mock
/* eslint-disable */
console.log(process.env);
if (process.env.NODE_ENV !== 'production' && process.env.VUE_APP_PROXY !== 'Y') {
  require('@/mock')
}
Vue.use(preventReClick)
Vue.use(iView, {
  i18n: (key, value) => i18n.t(key, value)
})
Vue.use(TreeTable)




/**
 * @description 注册admin内置插件
 */
installPlugin(Vue)
/**
 * @description 生产环境关掉提示
 */
Vue.config.productionTip = false
/**
 * @description 全局注册应用配置
 */
Vue.prototype.$config = config

/**
 * 注册指令
 */
importDirective(Vue)

/* eslint-disable no-new */
 new Vue({
  el: '#app',
  router,
  i18n,
  store,
  render: h => h(App),
   components: {
     App
   }
}).$mount('#app')
