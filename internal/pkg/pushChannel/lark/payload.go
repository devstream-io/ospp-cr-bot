package lark

type CallBackReq struct {
	Encrypt string `json:"encrypt"`
}

type DecodeCallBackReq struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type CallBackResp struct {
	Challenge string `json:"challenge"`
}

type MessageCallBackResp struct {
	Schema string `json:"schema"`
	Header struct {
		EventID    string `json:"event_id"`
		EventType  string `json:"event_type"`
		CreateTime string `json:"create_time"`
		Token      string `json:"token"`
		AppID      string `json:"app_id"`
		TenantKey  string `json:"tenant_key"`
	} `json:"header"`
	Event struct {
		Sender struct {
			SenderID struct {
				UnionID string `json:"union_id"`
				UserID  string `json:"user_id"`
				OpenID  string `json:"open_id"`
			} `json:"sender_id"`
			SenderType string `json:"sender_type"`
			TenantKey  string `json:"tenant_key"`
		} `json:"sender"`
		Message struct {
			MessageID   string `json:"message_id"`
			RootID      string `json:"root_id"`
			ParentID    string `json:"parent_id"`
			CreateTime  string `json:"create_time"`
			ChatID      string `json:"chat_id"`
			ChatType    string `json:"chat_type"`
			MessageType string `json:"message_type"`
			Content     string `json:"content"`
			Mentions    []struct {
				Key string `json:"key"`
				ID  struct {
					UnionID string `json:"union_id"`
					UserID  string `json:"user_id"`
					OpenID  string `json:"open_id"`
				} `json:"id"`
				Name      string `json:"name"`
				TenantKey string `json:"tenant_key"`
			} `json:"mentions"`
		} `json:"message"`
	} `json:"event"`
}

type SendMessageReq struct {
	ReceiveID string `json:"receive_id"`
	Content   string `json:"content"`
	MsgType   string `json:"msg_type"`
}

type CardPostReq struct {
	OpenId        string `json:"open_id"`
	UserId        string `json:"user_id"`
	OpenMessageId string `json:"open_message_id"`
	TenantKey     string `json:"tenant_key"`
	Token         string `json:"token"`
	Action        struct {
		Value struct {
			Type    string `json:"type"`
			Numbers string `json:"numbers"`
		} `json:"value"`
		Tag string `json:"tag"`
	} `json:"action"`
}

/*
	Message API Model
*/

type Header struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
}

type CreateMessageRequest struct {
	ReceiveID string `json:"receive_id"`
	Content   string `json:"content"`
	MsgType   string `json:"msg_type"`
}

type CreateMessageResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *MessageItem `json:"data"`
}

type MessageItem struct {
	MessageID  string         `json:"message_id,omitempty"`
	RootID     string         `json:"root_id,omitempty"`
	ParentID   string         `json:"parent_id,omitempty"`
	MsgType    string         `json:"msg_type,omitempty"`
	CreateTime string         `json:"create_time,omitempty"`
	UpdateTime string         `json:"update_time,omitempty"`
	Deleted    bool           `json:"deleted,omitempty"`
	ChatID     string         `json:"chat_id,omitempty"`
	Sender     *MessageSender `json:"sender,omitempty"`
	Body       *MessageBody   `json:"body,omitempty"`
}

type MessageBody struct {
	Content string `json:"content,omitempty"`
}

type MessageSender struct {
	ID         string `json:"id,omitempty"`
	IDType     string `json:"id_type,omitempty"`
	SenderType string `json:"sender_type"`
	TenantKey  string `json:"tenant_key"`
}

type ReceiveEventEncrypt struct {
	Encrypt string `json:"encrypt" form:"encrypt"`
}

type DecryptToken struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type Event struct {
	Schema string      `json:"schema"`
	Header Header      `json:"header"`
	Event  interface{} `json:"event"`
}

type ReceiveMessageEvent struct {
	Schema string       `json:"schema"`
	Header Header       `json:"header"`
	Event  MessageEvent `json:"event"`
}

type MessageEvent struct {
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
}

type Sender struct {
	SenderID   map[string]string `json:"sender_id"`
	SenderType string            `json:"sender_type"`
	TenantKey  string            `json:"tenant_key"`
}

type Message struct {
	MessageID   string     `json:"message_id"`
	RootID      string     `json:"root_id"`
	ParentID    string     `json:"parent_id"`
	CreateTime  string     `json:"create_time"`
	ChatID      string     `json:"chat_id"`
	ChatType    string     `json:"chat_type"`
	MessageType string     `json:"message_type"`
	Content     string     `json:"content"`
	Mentions    []*Mention `json:"mentions,omitempty"`
}

type Mention struct {
	Key       string  `json:"key,omitempty"`
	ID        *UserID `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	TenantKey string  `json:"tenant_key,omitempty"`
}

type MentionV1 struct {
	Key       string `json:"key,omitempty"`
	ID        string `json:"id,omitempty"`
	IDType    string `json:"id_type,omitempty"`
	Name      string `json:"name,omitempty"`
	TenantKey string `json:"tenant_key,omitempty"`
}

type UserID struct {
	UserID  string `json:"user_id,omitempty"`
	OpenID  string `json:"open_id,omitempty"`
	UnionID string `json:"union_id,omitempty"`
}

type GetMessageHistoryResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *GetMessageHistoryBody `json:"data"`
}

type GetMessageHistoryBody struct {
	HasMore   bool                            `json:"has_more"`
	PageToken string                          `json:"page_token,omitempty"`
	Items     []*GetMessageHistoryMessageItem `json:"items,omitempty"`
}

type GetMessageHistoryMessageItem struct {
	MessageID      string         `json:"message_id,omitempty"`
	RootID         string         `json:"root_id,omitempty"`
	ParentID       string         `json:"parent_id,omitempty"`
	MsgType        string         `json:"msg_type,omitempty"`
	CreateTime     string         `json:"create_time,omitempty"`
	UpdateTime     string         `json:"update_time,omitempty"`
	Deleted        bool           `json:"deleted,omitempty"`
	ChatID         string         `json:"chat_id,omitempty"`
	Sender         *MessageSender `json:"sender,omitempty"`
	Body           *MessageBody   `json:"body,omitempty"`
	Mentions       []*MentionV1   `json:"mentions,omitempty"`
	UpperMessageID string         `json:"upper_message_id,omitempty"`
}

type UploadImageResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    *UploadImageResponseBody `json:"data,omitempty"`
}

type UploadImageResponseBody struct {
	ImageKey string `json:"image_key"`
}
