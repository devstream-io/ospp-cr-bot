package constants

import "time"

type MsgType int

const (
	READ    MsgType = 0
	Unread1 MsgType = 1
	Unread2 MsgType = 2
	Unread3 MsgType = 3
)

//const (
//	TimeUnread1   = 30 * time.Minute
//	TimeUnread2   = 30 * time.Minute
//	TimeUnread3   = 11 * time.Hour
//	CommentUnread = 1 * time.Hour
//)

const (
	TimeUnread1   = 3 * time.Second
	TimeUnread2   = 3 * time.Second
	TimeUnread3   = 10 * time.Second
	CommentUnread = 20 * time.Second
)
