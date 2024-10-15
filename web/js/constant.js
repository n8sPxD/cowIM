// js/constant.js

// 消息类型
export const MSG_COMMON_MSG = 1;
export const MSG_REPLY_MSG = 2;
export const MSG_IMAGE_MSG = 3;
export const MSG_VIDEO_MSG = 4;
export const MSG_FILE_MSG = 5;

// 系统相关 消息类型
export const MSG_SYSTEM_MSG = 100;
export const MSG_ACK_MSG = 101;

// 群聊或者单聊
export const SINGLE_CHAT = 1;
export const GROUP_CHAT = 2;
export const BIG_GROUP_CHAT = 3;
export const SYSTEM_INFO = 99;

// 用户类型
export const USER_SYSTEM = 1;
export const USER_COMMON = 2;

// 群聊权限
export const GROUP_MASTER = 1;
export const GROUP_ADMIN = 2;
export const GROUP_COMMON = 3;

// Websocket 关闭代码
export const DUP_CLIENT_CODE = 4001;