package webot

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cutesdk/webot/wxbizmsgcrypt"
	"github.com/idoubi/goutils"
)

// Bot 企业微信机器人
type Bot struct {
	Token          string      // 接入验证的token
	AesKey         string      // 接入验证的encodingAesKey
	WebhookURL     string      // 主动推送消息的地址
	chatids        []string    // 聊天对象ID
	mentionedUsers []string    // 提到的用户
	visibleUsers   []string    // 消息可见的用户
	actions        []MsgAction // 消息操作
	callbackID     string      // 消息操作回调ID
	reqMsg         *ReqMsg     // 请求消息
	respMsg        *RespMsg    // 响应消息
}

// ReqMsg 请求消息
type ReqMsg struct {
	From           ReqMsgFrom    `xml:"From"`
	WebhookURL     string        `xml:"WebhookUrl"`
	ChatID         string        `xml:"ChatId"`
	ChatType       string        `xml:"ChatType"`
	GetChatInfoURL string        `xml:"GetChatInfoUrl"`
	MsgID          string        `xml:"MsgId"`
	MsgType        string        `xml:"MsgType"`
	Text           ReqMsgText    `xml:"Text"`
	Attachment     MsgAttachment `xml:"Attachment"`
	Event          Event         `xml:"Event"`
}

// Event 请求消息来源
type Event struct {
	EventType string `xml:"EventType"`
}

// ReqMsgFrom 请求消息来源
type ReqMsgFrom struct {
	UserID string `xml:"UserId"`
	Name   string `xml:"Name`
	Alias  string `xml:"Alias"`
}

// ReqMsgText 请求消息内容
type ReqMsgText struct {
	Content string `xml:"Content"`
}

// RespMsg 响应消息
type RespMsg struct {
	MsgType       string          `xml:"MsgType" json:"msgtype"`
	ChatID        string          `xml:"-" json:"chatid,omitempty"`
	VisibleToUser string          `xml:"VisibleToUser" json:"visible_to_user"`
	Text          RespMsgText     `xml:"Text" json:"text"`
	Markdown      RespMsgMarkdown `xml:"Markdown" json:"markdown"`
	News          RespMsgNews     `xml:"-" json:"news"`
	Image         RespMsgImage    `xml:"-" json:"image"`
}

// RespMsgText 文本消息类型
type RespMsgText struct {
	Content       string         `xml:"Content" json:"content"`
	MentionList   *MentionedList `xml:"MentionedList,omitempty" json:"-"`
	MentionedList []string       `xml:"-" json:"mentioned_list"`
}

// RespMsgMarkdown markdown消息类型
type RespMsgMarkdown struct {
	Content     string           `xml:"Content" json:"content"`
	Attachment  *MsgAttachment   `xml:"Attachment,omitempty" json:"-"`
	Attachments []*MsgAttachment `xml:"-" json:"attachments,omitempty"`
}

// RespMsgNews 图文消息
type RespMsgNews struct {
	Articles []NewsArticle `xml:"-" json:"articles"`
}

// NewsArticle 图文项
type NewsArticle struct {
	Title       string `xml:"-" json:"title"`
	Description string `xml:"-" json:"description"`
	URL         string `xml:"-" json:"url"`
	PicURL      string `xml:"-" json:"picurl"`
}

// RespMsgImage 图片消息
type RespMsgImage struct {
	Base64 string `xml:"-" json:"base64"`
	MD5    string `xml:"-" json:"md5"`
}

// MsgAttachment markdown附加数据
type MsgAttachment struct {
	CallbackID string      `xml:"CallbackId" json:"callback_id"`
	Actions    []MsgAction `xml:"Actions" json:"actions"`
}

// MsgAction 操作
type MsgAction struct {
	Name        string `xml:"Name" json:"name"`
	Value       string `xml:"Value" json:"value"`
	Text        string `xml:"Text" json:"text"`
	Type        string `xml:"Type" json:"type"`
	BorderColor string `xml:"BorderColor" json:"border_color"`
	TextColor   string `xml:"TextColor" json:"text_color"`
	ReplaceText string `xml:"ReplaceText" json:"replace_text"`
}

// MentionedList 提到的人
type MentionedList struct {
	Item []string `xml:"Item" json:"mentioned_list"`
}

// Valid 接入验证
func (b *Bot) Valid(msgSign, timestamp, nonce, echoStr string) (string, error) {
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(b.Token, b.AesKey, "", wxbizmsgcrypt.XmlType)

	bt, err := wxcpt.VerifyURL(msgSign, timestamp, nonce, echoStr)
	if err != nil {
		return "", err
	}

	return string(bt), nil
}

