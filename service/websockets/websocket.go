package websockets

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var MessageChan = make(chan string)

func UpdateMsg(msg string) {
	MessageChan <- msg
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error %s", err)
	}

	done := make(chan struct{})
	go writer(ws, done)
	go reader(ws, done)
}

func writer(conn *websocket.Conn, done chan struct{}) {
	defer conn.Close()
	for {
		select {
		case <-done:
			// the reader is done, so return
			return
		case singleset := <-MessageChan: // get data from channel
			msgErr := conn.WriteMessage(1, []byte(singleset))
			if msgErr != nil {
				logrus.Error(msgErr)
				return
			}
			logrus.Info("ws write===>", singleset)
		}
	}
}

func reader(conn *websocket.Conn, done chan struct{}) {
	defer conn.Close()
	defer close(done)
	for {
		_, _, err := conn.ReadMessage() //what is message type?
		if err != nil {
			logrus.Errorf("there is errors%s", err)
			return
		}
		//logrus.Info("ws read ====>",string(readMsg))
	}
}

func WsServer(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//升级get请求为webSocket协议
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	logrus.Info("----------------初始化WsServer--------------", ws.RemoteAddr())
	defer func() {
		ws.Close()
		logrus.Error("关闭WsServer====>")
	}()

	for {
		messageType, message, readErr := ws.ReadMessage()
		if readErr != nil || messageType == websocket.CloseMessage {
			logrus.Error("ws home read err====>", readErr)
		}
		//logrus.Info("ws home read ====>",string(message))

		msgErr := ws.WriteMessage(websocket.TextMessage, message)
		if msgErr != nil {
			logrus.Error("ws home write err====>", msgErr)
		}
		//logrus.Info("ws home write====>", string(message))
	}

}
