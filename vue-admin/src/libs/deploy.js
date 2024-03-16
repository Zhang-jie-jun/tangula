/**
* @param {stateText}
* @description 发布状态颜色
*/
export const choiceDeployStateColor = (stateText) => {
  let color = 'warning'
  switch (stateText) {
    case '挂载中':
      color = 'warning'
      break
    case '卸载中':
      color = 'warning'
      break
    case '空闲':
      color = '#5cadff'
      break
    case '失败':
      color = 'error'
      break
    case '已挂载':
      color = 'success'
      break
    default:
      color = '#5cadff'
  }
  return color
}

export const choiceDeployState = (stateText) => {
  let color = 'warning'
  switch (stateText) {
    case '待发布':
      color = 'warning'
      break
    case '发布中':
      color = 'primary'
      break
    case '失败':
      color = 'error'
      break
    case '成功':
      color = 'success'
      break
    case '上传中':
      color = 'primary'
      break
    default:
      color = 'warning'
  }
  return color
}