// DecryptMsg 消息解密
func (b *Bot) DecryptMsg(msgSign, timestamp, nonce string, rawMsg []byte) (*ReqMsg, error) {
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(b.Token, b.AesKey, "", wxbizmsgcrypt.XmlType)

	bt, err := wxcpt.DecryptMsg(msgSign, timestamp, nonce, rawMsg)
	if err != nil {
		return nil, err
	}

	reqMsg := &ReqMsg{}
	if err := xml.Unmarshal(bt, &reqMsg); err != nil {
		return nil, err
	}

	b.reqMsg = reqMsg

	return reqMsg, nil
}

// Reset 重置
func (b *Bot) Reset() *Bot {
	return &Bot{
		Token:      b.Token,
		AesKey:     b.AesKey,
		WebhookURL: b.WebhookURL,
	}
}

// Target 发送消息对象
func (b *Bot) Target(chatids []string) *Bot {
	b.chatids = chatids
	return b
}

// Mention 艾特用户
func (b *Bot) Mention(users []string) *Bot {
	b.mentionedUsers = users
	return b
}

// Visible 指定用户可见
func (b *Bot) Visible(users []string) *Bot {
	b.visibleUsers = users
	return b
}

// Actions 设置操作
func (b *Bot) Actions(callbackID string, actions []MsgAction) *Bot {
	b.callbackID = callbackID
	b.actions = actions
	return b
}

// Text 构建文本消息
func (b *Bot) Text(text string) *Bot {
	respMsg := &RespMsg{
		MsgType: "text",
		Text: RespMsgText{
			Content:     text,
			MentionList: nil,
		},
	}
	if len(b.mentionedUsers) > 0 {
		respMsg.Text.MentionList = &MentionedList{
			Item: b.mentionedUsers,
		}
		respMsg.Text.MentionedList = b.mentionedUsers
	}

	b.respMsg = respMsg

	return b
}

// Markdown 构建markdown消息
func (b *Bot) Markdown(markdown string) *Bot {
	respMsg := &RespMsg{
		MsgType: "markdown",
		Markdown: RespMsgMarkdown{
			Content:    markdown,
			Attachment: nil,
		},
	}

	if b.callbackID != "" && len(b.actions) > 0 {
		respMsg.Markdown.Attachment = &MsgAttachment{
			CallbackID: b.callbackID,
			Actions:    b.actions,
		}
		respMsg.Markdown.Attachments = []*MsgAttachment{respMsg.Markdown.Attachment}
	}

	b.respMsg = respMsg

	return b
}

// News 构建图片消息
func (b *Bot) News(articles []NewsArticle) *Bot {
	respMsg := &RespMsg{
		MsgType: "news",
		News: RespMsgNews{
			Articles: articles,
		},
	}

	b.respMsg = respMsg

	return b
}

// Image 构建图片消息
func (b *Bot) Image(imgurl string) *Bot {
	res, err := http.Get(imgurl)
	if err != nil {
		fmt.Println("get image error:", err)
		return b
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read image error:", err)
		return b
	}

	// 计算文件base64
	base64Str := base64.StdEncoding.EncodeToString(body)

	// 计算文件md5
	md5Handler := md5.New()
	_, err = io.Copy(md5Handler, bytes.NewReader(body))
	if err != nil {
		fmt.Println("get md5 error:", err)
		return b
	}
	md5Value := md5Handler.Sum(nil)
	md5Str := fmt.Sprintf("%x", md5Value)

	respMsg := &RespMsg{
		MsgType: "image",
		Image: RespMsgImage{
			Base64: base64Str,
			MD5:    md5Str,
		},
	}

	b.respMsg = respMsg

	return b
}

// Reply 被动回复消息
func (b *Bot) Reply() ([]byte, error) {
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(b.Token, b.AesKey, "", wxbizmsgcrypt.XmlType)

	if len(b.chatids) > 0 {
		b.respMsg.ChatID = strings.Join(b.chatids, "|")
	}
	if len(b.visibleUsers) > 0 {
		b.respMsg.VisibleToUser = strings.Join(b.visibleUsers, "|")
	}

	bt, err := xml.Marshal(b.respMsg)
	if err != nil {
		return nil, err
	}

	b.Reset()

	timestamp := goutils.TimestampStr()
	nonce := goutils.NonceStr(16)
	msg, err := wxcpt.EncryptMsg(string(bt), timestamp, nonce)
	if len(msg) == 0 {
		return nil, errors.New("no reply msg")
	}

	return msg, nil
}

// Send 主动发送消息
func (b *Bot) Send() ([]byte, error) {
	if len(b.chatids) > 0 {
		b.respMsg.ChatID = strings.Join(b.chatids, "|")
	}
	if len(b.visibleUsers) > 0 {
		b.respMsg.VisibleToUser = strings.Join(b.visibleUsers, "|")
	}

	bt, err := json.Marshal(b.respMsg)
	if err != nil {
		return nil, err
	}

	b.Reset()

	res, err := http.Post(b.WebhookURL, "application/json", bytes.NewReader(bt))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
