package apiv1

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// type Notification struct {
// 	Message string `json:"message"`
// }
type Notification interface{}
type Connection struct {
	wsConn *websocket.Conn
	c      notifChan
}
type notifChan chan Notification

// global state
const (
	maxQueueLen = 100
)

var (
	conns  = connMap{v: map[string]Connection{}}
	addrs  = addrMap{v: map[string]uint{}}
	notifs = notifMap{v: map[uint][]*Notification{}}
)

type connMap struct {
	v   map[string]Connection // IP address : websocket conn
	mux sync.Mutex
}
type addrMap struct {
	v   map[string]uint // IP address : ID
	mux sync.Mutex
}
type notifMap struct {
	v   map[uint][]*Notification // ID : notif struct
	mux sync.Mutex
}

func addConn(uid uint, uaddr string, conn *websocket.Conn) notifChan {
	conns.mux.Lock()
	addrs.mux.Lock()
	defer conns.mux.Unlock()
	defer addrs.mux.Unlock()

	c := make(notifChan, maxQueueLen)

	conns.v[uaddr] = Connection{
		wsConn: conn,
		c:      c,
	}
	addrs.v[uaddr] = uid

	return c
}
func removeConn(uid uint, uaddr string) {
	conns.mux.Lock()
	addrs.mux.Lock()
	defer conns.mux.Unlock()
	defer addrs.mux.Unlock()

	delete(conns.v, uaddr)
	delete(addrs.v, uaddr)
}
func sendNotif(uid uint, n Notification) {
	for ipAddr, id := range addrs.v {
		if uid == id {
			conn := conns.v[ipAddr]
			conn.c <- n
		}
	}
}

type resultWS struct {
	Message string `json:"message"`
}

// TODO: fix CheckOrigin function
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func PostWs(context *Context, c echo.Context) error {
	user, err := getClientUser(context, c)
	fmt.Printf("->>>>>> %s\n", c.RealIP())
	if err != nil {
		// TODO
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("WebSocket connection failed: %s", err.Error()))
	}
	var postData = jsonmap{}
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// sendNotif(user.ID, Notification{Message: "hihihi"})
	sendNotif(user.ID, &postData)
	return c.JSON(http.StatusOK, resultWS{Message: "connection end"})
}
func GetWs(context *Context, c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("WebSocket upgrade failed: ", err.Error()))
	}
	defer ws.Close()

	user, err := getClientUser(context, c)
	fmt.Printf("->>>>>> %s\n", c.RealIP())
	if err != nil {
		// TODO
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("WebSocket connection failed: %s", err.Error()))
	}
	uid := user.ID
	ip := c.RealIP()

	ch := addConn(uid, ip, ws)
	stop := make(chan int)

	ws.SetCloseHandler(func(code int, text string) error {
		stop <- 1
		return nil
	})

	log.Infof("ws connection: %s", ip)

	for {
		// Write
		select {
		case notif := <-ch:
			err = ws.WriteJSON(notif)
			if err != nil {
				c.Logger().Error(err)
			}
		case <-stop:
			removeConn(uid, ip)
			break
		}

		// Read
		// _, msg, err := ws.ReadMessage()
		// if err != nil {
		// 	c.Logger().Error(err)
		// }
		// fmt.Printf("%s\n", msg)
	}

	return c.JSON(http.StatusOK, resultWS{Message: "connection end"})
}
