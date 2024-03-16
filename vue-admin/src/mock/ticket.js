import Mock from 'mockjs'

const Random = Mock.Random
const getParamObj = (url) => {
  const search = decodeURIComponent(url).split('?')[1]
  const paramOjb = {}
  search.split('&').forEach(d => {
    const [key, value] = d.split('=')
    paramOjb[key] = value
  })
  return paramOjb
}

/**
 * 查询发布计划列表
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1564
 */
export const getTicketList = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'code': '0000',
    'msg': '',
    'ticketList': [{
      'ticketId': 39,
      'title': '市场上线工单190419-3',
      'targetVersion': 190418,
      'status': 0,
      'version': 14,
      'creater': 'QA',
      'lastModified_by': 'zzw',
      'assignTo': 'SOR',
      'ticketType': 1,
      'deployEnv': '0'
    }, {
      'ticketId': 12,
      'title': '市场上线工单190422-3',
      'targetVersion': 190511,
      'status': 1,
      'version': 14,
      'creater': 'QA',
      'lastModified_by': 'zzw',
      'assignTo': 'SOR',
      'ticketType': 1,
      'deployEnv': '1'
    }, {
      'ticketId': 41,
      'title': '市场上线工单190422-3',
      'targetVersion': 190311,
      'status': 2,
      'version': 14,
      'creater': 'QA',
      'lastModified_by': 'zzw',
      'assignTo': 'SOR',
      'ticketType': 1,
      'deployEnv': '3'
    }, {
      'ticketId': 42,
      'title': '市场上线工单190512-9',
      'targetVersion': 190211,
      'status': 3,
      'version': 14,
      'creater': 'QA',
      'lastModified_by': 'zzw',
      'assignTo': 'SOR',
      'ticketType': 1,
      'deployEnv': '3'
    }, {
      'ticketId': 51,
      'title': '市场上线工单190512-9',
      'targetVersion': 190321,
      'status': 4,
      'version': 14,
      'creater': 'QA',
      'lastModified_by': 'zzw',
      'assignTo': 'SOR',
      'ticketType': 1,
      'deployEnv': '3'
    }, {
      'ticketId': 42,
      'title': '市场上线工单190512-9',
      'targetVersion': 190319,
      'status': 5,
      'version': 14,
      'creater': 'QA',
      'lastModified_by': 'zzw',
      'assignTo': 'SOR',
      'ticketType': 1,
      'deployEnv': '3'
    }],
    count: 40
  })
}

/**
 * 获取需求列表
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1551
 * @param {*} biz 产线
 * @param {*} datetime 目标版本日期
 */
export const getDemandList = (req) => {
  const params = getParamObj(req.url)
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'status': 200,
    'code': '0000',
    'msg': '',
    'result|5-14': [{
      'biz': params.biz,
      'targetVersion': params.fixed_date,
      'reqId|40000-60000': 40011,
      'reqSubject|1-3': Mock.mock('@cword'),
      'issueStatus|+1': ['新建', '转测试', '测试完成', '上线完成', '已排期', '任务完成'],
      'srvMissFlag': 0,
      'missedApp': '',
      'project|1': ['项目1', '项目2'],
      devDuty: '王小明',
      testDuty: '李小狼',
      'relatedTicket|1': [['41', ''], ['52'], ['']],
      'apps|1': [
        ['biz-cam-hotel-srv#43286'],
        ['biz-cam-hotel-srv#43286', 'biz-cam-merchant-srv#43285'],
        ['biz-cam-hotel-srv#43286', 'biz-cam-hotel-web#43293'],
        ['biz-cam-trade-center-web#43179']
      ],
      'sqlScripts|1': [
        ['tcbiz_cam_hotel_48324_ddl.sql'],
        ['tcbiz_cam_hotel_48184_ddl.sql', 'tcbiz_cam_merchant_48184_ddl.sql'],
        ['tcbiz_cam_merchant_48184_update.sql', 'tcbiz_cam_hotel_48184_ddl.sql', 'tcbiz_cam_merchant_48184_ddl.sql']
      ],
      'config|0-3': [
        Random.url(),
        Random.url(),
        Random.url()
      ],
      'frontPkg|0-3': [
        Random.url(),
        Random.url(),
        Random.url()
      ]
    }]
  })
}

/**
 * 获取单条 发布计划/工单，因为两者数据一致，因此共用
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1552
 * @param {*} ticketId 工单ID
 *
 */
