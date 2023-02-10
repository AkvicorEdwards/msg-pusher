# Args

使用`-h`查看内置参数

- `-c` 指定配置文件路径
- `-i` 创建并初始化数据库

# API

## `/send`

### 使用`POST`方式发送请求

- `caller`: 调用者名称，用于区分调用者，可随意填写。
- `target`: 接收者的ID，多个目标时使用`,`分隔。如`1,2,3`
- `secret`: 密钥
- `data`: 信息内容，为json字符串，格式见下方

#### `data`

```text
{
   "title":"t_title",  // 信息标题
   "content":"t_content",  // 信息内容
   "sender":"t_sender",  // 发送者
   "urgency":1,  // 信息紧急程度（0：指令，1：紧急，2：重要，3：普通
   "time_create":2,  // 信息创建时间（秒为单位的时间戳
   "time_send":3,  // 信息期望的发送时间（秒为单位的时间戳
   "extra":{  // 各模块需要的额外信息
      "t_mod":{  // 模块的key
         "t_key":"t_value"  // 模块需要的信息
      }
   }
}
```

# Mod API

## wecom

### 额外信息

额外信息需要填充于`data`的`extra`中。

```text
"extra":{
   "wecom":{
      "type":"textcard", // 信息类型，默认为textcard。textcard,text
      "touser":"Akvicor", // 接收者，以|分割（以企业微信官方要求为准），默认为Akvicor
      "toparty":"", // 接收者party，以|分割（以企业微信官方要求为准，默认为空
      "totag":"", // 接收者tag，以|分割（以企业微信官方要求为准，默认为空
      "url":"url", // 消息点击后的url（以企业微信官方要求为准，默认为url
      "btn":"BUTTON" // 按钮文字（以企业微信官方要求为准，默认为BUTTON
   }
}
```

