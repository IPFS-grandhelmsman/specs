---
title: Introduction - 介绍
entries:
#  - specmap
- arch - 架构
- concepts - 概念
- filecoin_vm - VM
- process - 进程
- changelog - 变更日志
- system - 系统
---

<center><img src="./docs/intro/underconstruction.gif" height="128px" /></center>

{{% notice warning %}}
**Warning:** This draft of the Filecoin protocol specification is a work in progress.
It is intended to establish the rough overall structure of the document,
enabling experts to fill in different sections in parallel.
However, within each section, content may be out-of-order, incorrect, and/or incomplete.
The reader is advised to refer to the
[official Filecoin spec document](https://filecoin-project.github.io/specs/)
for specification and implementation questions.
{{% /notice %}}

{{% notice warning %}}
**警告:** 
Filecoin协议规范的草案正在进行中。
它的目的是建立文件的大致总体结构，
使专家能够并行地填写不同的部分。
但是，在每个部分中，内容可能是无序的、不正确的和/或不完整的。
建议读者参阅
[official Filecoin spec document](https://filecoin-project.github.io/specs/)
关于规范和实现的问题。
{{% /notice %}}


Filecoin is a distributed storage network based on a blockchain mechanism.
Filecoin *miners* can elect to provide storage capacity for the network, and thereby
earn units of the Filecoin cryptocurrency (FIL) by periodically producing
cryptographic proofs that certify that they are providing the capacity specified.
In addition, Filecoin enables parties to exchange FIL currency
through transactions recorded in a shared ledger on the Filecoin blockchain.
Rather than using Nakamoto-style proof of work to maintain consensus on the chain, however,
Filecoin uses proof of storage itself: a miner's power in the consensus protocol
is proportional to the amount of storage it provides.

Filecoin是一个基于区块链机制的分布式存储网络。
Filecoin *miner* 可以选择为网络提供存储容量，从而通过定期生成来获得Filecoin加密货币(FIL)的单位证明它们提供指定容量的密码证明。
此外，Filecoin允许各方交换FIL货币通过记录在Filecoin区块链上的共享分类帐中的交易。
然而，与其使用nakamoto风格的工作证明来维持链上的一致，Filecoin使用本身的存储证明:共识协议中的矿工的能力与它提供的存储总量成正比。

The Filecoin blockchain not only maintains the ledger for FIL transactions and
accounts, but also implements the Filecoin VM, a replicated state machine which executes
a variety of cryptographic contracts and market mechanisms among participants
on the network.
These contracts include *storage deals*, in which clients pay FIL currency to miners
in exchange for storing the specific file data that the clients request.
Via the distributed implementation of the Filecoin VM, storage deals
and other contract mechanisms recorded on the chain continue to be processed
over time, without requiring further interaction from the original parties
(such as the clients who requested the data storage).

Filecoin区块链不仅为FIL交易和维护分类帐但也实现了Filecoin VM，一个执行的复制状态机参与者之间的各种加密契约和市场机制在网络上。
这些合约包括 *储存交易* 客户向矿工支付外汇作为交换，存储客户端请求的特定文件数据。
通过分布式实现的Filecoin虚拟机，存储交易链上记录的其他契约机制随着时间的推移继续被处理，不需要原始方的进一步交互(例如请求数据存储的客户机)。

