# Message Pusher

对各种消息推送方式的封装，提供统一的接口和定时发送

- [x] 企业微信自建应用（支持`text`和`textcard`类型
- [ ] 自编写的消息推送
- [ ] 邮件
- [ ] 短信

# 1. Args

使用`-h`查看内置参数

- `-c` 指定配置文件路径
- `-i` 创建并初始化数据库

# 2. API

## 2.1 `/send`

### 2.1.1 使用`POST`方式发送请求

- `caller`: 调用者名称，用于区分调用者，可随意填写。
- `target`: 接收者的ID，多个目标时使用`,`分隔。如`1,2,3`
- `secret`: 密钥
- `data`: 信息内容，为json字符串，格式见下方

#### 2.1.1.1 `data`

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

# 3. Mod API

## 3.1 wecom

### 3.1.1 额外信息

额外信息需要填充于`data`的`extra`中。默认值定义在`mod/wecom/def.go`文件中。

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

# 4. Develop

## 4.1 根服务

用于提供基础的服务

### 4.1.1 模块可调用的方法

```go
SessionVerify(r *http.Request) bool // 检查当前是否登录
Login(w http.ResponseWriter, r *http.Request) // 提供登录界面
InsertMod(mod mod.Model) bool // 启用传入的模块，返回是否成功
GetIP(r *http.Request) string // 返回访问者的ip地址
LastPage(w http.ResponseWriter, r *http.Request) // 重定向到上一级页面
```

## 4.2 启用模块

### 4.2.1 创建数据库表

编写一个用于创建数据库的函数，并在`main.go`文件中`initDatabase()`中调用。可参照以下模板

```go
func InsertDatabaseTable() bool {
	d := db.Connect()
	if d == nil {
		glog.Fatal("can not connect database")
		return false
	}
	err := d.AutoMigrate(&ModModel{})
	if err != nil {
		glog.Fatal("failed to create table [%s]", err.Error())
		return false
	}
	glog.Info("database table insert finished [%s]", ModModel{}.TableName())
	return true
}
```

`ModModel`应继承`gorm.Tabler`接口，如以下模板

```go
type ModModel struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Secret         string `gorm:"column:secret"`
	Name           string `gorm:"column:name"`
	ValidityPeriod int64  `gorm:"column:validity_period"`
	CreateTime     int64  `gorm:"column:create_time;autoCreateTime"`
	LastUsed       int64  `gorm:"column:last_used"`
	Expired        int64  `gorm:"column:expired"`
}

func (ModModel) TableName() string {
	return tableName
}
```

### 4.2.2 启用模块

编写一个用于启用模块的函数，并在`main.go`文件中`initApp()`中调用。可参照以下模板

```go
func Load() bool {
	self = new(Model)
	self.handler = make(map[string]func(w http.ResponseWriter, r *http.Request))
	self.handler["/"] = index
	return app.InsertMod(self)
}
```

`Model`应继承`Model`接口，接口详细信息请看对应的接口介绍。

## 4.3 新模块需要继承的相关接口

### 4.3.1 Package

当信息需要发送时，会调用Send()函数尝试发送，当发送失败时会间隔一定时间后再次发送。

```go
type Package interface {
// Send will try to send information and return whether it is successful
Send() bool
}
```

### 4.3.2 Target

由于每个模块的数据库表中储存的信息不相同，因此需要通过继承此接口来提供统一的获取方法

```go
type Target interface {
	// GetName return secret name
	GetName() string
	// GetKey return mod key
	GetKey() string
	// GetSecret return secret
	GetSecret() string
	GetValidityPeriod() int64
	GetCreateTime() int64
	GetLastUsed() int64
	GetExpired() int64
}
```

### 4.3.3 Model

将自身的Model通过`app.InsertMod(mod mod.Model) bool `插入根服务

```go
type Model interface {
	// URL is url path prefix
	URL() string
	// Name is mod name
	Name() string
	// Table is the table name used in the database
	Table() string
	// Key is the key for extra
	//     The message sender can use the extra in data to carry the information required by the specified mod.
	Key() string
	// Handle is HTTP handler, handle http requests
	Handle(http.ResponseWriter, *http.Request)
	// Prepare is used to process the general information passed in and return the processed data
	Prepare(id int64, data *MessageModel) Package
	// GetTarget returns the information of the key corresponding to the id
	GetTarget(id int64) Target
}
```

## 4.4 新模块注意事项

### 4.4.1 新密钥插入

当模块所属表中插入新的记录时，需要手动调用`db.InsertTarget`函数，以建立对应target的唯一全局id。否则无法通过全局发送服务为此密钥提供对应的发送服务，也无法在`/target`中查询到对应的项。

```go
type TargetFunc func(d *gorm.DB) int64
func InsertTarget(fun TargetFunc, table string) bool
```

在`InsertTarget`对数据库的操作以事务方式进行，包括调用`TargetFunc`时传入的`d`。在发生数据库操作失败时，数据库会回滚到操作之前的状态。

在`TargetFunc`中，需要借助传入的`d`对数据库表进行操作，返回值为插入的项在表中的id，插入失败时返回负数。模板如下

```go
func insertModName(data *ModModel) bool {
    d := db.Connect()
    if d == nil {
        return false
    }
    dbLock.Lock()
    defer dbLock.Unlock()

    fun := func(d *gorm.DB) int64 {
        res := d.Model(&ModModel{}).Create(data)
        if res.Error != nil || res.RowsAffected != 1 {
            glog.Warning("failed to insert mod [%s]", res.Error)
            return -1
        }
        return data.ID
    }

    if db.InsertTarget(fun, data.TableName()) {
        return true
    }

    return false
}
```

