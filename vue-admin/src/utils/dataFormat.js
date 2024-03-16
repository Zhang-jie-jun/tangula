import { envReverseMap } from './demandProtocol'

const btnTypeMap = {
  '0': 'info',
  '1': 'success',
  '2': 'warning',
  '3': 'error',
  '4': 'error',
  '5': 'warning',
  '6': 'success',
  '8': 'success',
  '9': 'error',
  '10': 'warning'
}

const sqlTypeMap = {
  '0': 'info',
  '1': 'success',
  '2': 'warning',
  '3': 'error',
  '4': 'primary',
  '5': 'warning',
  '6': 'error',
  '7': 'error'
}

const sqlText = {
  '0': '确定发起审核吗？',
  '1': '已发布',
  '2': 'SQL工单发布中，若长时间无响应请联系DBA人工处理',
  '3': '确定重新发布SQL工单吗？',
  '5': '等待审核中……',
  '6': '请申请新的SQL工单',
  '7': '确定发起审核吗'
}
export default class DataFormat {
  /**
   * @param {*} rawData
   * @param {Array} config // 配置
   * @param {Array} sqlScripts // sql
   * @param {Array} apps
   * @param {Array} frontPkg
   * @param {string} config // 配置表
   * @param {string} frontPkg // 前端包
   *
   */
  static demandDataFormat (rawData, forceUpdate = true) {
    console.log('rawData', rawData)
    if (!rawData.deploySqlList_new || !rawData.apps) {
      return rawData
    }

    const raw = Object.assign({}, rawData)
    console.log('raw', raw)
    const orderTemplate = {
      order: 'pre'
    }
    const envTemplate = {
      env: '3' // 默认为“预发+生产”
    }
    const sqlEnvTemplate = {
      env: '1' // 默认为预发
    }
    // 已有数据

    if (raw.modified) {
      return raw
    }

    if (!forceUpdate) {
      return raw
    }

    raw['modified'] = true
    // config
    // raw['modifyConfig'] = raw.config.map(d => (Object.assign({
    //   itemName: d
    // }, orderTemplate, envTemplate)))
    raw['modifyConfig'] = raw.deployConfigListNew.map(d => (Object.assign({
      itemName: d.split('|*|')[0] + '=>' + d.split('|*|')[1]
    }, orderTemplate, envTemplate)))
    // sql
    raw['modifySql'] = raw.deploySqlList_new.map(d => (Object.assign({
      itemName: d.split('|*|')[0],
      sql: d.split('|*|')[1]
    }, orderTemplate, sqlEnvTemplate)))
    // raw['modifySql'] = raw.deploySqlList_new.map(d => {
    //   const [sqlName, sqlDetails] = d.split('|*|')
    //   return Object.assign({
    //     sqlName,
    //     sqlDetails
    //   }, orderTemplate, sqlEnvTemplate)
    // })
    // app
    raw['modifyApp'] = raw.apps.map(d => {
      const [itemName, srvVer] = d.split('#')
      return Object.assign({
        itemName,
        srvVer
      }, envTemplate)
    })
    // front
    // raw['modifyFront'] = raw.frontPkg.map(d => (
    //   Object.assign({
    //     itemName: d
    //   }, envTemplate)
    // ))
    //
    raw['modifyFront'] = raw.appFronts.map(d => {
      const [itemName, srvVer] = d.split('#')
      return Object.assign({
        itemName,
        srvVer
      }, envTemplate)
    })
    return raw
  }

  static demandDataListFormat (arr, forceUpdate = true) {
    if (!(arr instanceof Array)) {
      console.error('TypeError: demandDataListFormate needs Array, but get:')
      console.log(arr)
      return
    }
    const list = arr.map(d => d)
    console.log('list', list)
    return list.map(d => this.demandDataFormat(d, forceUpdate))
  }

  static duplicateRemove (arr) {
    return arr.reduce((all, next, index) => {
      const nameList = all.map(d => d.itemName)
      const duplicatePosition = nameList.indexOf(next.itemName)
      if (duplicatePosition !== -1) {
        const dupulicate = all[duplicatePosition]
        dupulicate['duplicated'] = true
        const dualIssueIds = dupulicate.issueId

        if (next.type === 'srv') {
          const target = all[duplicatePosition]
          const srvVer = target.srvVer
          target.srvVer = Number(srvVer).toFixed(0) > Number(next.srvVer).toFixed(0) ? srvVer : next.srvVer
        }

        if (!dualIssueIds.includes(next.issueId)) {
          const issues = new Set(dupulicate.issueId.concat(next.issueId))
          dupulicate.issueId = Array.from(issues)
        }

        return all
      }
      all.push(next)
      return all
    }, [])
  }

