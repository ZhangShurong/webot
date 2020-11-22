## Webot 

企业微信机器人交互SDK

## 创建机器人

在企业微信选择一个群聊，右键点击添加机器人，即可创建一个具备基本消息推送能力的机器人。

![](https://imgkr2.cn-bj.ufileos.com/cf7d98da-78b8-49df-a9c8-d4e21a60bb00.png?UCloudPublicKey=TOKEN_8d8b72be-579a-4e83-bfd0-5f6ce1546f13&Signature=LvWEMJ9%252FqZsjzVl6ykcP2Wc0pFM%253D&Expires=1606147094)

## 下载SDK


```shell
go get -u github.com/cutesdk/webot
```

## 机器人推送消息

使用机器人推送消息非常简单，只需要在创建完机器人之后，拿到机器人的 Webhook地址，使用 SDK 创建一个 Bot 对象，通过链式调用完成消息推送。

### 基本使用

```go
package main

import (
	"github.com/cutesdk/webot"
)

func main() {
	bot := &webot.Bot{
		WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=b8f1e424-d48d-46cc-a2c7-d360c8e98b3d",
	}

	bot.Text("你好，我是机器人").Send()
}
```

运行上面的代码，添加了机器人的群就会收到消息：

![](https://imgkr2.cn-bj.ufileos.com/7b79cac0-b95c-4cca-95cd-b63747e67af5.png?UCloudPublicKey=TOKEN_8d8b72be-579a-4e83-bfd0-5f6ce1546f13&Signature=ayWUENfkxbzxMgHUCGl8jMeVsSM%253D&Expires=1606146549)

### 发送 markdown 类型消息

```go
	mdmsg := `*这是一条Markdown消息*
> 这是引用文本
- 这是列表1
- 这是列表2	
`
	bot.Markdown(mdmsg).Send()
```

运行上面的代码，机器人会发送 markdown 消息：

![](https://imgkr2.cn-bj.ufileos.com/31e121dc-1a5e-403a-9f67-0bb739567352.png?UCloudPublicKey=TOKEN_8d8b72be-579a-4e83-bfd0-5f6ce1546f13&Signature=1zMS8jk74UFUKKyCzoc0%252Fibc5J4%253D&Expires=1606146568)


### 发送带操作按钮的 markdown 消息

```go
mdmsg := `这是一条带操作按钮的 Markdown 消息
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

bot.Actions(callbackID, actions).Markdown(mdmsg).Send()
```

运行上面的代码，机器人发送带有操作按钮的 markdown 消息：

![](https://imgkr2.cn-bj.ufileos.com/a5ee8428-6eee-4cb2-863f-9b61a7f44970.png?UCloudPublicKey=TOKEN_8d8b72be-579a-4e83-bfd0-5f6ce1546f13&Signature=LjpS%252BRDe6dCOytmgRrJTc0FAuVs%253D&Expires=1606146586)


点击操作按钮，机器人会给回调地址推送事件消息，开发者可以在回调事件里面进行相应的处理，实现复杂的机器人交互逻辑。

### 提醒谁看

```go
msg := "这是一条要提示用户<@adam>的文本消息"
bot.Mention([]string{"@all"}).Text(msg).Send()
```

运行上面代码，机器人会发送消息，并艾特指定的用户：

![](https://imgkr2.cn-bj.ufileos.com/764d4c09-0551-4063-a265-907ffb4adb91.png?UCloudPublicKey=TOKEN_8d8b72be-579a-4e83-bfd0-5f6ce1546f13&Signature=yQ1omK%252B65cnwT3ZlWg1XdSbapVI%253D&Expires=1606146820)


在文本中艾特人用 `<@userid>`，艾特全体成员用 `@all`

### 指定用户可见

```go
msg := "这是一条仅指定用户可见的文本消息"
bot.Target([]string{"wrkSFfCgAAIfjQ7pSijsAAx_KpflhDBw"}).Visible([]string{"adam"}).Text(msg).Send()
```

运行上面的代码，机器人会发送消息到群聊，只有指定的群聊中的指定用户可以看到消息：

![](https://imgkr2.cn-bj.ufileos.com/d8c2a9c1-cb6d-4b08-a7bd-dace3e5a4182.png?UCloudPublicKey=TOKEN_8d8b72be-579a-4e83-bfd0-5f6ce1546f13&Signature=9AzHtmFIo5aDSoynL84TuKVjsgE%253D&Expires=1606146846)


`Target` 用于指定一个或多个群聊的 `chatid` ，每个群有唯一的 `chatid` ，在配置了机器人接收消息后，才能获取到群聊的 `chatid`。

## 机器人接收消息

目前公网版本的企业微信机器人并不支持接收消息配置，暂不能使用机器人完成复杂的交互功能。请关注企业微信官网查看更新动态。