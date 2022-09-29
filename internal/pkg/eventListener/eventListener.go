package eventListener

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/config"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/constants"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/pushChannel/lark/messageTemplate"
	"github.com/lyleshaw/ospp-cr-bot/pkg/utils/log"
	"github.com/valyala/fasttemplate"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/pushChannel/lark"
)

// Event defines a GitHub hook event type
type Event string

// GitHub hook types
const (
	IssueCommentEvent             Event = "issue_comment"
	IssuesEvent                   Event = "issues"
	PullRequestEvent              Event = "pull_request"
	PullRequestReviewEvent        Event = "pull_request_review"
	PullRequestReviewCommentEvent Event = "pull_request_review_comment"
)

func getAtPeopleFromComment(comment string) (res []string) {
	compileRegex := regexp.MustCompile("@(\\S*?) ")
	data := compileRegex.FindAllStringSubmatch(comment, -1)
	for _, v := range data {
		res = append(res, v[1])
	}
	return
}

// GitHubWebHook .
// @router /github/webhook [POST]
func GitHubWebHook(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")
	if event == "" {
		return
	}

	body, _ := ioutil.ReadAll(c.Request.Body)
	bodyStr := string(body)
	bodyStr = bodyStr[8:]
	bodyStr, _ = url.QueryUnescape(bodyStr)

	gitHubEvent := Event(event)

	log.Infof("Event Type: %s\n", gitHubEvent)
	switch gitHubEvent {
	case IssueCommentEvent:
		var req IssueCommentPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)

		t := fasttemplate.New(messageTemplate.IssueCommentMsg, "{{", "}}")

		// 正则拿到提及人
		var requestedReviewers string
		atPeople := getAtPeopleFromComment(req.Comment.Body)
		for _, v := range atPeople {
			if _, ok := config.LarkMaps[v]; ok {
				requestedReviewers += "<at id=" + config.LarkMaps[v].Lark + "></at>"
			}
		}

		receiveID, err := config.QueryReceiveIDByRepo(req.Repository.FullName)
		if err != nil {
			return
		}

		var issueCommentMsg string
		// 如果没有 Assignees 直接发到群
		if len(atPeople) == 0 {
			issueCommentMsg = t.ExecuteString(map[string]interface{}{
				"Title":              messageTemplate.IssueCommentTitle2,
				"Time":               req.Issue.CreatedAt.Format("2006年01月02日 15:04:05"),
				"IssueTitle":         req.Issue.Title,
				"Login":              req.Issue.User.Login,
				"IssueNumber":        strconv.FormatInt(req.Issue.Number, 10),
				"IssueContent":       req.Comment.Body,
				"RequestedReviewers": "无",
				"IssueURL":           req.Issue.HTMLURL,
			})

			_, err = lark.SendGroupMessage(receiveID, issueCommentMsg)
			if err != nil {
				log.Errorf("send message error: %+v", err)
			}
			return
		}

		issueCommentMsg = t.ExecuteString(map[string]interface{}{
			"Title":              messageTemplate.IssueCommentTitle1,
			"Time":               req.Issue.CreatedAt.Format("2006年01月02日 15:04:05"),
			"IssueTitle":         req.Issue.Title,
			"Login":              req.Issue.User.Login,
			"IssueNumber":        strconv.FormatInt(req.Issue.Number, 10),
			"IssueContent":       req.Comment.Body,
			"RequestedReviewers": requestedReviewers,
			"IssueURL":           req.Issue.HTMLURL,
		})

		number := strconv.FormatInt(req.Issue.Number, 10)

		// 简化版的 SendMessages
		for _, receiver := range atPeople {
			if _, ok := config.MsgQueue[receiver+number+string(gitHubEvent)]; !ok {
				config.MsgQueue[receiver+number+string(gitHubEvent)] = constants.Unread1
				msg, err := lark.SendMessage(config.LarkMaps[receiver].Lark, issueCommentMsg)
				if err != nil {
					log.Errorf("send message failed, msg=%v", msg)
					log.Errorf("send message failed, err=%v", err)
				}
				log.Infof("msgQ=%v", config.MsgQueue)
			} else {
				if config.MsgQueue[receiver+number+string(gitHubEvent)] == constants.READ {
					// 已读，从队列中删除
					delete(config.MsgQueue, receiver+number+string(gitHubEvent))
					log.Infof("msgQ=%v", config.MsgQueue)
				}
			}
		}

		// 定时未读则转发群聊
		time.AfterFunc(config.CommentUnread, func() {
			issueCommentMsg = lark.MsgTemplateUpgrade(issueCommentMsg)
			for _, receiver := range atPeople {
				if _, ok := config.MsgQueue[receiver+number+string(gitHubEvent)]; ok {
					if config.MsgQueue[receiver+number+string(gitHubEvent)] == constants.Unread1 {
						_, err = lark.SendGroupMessage(receiveID, issueCommentMsg)
						if err != nil {
							log.Errorf("send message error: %+v", err)
						}
						delete(config.MsgQueue, receiver+number+string(gitHubEvent))
						log.Infof("msgQ=%v", config.MsgQueue)
					}
				}
			}
		})
	case IssuesEvent:
		var req IssuesPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)

		t := fasttemplate.New(messageTemplate.IssueMsg, "{{", "}}")

		receiveID, err := config.QueryReceiveIDByRepo(req.Repository.FullName)
		if err != nil {
			return
		}

		var assigneesToAt string
		assignees := make([]string, len(req.Issue.Assignees))
		for i, v := range req.Issue.Assignees {
			if _, ok := config.LarkMaps[v.Login]; ok {
				assigneesToAt += "<at id=" + config.LarkMaps[v.Login].Lark + "></at>"
				assignees[i] = v.Login
			}
		}

		var issueMsg string
		// 如果没有 Assignees 直接发到群
		if len(assignees) == 0 {
			issueMsg = t.ExecuteString(map[string]interface{}{
				"Title":              messageTemplate.IssueTitle1,
				"Time":               req.Issue.CreatedAt.Format("2006年01月02日 15:04:05"),
				"IssueTitle":         req.Issue.Title,
				"Login":              req.Issue.User.Login,
				"IssueNumber":        strconv.FormatInt(req.Issue.Number, 10),
				"IssueContent":       req.Issue.Body,
				"RequestedReviewers": "无",
				"IssueURL":           req.Issue.HTMLURL,
			})

			_, err = lark.SendGroupMessage(receiveID, issueMsg)
			if err != nil {
				log.Errorf("send message error: %+v", err)
			}
			return
		}

		// 如果有，则先发送到人，等待超时后再发送到群
		issueMsg = t.ExecuteString(map[string]interface{}{
			"Title":              messageTemplate.IssueTitle1,
			"Time":               req.Issue.CreatedAt.Format("2006年01月02日 15:04:05"),
			"IssueTitle":         req.Issue.Title,
			"Login":              req.Issue.User.Login,
			"IssueNumber":        strconv.FormatInt(req.Issue.Number, 10),
			"IssueContent":       req.Issue.Body,
			"RequestedReviewers": assigneesToAt,
			"IssueURL":           req.Issue.HTMLURL,
		})

		lark.SendMessages(receiveID, assignees, issueMsg, strconv.FormatInt(req.Issue.Number, 10), string(gitHubEvent))
	case PullRequestEvent:
		var req PullRequestPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)

		t := fasttemplate.New(messageTemplate.PrMsg, "{{", "}}")

		receiveID, err := config.QueryReceiveIDByRepo(req.Repository.FullName)
		if err != nil {
			return
		}

		var assigneesToAt string
		assignees := make([]string, len(req.PullRequest.Assignees))
		for i, v := range req.PullRequest.Assignees {
			if _, ok := config.LarkMaps[v.Login]; ok {
				assigneesToAt += "<at id=" + config.LarkMaps[v.Login].Lark + "></at>"
				assignees[i] = v.Login
			}
		}

		var prMsg string
		// 如果没有 Assignees 直接发到群
		if len(assignees) == 0 {
			prMsg = t.ExecuteString(map[string]interface{}{
				"Title":              messageTemplate.PrTitle2,
				"Time":               req.PullRequest.CreatedAt.Format("2006年01月02日 15:04:05"),
				"PRTitle":            req.PullRequest.Title,
				"Login":              req.PullRequest.User.Login,
				"PRNumber":           strconv.FormatInt(req.PullRequest.Number, 10),
				"PRContent":          req.PullRequest.Body,
				"RequestedReviewers": "无",
				"PRURL":              req.PullRequest.HTMLURL,
			})

			_, err = lark.SendGroupMessage(receiveID, prMsg)
			if err != nil {
				log.Errorf("send message error: %+v", err)
			}
			return
		}

		// 如果有，则先发送到人，等待超时后再发送到群
		prMsg = t.ExecuteString(map[string]interface{}{
			"Title":              messageTemplate.PrTitle1,
			"Time":               req.PullRequest.CreatedAt.Format("2006年01月02日 15:04:05"),
			"PRTitle":            req.PullRequest.Title,
			"Login":              req.PullRequest.User.Login,
			"PRNumber":           strconv.FormatInt(req.PullRequest.Number, 10),
			"PRContent":          req.PullRequest.Body,
			"RequestedReviewers": assigneesToAt,
			"PRURL":              req.PullRequest.HTMLURL,
		})
		lark.SendMessages(receiveID, assignees, prMsg, strconv.FormatInt(req.PullRequest.Number, 10), string(gitHubEvent))
	case PullRequestReviewEvent:
		var req PullRequestReviewPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)

		t := fasttemplate.New(messageTemplate.PrReviewMsg, "{{", "}}")

		receiveID, err := config.QueryReceiveIDByRepo(req.Repository.FullName)
		if err != nil {
			return
		}

		// 正则拿到提及人
		var requestedReviewers string
		atPeople := getAtPeopleFromComment(req.Review.Body)
		for _, v := range atPeople {
			if _, ok := config.LarkMaps[v]; ok {
				requestedReviewers += "<at id=" + config.LarkMaps[v].Lark + "></at>"
			}
		}

		var prReviewMsg string
		// 如果没有 Assignees 直接发到群
		if len(atPeople) == 0 {
			prReviewMsg = t.ExecuteString(map[string]interface{}{
				"Title":              messageTemplate.PrReviewTitle1,
				"Time":               req.PullRequest.CreatedAt.Format("2006年01月02日 15:04:05"),
				"PRTitle":            req.PullRequest.Title,
				"Login":              req.PullRequest.User.Login,
				"PRNumber":           strconv.FormatInt(int64(req.PullRequest.Number), 10),
				"PRContent":          req.PullRequest.Body,
				"RequestedReviewers": "无",
				"PRURL":              req.PullRequest.HTMLURL,
			})

			_, err = lark.SendGroupMessage(receiveID, prReviewMsg)
			if err != nil {
				log.Errorf("send message error: %+v", err)
			}
			return
		}

		prReviewMsg = t.ExecuteString(map[string]interface{}{
			"Title":              messageTemplate.PrReviewTitle1,
			"Time":               req.PullRequest.CreatedAt.Format("2006年01月02日 15:04:05"),
			"PRTitle":            req.PullRequest.Title,
			"Login":              req.PullRequest.User.Login,
			"PRNumber":           strconv.FormatInt(int64(req.PullRequest.Number), 10),
			"PRContent":          req.Review.Body,
			"RequestedReviewers": requestedReviewers,
			"PRURL":              req.PullRequest.HTMLURL,
		})

		number := strconv.FormatInt(int64(req.PullRequest.Number), 10)

		// 简化版的 SendMessages
		for _, receiver := range atPeople {
			if _, ok := config.MsgQueue[receiver+number+string(gitHubEvent)]; !ok {
				config.MsgQueue[receiver+number+string(gitHubEvent)] = constants.Unread1
				msg, err := lark.SendMessage(config.LarkMaps[receiver].Lark, prReviewMsg)
				if err != nil {
					log.Errorf("send message failed, msg=%v", msg)
					log.Errorf("send message failed, err=%v", err)
				}
				log.Infof("msgQ=%v", config.MsgQueue)
			} else {
				if config.MsgQueue[receiver+number+string(gitHubEvent)] == constants.READ {
					// 已读，从队列中删除
					delete(config.MsgQueue, receiver+number+string(gitHubEvent))
					log.Infof("msgQ=%v", config.MsgQueue)
				}
			}
		}

		// 定时未读则转发群聊
		time.AfterFunc(config.CommentUnread, func() {
			prReviewMsg = lark.MsgTemplateUpgrade(prReviewMsg)
			for _, receiver := range atPeople {
				if _, ok := config.MsgQueue[receiver+number+string(gitHubEvent)]; ok {
					if config.MsgQueue[receiver+number+string(gitHubEvent)] == constants.Unread1 {
						_, err = lark.SendGroupMessage(receiveID, prReviewMsg)
						if err != nil {
							log.Errorf("send message error: %+v", err)
						}
					}
				}
			}
		})
	case PullRequestReviewCommentEvent:
		var req PullRequestReviewCommentPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)

		t := fasttemplate.New(messageTemplate.PrCommentMsg, "{{", "}}")

		receiveID, err := config.QueryReceiveIDByRepo(req.Repository.FullName)
		if err != nil {
			return
		}

		// 正则拿到提及人
		var requestedReviewers string
		atPeople := getAtPeopleFromComment(req.Comment.Body)
		for _, v := range atPeople {
			if _, ok := config.LarkMaps[v]; ok {
				requestedReviewers += "<at id=" + config.LarkMaps[v].Lark + "></at>"
			}
		}

		var prCommentMsg string
		// 如果没有 Assignees 直接发到群
		if len(atPeople) == 0 {
			prCommentMsg = t.ExecuteString(map[string]interface{}{
				"Title":              messageTemplate.PrCommentTitle1,
				"Time":               req.PullRequest.CreatedAt.Format("2006年01月02日 15:04:05"),
				"PRTitle":            req.PullRequest.Title,
				"Login":              req.PullRequest.User.Login,
				"PRNumber":           strconv.FormatInt(req.PullRequest.Number, 10),
				"PRContent":          req.PullRequest.Body,
				"RequestedReviewers": "无",
				"PRURL":              req.PullRequest.HTMLURL,
			})

			_, err = lark.SendGroupMessage(receiveID, prCommentMsg)
			if err != nil {
				log.Errorf("send message error: %+v", err)
			}
			return
		}

		prCommentMsg = t.ExecuteString(map[string]interface{}{
			"Title":              messageTemplate.PrCommentTitle1,
			"Time":               req.PullRequest.CreatedAt.Format("2006年01月02日 15:04:05"),
			"PRTitle":            req.PullRequest.Title,
			"Login":              req.PullRequest.User.Login,
			"PRNumber":           strconv.FormatInt(req.PullRequest.Number, 10),
			"PRContent":          req.Comment.Body,
			"RequestedReviewers": requestedReviewers,
			"PRURL":              req.PullRequest.HTMLURL,
		})

		number := strconv.FormatInt(req.PullRequest.Number, 10)

		// 简化版的 SendMessages
		for _, receiver := range atPeople {
			if _, ok := config.MsgQueue[receiver+number+string(gitHubEvent)]; !ok {
				config.MsgQueue[receiver+number+string(gitHubEvent)] = constants.Unread1
				msg, err := lark.SendMessage(config.LarkMaps[receiver].Lark, prCommentMsg)
				if err != nil {
					log.Errorf("send message failed, msg=%v", msg)
					log.Errorf("send message failed, err=%v", err)
				}
				log.Infof("msgQ=%v", config.MsgQueue)
			} else {
				if config.MsgQueue[receiver+number+string(gitHubEvent)] == constants.READ {
					// 已读，从队列中删除
					delete(config.MsgQueue, receiver+number+string(gitHubEvent))
					log.Infof("msgQ=%v", config.MsgQueue)
				}
			}
		}

		// 定时未读则转发群聊
		time.AfterFunc(config.CommentUnread, func() {
			prCommentMsg = lark.MsgTemplateUpgrade(prCommentMsg)
			for _, receiver := range atPeople {
				if _, ok := config.MsgQueue[receiver+number+string(gitHubEvent)]; ok {
					if config.MsgQueue[receiver+number+string(gitHubEvent)] == constants.Unread1 {
						_, err = lark.SendGroupMessage(receiveID, prCommentMsg)
						if err != nil {
							log.Errorf("send message error: %+v", err)
						}
					}
				}
			}
		})
	}
	c.JSON(200, nil)
}
