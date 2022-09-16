package messageTemplate

const PrReviewTitle1 = "有人在 Review 中 @ 你!"
const PrReviewTitle2 = "有 Review 待处理!"
const PrReviewMsg = `{
  "config": {
    "wide_screen_mode": true
  },
  "header": {
    "template": "blue",
    "title": {
      "content": "⏰ GitHub PullRequest Review 提醒",
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
            "type": "pull_request_review",
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
