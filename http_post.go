package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//用于对钉钉自定义机器人（加签安全认证）推送Text格式的消息

const (
	secret      = "your secret"
	baseWebhook = "your dingtalk-rebot webhook"
)

func main() {
	httpDo("your message")
}

func httpDo(message string) {
	//生成sign
	timestamp := time.Now().UnixMilli()
	string_sign := fmt.Sprintf("%d\n%s", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(string_sign))
	hashData := mac.Sum(nil)
	sign := url.QueryEscape(base64.StdEncoding.EncodeToString(hashData))

	client := &http.Client{}
	content := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": message,
		},
	}
	contentData, _ := json.Marshal(content)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s&timestamp=%d&sign=%s", baseWebhook, timestamp, sign), strings.NewReader(string(contentData)))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}
