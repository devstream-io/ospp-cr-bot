package messageTemplate

const IssueTitle1 = "您有一个 Issue 待处理!"
const IssueTitle2 = "%s 未处理 Issue!"
const IssueMsg = `{
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
            "content": "**📋 Issue 标题：**\n{{IssueTitle}}",
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
            "content": "**🔢 Issue 号：**\n{{IssueNumber}}",
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
            "content": "**📜 Issue 内容：**\n{{IssueContent}}",
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
            "type": "issues",
            "numbers": "{{IssueNumber}}"
          }
        },
        {
          "tag": "button",
          "text": {
            "tag": "plain_text",
            "content": "跟进 Issue"
          },
          "url": "{{IssueURL}}",
          "type": "primary"
        }
      ]
    }
  ]
}`
