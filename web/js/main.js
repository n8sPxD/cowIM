// js/main.js

import {addMessage, getLatestTimestamp, initDB} from './db.js';
import {connectWebSocket, onMessageReceived} from './websocket.js';
import {loadProto} from './message.js';

// 页面加载完成后初始化
export async function initializeMain() {
    const jwtToken = sessionStorage.getItem('jwtToken');
    const cowID = sessionStorage.getItem('CowID');

    if (!jwtToken || !cowID) {
        // 未登录，跳转到登录页
        window.location.href = 'login.html';
        return;
    }

    try {
        // 初始化 IndexedDB
        await initDB();

        // 加载 Protobuf
        await loadProto('/proto/message.proto');

        // 获取 WebSocket 服务器地址
        const wsServerIP = await getWebSocketServerIP();

        // 连接 WebSocket
        await connectWebSocket(wsServerIP, jwtToken);

        // 同步初始数据
        await fetchInitialData();

        // 初始化 UI
        initializeUI();

        // 处理接收到的消息
        onMessageReceived(handleIncomingMessage);

    } catch (error) {
        console.error('初始化失败:', error);
        alert('初始化失败，请稍后重试');
    }
}

// 获取 WebSocket 服务器地址
async function getWebSocketServerIP() {
    const jwtToken = sessionStorage.getItem('jwtToken');

    const response = await fetch('http://localhost:8080/wsget', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${jwtToken}`
        }
    });

    if (!response.ok) throw new Error('获取 WebSocket 服务器地址失败');

    const data = await response.json();
    const wsIP = data.content.ip; // 假设返回的 JSON 中包含 wsIP 字段

    if (!wsIP) throw new Error('无效的 WebSocket 服务器地址');

    return wsIP;
}

// 同步初始数据
async function fetchInitialData() {
    const jwtToken = sessionStorage.getItem('jwtToken');

    // 获取最新的本地消息 timestamp
    const latestTimestamp = await getLatestTimestamp();

    // 获取离线消息 /timelinesync，并传递 timestamp 参数
    const timelineResponse = await fetch(`http://localhost:8080/timelinesync?timestamp=${latestTimestamp}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${jwtToken}`
        }
    });

    if (!timelineResponse.ok) throw new Error('获取离线消息失败');

    const timelineData = await timelineResponse.json();
    console.log("timelineData: ", timelineData)
    for (const timeline of timelineData.infos) {
        const {msgForward} = timeline;
        // 将后端返回的消息转换为 IndexedDB 所需的消息格式
        const formattedMessage = {
            id: msgForward.id.toString(),         // 消息 ID，转换为字符串
            from: timeline.senderID,               // 发送者 ID
            to: timeline.receiverID,               // 接受者 ID
            group: msgForward.msgType === 1 ? msgForward.group : undefined, // 如果是群聊则有群组 ID
            content: msgForward.content,           // 消息内容
            type: msgForward.type,              // 消息类型，0 表示单聊，1 表示群聊
            msg_type: msgForward.msgType,          // 聊天类型
            timestamp: new Date(msgForward.timestamp).getTime() // 时间戳转换为时间戳（毫秒）
        };
        await addMessage(timeline.msgForward); // 假设 addMessage 函数会处理消息存储
    }

    // 获取群组信息 /groups
    const groupsResponse = await fetch('http://localhost:8080/groups', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${jwtToken}`
        }
    });

    if (!groupsResponse.ok) throw new Error('获取群组信息失败');

    const groupsData = await groupsResponse.json();
    for (const group of groupsData) {
        await addGroup(group);
    }

    // 获取好友列表 /friends
    const friendsResponse = await fetch('http://localhost:8080/friends', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${jwtToken}`
        }
    });

    if (!friendsResponse.ok) throw new Error('获取好友列表失败');

    const friendsData = await friendsResponse.json();
    for (const friend of friendsData) {
        await addFriend(friend);
    }
}

// 初始化 UI 组件和事件监听
function initializeUI() {
    // 显示最近会话列表
    displayRecentConversations();

    // 设置侧边栏按钮事件
    document.getElementById('recentButton').addEventListener('click', displayRecentConversations);
    document.getElementById('friendsButton').addEventListener('click', displayFriendsList);
    document.getElementById('groupsButton').addEventListener('click', displayGroupsList);

    // 设置发送按钮事件
    document.getElementById('sendButton').addEventListener('click', handleSendMessage);

    // 设置消息输入框回车发送事件
    document.getElementById('messageInput').addEventListener('keypress', (event) => {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            handleSendMessage();
        }
    });
}

document.addEventListener('DOMContentLoaded', async () => {
    await initializeMain();
});

// 其他函数如 displayRecentConversations, displayFriendsList, displayGroupsList, selectConversation, handleSendMessage, appendMessageToChatHistory, handleIncomingMessage 等保持不变

// 生成 UUID 已在 utils.js 中定义
