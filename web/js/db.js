// js/db.js

const DB_NAME = 'CowIMDB';
const DB_VERSION = 1;
let db;

// 初始化 IndexedDB
export function initDB() {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, DB_VERSION);

        request.onerror = (event) => {
            console.error('IndexedDB 打开失败:', event.target.error);
            reject(event.target.error);
        };

        request.onsuccess = (event) => {
            db = event.target.result;
            resolve();
        };

        request.onupgradeneeded = (event) => {
            db = event.target.result;

            // 创建 messages 表
            if (!db.objectStoreNames.contains('messages')) {
                const messagesStore = db.createObjectStore('messages', {keyPath: 'id'});
                messagesStore.createIndex('from', 'from', {unique: false});
                messagesStore.createIndex('to', 'to', {unique: false});
                messagesStore.createIndex('group', 'group', {unique: false});
                messagesStore.createIndex('timestamp', 'timestamp', {unique: false});
            }

            // 创建 friends 表
            if (!db.objectStoreNames.contains('friends')) {
                db.createObjectStore('friends', {keyPath: 'friendID'});
            }

            // 创建 groups 表
            if (!db.objectStoreNames.contains('groups')) {
                db.createObjectStore('groups', {keyPath: 'groupID'});
            }
        };
    });
}

// 添加消息
export function addMessage(message) {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['messages'], 'readwrite');
        const store = transaction.objectStore('messages');
        const request = store.put(message);

        request.onsuccess = () => resolve();
        request.onerror = (event) => {
            console.error('添加消息失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 更新消息
export function updateMessage(message) {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['messages'], 'readwrite');
        const store = transaction.objectStore('messages');
        const request = store.put(message);

        request.onsuccess = () => resolve();
        request.onerror = (event) => {
            console.error('更新消息失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 获取所有消息
export function getAllMessages() {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['messages'], 'readonly');
        const store = transaction.objectStore('messages');
        const request = store.getAll();

        request.onsuccess = (event) => resolve(event.target.result);
        request.onerror = (event) => {
            console.error('获取所有消息失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 通过 ID 获取特定消息
export function getMessageByID(messageID) {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['messages'], 'readonly');
        const store = transaction.objectStore('messages');
        const request = store.get(messageID);

        request.onsuccess = (event) => {
            resolve(event.target.result);
        };
        request.onerror = (event) => {
            console.error('获取消息失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 获取特定会话的消息
export function getChatMessages(chatID) {
    return new Promise(async (resolve, reject) => {
        const transaction = db.transaction(['messages'], 'readonly');
        const store = transaction.objectStore('messages');

        if (typeof chatID === 'number') { // 单聊
            const cowID = Number(sessionStorage.getItem('CowID'));
            const indexFrom = store.index('from');
            const indexTo = store.index('to');

            const fromRequest = indexFrom.getAll(cowID);
            const toRequest = indexTo.getAll(cowID);

            Promise.all([
                new Promise((res, rej) => {
                    fromRequest.onsuccess = (e) => res(e.target.result);
                    fromRequest.onerror = (e) => rej(e.target.error);
                }),
                new Promise((res, rej) => {
                    toRequest.onsuccess = (e) => res(e.target.result);
                    toRequest.onerror = (e) => rej(e.target.error);
                })
            ]).then(([fromMsgs, toMsgs]) => {
                const combined = fromMsgs.concat(toMsgs).filter(msg => {
                    if (msg.type === 0) { // SINGLE_CHAT
                        return (msg.from === cowID && msg.to === chatID) || (msg.from === chatID && msg.to === cowID);
                    }
                    return false;
                });

                combined.sort((a, b) => a.timestamp - b.timestamp);
                resolve(combined);
            }).catch(reject);

        } else if (typeof chatID === 'string' && chatID.startsWith('group_')) { // 群聊
            const groupID = Number(chatID.split('_')[1]);
            const index = store.index('group');
            const groupRequest = index.getAll(groupID);

            groupRequest.onsuccess = (event) => {
                const groupMsgs = event.target.result.filter(msg => msg.type === 1); // GROUP_CHAT
                groupMsgs.sort((a, b) => a.timestamp - b.timestamp);
                resolve(groupMsgs);
            };
            groupRequest.onerror = (event) => {
                console.error('获取群消息失败:', event.target.error);
                reject(event.target.error);
            };
        } else {
            resolve([]);
        }
    });
}

// 添加好友
export function addFriend(friend) {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['friends'], 'readwrite');
        const store = transaction.objectStore('friends');
        const request = store.put(friend);

        request.onsuccess = () => resolve();
        request.onerror = (event) => {
            console.error('添加好友失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 获取所有好友
export function getAllFriends() {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['friends'], 'readonly');
        const store = transaction.objectStore('friends');
        const request = store.getAll();

        request.onsuccess = (event) => resolve(event.target.result);
        request.onerror = (event) => {
            console.error('获取好友列表失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 通过ID获取好友
export function getFriendByID(friendID) {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['friends'], 'readonly');
        const store = transaction.objectStore('friends');
        const request = store.get(friendID);

        request.onsuccess = (event) => resolve(event.target.result);
        request.onerror = (event) => {
            console.error('获取好友失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 添加群组
export function addGroup(group) {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['groups'], 'readwrite');
        const store = transaction.objectStore('groups');
        const request = store.put(group);

        request.onsuccess = () => resolve();
        request.onerror = (event) => {
            console.error('添加群组失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 获取所有群组
export function getAllGroups() {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['groups'], 'readonly');
        const store = transaction.objectStore('groups');
        const request = store.getAll();

        request.onsuccess = (event) => resolve(event.target.result);
        request.onerror = (event) => {
            console.error('获取群组列表失败:', event.target.error);
            reject(event.target.error);
        };
    });
}

// 通过ID获取群组
export function getGroupByID(groupID) {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['groups'], 'readonly');
        const store = transaction.objectStore('groups');
        const request = store.get(groupID);

        request.onsuccess = (event) => resolve(event.target.result);
        request.onerror = (event) => {
            console.error('获取群组失败:', event.target.error);
            reject(event.target.error);
        };
    });
}


// 获取最新的消息 timestamp
export function getLatestTimestamp() {
    return new Promise((resolve, reject) => {
        const transaction = db.transaction(['messages'], 'readonly');
        const store = transaction.objectStore('messages');

        // 打开游标，按时间戳倒序获取最新的记录
        const request = store.index('timestamp').openCursor(null, 'prev');

        request.onsuccess = function (event) {
            const cursor = event.target.result;
            if (cursor) {
                const latestTimestamp = new Date(cursor.value.timestamp).getTime(); // 获取最新的 timestamp
                console.log("timestamp: ", latestTimestamp)
                resolve(latestTimestamp);
            } else {
                resolve(0); // 如果没有记录，返回 null
            }
        };

        request.onerror = function (event) {
            console.error('获取最新 timestamp 失败:', event.target.error);
            reject(event.target.error);
        };
    });
}