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
