---
title: Message Pool - 消息池
statusIcon: 🛑
entries:
- message_syncer
- message_storage
---

{{<label message_pool>}}
The Message Pool is a subsystem in the Filecoin blockchain system. The message pool is acts as the interface between Filecoin nodes and a peer-to-peer network used for off-chain message transmission. It is used by nodes to maintain a set of messages to transmit to the Filecoin VM (for "on-chain" execution).

消息池是Filecoin区块链系统中的一个子系统。消息池充当Filecoin节点和用于脱链消息传输的对等网络之间的接口。节点使用它来维护一组要传输到Filecoin VM的消息(用于"链上"执行)。

{{< readfile file="message_pool_subsystem.id" code="true" lang="go" >}}

Clients that use a message pool include:

使用消息池的客户端包括:

- storage market provider and client nodes - for transmission of deals on chain
- 存储市场供应商和客户节点 - 用于传输链上的交易
- storage miner nodes - for transmission of PoSts, sector commitments, deals, and other operations tracked on chain
- 存储矿工节点 - 用于时空证明的传输，扇区承诺，交易，和其他操作的链上跟踪
- verifier nodes - for transmission of potential faults on chain
- 验证节点 - 用于潜在故障在链上的传输
- relayer nodes - for forwarding and discarding messages appropriately.
- 中继器节点 - 用于适当地转发和丢弃消息。

The message pool subsystem is made of two components:

消息池子系统由两个组件组成：

- The message syncer {{<sref message_syncer>}} -- which receives and propagates messages.
- 消息同步器{{<sref message_syncer "消息同步器">}} -- 接收和传播消息。
- Message storage {{<sref message_storage>}} -- which caches messages according to a given policy.
- 消息存储{{<sref message_storage "消息存储">}} -- 根据给定的策略缓存消息。

TODOs:

计划：

- discuss how messages are meant to propagate slowly/async
- 讨论消息是如何缓慢/异步传播的
- explain algorithms for choosing profitable txns
- 解释选择可盈利的txns的算法

