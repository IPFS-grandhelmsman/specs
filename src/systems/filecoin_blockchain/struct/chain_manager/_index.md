---
title: Chain Manager - 链管理器
---

The _Chain Manager_ is a central component in the blockchain system. It tracks and updates competing subchains received by a given node in order to select the appropriate blockchain head: the latest block of the heaviest subchain it is aware of in the system.

_链管理器_ 是区块链系统中的一个中心组件。它跟踪并更新一个给定节点接收到的竞争子链，以选择适当的区块链头:它在系统中知道的最重子链的最新块。

In so doing, the _chain manager_ is the central subsystem that handles bookkeeping for numerous other systems in a Filecoin node and exposes convenience methods for use by those systems, enabling systems to sample randomness from the chain for instance, or to see which block has been finalized most recently.

这样做时，_链管理器_ 是中心子系统，它处理Filecoin节点中许多其他系统的记录本，并公开这些系统使用的方便方法，使系统能够从链中取样随机性，或者查看最近完成了哪个块。

The chain manager interfaces and functions are included here, but we expand on important details below for clarity.

链管理器的接口和功能包括在这里，但是为了清晰起见，我们在下面展开重要的细节。

# Chain Expansion - 链扩张

## Incoming blocks and semantic validation 区块传入和语义校验

Once a block has been received and syntactically validated by the {{<sref chainsync>}}, it must be semantically validated by the chain manager for inclusion on a given chain.

一旦接收到一个区块并通过{{<sref chainsync "区块同步">}}进行语法验证，它必须经过链管理器的语义验证才能包含在给定的链上。

A semantically valid block:

一个语义上有效的区块:

- must be from a valid miner with power on the chain
- 必须是一个来自链上有存储能力的有效矿工
- must only have valid parents in the tipset, meaning
- 在Tipset中必须仅有且有效的父级，这意味着什么
  - that each parent itself must be a valid block
  - 每个父级本身必须是一个有效的区块
  - that they all have the same parents themselves
  - 他们自己都有相同的父级
  - that they are all at the same height (i.e. include the same number of tickets)
  - 他们都在同一高度(即包括相同数量的票据)
- must have a valid tickets generated from the minTicket in its parent tipset.
- 必须具有从其父tipset中的最小票据生成的有效票据。
- must only have valid state transitions:
- 必须只有有效的状态转换:
  - all messages in the block must be valid
  - 在该区块中所有消息必须被校验
  - the execution of each message, in the order they are in the block, must produce a receipt matching the corresponding one in the receipt set of the block.
  - 每个消息的执行，按照它们在块中的顺序，必须生成一个收据匹配到对应的块收据集合中的一个。
- the resulting state root after all messages are applied, must match the one in the block
- 所有消息生效后的结果状态根，必须匹配块中的一个。


{{% notice info %}}
Once the block passes validation, it must be added to the local datastore, regardless whether it is understood as the best tip at this point. Future blocks from other miners may be mined on top of it and in that case we will want to have it around to avoid refetching.

一旦块通过验证，就必须将它添加到本地数据存储中，不管它是否被理解为此时的最佳的Tip。未来来自其他矿工的区块可能会在它上面开采，在这种情况下，我们希望它在附近以避免重新抓取。
{{% /notice %}}

{{% notice info %}}
To make certain validation checks simpler, blocks should be indexed by height and by parent set. That way sets of blocks with a given height and common parents may be quickly queried. It may also be useful to compute and cache the resultant aggregate state of blocks in these sets, this saves extra state computation when checking which state root to start a block at when it has multiple parents.

为了使某些验证检查更简单，块应该根据高度和父集建立索引。这样，具有给定高度和公共父集的块集可以快速查询。在这些集合中计算和缓存块的最终聚合状态可能也很有用，这在检查当一个块有多个父块时在哪个状态根启动块时节省了额外的状态计算。
{{% /notice %}}

The following requires having and processing (executing) the messages

以下要求拥有和处理(执行)消息

