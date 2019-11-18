---
title: Block Producer - 区块产生器
---

{{<label block_producer>}}

# Mining Blocks - 挖块

Having registered as a miner, it's time to start making and checking tickets. At this point, the miner should already be running chain validation, which includes keeping track of the latest tipsets seen on the network.

注册为矿工后，是时候开始制作和检查票据了。此时，矿工应该已经在运行链验证，其中包括跟踪网络上看到的最新Tipset。

For additional details around how consensus works in Filecoin, see {{<sref expected_consensus>}}. For the purposes of this section, there is a consensus protocol (Expected Consensus) that guarantees a fair process for determining what blocks have been generated in a round, whether a miner is eligible to mine a block itself, and other rules pertaining to the production of some artifacts required of valid blocks (e.g. Tickets, ElectionProofs).

有关共识如何在Filecoin中工作的更多细节，请参见{{<sref expected_consensus "预期共识">}}。对于本节的目的，都有一个共识协议(预期共识)，保证一个公平的过程确定块生成在一个回合中的，一个矿工本身是否符合挖掘一个块，和其他规则与有效块的生产所需的一些工件(例如票据，选举证明)。

## Mining Cycle - 挖掘周期

At any height `H`, there are three possible situations:

在任何高度的“H”，它们都有三种可能的情况：

- The miner is eligible to mine a block: they produce their block and propagate it. They then resume mining at the next height `H+1`.
- 矿工有资格开采一个区块：他们生产自己的区块并传播它。然后在下一个高度`H+1`继续开采。
- The miner is not eligible to mine a block but has received blocks: they form a Tipset with them and resume mining at the next height `H+1`.
- 矿工没有资格开采一个块，但收到了块：他们形成一个Tipset与他们的下一个高度`H+1`恢复开采。
- The miner is not eligible to mine a block and has received no blocks: prompted by their clock they run leader election again, incrementing the epoch number.
- 矿工没有资格开采一个区块，也没有收到任何区块：在他们的时钟推动下，他们再次进行领导人选举，并增加纪元数。

This process is repeated until either a winning ticket is found (and block published) or a new valid Tipset comes in from the network.

此过程重复进行直到找到一个赢了的票据(和发布了区块)或从网络中传入新的有效Tipset。

Let's illustrate this with an example.

让我们用一个例子来说明这一点。

Miner M is mining at Height H.
Heaviest tipset at H-1 is {B0}

矿工M在高度H处采矿。
H-1处最重的tipset是{B0}

- New Round:
- 新的一轮：
    - M produces a ticket at H, from B0's ticket (the min ticket at H-1)
    - M在H处生成票据，来自B0的票据(在H-1处的最小票据)
    - M draws the ticket from height H-K to generate an ElectionProof
    - M从高度H-K处取出票据去生成选举证明
    - That ElectionProof is invalid
    - 那个选举证明是无效的
    - M has not heard about other blocks on the network.
    - M没有听说过网络上的其他区块。
- New Round:
- 新的一轮：
    - Epoch/Height is incremented to H + 1.
    - 纪元/高度增加到`H + 1`。
    - M generates a new ElectionProof with this new epoch number.
    - M用这个新的纪元数生成一个新的选举证明。
    - That ElectionProof is valid
    - 那个选举证明是有效的
    - M generates a block B1 using the new ElectionProof and the ticket drawn last round.
    - M使用新的选举证明和最后一轮抽到的票据生成B1块。
    - M has received blocks B2, B3 from the network with the same parents and same height.
    - M已经从具有相同父级和相同高度的网络中接收到了块B2, B3。
    - M forms a tipset {B1, B2, B3}
    - M形成一个tipset {B1, B2, B3}

Anytime a miner receives new blocks, it should evaluate what is the heaviest Tipset it knows about and mine atop it.

无论何时矿工收到新的块，它应该评估哪个是最重的Tipset，它知道关于和在它之上的挖矿。

## Block Creation - 区块创建

Scratching a winning ticket, and armed with a valid `ElectionProof`, a miner can now publish a new block!

抓一张赢的票据，并配备一个有效的`选举证明`，矿工现在可以发布一个新的区块了!

