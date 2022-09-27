package constants

import (
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/config"
	"time"
)

type MsgType int

const (
	READ    MsgType = 0
	Unread1 MsgType = 1
	Unread2 MsgType = 2
	Unread3 MsgType = 3
)

var (
	TimeUnread1   = time.Duration(config.Cfg.Scheduler.TimeUnread1) * time.Minute   // PR/Issue 消息第一次发送消息后若未读，经过 TimeUnread1 后重发
	TimeUnread2   = time.Duration(config.Cfg.Scheduler.TimeUnread2) * time.Minute   // PR/Issue 消息第二次发送消息后若未读，经过 TimeUnread2 后发送给上级
	TimeUnread3   = time.Duration(config.Cfg.Scheduler.TimeUnread3) * time.Minute   // PR/Issue 消息第三次发送消息后若未读，经过 TimeUnread3 后抄送群聊
	CommentUnread = time.Duration(config.Cfg.Scheduler.CommentUnread) * time.Minute // Comment 消息第一次发送后若未读，经过 CommentUnread 后抄送群聊
)
