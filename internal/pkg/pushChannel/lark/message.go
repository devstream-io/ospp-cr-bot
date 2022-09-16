package lark

import (
	"encoding/json"
	"fmt"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/config"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/constants"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/pushChannel/lark/messageTemplate"
	"github.com/lyleshaw/ospp-cr-bot/pkg/utils/log"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	createMessageURL = "https://open.feishu.cn/open-apis/im/v1/messages"
)

func MsgTemplateUpgrade(msgTemplate string, receiver string) string {
	msgTemplate = strings.Replace(msgTemplate, messageTemplate.PrTitle1, fmt.Sprintf(messageTemplate.PrTitle2, receiver), -1)
	msgTemplate = strings.Replace(msgTemplate, messageTemplate.IssueTitle1, fmt.Sprintf(messageTemplate.IssueTitle2, receiver), -1)
	msgTemplate = strings.Replace(msgTemplate, messageTemplate.PrCommentTitle1, messageTemplate.PrCommentTitle2, -1)
	msgTemplate = strings.Replace(msgTemplate, messageTemplate.IssueCommentTitle1, messageTemplate.IssueCommentTitle2, -1)
	msgTemplate = strings.Replace(msgTemplate, messageTemplate.PrReviewTitle1, messageTemplate.PrReviewTitle2, -1)
	return msgTemplate
}

