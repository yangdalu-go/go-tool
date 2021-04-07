package dingmessage

import (
	"fmt"
	"testing"
)

func TestDingClient_SendMarkdownMessage(t *testing.T) {
	cli := NewDingClient("")
	cli.SendMarkdownMessage("自动部署消息测试", fmt.Sprintf("提交者：%s  \n服务：%s  \n分支：%s  ", "测试", "mrk-yudong-message", "production"), false)
}