To create a block, the eligible miner must compute a few fields:

要创建一个块，符合条件的矿工必须计算几个字段:

- `Ticket` - new ticket generated from that in the prior epoch (see {{<sref ticket_generation>}}).
- `票据` - 在前一个纪元中生成的新票据(参见{{<sref ticket_generation "票据生成">}})。
- `ElectionProof` - A specific signature over the min_ticket from `randomness_lokkback` epochs back (see {{<sref leader_election>}}).
- `选举证明` - 从`randomness_lokkback`纪元返回的最小票据上的特定签名(参见{{<sref leader_election "领导人选举">}}})。
- `ParentWeight` - The parent chain's weight (see {{<sref chain_selection>}}).
- `父级权重` - 父链的权值(参见{{<sref chain_selection "链选择">}})。
- `Parents` - the CIDs of the parent blocks.
- `父级` - 父级区块CID们
- `ParentState` - Note that it will not end up in the newly generated block, but is necessary to compute to generate other fields. To compute this:
- `父状态` - 注意，它不会出现在新生成的块中，但是对于计算生成其他字段是必需的。计算:
  - Take the `ParentState` of one of the blocks in the chosen parent set (invariant: this is the same value for all blocks in a given parent set).
  - 以所选父集合中的一个块的`父状态`为例(不变式:对于给定父集合中的所有块该值都是相同的)。
  - For each block in the parent set, ordered by their tickets:
  - 对于父集合中的每个块，根据它们的票据排序:
    - Apply each message in the block to the parent state, in order. If a message was already applied in a previous block, skip it.
    - 将块中的每个消息应用到父状态，按顺序的。如果一个消息已经在前一个块中应用，那么跳过它。
    - Transaction fees are given to the miner of the block that the first occurance of the message is included in. If there are two blocks in the parent set, and they both contain the exact same set of messages, the second one will receive no fees.
    - 交易费被给予对应区块的矿工, 那含有该消息的第一次的区块事件。如果在父集合中有两个块，并且它们都包含完全相同的消息集合，则第二个块将不收取任何费用。
    - It is valid for messages in two different blocks of the parent set to conflict, that is, A conflicting message from the combined set of messages will always error.  Regardless of conflicts all messages are applied to the state.
    - 对于父集的两个不同块中的消息冲突是有效的，也就是说，一个来自合并消息集的冲突消息将总是错误的。不管如何冲突所有消息都被应用于状态树。
    - TODO: define message conflicts in the state-machine doc, and link to it from here
    - 计划：在状态机文档中定义消息冲突，并从这里链接到它
- `MsgRoot` - To compute this:
- `消息根` - 来计算这个：
  - Select a set of messages from the mempool to include in the block.
  - 从内存池中选择一个要包含在块中的消息的集合。
  - Separate the messages into BLS signed messages and secpk signed messages
  - 将消息分为BLS签名消息和secpk签名消息
  - For the BLS messages:
  - 有关BLS消息：
    - Strip the signatures off of the messages, and insert all the bare `Message`s for them into a sharray.
    - 去掉消息上的签名，然后将所有的空`消息`插入到sharray中。
    - Aggregate all of the bls signatures into a single signature and use this to fill out the `BLSAggregate` field
    - 将所有bls签名聚合到一个签名中，并使用该签名填写`BLSAggregate`字段
  - For the secpk messages:
  - 有关secpk消息：
    - Insert each of the secpk `SignedMessage`s into a sharray
    - 将每个secpk的`已签消息`插入到sharray中
  - Create a `TxMeta` object and fill each of its fields as follows:
  - 创建一个`TxMeta`对象并填写其每个字段如下:
    - `blsMessages`: the root cid of the bls messages sharray
    - `secpkMessages`: the root cid of the secp messages sharray
  - The cid of this `TxMeta` object should be used to fill the `MsgRoot` field of the block header.
  - 应该使用这个`TxMeta`对象的cid来填充块头的`MsgRoot`字段。
