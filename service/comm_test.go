package service

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestMqttclient(t *testing.T) {
	var m = &MQClient{
		Url:      "127.0.0.1:1883",
		UserName: "garden",
		Password: "smartGarden@#",
		ClientID: "xxxx",
	}
	opts := MQTT.NewClientOptions()
	opts.AddBroker(m.Url)
	opts.SetClientID(m.ClientID)
	opts.SetUsername(m.UserName)
	opts.SetPassword(m.Password)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetKeepAlive(30 * time.Second) // 设置心跳间隔
	opts.SetAutoReconnect(true)
	//全局MQTT pub消息处理
	opts.SetDefaultPublishHandler(handleFunc)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLossHandler
	opts.OnReconnecting = reconnectHandler
	opts.ConnectTimeout = 5 * time.Second
	m.Client = MQTT.NewClient(opts)
	token := m.Client.Connect()
	fmt.Println("wait ", time.Now())
	if token.Wait() {
		fmt.Println("wait finish ", time.Now())
		return
	}

	tk := time.NewTicker(5 * time.Second)
	for range tk.C {
		m.Pub("sg_utek/xxx/from_dev", "")

	}

}
func TestWebscoketClient(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8077", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("received: %s\n", message)
		}
	}()
	for {
		println(111111111)
		err := c.WriteMessage(websocket.TextMessage, []byte("xxxxx"))
		if err != nil {
			log.Println("write:", err)
			return
		}
		println(22222)
		time.Sleep(2 * time.Second)
	}
}