func SendGroupMessage(receiveID string, msgTemplate string) (*MessageItem, error) {
	var err error
	token, err := Atm.GetAccessToken()
	cli := &http.Client{}

	createReq := GenCreateMessageRequest(receiveID, msgTemplate, "interactive")
	reqBytes, err := json.Marshal(createReq)
	if err != nil {
		logrus.WithError(err).Errorf("failed to marshal")
		return nil, err
	}

	log.Infof("req: %+v", string(reqBytes))
	req, err := http.NewRequest("POST", createMessageURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		logrus.WithError(err).Errorf("new request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("receive_id_type", "chat_id")
	req.URL.RawQuery = q.Encode()

	var logID string
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create message failed, err=%v", err)
	}
	if resp != nil && resp.Header != nil {
		logID = resp.Header.Get("x-tt-logid")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("read body failed")
		return nil, err
	}

	createMessageResp := &CreateMessageResponse{}
	err = json.Unmarshal(body, createMessageResp)
	if err != nil {
		logrus.WithError(err).Errorf("failed to unmarshal")
		return nil, err
	}
	if createMessageResp.Code != 0 {
		logrus.Errorf("failed to create message, code: %v, msg: %v, log_id: %v", createMessageResp.Code, createMessageResp.Message, logID)
		return nil, fmt.Errorf("create message failed")
	}
	logrus.Infof("succeed create message, msg_id: %v", createMessageResp.Data.MessageID)
	return createMessageResp.Data, nil
}

func SendMessage(receiveID string, msgTemplate string) (*MessageItem, error) {
	var err error
	token, err := Atm.GetAccessToken()
	cli := &http.Client{}

	createReq := GenCreateMessageRequest(receiveID, msgTemplate, "interactive")
	reqBytes, err := json.Marshal(createReq)
	if err != nil {
		logrus.WithError(err).Errorf("failed to marshal")
		return nil, err
	}

	log.Infof("receiveID=%s", receiveID)

	req, err := http.NewRequest("POST", createMessageURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		logrus.WithError(err).Errorf("new request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("receive_id_type", "open_id")
	req.URL.RawQuery = q.Encode()

	var logID string
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create message failed, err=%v", err)
	}
	if resp != nil && resp.Header != nil {
		logID = resp.Header.Get("x-tt-logid")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("read body failed")
		return nil, err
	}

	createMessageResp := &CreateMessageResponse{}
	err = json.Unmarshal(body, createMessageResp)
	if err != nil {
		logrus.WithError(err).Errorf("failed to unmarshal")
		return nil, err
	}
	if createMessageResp.Code != 0 {
		logrus.Errorf("failed to create message, code: %v, msg: %v, log_id: %v", createMessageResp.Code, createMessageResp.Message, logID)
		return nil, fmt.Errorf("create message failed")
	}
	logrus.Infof("succeed create message, msg_id: %v", createMessageResp.Data.MessageID)
	return createMessageResp.Data, nil
}

func SendMessages(receiveID string, receivers []string, msgTemplate string, number string, gitHubEvent string) {
	for _, receiver := range receivers {
		if _, ok := config.MsgQueue[receiver+number+gitHubEvent]; !ok {
			config.MsgQueue[receiver+number+gitHubEvent] = constants.Unread1
			_, err := SendMessage(config.LarkMaps[receiver].Lark, msgTemplate)
			if err != nil {
				log.Errorf("send message failed, msgTemplate=%s, receivers=%+v, gitHubEvent=%s", msgTemplate, receivers, gitHubEvent)
				log.Errorf("send message failed, err=%v", err)
			}
			log.Infof("msgQ=%v", config.MsgQueue)
		} else {
			if config.MsgQueue[receiver+number+gitHubEvent] == constants.READ {
				// 已读，从队列中删除
				delete(config.MsgQueue, receiver+number+gitHubEvent)
				log.Infof("msgQ=%v", config.MsgQueue)
			}
		}
	}
	// 设置定时任务
	time.AfterFunc(constants.TimeUnread1, func() {
		TimeCheck1(receiveID, receivers, msgTemplate, number, gitHubEvent)
	})
}

func TimeCheck1(receiveID string, receivers []string, msgTemplate string, number string, gitHubEvent string) {
	for _, receiver := range receivers {
		if _, ok := config.MsgQueue[receiver+number+gitHubEvent]; ok {
			if config.MsgQueue[receiver+number+gitHubEvent] == constants.Unread1 {
				// 如果第一次未读，则调整状态为第二次未读，并重新发送消息
				config.MsgQueue[receiver+number+gitHubEvent] = constants.Unread2
				_, err := SendMessage(config.LarkMaps[receiver].Lark, msgTemplate)
				if err != nil {
					log.Errorf("send message failed, msgTemplate=%s, receivers=%+v, gitHubEvent=%s", msgTemplate, receivers, gitHubEvent)
					log.Errorf("send message failed, err=%v", err)
				}
				log.Infof("msgQ=%v", config.MsgQueue)
			} else if config.MsgQueue[receiver+number+gitHubEvent] == constants.READ {
				// 已读，从队列中删除
				delete(config.MsgQueue, receiver+number+gitHubEvent)
				log.Infof("msgQ=%v", config.MsgQueue)
			}
		}
	}
	// 设置定时任务
	time.AfterFunc(constants.TimeUnread2, func() {
		TimeCheck2(receiveID, receivers, msgTemplate, number, gitHubEvent)
	})
}

func TimeCheck2(receiveID string, receivers []string, msgTemplate string, number string, gitHubEvent string) {
	for _, receiver := range receivers {
		if _, ok := config.MsgQueue[receiver+number+gitHubEvent]; ok {
			if config.MsgQueue[receiver+number+gitHubEvent] == constants.Unread2 {
				// 如果第二次未读，则调整状态为第三次未读，并发送消息给上级
				config.MsgQueue[receiver+number+gitHubEvent] = constants.Unread3
				msgTemplate = MsgTemplateUpgrade(msgTemplate, receiver)
				if config.LarkMaps[receiver].Boss == "0" {
					_, err := SendGroupMessage(receiveID, msgTemplate)
					if err != nil {
						log.Errorf("send message error: %+v", err)
					}
					continue
				}
				_, err := SendMessage(config.LarkMaps[config.LarkMaps[receiver].Boss].Lark, msgTemplate)
				if err != nil {
					log.Errorf("send message failed, msgTemplate=%s, receivers=%+v, gitHubEvent=%s", msgTemplate, receivers, gitHubEvent)
					log.Errorf("send message failed, err=%v", err)
				}
				log.Infof("msgQ=%v", config.MsgQueue)
			} else if config.MsgQueue[receiver+number+gitHubEvent] == constants.READ {
				// 已读，从队列中删除
				delete(config.MsgQueue, receiver+number+gitHubEvent)
				log.Infof("msgQ=%v", config.MsgQueue)
			}
		}
	}
	// 设置定时任务
	time.AfterFunc(constants.TimeUnread3, func() {
		TimeCheck3(receiveID, receivers, msgTemplate, number, gitHubEvent)
	})
}

func TimeCheck3(receiveID string, receivers []string, msgTemplate string, number string, gitHubEvent string) {
	for _, receiver := range receivers {
		if _, ok := config.MsgQueue[receiver+number+gitHubEvent]; ok {
			if config.MsgQueue[receiver+number+gitHubEvent] == constants.Unread3 {
				msgTemplate = MsgTemplateUpgrade(msgTemplate, receiver)
				// 如果第三次未读，则从队列清除，并发送群消息
				_, err := SendGroupMessage(receiveID, msgTemplate)
				if err != nil {
					log.Errorf("send message error: %+v", err)
				}
				delete(config.MsgQueue, receiver+number+gitHubEvent)
				log.Infof("msgQ=%v", config.MsgQueue)
			} else if config.MsgQueue[receiver+number+gitHubEvent] == constants.READ {
				delete(config.MsgQueue, receiver+number+gitHubEvent)
				log.Infof("msgQ=%v", config.MsgQueue)
			}
		}
	}
}

func GenCreateMessageRequest(chatID, content, msgType string) *CreateMessageRequest {
	return &CreateMessageRequest{
		ReceiveID: chatID,
		Content:   content,
		MsgType:   msgType,
	}
}
