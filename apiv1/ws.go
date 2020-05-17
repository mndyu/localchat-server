package apiv1

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{}

// global state
var (
	conns       connMap
	userIPAddrs addrMap
	notifQueues notifMap
)

type notif struct {
	message string
}

type connMap struct {
	conns map[string]*websocket.Conn // IP address : websocket conn
	lock  sync.Mutex
}
type addrMap struct {
	userIPAddrs map[string][]uint // IP address : IDs
	lock        sync.Mutex
}
type notifMap struct {
	notifQueues map[uint][]*notif // ID : notif struct
	lock        sync.Mutex
}

func addConn() {

}
func sendNotif() {

}

func PostWs(context *Context, c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
