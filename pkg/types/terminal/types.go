package terminal

import (
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalSession struct {
	WsConn   *websocket.Conn
	SizeChan chan remotecommand.TerminalSize
	DoneChan chan struct{}
}

type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}
