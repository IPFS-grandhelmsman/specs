---
title: Message Syncer 消息同步器
---

{{<label message_syncer>}}

TODO:

计划

- explain message syncer works
- 解释消息同步器工作
- include the message syncer code
- 包括消息同步器代码

# Message Propagation - 消息传播

Messages are propagated over the libp2p pubsub channel `/fil/messages`. On this channel, every [serialised `SignedMessage`](data-structures.md#messages) is announced.

消息通过libp2p 发布订阅通道`/fil/ Messages`传播。在这个频道，每一条[序列化的`已签名消息`](data-structures.md#messages)都会被公布。

Upon receiving the message, its validity must be checked: the signature must be valid, and the account in question must have enough funds to cover the actions specified. If the message is not valid it should be dropped and must not be forwarded.

在接收到消息时，必须检查其有效性：签名必须有效，相关帐户必须有足够的资金来支付指定的操作。如果消息无效，则它应被丢弃，并被禁止转发。

{{% notice todo %}}
discuss checking signatures and account balances, some tricky bits that need consideration. Does the fund check cause improper dropping? E.g. I have a message sending funds then use the newly constructed account to send funds, as long as the previous wasn't executed the second will be considered "invalid" ... though it won't be at the time of execution.

讨论检查签名和帐户余额，一些需要考虑的棘手的细节。资金检查是否造成不当的丢失?例如我有一条信息正在发送资金，然后使用新创建的账户发送资金，只要前一个没有执行，第二个将被视为“无效”……虽然不会在执行的时候。
{{% /notice %}}
