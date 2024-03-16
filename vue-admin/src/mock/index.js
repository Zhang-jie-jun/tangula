import Mock from 'mockjs'
import { checkUserToken, auth } from './login'
import { messageCount, userPermission } from './user'
// import { getTableData, getDragList, uploadImage } from './data'
// import { getMessageInit, getContentByMsgId, hasRead, removeReaded, restoreTrash, messageCount } from './user'
import {
  getDemandList,
  getPlan,
  savePlan,
  createEmptyTicket,
  getTicketPassStatus,
  getTicketList,
  confrimTicket,
  // createTicket,
  addTicketItem,
  delTicketItem,
  modifyTicketItemStatus,
  modifySrvVersion,
  getTargetVerDate,
  cancelTicket } from './ticket'

// 配置Ajax请求延时，可用来测试网络延迟大时项目中一些效果
Mock.setup({
  timeout: 1000
})

// 登录相关和获取用户信息
Mock.mock(/\/account\/auth/, auth)
Mock.mock(/\/user_info/, checkUserToken)
Mock.mock(/message\/count/, messageCount)
Mock.mock(/rd\/permission/, userPermission)
// Mock.mock(/\/logout/, logout)
// Mock.mock(/\/get_table_data/, getTableData)
// Mock.mock(/\/get_drag_list/, getDragList)
// Mock.mock(/\/save_error_logger/, 'success')
// Mock.mock(/\/image\/upload/, uploadImage)
// Mock.mock(/\/message\/init/, getMessageInit)
// Mock.mock(/\/message\/content/, getContentByMsgId)
// Mock.mock(/\/message\/has_read/, hasRead)
// Mock.mock(/\/message\/remove_readed/, removeReaded)
// Mock.mock(/\/message\/restore/, restoreTrash)
// Mock.mock(/\/message\/count/, messageCount)

// 工单相关
Mock.mock(/\/rd\/getRdIssues/, getDemandList)
Mock.mock(/\/rd\/getTickets/, 'get', getPlan)
Mock.mock(/\/rd\/saveTicket/, savePlan)
Mock.mock(/\/rd\/getTickets/, 'post', createEmptyTicket)
Mock.mock(/\/rd\/getIssuesInfo/, getTicketPassStatus)
Mock.mock(/\/rd\/getTicketList/, getTicketList)
Mock.mock(/\/rd\/verifyTicket/, confrimTicket)
// Mock.mock(/\/rd\/saveDeployTicket/, createTicket)
Mock.mock(/\/rd\/addItem/, addTicketItem)
Mock.mock(/\/rd\/delItem/, delTicketItem)
Mock.mock(/\/rd\/modifyStatus/, modifyTicketItemStatus)
Mock.mock(/\/rd\/modifySrvVersion/, modifySrvVersion)
Mock.mock(/\/rd\/getTargetVersion/, getTargetVerDate)
Mock.mock(/\/rd\/cancelTicket/, cancelTicket)

Mock.mock(/\/rd\/inte_build/, getTicketList)
Mock.mock(/\/sch\/getScheduleList/, getTicketList)
Mock.mock(/\/sch\/getSchReqs/, getTicketList)
Mock.mock(/\/sch\/getRecentTarget/, getTicketList)
Mock.mock(/\/rd\/getBizs/, getTicketList)

export default Mock
