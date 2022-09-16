package messageTemplate

const SendIdMsg = `{
  "config": {
    "wide_screen_mode": true
  },
  "header": {
    "template": "blue",
    "title": {
      "content": " 飞书 ID",
      "tag": "plain_text"
    }
  },
  "i18n_elements": {
    "zh_cn": [
      {
        "tag": "div",
        "text": {
          "content": "号码为：{{ChatID}}",
          "tag": "lark_md"
        }
      }
    ]
  }
}`
