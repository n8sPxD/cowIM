# cowIM

个人绞尽脑汁设计的即时通讯系统

### 架构图
<img src="docs/cowIm架构v2.png" height="719" alt="架构图v2">

### 目前实现与进度
#### 底层及架构
- [x] 自定义应用层协议
<br><br/>
- [x] HTTP 网关 反向代理
- [x] HTTP 网关 鉴权
- [ ] HTTP 网关 用户黑名单
- [x] HTTP 网关 CORS
<br><br/>
- [x] TCP 网关 反向代理
- [ ] TCP 网关 鉴权
- [ ] TCP 网关 客户端连接信息保存

#### 业务
- [x] 注册
- [ ] 登陆
- [ ] 好友列表获取
- [ ] 消息记录
- [ ] 群聊列表获取
- [ ] 最近会话列表
- [ ] 单聊
- [ ] 群聊