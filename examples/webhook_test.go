package examples

import (
	"fmt"

	"github.com/cutesdk/webot"
)

var bot *webot.Bot

// ExampleSendText 发送文本消息
func ExampleSendText() {
	msg := "这是一条需要<@adam>关注的文本消息"
	res, err := bot.Text(msg).Send()
	fmt.Printf("res:%s, err:%v\n", res, err)
	// Output: xxx
}

// ExampleSendTextMetion 发送文本消息提醒指定用户
func ExampleSendTextMethon() {
	msg := "这是一条要提示用户的文本消息"
	res, err := bot.Mention([]string{"adam"}).Text(msg).Send()
	fmt.Printf("res:%s, err:%v\n", res, err)
	// Output: xxx
}

// ExampleSendTextVisible 发送文本消息仅指定用户可见
func ExampleSendTextVisible() {
	msg := "这是一条仅指定用户可见的文本消息"
	res, err := bot.Target([]string{"wrkSFfCgAAIfjQ7pSijsAAx_KpflhDBw"}).Visible([]string{"adam"}).Text(msg).Send()
	fmt.Printf("res:%s, err:%v\n", res, err)
	// Output: xxx
}

// ExampleSendMarkdown 发送markdown消息
func ExampleSendMarkdown() {
	msg := `这是一条 Markdown 消息
	> 机器人主动推送测试
	*推送结束*`
	res, err := bot.Markdown(msg).Send()
	fmt.Printf("推送结果: %s, %v\n", res, err)
	// Output: xxx
}

// ExampleSendMarkdownWithActions
func ExampleSendMarkdownWithActions() {
	msg := `这是一条带操作按钮的 Markdown 消息
	> 请确认你要提交吗？`

	callbackID := "testwebot"
	actions := []webot.MsgAction{}
	actions = append(actions, webot.MsgAction{
		Text:        "确认",
		Name:        "acts",
		Value:       "confirm",
		ReplaceText: "已确认",
		Type:        "button",
		TextColor:   "2EAB49",
		BorderColor: "2EAB49",
	})
	actions = append(actions, webot.MsgAction{
		Text:        "取消",
		Name:        "acts",
		Value:       "cancel",
		ReplaceText: "已取消",
		Type:        "button",
	})
	res, err := bot.Actions(callbackID, actions).Markdown(msg).Send()
	fmt.Printf("推送结果: %s, %v\n", res, err)
	// Output: xxx
}

// ExampleSendNews 发送图文消息
func ExampleSendNews() {
	articles := []webot.NewsArticle{}
	articles = append(articles, webot.NewsArticle{
		Title:       "这是图文标题",
		Description: "这是图文描述",
		URL:         "http://mmbiz.qpic.cn/mmbiz_png/n3aTUWcsxPlJFSUwib1yNVtWrNqPEbHsDt32Px4FfhRQutAy8on7nianPR0iarNic5BZv2Kq9C47CeBYb3VXrCMWZw/0?wx_fmt=png",
		PicURL:      "http://mmbiz.qpic.cn/mmbiz_png/n3aTUWcsxPlJFSUwib1yNVtWrNqPEbHsDt32Px4FfhRQutAy8on7nianPR0iarNic5BZv2Kq9C47CeBYb3VXrCMWZw/0?wx_fmt=png",
	})

	res, err := bot.News(articles).Send()
	fmt.Printf("推送结果：%s, %v\n", res, err)
	// Output: xxx
}

// ExampleSendImage 发送图片
func ExampleSendImage() {
	imgurl := "http://mmbiz.qpic.cn/mmbiz_png/n3aTUWcsxPlJFSUwib1yNVtWrNqPEbHsDt32Px4FfhRQutAy8on7nianPR0iarNic5BZv2Kq9C47CeBYb3VXrCMWZw/0?wx_fmt=png"

	bot.Text("请点击查看大图").Send()
	res, err := bot.Image(imgurl).Send()
	fmt.Printf("推送结果：%s, %v\n", res, err)
	// Output: xxx
}

func init() {
	bot = &webot.Bot{
		WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=b8f1e424-d48d-46cc-a2c7-d360c8e98b3d",
	}
}
