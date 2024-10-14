# cowIM
个人绞尽脑汁设计的即时通讯系统，目前处于开发阶段，勉强能跑，能够在前端进行一对一聊天

### 架构图
<img src="docs/pics/cowIm架构v7.png" alt="架构图">

### 如何运行
当前未开发完毕，构建处于测试阶段，比较复杂，后续采取服务容器化后会改进
##### 基础服务
```sh
% docker-compose -f deployments/base/docker-compose.yaml up -d
```
##### 即时通讯系统
```sh
% cd build/run

# 确保运行之前给予执行权限
# chmod +x *.sh

% ./auth_api.sh
% ./auth_rpc.sh
...
```

### 目前实现与进度
#### 业务(后端)
 - [x] 单聊
 - [x] 群聊
 - [x] 好友添加
 - [x] 好友查询
 - [ ] 好友管理
 - [x] 群聊创建
 - [x] 群聊查询

### 目录说明
```
.
├── build           // 构建脚本
├── cmd             // 初始化文件
├── deployments     // 服务容器部署
├── docs            // 文档
├── internal        // 系统内部代码
├── pkg             // 系统公共代码
├── test            // 测试代码
└── web             // 前端页面

```



### 参考与启发
- https://www.bilibili.com/video/BV1KM411S7WT (架构)
- http://www.52im.net/thread-4257-1-1.html (架构)
- https://space.bilibili.com/30625295/channel/collectiondetail?sid=3179321 (架构)
- https://www.bilibili.com/video/BV1rU4y17769 (消息收发)
- https://github.com/zeromicro/zero-examples/tree/main/chat (websocket)
- https://www.bilibili.com/video/BV1se4ReWEHL (架构 timeline 读写扩散)
- http://www.52im.net/thread-1616-1-1.html (读写扩散)
- https://xie.infoq.cn/article/19e95a78e2f5389588debfb1c (推拉)
- https://zhuanlan.zhihu.com/p/65032348 (timeline)