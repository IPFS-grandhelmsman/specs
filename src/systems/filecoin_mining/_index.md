---
menuIcon: ⛏
menuTitle: "**Storage Mining - 存储挖掘**"
title: "Storage Mining System - proving storage for producing blocks - 存储挖掘系统 - 为生产区块校验存储"
entries:
- storage_mining
- sector
- sector_index
- storage_proving
---

{{< incTocMap "/docs/systems/filecoin_mining" 2 >}}


The Storage Mining System is the part of the Filecoin Protocol that deals with storing Client's
data, producing proof artifacts that demonstrate correct storage behavior, and managing the work
involved.

存储挖掘系统是Filecoin协议的一部分，用于处理存储客户端的数据、生成演示正确存储行为的证明构件和管理相关工作。

Storing data and producing proofs is a complex, highly optimizable process, with lots of tunable
choices. Miners should explore the design space to arrive at something that (a) satisfies protocol
and network-wide constraints, (b) satisfies clients' requests and expectations (as expressed in
`Deals`), and \(c) gives them the most cost-effective operation. This part of the Filecoin Spec
primarily describes in detail what MUST and SHOULD happen here, and leaves ample room for
various optimizations for implementers, miners, and users to make. In some parts, we describe
algorithms that could be replaced by other, more optimized versions, but in those cases it is
important that the **protocol constraints** are satisfied. The **protocol constraints** are
spelled out in clear detail (an unclear, unmentioned constraint is a "spec error").  It is up
to implementers who deviate from the algorithms presented here to ensure their modifications
satisfy those constraints, especially those relating to protocol security.

存储数据和生成证明是一个复杂的、高度优化的过程，有许多可调选项。采矿者应该探索设计空间，以达到以下目的：(a)满足协议和网络范围的约束，(b)满足客户的要求和期望(如“交易”中所表达的)，\(c)为他们提供最具成本效益的操作。Filecoin规范的这一部分主要详细描述了这里必须发生和应该发生的事情，并为实现者、矿工和用户提供了大量的优化空间。在某些部分，我们描述了可以被其他更优化的版本替换的算法，但是在这些情况下，满足 **协议约束** 是很重要的。**协议约束** 被详细的描述出来(一个不清楚的，未被提及的约束是一个“规范错误”)。要由偏离这里提供的算法的实现者来确保他们的修改满足这些约束，特别是那些与协议安全性相关的约束。
