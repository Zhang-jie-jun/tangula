import config from '@/config'
import {createSocket} from "@/libs/new-socket";
export default class JFSocketIo {
  constructor (channel) {
    this.websock = null
    this.setIntervalWesocketPush=null
    this.channel = channel
    this.url = process.env.NODE_ENV === 'production' ? config.socketIoUrl.pro : config.socketIoUrl.dev
    this.initWebSocket()
  }

  /**
   * 初始化weosocket
   */
  initWebSocket = () => {
    if(this.websock){
      this.websock.close()
      console.log('断开websocket')
    }
    const wsuri = this.url + this.channel
    console.log('初始化WebSocket', wsuri)
    this.websock = new WebSocket(wsuri)
    this.websock.onmessage=this.websocketonMessage
    this.websock.onopen = this.websocketonopen
    this.websock.onerror = this.websocketonerror
    this.websock.onclose = this.websocketclose
  }

  websocketonMessage = (e) => {
    window.dispatchEvent(new CustomEvent('onmessageWS', {
      detail: {
        data: e.data
      }
    }))
  }

  /**
   * 连接建立之后执行send方法发送数据
   */
  websocketonopen = () => {
    this.sendPing()
  }

  /**
   * 关闭连接
   */
  closeConnection = () => {
    this.websock.close()
  }

  /**
   * 连接建立失败重连
   */
  websocketonerror = () => {
    this.websock.close()
    clearInterval(this.setIntervalWesocketPush)
    console.log('连接失败重连中')
    if (this.websock.readyState !== 3) {
      this.websock = null
      this.initWebSocket()
    }
  }


  connecting = (message) => {
    setTimeout(() => {
      console.log('connecting……')
      if (this.websock.readyState === 0) {
        this.connecting(message)
      } else {
        this.websock.send(JSON.stringify(message))
        console.log('connecting success!')
      }
    }, 1000)
  }
  /**
   * 关闭
   */
  websocketclose = (e) => {
    clearInterval(this.setIntervalWesocketPush)
    this.websock.close()
    console.log('断开连接', e)
  }

  /**发送心跳
   * @param {number} time 心跳间隔毫秒 默认5000
   * @param {string} ping 心跳名称 默认字符串ping
   */
  sendPing = (time = 10000, ping = 'heart') => {
    clearInterval(this.setIntervalWesocketPush)
    this.websock.send(ping)
    this.setIntervalWesocketPush = setInterval(() => {
      this.websock.send(ping)
      console.log('发送心跳',ping)
    }, time)
  }
  /**
   * 发送数据
   * @param {any} message 需要发送的数据
   */
  sendWS = message => {
    if (this.websock !== null && this.websock.readyState === 3) {
      this.websock.close()
      this.initWebSocket()
    } else if (this.websock.readyState === 1) {
      this.websock.send(JSON.stringify(message))
      console.log('sendMessge',message)
    } else if (this.websock.readyState === 0) {
      this.connecting(message)
    }
  }
}