  static extractDemandDataFromExist = (arr, env) => { // 已有发布计划数据 格式化
    console.log('dataFormat', arr, env)
    let curEnv = ''
    if (env) {
      curEnv = env
    } else {
      curEnv = 'inte'
    }
    const [rcEnv, prodEnv, bothEnv] = Object.keys(envReverseMap)
    const arrMap = arr.map(d => {
      const data = {
        id: d.id,
        itemName: d.item,
        order: String(d.before_deploy) === '1' ? 'pre' : 'post', // 0发布后, 1发布前
        env: String(d.deployFlag),
        curEnv: curEnv,
        type: d.type,
        issueId: d.reqId,
        srvVer: d.app_version,
        serialNo: String(d.serial_no),
        coverage: String(d.coverage),
        branch: String(d.branch),
        sql: String(d.sql), // 记录sql详情
        status: (String(d.rcFlag) === '1' && env === 'rc') || (String(d.prodFlag) === '1' && env === 'prod') || (String(d.prodFlag) === '8' && env === 'prod'), // 是否已发布, true 已发布，false，未发布
        deployStatus: env === 'rc' ? String(d.rcFlag) : String(d.prodFlag),
        batchStatus: env === 'rc' ? String(d.rcBatchStatus) : String(d.prodBatchStatus),
        rcFlag: String(d.rcFlag),
        prodFlag: String(d.prodFlag),
        reviewer: String(d.reviewer),
        denyReason: String(d.denyReason),
        affected_rows: String(d.affected_rows),
        deploy_mode: String(d.deploy_mode),
        nodeStatusList: env === 'rc' ? d.rcNodeStatus : d.prodNodeStatus,
        frontNodeStatusList: env === 'rc' ? d.frontRcNodeStatus : d.frontProdNodeStatus,
        deployText: (String(d.rcFlag) === '2' && env === 'rc') || (String(d.prodFlag) === '2' && env === 'prod') || (String(d.prodFlag) === '10' && env === 'prod') ? '是否将应用置为发布失败?' : '是否发布该应用？'
      }
      return Object.assign(data, {
        srvButtonType: btnTypeMap[data.deployStatus],
        batchButtonType: btnTypeMap[data.batchStatus],
        sqlButtonType: sqlTypeMap[data.deployStatus],
        sqlButtonText: sqlText[data.deployStatus]
      })
    })

    const rc = {
      listConfigPre: [],
      listConfigPost: [],
      listSqlPre: [],
      listSqlPost: [],
      listSrv: [],
      listFront: []
    }
    const prod = {
      listConfigPre: [],
      listConfigPost: [],
      listSqlPre: [],
      listSqlPost: [],
      listSrv: [],
      listFront: []
    }

    rc.listConfigPre = arrMap.filter(f => f.type === 'config' && f.order === 'pre' && (f.env === rcEnv || f.env === bothEnv))
    rc.listConfigPost = arrMap.filter(f => f.type === 'config' && f.order === 'post' && (f.env === rcEnv || f.env === bothEnv))
    rc.listSqlPre = arrMap.filter(f => f.type === 'sql' && f.order === 'pre' && (f.env === rcEnv || f.env === bothEnv))
    rc.listSqlPost = arrMap.filter(f => f.type === 'sql' && f.order === 'post' && (f.env === rcEnv || f.env === bothEnv))
    rc.listSrv = arrMap.filter(f => f.type === 'srv' && (f.env === rcEnv || f.env === bothEnv))
    rc.listFront = arrMap.filter(f => f.type === 'front' && (f.env === rcEnv || f.env === bothEnv))
    rc.listSrv = DataFormat.sortSrvList(rc.listSrv)

    prod.listConfigPre = arrMap.filter(f => f.type === 'config' && f.order === 'pre' && (f.env === prodEnv || f.env === bothEnv))
    prod.listConfigPost = arrMap.filter(f => f.type === 'config' && f.order === 'post' && (f.env === prodEnv || f.env === bothEnv))
    prod.listSqlPre = arrMap.filter(f => f.type === 'sql' && f.order === 'pre' && (f.env === prodEnv || f.env === bothEnv))
    prod.listSqlPost = arrMap.filter(f => f.type === 'sql' && f.order === 'post' && (f.env === prodEnv || f.env === bothEnv))
    prod.listSrv = arrMap.filter(f => f.type === 'srv' && (f.env === prodEnv || f.env === bothEnv))
    prod.listFront = arrMap.filter(f => f.type === 'front' && (f.env === prodEnv || f.env === bothEnv))
    prod.listSrv = DataFormat.sortSrvList(prod.listSrv)
    console.log('listSqlPre', rc.listSqlPre)

    return {
      rc,
      prod
    }
  }
  static sortSrvList = (arr) => {
    var result = []
    var sno_arr = arr.map(d => d.serialNo)
    sno_arr = new Set(sno_arr)
    sno_arr = Array.from(sno_arr)
    sno_arr.sort()
    console.log('array', sno_arr)
    sno_arr.forEach(d => {
      var list = arr.filter(f => f.serialNo === d)
      const item = { batIndex: d, batList: list }
      result.push(item)
    })
    return result
  }
  static genRandom = (min, max) => {
    return (Math.random() * (max - min + 1) | 0.0) + min
  }
  // 时间格式化 safari失效
  static formatTime = (date, fmt) => {
    let o = {
      'M+': date.getMonth() + 1, // 月份
      'd+': date.getDate(), // 日
      'h+': date.getHours(), // 小时
      'm+': date.getMinutes(), // 分
      's+': date.getSeconds(), // 秒
      'S': date.getMilliseconds() // 毫秒
    }
    if (/(y+)/.test(fmt)) {
      fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
    }
    for (var k in o) {
      if (new RegExp('(' + k + ')').test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (('00' + o[k]).substr(('' + o[k]).length)))
      }
    }
    return fmt
  }
  // 时间格式化
  static getDayTime = (date) => {
    let formatTime = date.substr(5, 11)
    return formatTime
  }
  static orderDataFormat () {}
}