- `BLSAggregate` - The aggregated signatures of all messages in the block that used BLS signing.
- `BLS总计` - 使用BLS签名的块中所有消息的聚合签名
- `StateRoot` - Apply each chosen message to the `ParentState` to get this.
- `状态根` - 将每个选择的消息应用到`ParentState`来得到这个。
  - Note: first apply bls messages in the order that they appear in the blsMsgs sharray, then apply secpk messages in the order that they appear in the secpkMessages sharray.
  - 注意：首先按在blsMsgs sharray中出现的顺序应用bls消息，然后按在secpkMessages sharray中出现的顺序应用secpk消息。
- `ReceiptsRoot` - To compute this:
- `收据根` - 来计算这个:
  - Apply the set of messages to the parent state as described above, collecting invocation receipts as this happens.
  - 如上所述，将消息集应用到父状态，并在此过程中收集调用收据。
  - Insert them into a sharray and take its root.
  - 将它们插入到sharray中并取其根
- `Timestamp` - A Unix Timestamp generated at block creation. We use an unsigned integer to represent a UTC timestamp (in seconds). The Timestamp in the newly created block must satisfy the following conditions:
- `时间戳` - 块创建时生成的Unix时间戳。我们使用无符号整数来表示UTC时间戳(以秒为单位)。新创建的块中的时间戳必须满足以下条件:
  - the timestamp on the block corresponds to the current epoch (it is neither in the past nor in the future) as defined by the clock subsystem.
  - 块上的时间戳对应于时钟子系统定义的当前纪元(既不在过去也不在将来)。
- `BlockSig` - A signature with the miner's private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after it propagates to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.
- `块签名` - 在整个块上用矿工的私人密钥签名(也必须匹配的相关票据签名)。这是为了确保在块传播到网络后没有人篡改它，因为与普通的PoW区块链不同，一个赢了的票据存在于区块生成。

An eligible miner can start by filling out `Parents`, `Tickets` and `ElectionProof` with values from the ticket checking process.
- 符合条件的矿工可以从填写`父级`、`票据`和`选举证明`开始，填上验票过程中的值。

Next, they compute the aggregate state of their selected parent blocks, the `ParentState`. This is done by taking the aggregate parent state of the blocks' parent Tipset, sorting the parent blocks by their tickets, and applying each message in each block to that state. Any message whose nonce is already used (duplicate message) in an earlier block should be skipped (application of this message should fail anyway). Note that re-applied messages may result in different receipts than they produced in their original blocks, an open question is how to represent the receipt trie of this tipset's messages (one can think of a tipset as a 'virtual block' of sorts).

接下来，他们计算所选父块的聚合状态，即`父状态`。这是通过获取块的父Tipset的聚合父状态，根据它们的票据对父块进行排序，并将每个块中的每个消息应用到该状态来实现的。在前面的块中已经使用了nonce的任何消息(重复消息)都应该被跳过(该消息的应用程序应该失败)。请注意，重新应用的消息可能会产生与原始块中不同的收据，一个有待解决的问题是如何表示这个tipset“虚拟块”的收据。有关消息执行和状态转换的详细信息，请参阅[Filecoin状态机](state-machine.md)文档。

Once the miner has the aggregate `ParentState`, they must apply the block reward. This is done by adding the correct block reward amount to the miner owner's account balance in the state tree. The reward will be spendable immediately in this block.

一旦矿工有了聚合的`父级状态`，他们必须应用块奖励。这是通过在状态树中将正确的块奖励合计到矿工所有者的帐户余额来实现的。奖励将会在这个区块内被立即可消费的。

Now, a set of messages is selected to put into the block. For each message, the miner subtracts `msg.GasPrice * msg.GasLimit` from the sender's account balance, returning a fatal processing error if the sender does not have enough funds (this message should not be included in the chain).

现在，将选择一组消息放入块中。对于每条消息，矿工从发送者的帐户余额减去`Gas价格乘以Gas限制值`，如果发送者没有足够的资金(这个消息不应该包括在链)返回一个致命的处理错误。

They then apply the messages state transition, and generate a receipt for it containing the total gas actually used by the execution, the executions exit code, and the return value . Then, they refund the sender in the amount of `(msg.GasLimit - GasUsed) * msg.GasPrice`. In the event of a message processing error, the remaining gas is refunded to the user, and all other state changes are reverted. (Note: this is a divergence from the way things are done in Ethereum)