- Messages can be checked by verifying the messages hash correctly to the value.
- 消息可以通过验证消息哈希值是否正确来检查。
- MessageAggregateSig can be checked by verifying the messages sign correctly
- 消息集合签名能够被通过验证消息签名的正确性来检查
- MessageReceipts can only be checked by executing the messages
- 消息收据只能通过执行消息来检查
- StateRoot is the result of the execution of the messages, and can only be verified by executing them
- 状态根是消息执行的结果，只能通过执行它们来验证

## Block reception algorithm - 区块接收算法

Chain selection is a crucial component of how the Filecoin blockchain works. Every chain has an associated weight accounting for the number of blocks mined on it and so the power (storage) they track. It is always preferable to mine atop a heavier Tipset rather than a lighter one. While a miner may be foregoing block rewards earned in the past, this lighter chain is likely to be abandoned by other miners forfeiting any block reward earned as miners converge on a final chain. For more on this, see [chain selection](expected-consensus.md#chain-selection) in the Expected Consensus spec.

链选择是Filecoin区块链如何工作的一个关键组件。每条链都有一个相关的权重，用于计算其上开采的块的数量以及它们跟踪的能力(存储)。它总是在一个更重的Tipset挖掘会更好于在一个更轻的上。当一个矿工放弃过去获得的区块奖励时，这个较轻的链条可能会被其他矿工放弃，当矿工聚集在最后一条链条上时，任何区块奖励都会被放弃。有关这方面的更多信息，请参见预期共识规范中的[链选择](expected-consensus.md#chain-selection)。

However, ahead of finality, a given subchain may be abandoned in order of another, heavier one mined in a given round. In order to rapidly adapt to this, the chain manager must maintain and update all subchains being considered up to finality.

然而，在结束之前，一个给定的子链因另一个的排序而可能会被放弃，更重的一个子链在给定的回合中会被开采。为了快速适应这一点，链管理器必须维护和更新所有被认为是最终的子链。

That is, for every incoming block, even if the incoming block is not added to the current heaviest tipset, the chain manager should add it to the appropriate subchain it is tracking, or keep track of it independently until either:

也就是说，对于每一个传入块，即使传入块没有添加到当前最重的tipset，链管理器也应该将它添加到它正在跟踪的适当子链中，或者保持独立跟踪它直到:

- it is able to do so, through the reception of another block in that subchain
- 所以它能这样做，通过接收子链中的另一个块
- it is able to discard it, as that block was mined before finality
- 它可以丢弃它，因为在最终结果之前已经挖掘了该块

We give an example of how this could work in the block reception algorithm.

我们给出了一个在块接收算法中如何工作的例子。

### ChainTipsManager - 区块Tips管理器

The Chain Tips Manager is a subcomponent of Filecoin consensus that is technically up to the implementer, but since the pseudocode in previous sections reference it, it is documented here for clarity.

区块Tips管理器是Filecoin共识的一个子组件，它在技术上由实现者决定，但是由于前面几节中的伪代码引用了它，所以在这里对它进行了说明。

The Chain Tips Manager is responsible for tracking all live tips of the Filecoin blockchain, and tracking what the current 'best' tipset is.

区块Tips管理器负责跟踪Filecoin区块链上所有存活的Tips，并跟踪什么是当前的'最佳'的Tipset。

```go
// Returns the ticket that is at round 'r' in the chain behind 'head'
func TicketFromRound(head Tipset, r Round) {}

// Returns the tipset that contains round r (Note: multiple rounds' worth of tickets may exist within a single block due to losing tickets being added to the eventually successfully generated block)
func TipsetFromRound(head Tipset, r Round) {}

// GetBestTipset returns the best known tipset. If the 'best' tipset hasn't changed, then this
// will return the previous best tipset.
func GetBestTipset()

// Adds the losing ticket to the chaintips manager so that blocks can be mined on top of it
func AddLosingTicket(parent Tipset, t Ticket)
```
