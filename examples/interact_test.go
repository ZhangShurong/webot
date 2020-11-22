package examples

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cutesdk/webot"
)

// Example 机器人接入
func ExampleEntry() {
	http.HandleFunc("/webot/entry", func(w http.ResponseWriter, req *http.Request) {
		bot := &webot.Bot{
			Token:      "V6s4MaHTslnPj1EDb2TcGB",
			AesKey:     "LyoclRcRYQhmT8V6RS1zE5DtwHhaFT3fuePXiFjJuYw",
			WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=b8f1e424-d48d-46cc-a2c7-d360c8e98b3d",
		}

		// 获取请求参数
		query := req.URL.Query()

		msgSign := query.Get("msg_signature")
		timestamp := query.Get("timestamp")
		nonce := query.Get("nonce")
		echoStr := query.Get("echostr")

		// 接入验证
		if echoStr != "" {
			decryptStr, err := bot.Valid(msgSign, timestamp, nonce, echoStr)
			if err != nil {
				fmt.Fprintf(w, "机器人接入验证失败："+err.Error())
				return
			}

			// 接入成功
			fmt.Fprintf(w, decryptStr)
			return
		}

		// 接收消息
		data, _ := ioutil.ReadAll(req.Body)

		msg, err := bot.DecryptMsg(msgSign, timestamp, nonce, data)
		if err != nil {
			fmt.Fprintf(w, "消息验证失败："+err.Error())
			return
		}

		// TODO: 消息处理

		// 主动推送消息
		bot.Text("机器人主动推送：").Send()

		// 被动回复消息
		rmsg, _ := bot.Text(fmt.Sprintf("接收到消息：%+v", msg)).Reply()

		fmt.Fprintf(w, string(rmsg))
	})

	http.ListenAndServe(":8097", nil)

	// Output: xxx
}