然后，它们应用消息状态转换，并为其生成一个包含执行实际使用的总gas、执行退出代码和返回值的票据。然后，他们退还发送者的金额`(msg.GasLimit - GasUsed) * msg.GasPrice`。在消息处理出错的情况下，剩余的gas将退还给用户，所有其他状态更改将恢复。(注:这与以太坊做的方式不同)

Each message should be applied on the resultant state of the previous message execution, unless that message execution failed, in which case all state changes caused by that message are thrown out. The final state tree after this process will be the block's `StateRoot`.

每个消息都应该应用于前一个消息执行的结果状态，除非该消息执行失败，在这种情况下，由该消息引起的所有状态更改都将被抛出。此过程之后的最终状态树将是块的`状态根`。

The miner merklizes the set of messages selected, and put the root in `MsgRoot`. They gather the receipts from each execution into a set, merklize them, and put that root in `ReceiptsRoot`. Finally, they set the `StateRoot` field with the resultant state.

矿工merklizes选择的消息集，并把根放在`消息根`中。他们把每次执行的收据收集到一个集合，用merklize把它们分类，然后把这个放进`收据根`。最后，他们用结果状态设置`状态根`字段。

{{% notice info %}}
Note that the `ParentState` field from the expected consensus document is left out, this is to help minimize the size of the block header. The parent state for any given parent set should be computed by the client and cached locally.

请注意，预期共识文档中的`父状态`字段被省略了，这是为了帮助最小化区块头的大小。任何给定父集的父状态都应该由客户机计算并在本地缓存。
{{% /notice %}}

Finally, the miner can generate a Unix Timestamp to add to their block, to show that the block generation was appropriately delayed.

最后，矿工可以生成一个Unix时间戳来添加到他们的块中，以显示块的生成被适当地延迟了。

The miner will wait until BLOCK_DELAY has passed since the latest block in the parent set was generated to timestamp and send out their block. We recommend using NTP or another clock synchronization protocol to ensure that the timestamp is correctly generated (lest the block be rejected). While this timestamp does not provide a hard proof that the block was delayed (we rely on the VDF in the ticket-chain to do so), it provides some softer form of block delay by ensuring that honest miners will reject undelayed blocks.

当生成父集合中的最新块以进行时间戳并发送它们的块时，矿工将等待BLOCK_DELAY结束。我们建议使用NTP或其他时钟同步协议来确保正确生成时间戳(以免拒绝块)。虽然这个时间戳并不能提供块被延迟的确凿证据(我们依赖于票据链中的VDF来做这件事)，但是它通过确保诚实的采矿者拒绝未延迟的块，提供了一种更温和的块延迟形式。

Now the block is complete, all that's left is to sign it. The miner serializes the block now (without the signature field), takes the sha256 hash of it, and signs that hash. They place the resultant signature in the `BlockSig` field.

现在区块完成了，剩下的就是签名了。矿工现在序列化块(没有签名字段)，获取它的sha256散列，并对该散列进行签名。他们将生成的签名放在`BlockSig`字段中。

## Block Broadcast - 区块广播

An eligible miner broadcasts the completed block to the network and assuming everything was done correctly, the network will accept it and other miners will mine on top of it, earning the miner a block reward!

一个合格的采矿者将完成的块广播到网络(通过[区块传播](data-propagation.md))，并且假设一切都做对了，网络将接受它，其他采矿者将在它上面采矿，为采矿者赢得块奖励!

# Block Rewards - 区块奖励

Over the entire lifetime of the protocol, 1,400,000,000 FIL (`TotalIssuance`) will be given out to miners. The rate at which the funds are given out is set to halve every six years, smoothly (not a fixed jump like in Bitcoin). These funds are initially held by the network account actor, and are transferred to miners in blocks that they mine. Over time, the reward will eventually become close zero as the fractional amount given out at each step shrinks the network account's balance to 0.

在该协议书的整个生命周期内，将向矿工发放14亿FIL(`总发行量`)。这些资金的发放速度将每6年平稳减半(不像比特币那样出现固定的跃升)。这些资金最初由网络账户角色持有，并向矿工的开采出的区块转移。随着时间的推移，奖励最终将接近于零，因为每一步所提供的小金额将网络帐户的余额缩小到0。

