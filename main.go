package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	client := &http.Client{}
	SignIn(client)
}

// SignIn 签到
func SignIn(client *http.Client) bool {
	//生成要访问的url
	urlStr := "https://www.hifiti.com/sg_sign.htm"
	cookie := os.Getenv("COOKIE")
	if cookie == "" {
		fmt.Println("COOKIE不存在，请检查是否添加")
		return false
	}

	//提交请求
	formData := url.Values{}

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(formData.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Cookie", cookie)
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//处理返回结果
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	buf, _ := ioutil.ReadAll(response.Body)
	fmt.Println("签到结果：")
	fmt.Println(string(buf))
	hasDing := os.Getenv("DINGDING_WEBHOOK")
	if hasDing != "" {
		dingding(string(buf))
	} else {
		fmt.Println("DINGDING_WEBHOOK 环境变量未定义，跳过通知步骤")
	}
	return strings.Contains(string(buf), "成功")
}

func dingding(result string) {
	// 构造要发送的消息
	message := struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: "HiFiNi：\n" + result,
		},
	}

	// 将消息转换为JSON格式
	messageJson, _ := json.Marshal(message)
	DINGDING_WEBHOOK := os.Getenv("DINGDING_WEBHOOK")
	// 发送HTTP POST请求
	resp, err := http.Post(DINGDING_WEBHOOK,
		"application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
