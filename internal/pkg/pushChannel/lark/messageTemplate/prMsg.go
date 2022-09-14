package messageTemplate

const PrTitle1 = "您有一个 PullRequest 待处理!"
const PrTitle2 = "%s 未处理 PullRequest!"
const PrMsg = `{
  "config": {
    "wide_screen_mode": true
  },
  "header": {
    "template": "blue",
    "title": {
      "content": "⏰ {{Title}}",
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
        },
        {
          "is_short": true,
          "text": {
            "content": "**🔢 提醒人：**\n{{RequestedReviewers}}",
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
            "content": "**📜 PR 内容：**\n{{PRContent}}",
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
            "content": "确认收到"
          },
          "type": "primary",
          "value": {
            "type": "pull_request",
            "numbers": "{{PRNumber}}"
          }
        },
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
