package eventListener

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/url"
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
	fmt.Printf("%s\n", bodyStr)
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
		// Do whatever you want from here...
		fmt.Printf("%+v\n", req)
		fmt.Printf("%s\n", req.PullRequest.Head.Repo.FullName)
		fmt.Printf("%s\n", req.PullRequest.Head.User.Login)
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

func Init() error {
	router := gin.Default()
	router.POST("/github/webhook", GitHubWebHook)
	err := router.Run(":3000")
	if err != nil {
		return err
	}
	return nil
}
