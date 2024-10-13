// js/websocket.js

import {getMessageByID, updateMessage} from "./db.js";
import {deserializeMessage} from "./message.js";
import {MSG_ACK_MSG, SYSTEM_INFO, USER_SYSTEM} from "./constant.js";

let websocket;
let messageHandler = null;
let ackHandler = new Map(); // Map<messageID, { messageObj, serializedMessage, retries, timeout }>

const MAX_RETRIES = 3;
const ACK_TIMEOUT = 3000; // 2秒

// 连接 WebSocket 服务器
export function connectWebSocket(wsIP, jwtToken) {
    return new Promise((resolve, reject) => {
        const wsURL = `ws://${wsIP}/ws?token=${jwtToken}`; // 根据实际情况调整 URL

        websocket = new WebSocket(wsURL);

        websocket.binaryType = "arraybuffer";

        websocket.onopen = () => {
            console.log('WebSocket 连接成功');
            resolve();
        };

        websocket.onmessage = (event) => {
            handleIncomingData(event.data);
        };

        websocket.onerror = (event) => {
            console.error('WebSocket 连接错误:', event);
            reject(event);
        };

        websocket.onclose = (event) => {
            console.log('WebSocket 连接关闭:', event);
            // 可选：实现重连逻辑
        };
    });
}

// 处理接收到的 WebSocket 数据
function handleIncomingData(data) {
    try {
        console.log("data: ", data)
        const parsedData = deserializeMessage(data);

        console.log("接受消息: ", parsedData)

        if (parsedData.from === USER_SYSTEM && parsedData.type === SYSTEM_INFO && parsedData.msg_type === MSG_ACK_MSG) {
            console.log("系统消息")
            const messageID = parsedData.id;
            if (ackHandler.has(messageID)) {
                clearTimeout(ackHandler.get(messageID).timeout);
                ackHandler.delete(messageID);
                console.log(`消息 ${messageID} 已被确认`);

                // 从IndexedDB中获取消息，更新ID
                const messages = getMessageByID(messageID);
                const messageToUpdate = messages.find(msg => msg.id === messageID);

                if (messageToUpdate) {
                    messageToUpdate.id = parsedData.content;
                    updateMessage(messageToUpdate);
                    console.log(`消息 ID 更新为 ${parsedData.content}`);
                } else {
                    console.warn(`未找到 ID 为 ${messageID} 的消息`);
                }
            }
        } else {
            console.log("正常消息")
            if (messageHandler) {
                messageHandler(data);
            }
        }
    } catch (error) {
        console.error('处理 WebSocket 数据时出错:', error);
        // 如果不是 JSON 格式，假设是 Protobuf 消息
        if (messageHandler) {
            messageHandler(data);
        }
    }
}

// 发送消息并处理 ACK
export function sendMessageWithAck(messageObj, serializedMessage) {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
        websocket.send(serializedMessage);

        const messageID = messageObj.id;
        const retries = 0;

        const timeout = setTimeout(() => {
            handleAckTimeout(messageID);
        }, ACK_TIMEOUT);

        ackHandler.set(messageID, { messageObj, serializedMessage, retries, timeout });
    } else {
        console.error('WebSocket 未连接');
    }
}

// 处理 ACK 超时与重传
function handleAckTimeout(messageID) {
    if (!ackHandler.has(messageID)) return;

    const entry = ackHandler.get(messageID);

    if (entry.retries < MAX_RETRIES) {
        console.log(`重发消息 ${messageID}，尝试次数: ${entry.retries + 1}`);
        websocket.send(entry.serializedMessage);

        // 更新重试次数和重置超时
        entry.retries += 1;
        entry.timeout = setTimeout(() => {
            handleAckTimeout(messageID);
        }, ACK_TIMEOUT);

        ackHandler.set(messageID, entry);
    } else {
        console.error(`消息 ${messageID} 发送失败，超过最大重试次数`);
        // 可选：通知 UI 消息发送失败
        ackHandler.delete(messageID);
    }
}

// 注册接收消息的处理器
export function onMessageReceived(handler) {
    messageHandler = handler;
}
