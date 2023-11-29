package service

import (
	"fmt"
	"gin-web-api-ws-mqtt-frame/tools/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var MqServiceCli = &MQClient{}

// 消息质量
var qos byte = 1

type MQClient struct {
	Client   MQTT.Client
	Url      string
	UserName string
	Password string
	ClientID string
}

var msgChan chan []byte

var subTopic = "sg_utek/+/from_dev"
var pubToic = "sg_utek/%s/from_cloud"

type comm struct {
	Type string `json:"type"`
}

// 土壤数据
type turangData struct {
	comm
	Key  int `json:"key"`
	Data []struct {
		Smo  float32 `json:"smo"`
		Trwd float32 `json:"trwd"`
		Cond float32 `json:"cond"`
	}
}
type environData struct {
	comm
	Key  string `json:"key"`
	Data struct {
		Noise     float32 `json:"noise"`
		PM25      float32 `json:"PM2.5"`
		PM10      float32 `json:"PM10"`
		Temp      float32 `json:"temp"`
		Humi      float32 `json:"humi"`      //空气湿度 float
		Ill       float32 `json:"ill"`       //光照度 float
		Ppt       float32 `json:"ppt"`       //降雨量 float
		Windspeed float32 `json:"windspeed"` //风速 float
		Wdir      float32 `json:"wdir"`      //风向 float
	} `json:"data"`
}
type switchData struct {
	comm
	Key  int `json:"key"`  // 1-8
	Data int `json:"data"` //0为关，1为开 int
}

func MqttServerRun() {
	MqServiceCli = &MQClient{
		UserName: config.Get("MQ_USER"),
		Password: config.Get("MQ_PASSWORD"),
		Url:      fmt.Sprintf("tcp://%s:%s", config.Get("MQ_IP"), config.Get("MQ_PORT")),
		ClientID: fmt.Sprintf("garden_cloud_%d", time.Now().UnixNano()),
	}
	err := MqServiceCli.Init()
	if err != nil {
		fmt.Println("mqtt connect err ", err)
	}
	timer := time.NewTimer(5 * time.Second)
	msgChan = make(chan []byte, 10)

	for {
		select {
		case <-timer.C:
			if !IsMqttConnected() {
				fmt.Println("reconnect ", MqServiceCli.Url)
				err = MqServiceCli.Init()
				if err != nil {
					fmt.Println("mqtt reconnect err ", err)
					timer.Reset(5 * time.Second)
				}
			}
		case data := <-msgChan:
			MqServiceCli.Pub("", data)

		}
	}
}

// Init 初始化
func (m *MQClient) Init() error {
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
		return token.Error()
	}
	return nil
}

// 接收上报逻辑处理
var handleFunc MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Println(msg.Topic(), string(msg.Payload()))

}

var connectHandler MQTT.OnConnectHandler = func(_ MQTT.Client) {
	fmt.Println("mqtt connected!")
	MqServiceCli.Client.Subscribe(subTopic, 1, handleFunc)

}
var connectLossHandler MQTT.ConnectionLostHandler = func(_ MQTT.Client, err error) {
	fmt.Println("mqtt disconnected! ", err)
}
var reconnectHandler MQTT.ReconnectHandler = func(client MQTT.Client, options *MQTT.ClientOptions) {
	fmt.Println("reconnect !! ", time.Now())
}

func (m *MQClient) Pub(topic string, data interface{}) {
	if IsMqttConnected() {
		token := m.Client.Publish(topic, qos, false, data)
		if token.Wait() && token.Error() != nil {
		}
	}
}

func IsMqttConnected() bool {
	if MqServiceCli == nil || MqServiceCli.Client == nil {
		return false
	}
	return MqServiceCli.Client.IsConnectionOpen()
}

func (m *MQClient) SendToDcc(dccNumber string, data interface{}) {
	if IsMqttConnected() {
		token := m.Client.Publish(fmt.Sprintf(pubToic, dccNumber), qos, false, data)
		if token.Wait() && token.Error() != nil {
		}
	}
}
