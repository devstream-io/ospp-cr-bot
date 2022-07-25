package lark

const (
	SEND_GROUP_ID = `{"text":"本群的群号为：\n{{ChatID}}"}`
	SEND_PR_MSG   = `{
    "config": {
      "wide_screen_mode": true
    },
    "header": {
      "template": "red",
      "title": {
        "content": "⏰ GitHub PullRequest 提醒",
        "tag": "plain_text"
      }
    },
    "elements": [
      {
        "fields": [
          {
            "is_short": true,
            "text": {
              "content": "**🕐 时间：**\n{{Time}}",
              "tag": "lark_md"
            }
          },
          {
            "is_short": true,
            "text": {
              "content": "**📋 PR 标题：**\n{{PRTitle}}",
              "tag": "lark_md"
            }
          },
          {
            "is_short": false,
            "text": {
              "content": "",
              "tag": "lark_md"
            }
          },
          {
            "is_short": true,
            "text": {
              "content": "**👤 提出人：**\n{{Login}}",
              "tag": "lark_md"
            }
          },
          {
            "is_short": true,
            "text": {
              "content": "**🔢 PR 号：**\n{{PRNumber}}",
              "tag": "lark_md"
            }
          },
          {
            "is_short": false,
            "text": {
              "content": "",
              "tag": "lark_md"
            }
          }
        ],
        "tag": "div"
      },
      {
        "tag": "hr"
      },
      {
        "tag": "action",
        "actions": [
          {
            "tag": "button",
            "text": {
              "tag": "plain_text",
              "content": "跟进 PR"
            },
            "url": "{{PRURL}}",
            "type": "primary"
          }
        ]
      }
    ]
  }`
)
