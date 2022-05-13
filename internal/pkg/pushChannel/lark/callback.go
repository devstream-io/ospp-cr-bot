package lark

import (
	"encoding/json"
	"fmt"
	"github.com/fastwego/feishu"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasttemplate"
	"strings"
)

func Callback(c *gin.Context) {
	var req CallBackReq

	err := c.Bind(&req)
	larkCrypto := feishu.NewCrypto(FeishuConfig["EncryptKey"])
	decryptMsg, err := larkCrypto.GetDecryptMsg(req.Encrypt)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}
	decryptMsgStr := string(decryptMsg)
	if strings.Contains(decryptMsgStr, "url_verification") { // 这里判断一下 callback 究竟是初始的验证还是消息回调
		var decodeReq DecodeCallBackReq

		err = json.Unmarshal(decryptMsg, &decodeReq)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}
		callBackResp := CallBackResp{
			Challenge: decodeReq.Challenge,
		}
		fmt.Printf("callBackResp:%+v\n", callBackResp)
		c.JSON(200, callBackResp)
		return
	}
	var messageCallBackResp MessageCallBackResp
	err = json.Unmarshal(decryptMsg, &messageCallBackResp)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}
	c.JSON(200, nil)

	// 发送消息
	t := fasttemplate.New(SEND_GROUP_ID, "{{", "}}")
	groupIdMsg := t.ExecuteString(map[string]interface{}{
		"ChatID": messageCallBackResp.Event.Message.ChatID,
	})
	SendMessage(messageCallBackResp.Event.Message.ChatID, groupIdMsg)
	return
}
