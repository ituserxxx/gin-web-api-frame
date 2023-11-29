package router

import (
	"encoding/json"
	"fmt"
	"gin-web-api-ws-mqtt-frame/tools/config"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	AccessKey     = "047I4WS1-U51UBO6W-1J4BT21P-MF17IT99-92J8WIHU-944Q4KIW"
	Secret        = "035c73f7-bb6b-4889-a715-d9eb2d1925cc"
	MediaServerId = "U0hIvA6l1ubtIxVr"
)

type T2 struct {
	VideoChannelList []struct {
		Id                int    `json:"id"`
		MainId            string `json:"mainId"`
		MediaServerId     string `json:"mediaServerId"`
		Vhost             string `json:"vhost"`
		App               string `json:"app"`
		ChannelName       string `json:"channelName"`
		DepartmentId      string `json:"departmentId"`
		DepartmentName    string `json:"departmentName"`
		PDepartmentId     string `json:"pDepartmentId"`
		PDepartmentName   string `json:"pDepartmentName"`
		DeviceNetworkType string `json:"deviceNetworkType"`
		DeviceStreamType  string `json:"deviceStreamType"`
		MethodByGetStream string `json:"methodByGetStream"`
		VideoDeviceType   string `json:"videoDeviceType"`
		AutoVideo         bool   `json:"autoVideo"`
		AutoRecord        bool   `json:"autoRecord"`
		RecordSecs        int    `json:"recordSecs"`
		IpV4Address       string `json:"ipV4Address"`
		IpV6Address       string `json:"ipV6Address"`
		HasPtz            bool   `json:"hasPtz"`
		DeviceId          string `json:"deviceId"`
		ChannelId         string `json:"channelId"`
		RtpWithTcp        bool   `json:"rtpWithTcp"`
		DefaultRtpPort    bool   `json:"defaultRtpPort"`
		CreateTime        string `json:"createTime"`
		UpdateTime        string `json:"updateTime"`
		Enabled           bool   `json:"enabled"`
		NoPlayerBreak     bool   `json:"noPlayerBreak"`
		IsShareChannel    bool   `json:"isShareChannel"`
	} `json:"videoChannelList"`
	Request struct {
		Enabled   bool `json:"enabled"`
		PageIndex int  `json:"pageIndex"`
		PageSize  int  `json:"pageSize"`
		OrderBy   []struct {
			FieldName  string `json:"fieldName"`
			OrderByDir string `json:"orderByDir"`
		} `json:"orderBy"`
	} `json:"request"`
	Total int `json:"total"`
}
type CameraInfo struct {
	Id        string `json:"id" `
	Name      string `json:"name"`
	ChannelId string `json:"channel_id"`
	MainId    string `json:"main_id"`
	Status    int    `json:"status" dc:"摄像头状态  1：在线 2：离线" `
}

// GetCameraDevices 获取摄像头设备列表
func GetCameraDevices() (res []CameraInfo, err error) {
	res = make([]CameraInfo, 0)
	url := getGb28181Url() + "/MediaServer/GetVideoChannelList"
	resp, err := gb28181PostRequest(url, fmt.Sprintf(`{"secret":"%s","pageIndex":1,"pageSize":9999,"MediaSeverId":"%s","orderBy": [{"fieldName": "mediaServerId","orderByDir": 0}],"enabled":true}`, Secret, MediaServerId))
	if err != nil {
		return
	}
	var data T2
	err = json.Unmarshal(resp, &data)
	if err != nil {

		return
	}
	for _, d := range data.VideoChannelList {
		res = append(res, CameraInfo{
			MainId:    d.MainId,
			Id:        d.DeviceId,
			ChannelId: d.ChannelId,
			Name:      d.ChannelName,
			Status: func(mainId string) int {
				if GetCameraIsOnline(mainId) {
					return 1
				}
				return 0
			}(d.MainId),
		})
	}
	return
}

func GetCameraChannelState() (data map[string]int) {
	data = make(map[string]int, 0)
	devicesInfo, err := GetCameraDevices()

	if err != nil {
		return
	}
	for _, d := range devicesInfo {
		if d.Status == 1 {
			data[d.Id+d.ChannelId] = 1
		}
	}
	return
}

// CameraPtzControl command ：
// 上：1
// 左上：2
// 右上：3
// 下：4
// 左下：5
// 右下：6
// 左：7
// 右：8
// 聚焦+：9
// 聚焦-：10
func CameraPtzControl(deviceId, channelId, command string) {
	url := fmt.Sprintf("%s/SipGate/PtzCtrl", getGb28181Url())
	_, err := gb28181PostRequest(url, fmt.Sprintf(`{"secret":"%s","ptzCommandType":%s,"speed":100,"deviceId":"%s","channelId":"%s"}`,
		Secret, command, deviceId, channelId))
	if err != nil {
		fmt.Printf("CameraPtzControl err :%#v", err)
	}
	time.Sleep(500 * time.Millisecond)
	CameraPtzControlStop(deviceId, channelId)
	return
}

func CameraPtzControlStop(deviceId, channelId string) {
	url := fmt.Sprintf("%s/SipGate/PtzCtrl", getGb28181Url())
	_, _ = gb28181PostRequest(url, fmt.Sprintf(`{"secret":"%s","ptzCommandType":%d,"speed":100,"deviceId":"%s","channelId":"%s"}`,
		Secret, 0, deviceId, channelId))
	return
}

func gb28181PostRequest(url, body string) ([]byte, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	payload := strings.NewReader(body)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("AccessKey", AccessKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return io.ReadAll(resp.Body)
}

func GetVideoPlayUrl(mainId string) string {
	return fmt.Sprintf("%s/rtp/%s.live.flv", getGb28181Ws(), mainId)
}
func getGb28181Ws() string {
	return config.Get("Gb28181Ws")
}
func getGb28181Url() string {
	return config.Get("Gb28181Url")
}

func GetCameraIsOnline(mainId string) bool {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	url := fmt.Sprintf("%s/MediaServer/StreamLive?mediaServerId=%s&mainId=%s&secret=%s", getGb28181Url(), MediaServerId, mainId, Secret)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Add("AccessKey", AccessKey)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}
