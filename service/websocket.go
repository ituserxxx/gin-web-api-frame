package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Socket    *websocket.Conn
	Send      chan []byte
	Type      string
	Id        string
	DevNumber string
}

var Manager = ClientManager{
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}

// WSStart is  项目运行前, 协程开启start -> go Manager.Start()
func (manager *ClientManager) WSStart() {
	for {
		select {
		case conn := <-Manager.Register:
			Manager.Clients[conn.Id] = conn
		case conn := <-Manager.Unregister:
			if _, ok := Manager.Clients[conn.Id]; ok {
				close(conn.Send)
				delete(Manager.Clients, conn.Id)
			}
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			_ = c.Socket.Close()

			break
		}
		fmt.Println(string(message))
		if len(message) == 0 {
			c.Send <- message
			continue
		}
		fmt.Println(string(message))
	}
}
func (c *Client) Write() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
func WsHandler(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 校验参数

	//可以添加用户信息验证
	client := &Client{
		Socket: conn,
		Send:   make(chan []byte),
		Id:     fmt.Sprint(time.Now().UnixNano()),
		Type:   c.DefaultQuery("t", "web"),
	}

	Manager.Register <- client
	go client.Read()
	go client.Write()
}

// RealDataToWebData 推送给 web 的数据类型
type PushRealDataToWebData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func SendToWeb(dccNumber string, data PushRealDataToWebData) {
	for _, c := range Manager.Clients {
		if (c.Type == "web" && c.DevNumber == dccNumber) || data.Type == "garden_list" {
			b, _ := json.Marshal(data.Data)
			c.Send <- b
		}
	}
}