The equation for the current block reward is of the form:

当前区块奖励的公式如下:

```
Reward = (IV * RemainingInNetworkActor) / TotalIssuance

奖励 = （IV * 网络账户剩余量) / 总发行量
```

`IV` is the initial value, and is set to:

`IV` 是初始化的值，它设定为：

```
IV = 153856861913558700202 attoFIL // 153.85 FIL
```

IV was derived from:

IV源自:
```
// Given one block every 30 seconds, this is how many blocks are in six years
HalvingPeriodBlocks = 6 * 365 * 24 * 60 * 2 = 6,307,200 blocks
λ = ln(2) / HalvingPeriodBlocks
IV = TotalIssuance * (1-e^(-λ)) // Converted to attoFIL (10e18)
```

Note: Due to jitter in EC, and the gregorian calendar, there may be some error in the issuance schedule over time. This is expected to be small enough that it's not worth correcting for. Additionally, since the payout mechanism is transferring from the network account to the miner, there is no risk of minting *too much* FIL.

注:由于EC的抖动和公历的原因，随着时间的推移，可能会在发行日程上出现一些错误。这个值应该足够小，不值得修正。此外，由于支付机制是从网络账户转移到矿工，因此不会产生`过多`FIL的风险。

TODO: Ensure that if a miner earns a block reward while undercollateralized, then `min(blockReward, requiredCollateral-availableBalance)` is garnished (transfered to the miner actor instead of the owner).

要做的事情:确保如果一个矿工赚了区块奖励而抵押不足，然后`最小(块奖励，要求抵押品可用性余额)”被加装饰(转移到矿工角色而不是持有者)。

## Notes on Block Reward Application - 区块奖励申请说明

As mentioned above, every round, a miner checks to see if they have been selected as the leader for that particular round. Thus, it is possible that multiple miners may be selected as winners in a given round, and thus, that there will be multiple blocks with the same parents that are produced at the same block height (forming a Tipset). Each of the winning miners will apply the block reward directly to their actor's state in their state tree.

如上所述，每一轮，矿工都会检查他们是否被选为这一轮的领导者。因此，在给定的一轮中，有可能选择多个矿工作为获胜者，因此，将有多个具有相同父级的块在相同的块高度生成(形成一个Tipset)。每个获胜的矿工都将直接将区块奖励应用到他们状态树中的参与者状态。

Other nodes will receive these blocks and form a Tipset out of the eligible blocks (those that have the same parents and are at the same block height). These nodes will then validate the Tipset. To validate Tipset state, the validating node will, for each block in the Tipset, first apply the block reward value directly to the mining node's account and then apply the messages contained in the block.

其他节点将接收这些块并从符合条件的块(具有相同父级且具有相同块高度的块)中形成一个Tipset。然后，这些节点将验证Tipset。要验证Tipset状态，验证节点将对Tipset中的每个块首先将块奖励值直接应用到挖掘节点的帐户，然后应用块中包含的消息。

Thus, each of the miners who produced a block in the Tipset will receive a block reward. There will be no lockup. These rewards can be spent immediately.

因此，每个在Tipset中生成一个块的矿工都会得到一个块奖励。不会有禁闭。这些奖励可以立即使用。

Messages in Filecoin also have an associated transaction fee (based on the gas costs of executing the message). In the case where multiple winning miners included the same message in their blocks, only the first miner will be paid this transaction fee. The first miner is the miner with the lowest ticket value (sorted lexicographically).

Filecoin中的消息也有相关的交易费用(基于执行消息的gas成本)。如果多个成功的矿工在他们的区块中包含相同的信息，只有第一个矿商将获得这笔交易费。第一个矿工是票据值最低的矿工(按字典顺序排序)。

# Open Questions - 开放的问题

- How should receipts for tipsets be referenced? It is common for applications to provide the merkleproof of a receipt to prove that a transaction was successfully executed.
- 如何引用Tipset的收据？应用程序通常提供收据的merkleproof来证明交易成功执行。

