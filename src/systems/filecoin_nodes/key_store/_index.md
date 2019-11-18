---
title: Key Store - 密钥存储
statusIcon: 🛑
---

The `Key Store` is a fundamental abstraction in any full Filecoin node used to store the keypairs associated to a given miner's address and distinct workers (should the miner choose to run multiple workers).

`密钥存储`是任何完整Filecoin节点中的一个基本抽象，用于存储与给定矿工地址和不同工人(如果矿工选择运行多个工人)相关联的密钥对。

Node security depends in large part on keeping these keys secure. To that end we recommend keeping keys separate from any given subsystem and using a separate key store to sign requests as required by subsystems as well as keeping those keys not used as part of mining in cold storage.

节点安全性在很大程度上取决于这些密钥的安全性。为此，我们建议将密钥与任何给定的子系统分开，并使用一个单独的密钥存储库来签署子系统所需的请求，同时将那些密钥不作为挖掘的一部分保存在冷藏库中。

{{< readfile file="key_store.id" code="true" lang="go" >}}
{{< readfile file="key_store.go" code="true" lang="go" >}}

TODO:

- describe the different types of keys used in the protocol and their usage
- 描述协议中使用的不同类型的密钥及其用法
- clean interfaces for getting signatures for full filecoin mining cycles
- 干净的接口，用于获取完整的filecoin挖掘周期的签名
- potential reccomendations or clear disclaimers with regards to consequences of failed key security
- 关于密钥安全性失败的后果的潜在修订或明确的免责声明
- protocol for changing worker keys in filecoin
- 用于在filecoin中更改工作密钥的协议