export const getPlan = (req) => {
  const params = getParamObj(req.url)
  console.log(`请求路径：${req.url}`)
  const type = params.type

  return Mock.mock({
    'code': '0000',
    'msg': '',
    'ticket': {
      'ticketId': 36,
      'title': '市场上线工单190418-1',
      'biz': '市场',
      'targetVersion': '190418',
      // 'status': '4',
      'status': String(type) === '1' ? '2' : '3', // 0:新建/1:待确认/2:预发执行中/3:生产执行中/4:发布完成/5:取消
      'assignTo': '',
      'version': 7
    },
    'issues': [{
      'reqId': '46951',
      'reqSubject': '众筹转让专区',
      'apps': [
        'biz-mkt-oss-web#43586',
        'biz-mkt-tkt-srv#43693'
      ],
      'sqlScripts': [
        '/mysql_db_scripts/CF/20190327/tcbiz_cf_47610_ddl.sql'
      ],
      'config': [
        Random.url()
      ],
      'frontPkg': [
        Random.url()
      ],
      'issueStatus': '新建',
      'version': 7,
      'targetVersion': Mock.mock('@date'),
      devDuty: '王小明',
      testDuty: '李小狼'
    }, {
      'reqId': '31134',
      'reqSubject': '旅游基金专区',
      'apps': [
        'biz-mkt-oss-web#43586',
        'biz-mkt-tkt-srv#43693'
      ],
      'sqlScripts': [
        '/mysql_db_scripts/CF/20190327/tcbiz_cf_47610_ddl.sql'
      ],
      'config': [
        Random.url()
      ],
      'frontPkg': [
        Random.url()
      ],
      'issueStatus': '测试完成',
      'version': 7,
      'targetVersion': Mock.mock('@date'),
      devDuty: '王小明',
      testDuty: '李小狼'
    }],
    'items': [{
      'item': 'biz-mkt-oss-srv',
      'id': '41',
      'type': 'srv',
      'reqId': ['412311', '444411'],
      'app_version': '43587',
      'order': 0,
      'before_deploy': 1,
      'deployFlag': 3,
      'version': 7,
      'rcFlag': 0,
      'prodFlag': 0
    }, {
      'item': 'biz-mkt-fff-srv',
      'id': '55',
      'type': 'srv',
      'reqId': ['14412', '62132'],
      'app_version': '33123',
      'order': 1,
      'before_deploy': 1,
      'deployFlag': 3,
      'rcFlag': 1,
      'prodFlag': 0
    }, {
      'item': 'biz-mkt-oss-config',
      'id': '12',
      'type': 'config',
      'reqId': ['412311'],
      'app_version': '',
      'order': 2,
      'before_deploy': 1,
      'deployFlag': 3,
      'version': 7,
      'rcFlag': 1,
      'prodFlag': 0
    }, {
      'item': 'biz-mkt-oss-sql',
      'id': '4',
      'type': 'sql',
      'app_version': '',
      'reqId': ['412311'],
      'order': 3,
      'before_deploy': 1,
      'deployFlag': 1,
      'version': 7,
      'rcFlag': 1,
      'prodFlag': 0
    }, {
      'item': 'biz-mkt-oss-sql',
      'id': '44',
      'type': 'sql',
      'app_version': '',
      'reqId': ['312231'],
      'order': 3,
      'before_deploy': 2,
      'deployFlag': 2,
      'version': 7,
      'rcFlag': 0,
      'prodFlag': 0
    }, {
      'item': Random.url(),
      'id': '89',
      'type': 'front',
      'app_version': '',
      'reqId': ['412311'],
      'order': 4,
      'before_deploy': 1,
      'deployFlag': 3,
      'version': 7,
      'rcFlag': 0,
      'prodFlag': 1
    }]
  })
}

/**
 * 保存发布计划
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1553
 */
export const savePlan = (req) => {
  console.log(`请求路径：${req.url}`)
  console.warn(`url过长`)
  return Mock.mock({
    'status': 200,
    'code': '0000',
    'msg': '保存成功',
    'new_version|1-200': 25,
    reqCheckResult: [
      { '41211': '需求有误' },
      { '21233': '需求有误' }
    ],
    'items_list': [
      146,
      147,
      148
    ],
    'req_list': [
      70,
      71
    ]
  })
}

/**
 * 新建一个空的工单
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1554
 *
 * @param biz {string} 产线
 * @param targetVersion {string} 目标版本(yymmdd)
 * @param creater {string} 创建人
 *
 */
