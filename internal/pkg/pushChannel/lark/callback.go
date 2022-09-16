package lark

import (
	"encoding/json"
	"github.com/fastwego/feishu"
	"github.com/gin-gonic/gin"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/config"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/pushChannel/lark/messageTemplate"
	"github.com/lyleshaw/ospp-cr-bot/pkg/utils/log"
	"github.com/valyala/fasttemplate"
	"io/ioutil"
	"strings"
)

func Callback(c *gin.Context) {
	var req CallBackReq

	err := c.Bind(&req)
	larkCrypto := feishu.NewCrypto(FeishuConfig["EncryptKey"])
	decryptMsg, err := larkCrypto.GetDecryptMsg(req.Encrypt)
	if err != nil {
		log.Errorf("err: %s\n", err.Error())
		return
	}
	decryptMsgStr := string(decryptMsg)
	if strings.Contains(decryptMsgStr, "url_verification") { // 这里判断一下 callback 究竟是初始的验证还是消息回调
		var decodeReq DecodeCallBackReq

		err = json.Unmarshal(decryptMsg, &decodeReq)
		if err != nil {
			log.Errorf("err: %s\n", err.Error())
			return
		}
		callBackResp := CallBackResp{
			Challenge: decodeReq.Challenge,
		}
		log.Infof("callBackResp:%+v\n", callBackResp)
		c.JSON(200, callBackResp)
		return
	}
	var messageCallBackResp MessageCallBackResp
	err = json.Unmarshal(decryptMsg, &messageCallBackResp)
	if err != nil {
		log.Errorf("err: %s\n", err.Error())
		return
	}
	c.JSON(200, nil)

	// 发送消息
	var groupIdMsg string
	if messageCallBackResp.Event.Message.Content == "{\"text\":\"ID\"}" {
		t := fasttemplate.New(messageTemplate.SendIdMsg, "{{", "}}")
		groupIdMsg = t.ExecuteString(map[string]interface{}{
			"ChatID": messageCallBackResp.Event.Sender.SenderID.OpenID,
		})
	}
	if messageCallBackResp.Event.Message.ChatType == "group" {
		t := fasttemplate.New(messageTemplate.SendIdMsg, "{{", "}}")
		groupIdMsg = t.ExecuteString(map[string]interface{}{
			"ChatID": messageCallBackResp.Event.Message.ChatID,
		})
	}
	_, err = SendGroupMessage(messageCallBackResp.Event.Message.ChatID, groupIdMsg)
	if err != nil {
		log.Errorf("send message error: %+v", err)
	}
	return
}

func CardCallback(c *gin.Context) {
	reqStr, _ := ioutil.ReadAll(c.Request.Body)
	log.Infof("reqStr: %s", reqStr)
	if strings.Contains(string(reqStr), "url_verification") {
		var req DecodeCallBackReq
		_ = json.Unmarshal(reqStr, &req)
		log.Infof("req:%+v\n", req)
		c.JSON(200, req)
		return
	}

	var cardPostReq CardPostReq
	err := json.Unmarshal(reqStr, &cardPostReq)
	if err != nil {
		log.Errorf("err: %+v\n", err)
		log.Errorf("err: %s\n", err.Error())
		return
	}
	// 卡片消息
	log.Infof("cardPostReq:%+v\n", cardPostReq)
	log.Infof("pre msqQueue:%+v\n", config.MsgQueue)
	larkId, _ := config.QueryGithubIdByLarkId(cardPostReq.OpenId)
	if _, ok := config.MsgQueue[larkId+cardPostReq.Action.Value.Numbers+cardPostReq.Action.Value.Type]; ok {
		delete(config.MsgQueue, larkId+cardPostReq.Action.Value.Numbers+cardPostReq.Action.Value.Type)
	}
	log.Infof("after msqQueue:%+v\n", config.MsgQueue)
	c.JSON(200, nil)
}
