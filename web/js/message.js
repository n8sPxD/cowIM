// js/message.js

let MessageType;

// 加载并编译 Protobuf
export async function loadProto(protoPath) {
    const response = await fetch(protoPath);
    if (!response.ok) throw new Error('无法加载 proto 文件');

    const protoText = await response.text();

    const root = protobuf.parse(protoText).root;
    MessageType = root.lookupType('Message');
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

    const message = MessageType.decode(new Uint8Array(buffer));
    const object = MessageType.toObject(message, {
        longs: String,
        enums: String,
        bytes: String,
        defaults: true
    });
    return object;
}
