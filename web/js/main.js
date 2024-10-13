// js/main.js

import {
    addFriend,
    addGroup,
    addMessage,
    getAllFriends,
    getAllGroups,
    getAllMessages,
    getChatMessages,
    getFriendByID,
    getGroupByID,
    getLatestTimestamp,
    initDB
} from './db.js';
import {connectWebSocket, onMessageReceived, sendMessageWithAck} from './websocket.js';
import {deserializeMessage, loadProto, serializeMessage} from './message.js';
import {generateUUID} from "./utils.js";
import {GROUP_CHAT, MSG_COMMON_MSG, SINGLE_CHAT} from "./constant.js";

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
        await loadProto('../proto/message.proto');

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
    if (timelineData.infos && Array.isArray(timelineData.infos)) {
        for (const timeline of timelineData.infos) {
            const {msgForward} = timeline;
            // 将后端返回的消息转换为 IndexedDB 所需的消息格式
            const formattedMessage = {
                id: msgForward.id.toString(),         // 消息 ID，转换为字符串
                from: timeline.senderID,               // 发送者 ID
                to: timeline.receiverID,               // 接受者 ID
                group: timeline.groupID, // 如果是群聊则有群组 ID
                content: msgForward.content,           // 消息内容
                type: timeline.type,              // 消息类型
                msgType: msgForward.msgType,          // 聊天类型
                extend: msgForward.extend,
                timestamp: new Date(msgForward.timestamp).getTime() // 时间戳转换为时间戳（毫秒）
            };
            await addMessage(formattedMessage); // 假设 addMessage 函数会处理消息存储
        }
    } else {
        console.log("没有最新的同步消息")
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
    if (groupsData.content.infos && Array.isArray(groupsData.content.infos)) {
        for (const group of groupsData.content.infos) {
            const tmpGroup = {
                groupID: group.groupId,
                groupName: group.groupName,
                groupAvatar: group.groupAvatar,
            }
            await addGroup(tmpGroup);
        }
    } else {
        console.log("群组为空")
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
    if (friendsData.content.friends && Array.isArray(friendsData.content.friends)) {
        for (const friend of friendsData.content.friends) {
            const formattedFriend = {
                friendID: friend.friendId,
                friendName: friend.username,
                friendAvatar: friend.avatar,
            };
            await addFriend(formattedFriend);
        }
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

// 显示最近会话列表
async function displayRecentConversations(selectedChatID = null) {
    const conversationList = document.getElementById('conversationList');
    conversationList.innerHTML = ''; // 清空现有列表

    const messages = await getAllMessages();
    const recentChats = new Map(); // key: chatID, value: latest message

    messages.forEach(msg => {
        if (msg.type === SINGLE_CHAT) { // SINGLE_CHAT
            const chatID = msg.from === Number(sessionStorage.getItem('CowID')) ? msg.to : msg.from;
            if (!recentChats.has(chatID) || recentChats.get(chatID).timestamp < msg.timestamp) {
                recentChats.set(chatID, msg);
            }
        } else if (msg.type === GROUP_CHAT) { // GROUP_CHAT
            const groupID = msg.group;
            if (!recentChats.has(`group_${groupID}`) || recentChats.get(`group_${groupID}`).timestamp < msg.timestamp) {
                recentChats.set(`group_${groupID}`, msg);
            }
        }
    });

    // 转换为数组并按时间戳降序排序
    const sortedChats = Array.from(recentChats.entries()).sort((a, b) => b[1].timestamp - a[1].timestamp);

    for (const [chatID, msg] of sortedChats) {
        let displayName;
        if (String(chatID).startsWith('group_')) {
            const groupID = chatID.split('_')[1];
            const group = await getGroupByID(Number(groupID));
            displayName = group ? group.groupName + `(${groupID})` : `群组 ${groupID}`;
        } else {
            const friend = await getFriendByID(Number(chatID));
            displayName = friend ? friend.friendName + `(${chatID})` : `用户 ${chatID}`;
        }

        const chatItem = document.createElement('div');
        chatItem.textContent = displayName;
        chatItem.dataset.chatId = chatID;
        chatItem.classList.add('chat-item');
        chatItem.addEventListener('click', () => selectConversation(chatID));

        // 保持选中的会话
        if (chatID == selectedChatID) {
            chatItem.classList.add('selected');
        }

        conversationList.appendChild(chatItem);
    }
}


// 显示好友列表
async function displayFriendsList() {
    const conversationList = document.getElementById('conversationList');
    conversationList.innerHTML = ''; // 清空现有列表

    const friends = await getAllFriends();

    // 按用户名排序
    friends.sort((a, b) => a.friendName.localeCompare(b.friendName));

    for (const friend of friends) {
        const friendItem = document.createElement('div');
        friendItem.textContent = friend.friendName;
        friendItem.dataset.friendId = friend.friendID;
        friendItem.classList.add('friend-item');
        friendItem.addEventListener('click', () => selectConversation(friend.friendID));

        conversationList.appendChild(friendItem);
    }
}

// 显示群组列表
async function displayGroupsList() {
    const conversationList = document.getElementById('conversationList');
    conversationList.innerHTML = ''; // 清空现有列表

    const groups = await getAllGroups();

    // 按群名排序
    groups.sort((a, b) => a.groupName.localeCompare(b.groupName));

    for (const group of groups) {
        const groupItem = document.createElement('div');
        groupItem.textContent = group.groupName;
        groupItem.dataset.groupID = group.groupID;
        groupItem.classList.add('group-item');
        groupItem.addEventListener('click', () => selectConversation(`group_${group.groupID}`));

        conversationList.appendChild(groupItem);
    }
}

// 选择一个会话（好友或群组）
async function selectConversation(chatID) {
    // 标记选中的会话
    document.querySelectorAll('.chat-item, .friend-item, .group-item').forEach(item => {
        item.classList.remove('selected');
    });
    const selectedItem = document.querySelector(`[data-chat-id="${chatID}"], [data-friend-id="${chatID}"], [data-group-id="${chatID}"]`);
    if (selectedItem) {
        selectedItem.classList.add('selected');
    }

    const chatHeader = document.getElementById('chatHeader');
    const chatHistory = document.getElementById('chatHistory');

    // 设置聊天头部
    if (String(chatID).startsWith('group_')) {  // 确保 chatID 是字符串
        const groupID = Number(chatID.split('_')[1]);
        const group = await getGroupByID(Number(groupID));
        chatHeader.textContent = group ? group.groupName : `群组 ${groupID}`;
    } else {
        const friend = await getFriendByID(Number(chatID));
        chatHeader.textContent = friend ? friend.friendName : `用户 ${chatID}`;
    }


    // 加载聊天记录
    const messages = await getChatMessages(chatID);
    const cowID = Number(sessionStorage.getItem('CowID'));
    chatHistory.innerHTML = ''; // 清空现有记录

    messages.forEach(msg => {
        const messageElement = document.createElement('div');
        messageElement.style.textAlign = msg.from === cowID ? "right" : "left";
        messageElement.textContent = msg.content;
        chatHistory.appendChild(messageElement);
    });

    // 滚动到底部
    chatHistory.scrollTop = chatHistory.scrollHeight;
}

// 处理发送消息
async function handleSendMessage() {
    const messageInput = document.getElementById('messageInput');
    const content = messageInput.value.trim();
    if (!content) return;

    const chatID = getSelectedChatID();
    if (!chatID) {
        alert('请选择一个会话');
        return;
    }

    const from = Number(sessionStorage.getItem('CowID'));
    let to;
    let group = null;

    if (chatID.startsWith('group_')) {
        group = Number(chatID.split('_')[1]);
        to = group
    } else {
        to = Number(chatID);
    }

    // 创建消息对象
    // TODO: 细化消息对象
    const message = {
        id: generateUUID(),
        from: from,
        to: to,
        group: group,
        content: content,
        type: group ? GROUP_CHAT : SINGLE_CHAT,
        msgType: MSG_COMMON_MSG,
        extend: null,
        timestamp: Date.now()
    };

    try {
        // 序列化消息
        const serializedMessage = serializeMessage(message);

        // 发送消息并处理 ACK
        sendMessageWithAck(message, serializedMessage);

        // 存储消息到 IndexedDB
        await addMessage(message);

        // 更新 UI
        appendMessageToChatHistory(message);

        // 清空输入框
        messageInput.value = '';
    } catch (error) {
        console.error('发送消息失败:', error);
        alert('发送消息失败，请重试');
    }
}

// 获取选中的会话ID
function getSelectedChatID() {
    const selectedItem = document.querySelector('.chat-item.selected, .friend-item.selected, .group-item.selected');
    if (selectedItem) {
        return selectedItem.dataset.chatId || selectedItem.dataset.friendId || selectedItem.dataset.groupId;
    }
    return null;
}

// 将消息追加到聊天记录
function appendMessageToChatHistory(message) {
    const chatHistory = document.getElementById('chatHistory');
    const cowID = Number(sessionStorage.getItem('CowID'))

    const messageElement = document.createElement('div');
    messageElement.style.textAlign = message.from === cowID ? "right" : "left";
    messageElement.textContent = message.content
    chatHistory.appendChild(messageElement);

    // 滚动到底部
    chatHistory.scrollTop = chatHistory.scrollHeight;
}

// 处理接收到的消息
async function handleIncomingMessage(data) {
    // 假设接收到的消息是 Protobuf 二进制数据
    const message = deserializeMessage(data);

    // 存储消息到 IndexedDB
    await addMessage(message);

    // 获取当前选中的会话 ID
    const currentChatID = getSelectedChatID();

    // 如果消息属于当前会话，追加到聊天记录
    const messageChatID = message.group ? `group_${message.group}` : (message.from === Number(sessionStorage.getItem('CowID')) ? message.to : message.from);

    if (messageChatID == currentChatID) {
        appendMessageToChatHistory(message);
    }

    // 更新最近会话列表，并保持当前选中的会话
    await displayRecentConversations(currentChatID);
}


document.addEventListener('DOMContentLoaded', async () => {
    await initializeMain();
});

// 其他函数如 displayRecentConversations, displayFriendsList, displayGroupsList, selectConversation, handleSendMessage, appendMessageToChatHistory, handleIncomingMessage 等保持不变

// 生成 UUID 已在 utils.js 中定义
