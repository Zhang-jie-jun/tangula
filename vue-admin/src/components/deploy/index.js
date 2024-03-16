import BuildInfo from './build-info.vue'
import QaDeploySearch from './qa-deploy-search.vue'
// 导出是export default 取值不能用析构的方式
// import t from '_c/deploy'
// t.BuildInfo  通过这个方式就能获取到
export default {
  BuildInfo,
  QaDeploySearch
}
