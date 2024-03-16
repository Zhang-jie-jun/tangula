export default {
  /**
   * @description 配置显示在浏览器标签的title
   */
  title: 'TANGULA',
  /**
   * @description token在Cookie中存储的天数，默认1天
   */
  cookieExpires: 1,
  /**
   * @description 是否使用国际化，默认为false
   *              如果不使用，则需要在路由中给需要在菜单中展示的路由设置meta: {title: 'xxx'}
   *              用来在菜单中显示文字
   */
  useI18n: false,
  /**
   * @description api请求基础路径
   */
  baseUrl: { // 服务端地址
    // dev: 'http://10.4.32.208:8080',
    dev: 'http://10.4.117.167:8000',
    // pro: 'http://192.168.212.34:8000'
    pro: 'http://10.4.117.167:8000'
  },

  /**
   *
   */
  socketIoUrl: {
    dev: 'ws://localhost:8000/',
    pro: 'ws://10.4.117.167:8000/'

  },
  /**
   * @description 默认打开的首页的路由name值，默认为home
   */
  homeName: 'home',
  /**
   * @description 需要加载的插件
   */
  plugin: {
    'error-store': {
      showInHeader: true, // 设为false后不会在顶部显示错误日志徽标
      developmentOff: true // 设为true后在开发环境不会收集错误信息，方便开发中排查错误
    }
  },
  /**
   * 发布环境
   */
  deployEnv: {
    inte: '集测',
    rc: '预发',
    prod: '生产'
  },

  deployState: {
    1: '待部署',
    2: '部署中',
    3: '成功',
    4: '失败'
  },
  roleMap: {
    0: 'QA',
    1: 'SUPER_ADMIN',
    3: 'ADMIN'
  },
}
