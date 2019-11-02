---
title: "GossipSub - 闲聊订阅"
---

{{<label gossip_sub>}}

Messages and block headers along side the message references are propagated using the [gossipsub libp2p pubsub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub). Every full node must implement and run that protocol. All pubsub messages are authenticated and must be [syntactically validated](./validation.md#syntactical-validation) before being propagated further.

消息和消息旁边的块头使用[gossipsub libp2p pubsub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub)传播。每个完整节点必须实现并运行该协议。所有的发布/订阅消息都是经过身份验证的，在进一步传播之前必须进行[语法验证](./validation.md#syntactical-validation)。

Further more, every full node must implement and offer the bitswap protocol and provide all Cid Referenced objects, it knows of, through it. This allows any node to fetch missing pieces (e.g. `Message`) from any node it is connected to. However, the node should fan out these requests to multiple nodes and not bombard any single node with too many requests at a time. A node may implement throttling and DDoS protection to prevent such a bombardment.

此外，每个完整节点必须实现和提供位交换协议，并通过它提供它所知道的所有Cid引用对象。这允许任何节点获取丢失的片段(例如，`消息`)来自它所连接的任何节点。但是，节点应该将这些请求分散到多个节点，而不是一次向单个节点发出过多的请求。节点可以实现节流和DDoS保护，以防止这种轰击。

# Bitswap - 位交换

Run bitswap to fetch and serve data (such as blockdata and messages) to and from other filecoin nodes. This is used to fill in missing bits during block propagation, and also to fetch data during sync.

运行位交换向其他filecoin节点获取和提供数据(比如块数据和消息)。这用于在块传播期间填充丢失的位，也用于在同步期间获取数据。

There is not yet an official spec for bitswap, but [the protobufs](https://github.com/ipfs/go-bitswap/blob/master/message/pb/message.proto) should help in the interim.

位交换目前还没有正式的规范，但是[the protobufs](https://github.com/ipfs/go-bitswap/blob/master/message/pb/message.proto)应该会有帮助。
