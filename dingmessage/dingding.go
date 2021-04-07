// https://developers.dingtalk.com/document/app/custom-robot-access
package dingmessage

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type EnumMsgType string

const (
	EnumMsgTypeText       EnumMsgType = "text"
	EnumMsgTypeLink       EnumMsgType = "link"
	EnumMsgTypeMarkdown   EnumMsgType = "markdown"
	EnumMsgTypeActionCard EnumMsgType = "actionCard"
	EnumMsgTypeFeedCard   EnumMsgType = "feedCard"
)

type DingClient struct {
	token string
}

func NewDingClient(token string) *DingClient {
	return &DingClient{token: token}
}

func (c *DingClient) SendTextMessage(content string, isAtAll bool) (string, error) {
	data := MessageRequest{
		MsgType: EnumMsgTypeText,
		Text: &Text{
			Content: content,
		},
		At: At{IsAtAll: isAtAll},
	}

	return sendDingMsgHandle(c.token, data)
}

func (c *DingClient) SendLinkMessage(title, text, picUrl, messageUrl string, isAtAll bool) (string, error) {
	data := MessageRequest{
		MsgType: EnumMsgTypeLink,
		Link: &Link{
			Title:      title,
			Text:       text,
			PicUrl:     picUrl,
			MessageUrl: messageUrl,
		},
		At: At{IsAtAll: isAtAll},
	}

	return sendDingMsgHandle(c.token, data)
}

func (c *DingClient) SendMarkdownMessage(title, markdownText string, isAtAll bool) (string, error) {
	data := MessageRequest{
		MsgType: EnumMsgTypeMarkdown,
		Markdown: &Markdown{
			Title: title,
			Text:  fmt.Sprintf("## %s \n> %s", title, markdownText),
		},
		At: At{IsAtAll: isAtAll},
	}

	return sendDingMsgHandle(c.token, data)
}

func (c *DingClient) SendActionCardMessage(actionCard *ActionCard, isAtAll bool) (string, error) {
	data := MessageRequest{
		MsgType:    EnumMsgTypeActionCard,
		ActionCard: actionCard,
		At:         At{IsAtAll: isAtAll},
	}

	return sendDingMsgHandle(c.token, data)
}

func (c *DingClient) SendFeedCardMessage(feedCard *FeedCard, isAtAll bool) (string, error) {
	data := MessageRequest{
		MsgType:  EnumMsgTypeFeedCard,
		FeedCard: feedCard,
		At:       At{IsAtAll: isAtAll},
	}

	return sendDingMsgHandle(c.token, data)
}

func sendDingMsgHandle(token string, data MessageRequest) (string, error) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("转json时候出现异常:%+v", err)
		return "", err
	}

	webHookUrl := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token)
	req, err := http.NewRequest("POST", webHookUrl, bufio.NewReader(bytes.NewBufferString(string(payloadBytes))))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("请求钉钉服务失败" + err.Error())
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("请求钉钉服务失败" + err.Error())
		return "", err
	}
	return string(respBody), nil
}

type MessageRequest struct {
	MsgType    EnumMsgType `json:"msgtype"`
	Markdown   *Markdown   `json:"markdown,omitempty"`
	Text       *Text       `json:"text,omitempty"`
	Link       *Link       `json:"link,omitempty"`
	ActionCard *ActionCard `json:"actionCard,omitempty"`
	FeedCard   *FeedCard   `json:"feedCard,omitempty"`
	At         At          `json:"at"`
}

type Text struct {
	Content string `json:"content"`
}

type Link struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ActionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	hideAvatar     string `json:"hideAvatar"`
	BtnOrientation string `json:"btnOrientation"`
	Btns           []Btn  `json:"btns"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
}

type FeedCard struct {
	Links []Link `json:"links"`
}

type At struct {
	IsAtAll bool `json:"isAtAll"`
}

type Btn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}
