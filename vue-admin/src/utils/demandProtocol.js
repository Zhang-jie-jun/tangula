export const excuteOrder = {
  '服务发布前执行': 'pre',
  '服务发布后执行': 'post'
}

export const excuteReverseOrder = {
  'pre': '服务发布前执行',
  'post': '服务发布后执行'
}

export const envMap = {
  '预发': '1',
  '生产': '2',
  '预发+生产': '3'
}
export const envReverseMap = {
  '1': '预发',
  '2': '生产',
  '3': '预发+生产'
}

export const envList = [{
  name: '预发',
  value: '1'
}, {
  name: '生产',
  value: '2'
}, {
  name: '预发+生产',
  value: '3'
}]

export const typeList = {
  '配置项': 'config',
  'SQL工单': 'sql',
  '后端应用': 'srv',
  '前端应用': 'front'
}
