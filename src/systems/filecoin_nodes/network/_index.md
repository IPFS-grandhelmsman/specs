---
title: Network Interface - 网络接口
statusIcon: ⚠️
---

{{< readfile file="network.id" code="true" lang="go" >}}


Filecoin nodes use the libp2p protocol for peer discovery, peer routing, and message multicast, and so on. Libp2p is a set of modular protocols common to the peer-to-peer networking stack. Nodes open connections with one another and mount different protocols or streams over the same connection. In the initial handshake, nodes exchange the protocols that each of them supports and all Filecoin related protcols will be mounted under `/filecoin/...` protocol identifiers.

Filecoin节点使用libp2p协议进行对等发现、对等路由和消息多播等。Libp2p是一组对点对点网络堆栈通用的模块协议。节点彼此打开连接，并在同一连接上挂载不同的协议或流。在最初的握手过程中，节点交换各自支持的协议，所有与Filecoin相关的原协议都将挂载在`/filecoin/...`的协议标识符。

Here is the list of libp2p protocols used by Filecoin.

下面是Filecoin使用的libp2p协议列表。

- Graphsync: TODO
- 图形同步: TODO
- Bitswap:  TODO
- 比特交换:  TODO
- Gossipsub: block headers and messages are broadcasted through a Gossip PubSub protocol where nodes can subscribe to topics such as `NewBlock`, `BlockHeader`, `BlockMessage`, etc and receive messages in those topics. When receiving messages related to a topic, nodes processes the message and forwards it to its peers who also subscribed to the same topic.
- 订阅闲聊：块标头和消息通过闲聊PubSub协议广播，节点可以订阅诸如`NewBlock`、`BlockHeader`、`BlockMessage`等主题，并接收这些主题中的消息。当接收到与某个主题相关的消息时，节点将处理该消息并将其转发给同样订阅了该主题的节点
- Kad-DHT: Kademlia DHT is a distributed hash table with a logarithmic bound on the maximum number of lookups for a particular node. Kad DHT is used primarily for peer routing as well as peer discovery in the Filecoin protocol.
- Kademlia DHT: Kademlia DHT是一个分布式哈希表，对特定节点的最大查找次数有一个对数界。在Filecoin协议中，Kad DHT主要用于对等路由和对等发现。
- Bootstrap: Bootstrap is a list of nodes that a new node attempts to connect upon joining the network. The list of bootstrap nodes and their addresses are defined by the users.
- Bootstrap: 启动引导是新节点在加入网络时尝试连接的节点列表。启动引导节点及其地址的列表由用户定义。
