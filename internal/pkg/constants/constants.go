package constants

import "time"

type MsgType int

const (
	READ    MsgType = 0
	Unread1 MsgType = 1
	Unread2 MsgType = 2
	Unread3 MsgType = 3
)

const (
	TimeUnread1   = 30 * time.Minute // PR/Issue 消息第一次发送消息后若未读，经过 TimeUnread1 后重发
	TimeUnread2   = 30 * time.Minute // PR/Issue 消息第二次发送消息后若未读，经过 TimeUnread2 后发送给上级
	TimeUnread3   = 11 * time.Hour   // PR/Issue 消息第三次发送消息后若未读，经过 TimeUnread3 后抄送群聊
	CommentUnread = 1 * time.Hour    // Comment 消息第一次发送后若未读，经过 CommentUnread 后抄送群聊
)
