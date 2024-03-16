import {
  getParams
} from '@/libs/util'
const USER_MAP = {
  super_admin: {
    name: 'super_admin',
    user_id: '1',
    access: ['super_admin', 'admin'],
    token: 'super_admin',
    avator: 'https://file.iviewui.com/dist/a0e88e83800f138b94d2414621bd9704.png'
  },
  admin: {
    name: 'admin',
    user_id: '2',
    access: ['admin'],
    token: 'admin',
    avator: 'https://avatars0.githubusercontent.com/u/20942571?s=460&v=4'
  }
}

export const login = req => {
  req = JSON.parse(req.body)
  return {
    token: USER_MAP[req.userName].token
  }
}

export const getUserInfo = req => {
  const params = getParams(req.url)
  return USER_MAP[params.token]
}

export const logout = req => {
  return null
}

export const auth = req => {
  return {
    app_status_code: 200,
    code: 1000,
    detail: '身份认证信息未提供。',
    msg: '身份认证信息未提供。',
    status: 200,
    type: 'client_exception'
  }
}

export const checkUserToken = (req) => {
  return {
    'status': 200,
    'code': 2000,
    'msg': '',
    'user': 'admin',
    'user_id': 15,
    'access': [
      'ADMIN'
    ],
    'token': '3fa721cda77069d4ed88ddfab6c8a43801a8ea29',
    'permission': {
      'user': 15,
      'business': [],
      'project': [{
        'id': 1,
        'project': 'CF',
        'project_chn': '众筹业务'
      },
      {
        'id': 2,
        'project': 'MAM',
        'project_chn': '会员系统，原CRM'
      },
      {
        'id': 3,
        'project': 'CAM',
        'project_chn': '供应链|资产管理'
      },
      {
        'id': 4,
        'project': 'LOAN',
        'project_chn': '保理业务'
      },
      {
        'id': 5,
        'project': 'INS',
        'project_chn': '信贷产品'
      },
      {
        'id': 6,
        'project': 'OA',
        'project_chn': '内部办公平台'
      },
      {
        'id': 7,
        'project': 'ITS',
        'project_chn': '分期游'
      },
      {
        'id': 8,
        'project': 'FE',
        'project_chn': '前端'
      },
      {
        'id': 9,
        'project': 'BI',
        'project_chn': '商业智能'
      },
      {
        'id': 10,
        'project': 'INF',
        'project_chn': '基础架构'
      },
      {
        'id': 11,
        'project': 'FE-INS',
        'project_chn': '市场信贷前端'
      },
      {
        'id': 12,
        'project': 'MKT',
        'project_chn': '市场后端'
      },
      {
        'id': 13,
        'project': 'FE-MKT',
        'project_chn': '市场营销前端'
      },
      {
        'id': 14,
        'project': 'PAYMENT-MAS',
        'project_chn': '支付-商户收单'
      },
      {
        'id': 15,
        'project': 'PAYMENT-CS',
        'project_chn': '支付-清结算'
      },
      {
        'id': 16,
        'project': 'DATA-PLATFORM',
        'project_chn': '数据平台'
      },
      {
        'id': 17,
        'project': 'DATA-SERVICE',
        'project_chn': '数据服务'
      },
      {
        'id': 18,
        'project': 'NEWPAY',
        'project_chn': '海外支付'
      },
      {
        'id': 19,
        'project': 'FE-NEWPAY',
        'project_chn': '海外支付前端'
      },
      {
        'id': 20,
        'project': 'IGB',
        'project_chn': '爱购宝，原BX保险业务'
      },
      {
        'id': 21,
        'project': 'FNC',
        'project_chn': '理财'
      },
      {
        'id': 22,
        'project': 'QA',
        'project_chn': '质量管理'
      },
      {
        'id': 23,
        'project': 'SOR',
        'project_chn': '运维研发'
      },
      {
        'id': 24,
        'project': 'RCS',
        'project_chn': '风控合规'
      }
      ],
      'app': [],
      'role': 'ADMIN',
      'role_chn': '管理员',
      'remark': ''
    },
    'avator': ''
  }
}
