package eventListener

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/config"
	"github.com/valyala/fasttemplate"
	"io/ioutil"
	"net/url"
	"strconv"

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

	switch gitHubEvent {
	case IssueCommentEvent:
		var req IssueCommentPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)
		// Do whatever you want from here...
		fmt.Printf("%+v\n", req)
		fmt.Printf("%s\n", req.Repository.FullName)
	case IssuesEvent:
		var req IssuesPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)
		// Do whatever you want from here...
		fmt.Printf("%+v\n", req)
		fmt.Printf("%s\n", req.Repository.FullName)
	case PullRequestEvent:
		var req PullRequestPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)
		t := fasttemplate.New(lark.SEND_PR_MSG, "{{", "}}")
		prMsg := t.ExecuteString(map[string]interface{}{
			"Time":     req.PullRequest.CreatedAt.Format("2006年01月02日 15:04:05"),
			"PRTitle":  req.PullRequest.Title,
			"Login":    req.PullRequest.User.Login,
			"PRNumber": strconv.FormatInt(req.PullRequest.Number, 10),
			"PRURL":    req.PullRequest.HTMLURL,
		})
		receiveID, err := config.QueryReceiveIDByRepo(req.Repository.FullName)
		if err != nil {
			return
		}
		lark.SendMessage(receiveID, prMsg)
	case PullRequestReviewEvent:
		var req PullRequestReviewPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)
		// Do whatever you want from here...
		fmt.Printf("%+v\n", req)
		fmt.Printf("%s\n", req.PullRequest.Head.Repo.FullName)
		fmt.Printf("%s\n", req.PullRequest.Head.User.Login)
	case PullRequestReviewCommentEvent:
		var req PullRequestReviewCommentPayload
		_ = json.Unmarshal([]byte(bodyStr), &req)
		// Do whatever you want from here...
		fmt.Printf("%+v\n", req)
		fmt.Printf("%s\n", req.PullRequest.Head.Repo.FullName)
		fmt.Printf("%s\n", req.PullRequest.Head.User.Login)
	}
	c.JSON(200, nil)
}