export const createEmptyTicket = (req) => {
  console.log(`请求路径：${req.url}`)
  const params = JSON.parse(req.body)
  console.log(params)
  const { biz, targetVersion } = params
  return Mock.mock({
    'code': '0000',
    'msg': '工单保存成功',
    'ticketId': 14,
    'title': `${biz}${targetVersion}-3`
  })
}

/**
 * 发布计划工单状态检测
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1558
 */
export const getTicketPassStatus = () => {
  return Mock.mock({
    code: '0000',
    msg: '通过'
  })
}

/**
 * 确认工单已上线完成
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1565
 * @param ticketId 工单编号
 * @param status 上线阶段 1:预发发布完成/2:生产发布完成
 * @param creater 操作人
 * @param cur_version 工单当前版本号(编辑工单时会返回)
 */
export const confrimTicket = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'status': 200,
    'code': '0000',
    'msg': '确认成功'
  })
}

/**
 * 不使用
 * 生成发布工单, 同保存生产计划数据
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1566
 */
// export const createTicket = () => {
//   return Mock.mock({
//     'code': '0000',
//     'msg': '成功'
//   })
// }

/**
 * 新增一条工单元素
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1567
 * @param ticketId {string} 工单id
 * @param item {string} 工单元素名称
 * @param app_version {string} srv版本号
 * @param type {string} 元素类型(srv,sql,config,front)
 * @param deployFlag {int} /1预发/2生产/3都执行
 * @param before_deploy {int} 是否在发布应用前执行:0否/1是
 * @param creater {string} 创建人
 * @param cur_version {int} 工单当前版本号(编辑工单时会返回)
 *
 */
export const addTicketItem = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'code': '0000',
    'msg': '保存成功',
    'isModified': 0,
    'items_id': 147,
    'new_version|1-200': 22
  })
}

/**
 * 删除一条工单元素
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1567
 * @param item {string} 工单元素名称
 * @param ticketId {int} 工单编号
 * @param creater {string} 创建人
 * @param cur_version {int} 工单当前版本号(编辑工单时会返回)
 *
 */
export const delTicketItem = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'code': '0000',
    'msg': '保存成功',
    'isModified': 0,
    'items_id': 147,
    'new_version|1-200': 22
  })
}

/**
 * 修改工单某项任务的状态
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1569
 * @param itemList {list} 工单元素
 * @param id {int} 工单编号
 * @param status {int} 修改后的工单元素状态:0新建/1预发发布中/2预发发布完成/3生产发布中/4生产发布完成
 * @param isRcDeployed {int} 该元素是否在预发发布过0否/1是（只有修改元素状态为"生产发布中"时才会判断，其他状态一律返回0）
 *
 * @param creater {string} 创建人
 * @param ticketId {string} 工单编号
 * @param cur_version {int} 工单当前版本号(编辑工单时会返回)
 */
export const modifyTicketItemStatus = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'code': '0000',
    'msg': '保存成功',
    'new_version|1-200': 22,
    'items_list': [{
      id: '31',
      itemName: 'oss-srv',
      status: '2',
      'isRcDeployed|1': ['0', '1']
    }]
  })
}
/**
 * 修改应用版本号
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1570
 * @param id {string} 工单元素id
 * @param ticketId {int} 工单编号
 * @param app_version {int} 修改后的工单版本号
 * @param creater {string} 创建人
 * @param cur_version {int} 工单当前版本号(编辑工单时会返回)
 */
export const modifySrvVersion = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'code': '0000',
    'msg': '保存成功',
    'isModified': 0,
    'items_id': 147,
    'new_version|1-200': 22
  })
}

/**
 * 获取今天之后的目标版本日期（包括今天）
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1559
 */
export const getTargetVerDate = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'code': '0000',
    'msg': '保存成功',
    'result': [
      '2019-04-24',
      '2019-05-08',
      '2019-05-12',
      '2019-06-12',
      '2019-07-12',
      '2019-08-12',
      '2019-09-12'
    ]
  })
}

/**
 * 将工单置为取消状态（不可恢复）
 * http://10.10.0.155:9090/index.php?s=/74&page_id=1578
 */
export const cancelTicket = (req) => {
  console.log(`请求路径：${req.url}`)
  return Mock.mock({
    'code': '0000',
    'msg': '取消成功'
  })
}
