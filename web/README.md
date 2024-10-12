## CowIM Web 设计规范
<h4 id="constants">1. 常量定义 </h4>

```go
// 群聊或者单聊
const (
    SINGLE_CHAT = iota  // 单聊
    GROUP_CHAT          // 群聊
    BIG_GROUP_CHAT      // 大规模群聊
    SYSTEM_INFO = 99    // 系统消息
)

// Message 类型
const (
    MSG_COMMON_MSG = iota   // 一般文本消息
    MSG_REPLY_MSG           // 回复消息
    MSG_IMAGE_MSG           // 图片消息
    MSG_VIDEO_MSG           // 视频消息
    MSG_FILE_MSG            // 上传文件
    MSG_SYSTEM_MSG = 99     // 系统通知（和SYSTEM_INFO一起用）
)

// 群聊权限
const (
    GROUP_MASTER = iota     // 群主
    GROUP_ADMIN             // 管理
    GROUP_COMMON            // 普通群员
)
```
<br>

<h4 id="types">2. 消息传递、数据入库格式</h4>

1. 前端和后端通过Websocket发送消息、IndexedDB入库的通用格式

```protobuf
message Message {
  uint32 from = 1; // 发送者
  uint32 to = 2; // 接受者
  optional uint32 group = 3; // 群组ID
  string content = 4; // 消息内容
  uint32 type = 5; // 群组还是单聊，群组是500人以下群还是500人以上群？这个字段其实是冗余的，只是方便后端判断
  uint32 msg_type = 6; // 聊天类型，直接聊天、回复、文件信息？
  optional int64 extend = 7; // 消息拓展内容（回复消息的原消息id）
  int64 timestamp = 8; // 时间戳
}
```

2. 用户刚登陆时，从后端获取离线时收到的消息，收到的消息格式

```go
// UserTimeline 用户时间线
type UserTimeline struct {
    ID         int64       `bson:"_id"                json:"id"`
    ReceiverID uint32      `bson:"receiver_id"        json:"receiverID"`
    SenderID   uint32      `bson:"sender_id"          json:"senderID"`
    GroupID    uint32      `bson:"group_id,omitempty" json:"groupID,omitempty"`
    Message    MessageSync `bson:"msgForward"         json:"msgForward"`
    Timestamp  time.Time   `bson:"timestamp"          json:"timestamp"` // 用于删除过时消息 + 实现Timeline模型(用户消息按时间线排列)
}
// MessageSync 消息同步表，用于用户即时能查询的信息，不直接入库，由Timeline间接入库（做数据冗余）
type MessageSync struct {
    ID        int64     `bson:"_id"              json:"id"`
    MsgType   uint8     `bson:"msg_type"         json:"msgType"`
    Content   string    `bson:"content"          json:"content"`
    Extend    int64     `bson:"extend,omitempty" json:"extend,omitempty"`
    Timestamp time.Time `bson:"timestamp"        json:"timestamp"` // 用于删除过时消息
}

/*
调用接口获取的消息大概长这个样子
    {
        "id": 600437751853127,
        "receiverID": 2,
        "senderID": 3,
        "msgForward": {
            "id": 600437751853126,
            "msgType": 0,
            "content": "123",
            "timestamp": "2024-10-12T10:07:30.011Z"
        },
        "timestamp": "2024-10-12T10:07:30.011Z"
    },
*/
```

3. 存在Web端sessionStorage中的数据

```
jwtToken: xxx   // jwt令牌，登陆后获得
CowID:    xxx   // 当前页面登陆的用户CowID
```

<br>

<h4 id="logic">3. 逻辑流程</h4>

- 登陆逻辑
: 1. 用户访问网页端，显示**login.html**
: 2. 用户输入CowID和密码，点击登陆按钮
: 3. 携带CowID和密码，调用后端/login接口，判断登陆合法性
: 4. 如果正确，则携带/login接口分配的jwtToken存储进sessionStorage中，然后跳转到**main.html**
: 5. 如果错误，前端处理，弹框提示

- 注册逻辑
: 1. 用户访问网页端，显示**login.html**
: 2. 用户点击注册按钮，然后前端跳转到**register.html**
: 3. 用户输入需要的信息，点击注册按钮
: 4. 前端判断输入信息合法性
: 5. 携带注册信息，调用后端/register接口
: 6. 如果注册成功，接口返回用户分配好的CowID，前端弹窗返回CowID，用户确认后返回**login.html**
: 7. 如果注册失败，前端处理，弹框提示，停留在**register.html**页面

- 登陆完成后的逻辑
: 1. 用户完成**登陆逻辑**后，显示**main.html**
: 2. 调用/wsget接口，获取可用Websocket server的IP地址
: 3. 连接分配到的Websocket server，如果失败则前端报错
: 4. 首先初始化IndexedDB，创建CowIMDB数据库，以及规定的所有数据表
: 5. 获取当前账号的所有信息，包括同步消息(/timelinesync)、个人资料(/info)、加入过的群聊(/groups)、加入过的好友(/friends)，并存在IndexedDB中
: 6. 获取完毕后，显示页面，否则前端报错
: 7. 默认显示在最近会话列表界面，聊天框不显示

- 点击/显示最近会话列表的逻辑
: 1. 用户完成**登陆完成后的逻辑**后，从IndexedDB中读取messages表，根据messages表来显示最近会话的用户
: 2. 遍历messages表，检查type为SINGLE_CHAT([常量定义](#constant)中的群聊或单聊)，然后将所有from和to整合起来，去重，再排除掉自己的ID，通过CowID显示最近对话列表
: 3. 遍历messages表，检查type为GROUP_CHAT([常量定义](#constant)中的群聊或单聊)，然后将所有group整合起来，去重，通过GroupID显示最近对话列表
: 4. 通过timestamp进行会话列表的排序

- 点击/显示好友按钮后的逻辑
: 1. 从IndexedDB中读取friends表，根据friends表来显示好友列表
: 2. 将所有的friendID整合起来，去重
: 3. 通过好友的用户名进行排序

- 点击/显示群组按钮后的逻辑
: 1. 从IndexedDB中读取groups表，根据groups表来显示群列表
: 2. 将所有的groupsID整合起来
: 3. 通过群聊的群名进行排序

<br>

<h4 id="db">4. IndexedDB设计</h4>

- 消息表 (同[消息传递、入库格式](#types)中 1 的格式)

```
messages {
    uint32 from;
    uint32 to;
    uint32 group;
    string content;
    uint32 type;
    uint32 msg_type;
    int64  extend;
    time   timestamp;
}
```

- 好友表
```
friends {
    uint32 friendID; // 好友的CowID
    string friendName; // 好友的用户名
    string friendNote; // 给好友设置的备注
    user friendInfo; // 好友的具体用户信息
}
```

- 用户信息表 (不单独存储，作为嵌套存储)
```
user {
    uint32 id; // CowID
    string name; // 用户名
    string avatar; // 头像
    ... // 字段待定，例如个性签名
}
```

- 群组表
```
groups {
    uint32 groupID; // 群组的ID
    string groupName; // 群组的名字
    []user groupMembers; // 群组的成员们
}
```

<br>

<h4 id="details">5. 细节<h4/>

> 什么时候去更新群聊、好友等消息？如果有改动，Web端如何收到？

客户端在每次登陆的时候，去请求后端，进行账号相关信息的同步。如果信息有改动，基于本即时通讯系统的同时在能一个客户端上在线的限制，改动只会在一个客户端上进行，而不需要考虑多客户端在线的同步问题，所以有改动的话，直接在本地进行IndexedDB的更新就可以了。
