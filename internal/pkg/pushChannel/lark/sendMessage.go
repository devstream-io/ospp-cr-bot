package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fastwego/feishu"
	"log"
	"net/http"
)

func SendMessage(receiveID string, msgTemplate string) {
	var sendMessage SendMessageReq
	sendMessage.ReceiveID = receiveID
	sendMessage.Content = msgTemplate
	sendMessage.MsgType = "interactive"

	data, err := json.Marshal(sendMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantAccessToken, err := Atm.GetAccessToken()
	if err != nil {
		log.Println(err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, feishu.ServerUrl+"/open-apis/im/v1/messages?receive_id_type=chat_id", bytes.NewReader(data))
	_, err = FeishuClient.Do(request, tenantAccessToken)
}
