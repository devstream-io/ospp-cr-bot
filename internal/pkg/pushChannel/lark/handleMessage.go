package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fastwego/feishu"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type HandleMessageReq struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type HandleMessageResp struct {
	Challenge string `json:"challenge"`
}

func HandleMessage(c *gin.Context) {
	// 加解密处理器
	dingCrypto := feishu.NewCrypto(FeishuConfig["EncryptKey"])

	// Post Body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	log.Printf(string(body))

	msgJson := struct {
		Encrypt string `json:"encrypt"`
	}{}
	err = json.Unmarshal(body, &msgJson)
	if err != nil {
		return
	}

	decryptMsg, err := dingCrypto.GetDecryptMsg(msgJson.Encrypt)
	if err != nil {
		return
	}

	eventJson := struct {
		Type  string `json:"type"`
		Event struct {
			Type   string `json:"type"`
			OpenId string `json:"open_id"`
			Text   string `json:"text"`
		} `json:"event"`
	}{}

	err = json.Unmarshal(decryptMsg, &eventJson)
	if err != nil {
		return
	}

	switch eventJson.Type {
	case "url_verification":
		// 响应 challenge
		_, _ = c.Writer.Write(decryptMsg)
		log.Println(string(decryptMsg))
		return
	}

	switch eventJson.Event.Type {
	case "message":
		// 响应 消息
		_, _ = c.Writer.Write(decryptMsg)
		log.Println(string(decryptMsg))

		replyTextMsg := struct {
			OpenId  string `json:"open_id"`
			MsgType string `json:"msg_type"`
			Content struct {
				Text string `json:"text"`
			} `json:"content"`
		}{
			OpenId:  eventJson.Event.OpenId,
			MsgType: "text",
			Content: struct {
				Text string `json:"text"`
			}{Text: eventJson.Event.Text},
		}

		data, err := json.Marshal(replyTextMsg)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantAccessToken, err := Atm.GetAccessToken()
		if err != nil {
			log.Println(err)
			return
		}

		request, err := http.NewRequest(http.MethodPost, feishu.ServerUrl+"/open-apis/message/v4/send/", bytes.NewReader(data))
		resp, err := FeishuClient.Do(request, tenantAccessToken)
		fmt.Println(string(resp), err)
	}
}
