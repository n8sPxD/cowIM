// js/message.js

let MessageType;

// 加载并编译 Protobuf
export async function loadProto(protoPath) {
    protobuf.load(protoPath)
        .then(root => {
            MessageType = root.lookupType("Message");
        })
        .catch(err => {
            console.error("Failed to load protobuf schema:", err);
        });
}

// 序列化消息对象为 Protobuf 二进制
export function serializeMessage(messageObj) {
    if (!MessageType) throw new Error('Protobuf 未加载');

    const errMsg = MessageType.verify(messageObj);
    if (errMsg) throw Error(errMsg);

    const message = MessageType.create(messageObj);
    const buffer = MessageType.encode(message).finish();
    return buffer;
}

// 反序列化 Protobuf 二进制为消息对象
export function deserializeMessage(buffer) {
    if (!MessageType) throw new Error('Protobuf 未加载');

    const message = MessageType.decode(buffer instanceof Uint8Array ? buffer : new Uint8Array(buffer));
    return message
}
