---
title: "Expected Consensus - 预期共识"
---

{{<label expected_consensus>}}
## Algorithm - 算法

Expected Consensus (EC) is a probabilistic Byzantine fault-tolerant consensus protocol. At a high level, it operates by running a leader election every round in which, on expectation, one participant may be eligible to submit a block. EC guarantees that this winner will be anonymous until they reveal themselves by submitting a proof of their election (we call this proof an `Election Proof`). All valid blocks submitted in a given round form a `Tipset`. Every block in a Tipset adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on how to select the heaviest chain, see {{<sref chain_selection>}}.

预期共识(EC)是一种概率拜占庭容错一致性协议。在高层，它的运作方式是每一轮进行领导人选举，在预期中，一名参与者可能有资格提交一份提案。EC保证这个获胜者将是匿名的，直到他们提交他们的选举证明(我们称这个证明为`选举证明`)。所有有效的区块都以`Tipset`的形式提交。Tipset中的每个块都增加了其链的权重。“最好”的链是权重最高的链，也就是说分叉的选择规则是选择已知的最重的链。有关如何选择最重链的更多细节，请参见{{<sref chain_selection "区块选择">}}。

At a very high level, with every new block generated, a miner will craft a new ticket from the prior one in the chain appended with the current epoch number (i.e. parentTipset.epoch + 1 to start). While on expectation at least one block will be generated at every round, in cases where no one finds a block in a given round, a miner will increment the round number and attempt a new leader election (using the new input) thereby ensuring liveness in the protocol.

The {{<sref storage_power_consensus>}} subsystem uses access to EC to use the following facilities:

{{<sref storage_power_consensus "存储能力证明">}}子系统使用对EC的访问来使用以下设施:

- Access to verifiable randomness for the protocol, derived from {{<sref tickets>}}.
- 从{{<sref tickets>}}获取协议的可验证随机性。
- Running and verifying {{<sref leader_election "leader election">}} for block generation.
- 运行并验证{{<sref leader_election "领导人选举">}}用于生成块。
- Access to a weighting function enabling {{<sref chain_selection>}} by the chain manager.
- 通过链管理器访问支持{{}<sref chain_selection "链选择">}的权重函数
- Access to the most recently {{<sref finality "finalized tipset">}} available to all protocol participants.
- 对所有协议参与者可用的最近的{{<sref finality "定型的tipset">}}的访问。

{{<label tickets>}}
## Tickets - 票据

For leader election in EC, participants win in proportion to the power they have within the network.

在EC的领导人选举中，参与者根据他们在网络中拥有的权力比例获胜。

A ticket is drawn from the past at the beginning of each new round to perform leader election. EC also generates a new ticket in every round for future use. Tickets are chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any other data in the block.
On expectation, in Filecoin, every block header contains one ticket, though it could contain more if that block was generated over multiple rounds.

在每一轮领导人选举开始时，都会从过去抽一张票。EC也会在每一轮产生一张新的票以备将来使用。票据独立于主区块链链接。一个票据只依赖于它之前的票据，而不是块中的任何其他数据。
根据预期，在Filecoin中，每个块头都包含一个票据，但是如果该块是在多个回合中生成的，那么它可能包含更多票据。

Tickets are used across the protocol as sources of randomness:

票据作为随机性的来源在整个协议中使用:

- The {{<sref sector_sealer>}} uses tickets to bind sector commitments to a given subchain.
- {{<sref sector_sealer "扇区密封器">}}使用票据将扇区承诺绑定到给定的子链。
- The {{<sref post_generator>}} likewise uses tickets to prove sectors remain committed as of a given block.
- {{<sref post_generator "投递生成器">}}同样使用票据来证明扇区仍然是作为给定块提交的。
- EC uses them to run leader election and generates new ones for use by the protocol, as detailed below.
- EC使用它们来运行领导人选举，并生成新的领导人供协议使用，具体如下

You can find the Ticket data structure {{<sref data_structures "here">}}.

您可以在这里找到Ticket数据结构{{<sref data_structures "这里">}}

### Comparing Tickets in a Tipset - 在Tipset中比较票据

Whenever comparing tickets is evoked in Filecoin, for instance when discussing selecting the "min ticket" in a Tipset, the comparison is that of the little endian representation of the ticket's VFOutput bytes.

每当在Filecoin中调用比较票据时，例如在讨论在Tipset中选择“小票据”时，比较的是票据的VFOutput字节的小端表示。

## Tickets in EC - EC中的票据

Within EC, a miner generates a new ticket in their block for every ticket they use running leader election, thereby ensuring the ticket chain is always as long as the block chain.

在EC中，矿工为他们使用的每一张竞选领导人的票在他们的区块生成一张新票，从而确保票链总是和区块链一样长。

Tickets are used to achieve the following:

票据是用来实现以下功能的:

- Ensure leader secrecy -- meaning a block producer will not be known until they release their block to the network.
-  确保领导人是保密的 -- 意思是一个块的生产者将不会被知道，直到他们释放他们的块到网络。
- Prove leader election -- meaning a block producer can be verified by any participant in the network.
- 证明领导人选举 -- 这意味着块生成器可以被网络中的任何参与者验证。


In practice, EC defines two different fields within a block:

在实践中，EC在一个块中定义了两个不同的字段:

- A `Ticket` field — this stores the new ticket generated during this block generation attempt. It is from this ticket that miners will sample randomness to run leader election in `K` rounds.
- 一个`票据`字段 - 它存储在此块生成尝试期间生成的新票据。从这张票中，矿工们将随机抽样进行领导人选举的`K`轮。
- An `ElectionProof` — this stores a proof that a given miner has won a leader election using the appropriate ticket `K` rounds back appended with the current epoch number. It proves that the leader was validly elected in this epoch.

{{< readfile file="election.id" code="true" lang="go" >}}
{{< readfile file="election.go" code="true" lang="go" >}}

```
But why the randomness lookback?

但是为什么会有随机性呢?

The randomness lookback helps turn independent ticket generation from a block one round back
into a global ticket generation game instead. Rather than having a distinct chance of winning or losing
for each potential fork in a given round, a miner will either win on all or lose on all
forks descended from the block in which the ticket is sampled.

随机性回溯有助于将独立票证生成从第一轮的阻塞转回到全局票证生成游戏。
在给定的一轮中，矿工不是对每个潜在的分支有明显的赢或输的机会，而是对所有分支都有赢或输的机会。

This is useful as it reduces opportunities for grinding, across forks or sybil identities.

这很有用，因为它减少了在fork或sybil标识之间进行研磨的机会。

However this introduces a tradeoff:

然而，这带来了一个权衡:

- The randomness lookback means that a miner can know K rounds in advance that they will win,
decreasing the cost of running a targeted attack (given they have local predictability).
- 随机回眸意味着矿工可以提前知道K轮他们将会赢，降低了目标攻击的成本(考虑到他们的局部可预测性)。
- It means electionProofs are stored separately from new tickets on a block, taking up
more space on-chain.
- 这意味着选举证明与新票分开存储在一个块上，在链上占用更多空间。

How is K selected?

K是如何被选择的?

- On the one end, there is no advantage to picking K larger than finality.
- 一方面，选择K大于最终结果是没有好处的。
- On the other, making K smaller reduces adversarial power to grind.
- 另一方面，使K更小给令人疲劳的工作减少对立力量。
```

{{<label ticket_generation>}}
### Ticket generation- 票据产生

{{< readfile file="expected_consensus.id" code="true" lang="go" >}}
{{< readfile file="expected_consensus.go" code="true" lang="go" >}}

This section discusses how tickets are generated by EC for the `Ticket` field.

本节讨论EC如何为`Ticket`字段生成票。

At round `N`, a new ticket is generated using tickets drawn from the Tipset at round `N-1` (for more on how tickets are drawn see {{<sref ticket_chain>}}).

在第`N`轮，使用从第`N-1`轮的Tipset中提取的票据生成一个新票据(有关如何提取票据的更多信息，请参见{{<sref ticket_chain "票据链">}})。

The miner runs the prior ticket through a Verifiable Random Function (VRF) to get a new unique output.

采矿者通过一个可验证的随机函数(VRF)运行先前的票据，以获得一个新的惟一输出。

The VRF's deterministic output adds entropy to the ticket chain, limiting a miner's ability to alter one block to influence a future ticket (given a miner does not know who will win a given round in advance).

VRF的确定性输出增加了票证链的熵，限制了矿工改变一个区块以影响未来票证的能力(假设矿工不知道在给定的一轮中谁会赢)。

We use the VRF from {{<sref vrf>}} for ticket generation in EC (see the `PrepareNewTicket` method below).

我们使用来自{{<sref vrf>}}的VRF来生成EC中的票据(请参阅下面的`准备新票据`方法)。

{{< readfile file="storage_mining_subsystem.id" code="true" lang="go" >}}
{{< readfile file="storage_mining_subsystem.go" code="true" lang="go" >}}


### Ticket Validation - 票据验证

Each Ticket should be generated from the prior one in the ticket-chain and verified accordingly as shown in `validateTicket` below.

{{< readfile file="storage_power_consensus_subsystem.id" code="true" lang="go" >}}
{{< readfile file="storage_power_consensus_subsystem.go" code="true" lang="go" >}}

每一张票都应该从票据链中的前一张生成。

{{<label leader_election>}}
## Secret Leader Election -  机密领导人选举

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their power. In the case of Filecoin, participants and powers are drawn from the {{<sref power_table>}}, where power is equivalent to storage provided through time.

预期共识是一种共识协议，通过根据矿工的权力按比例加权选举产生。在Filecoin的情况下，参与者和权力都来自于存储[能力表](storage-market.md#the-power-table)，其中的权力等同于通过时间提供的存储。

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of randomness used to run the election. In the case of Filecoin's EC, the blockchain tracks an independent ticket chain. These tickets are used as randomness inputs for Leader Election. Every block generated references an `ElectionProof` derived from a past ticket. The ticket chain is extended by the miner who generates a new block for each successful leader election.

预期一致的领导人选举必须是秘密的、公平的和可核实的。这是通过使用随机用于运行选举。在Filecoin的EC中，区块链跟踪一个独立的票据链。这些票被用作领导人选举的随机输入。每一个区块都会产生一个从过去的选票中衍生出来的`选举证明`。矿工扩展了票证链，为每一次成功的领导人选举生成一个新的区块。

### Running a leader election - 运行一个领导人选举

Now, a miner must also check whether they are eligible to mine a block in this round. 

现在，矿工还必须检查他们是否有资格在这一轮开采一个区块。

Design goals here include:

这里的设计目标包括:
- There should be one block per miner per epoch at most (for simplicity)
- 每个矿工每个纪元最多应该有一个块(为简单起见)
- Miners should be rewarded proportional to their power in the system
- 矿工的奖励应该与他们在系统中的存力成比例
- The system should be able to tune how many blocks are put out per epoch on expectation (hence "expected consensus").
- 系统应该能够根据预期(因此称为“预期共识”)调整每个纪元输出的块的数量。

现在，矿工还必须检查他们是否有资格在这一轮开采一个区块。

To do so, the miner will use tickets from K rounds back as randomness to uniformly draw a value from 0 to 1. Comparing this value to their power, they determine whether they are eligible to mine. A user's `power` is defined as the ratio of the amount of storage they proved as of their last PoSt submission to the total storage in the network as of the current block.

为了做到这一点，矿工将使用K轮返回的票据作为随机性，均匀地从0到1绘制一个值。将这个值与他们的权力进行比较，他们决定是否有资格开采。用户的`能力`被定义为他们在最后一次提交时所证明的存储容量与当前块在网络中的总存储容量的比值。

We use the VRF from {{<sref vrf>}} to run leader election in EC.

我们使用{{<sref vrf>}}中的VRF来运行EC中的leader election。

If the miner wins the election in this round, it can use newEP, along with a newTicket to generate and publish a new block. Otherwise, it waits to hear of another block generated in this round.

如果矿工在这一轮的选举中获胜，它可以使用newEP和一个newTicket来生成和发布一个新的块。否则，它将等待这轮生成的另一个块。

In short, the process of crafting a new ElectionProof in round N is as follows in the `DrawElectionProof` function:

简而言之，在第N轮中制作新的选举证明的过程如下:

{{< readfile file="storage_mining_subsystem.go" code="true" lang="go" >}}

It is important to note that every block contains two artifacts: one, a ticket derived from last block's ticket to extend the ticket-chain, and two, an election proof derived from the ticket `K` rounds back used to run leader election.

需要注意的是，每个块都包含两个人工制品:一个是从上一个块的ticket派生出的票据，用于扩展票证链;另一个是从用于运行领导人选举的票据`K`轮派生出的选举证明。

Note: Miner power is drawn from the power table, accounting only for power that has been proven over time (see {{<sref power_table>}}).

注意:矿工的能力是从能力表中提取的，只考虑经过一段时间证明的能力(参见{{<sref power_table>}})。

The miner can then check whether they drew a winning election proof by comparing their power fraction to the `ElectionProof`s value, as follows:

然后，矿工可以通过比较他们的权力比例和`选举证明`的价值，来检查他们是否抽到了获胜的选举证明，如下所示:

`ElectionProof * TotalPower < e * MinerPower * MaxValue`

`选举证明 * 总存力 < e * 矿工存力 * 最大值`

with:
和：
- `ElectionProof`   the election proof's byte output
- `选举证明`        选举证明的字节输出
- `TotalPower`      the total power in the current power table
- `总存力`          当前存力表中的总存力
- `e`               the expected number of blocks crafted per epoch
- `e`               每个epoch的块的预期数量
- `MinerPower`      the current miner's power in the power table
- `矿工存力`        当前矿工在存力表的存力
- `MaxValue`        the maximum possible value for the election proof using SHA-256 with our vdf, we would have `maxValue = 2^256`.
- `最大值`          对于与我们的vdf一起使用SHA-256的选举证明的最大可能的值，我们将有`最大值 = 2^256`.

The above is the integer format of `ElectionProof/MaxValue < e * MinerPower/TotalPower` which makes it easy to see that a miner should win proportionally to the amount of power they have staked in the system. Note that miners with more than `1/e` fraction of power will win once at every round, meaning there may be fewer than e blocks per round with many large miners.

注意:我们从上一轮获取矿机能量。这意味着如果一个矿工在他们的验证期结束时赢得了一个区块，即使他们还没有重新提交帖子，他们仍然保留他们的权力(直到下一轮)。

If successful, the miner can craft a block, passing it to the block producer. If unsuccessful, it will wait to hear of another block mined this round to try again. In the case no other block was found in this round the miner can increment the epoch number and try leader election again using the same past ticket and new epoch number.

如果成功，矿工可以创建一个块，并将其传递给块生成器。如果不成功，它将等待听到另一个块被挖出这一轮再次尝试。如果在这一轮中没有发现其他块，矿工可以增加纪元号，并尝试使用相同的过去票据和新的纪元号再次选举领导人。

While a miner could try to run through multiple epochs in parallel in order to quickly generate a block, this effort will be futile as the rational majority of miners will reject blocks crafted with ElectionProofs whose epochs prove too high (i.e. in the future -- see below).

虽然一个矿商可以尝试同时运行多个纪元来快速生成一个块，但是这种努力是徒劳的，因为大多数合理的矿商将会拒绝那些纪元过高的选举证明(即在未来——见下文)。

### Election Validation - 选举校验

In order to determine that the mined block was generated by an eligible miner, one must check its `ElectionProof`'s validity and that its input was generated using the current epoch value as shown in the `ValidateElectionProof` and `IsWinningElectionProof` methods above.

为了确定被挖掘的块是由一个合格的矿工生成的，必须检查它的`选举证明`的有效性，并且它的输入是使用当前纪元值生成的，如上面的`有效的选举证明`和`赢了的选举证明`方法所示。

It is worth drawing attention to how mining rewards are split among miners. Specifically, for every miner publishing a block in a given round, reward is determined as follows:

值得注意的是，矿工奖励是如何在矿工之间分配的。具体来说，每一个矿工在给定的回合中发布一个区块，奖励如下:

`
draw = ElectionProof/maxValue
req = e * minerPower/totalPower
rewCount = ceil(req - draw)
reward = rewCount * targetRewardPerEpoch/e
`

`
抽取 = 选举证明/最大值
请求 = e * 矿工存力/总存力
奖励数 = 向上取整(请求 - 抽取)
奖励 = 奖励数 * 每纪元目标奖励/e
`

This process, shown in the `GetBlockRewards` method above, allows all miners to receive their fair share of rewards on expectation (over time), accounting for some miners winning at every round. Read more about design tradeoffs [here](https://github.com/filecoin-project/specs/issues/603).

为了确定被挖掘的区块是由一个合格的矿工所生成的，必须检查它的`选举证明`。

{{<label chain_selection>}}
## Chain Selection - 链选择

Just as there can be 0 miners win in a round, multiple miners can be elected in a given round. This in turn means multiple blocks can be created in a round, as seen above. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round.

就像可以有0个矿工在一轮中获胜一样，在给定的一轮中多个矿工可以被选举。这意味着在一轮中可以创建多个块。为了避免浪费矿工完成的有效工作，EC使用了一轮生成的所有有效块。

### Chain Weighting - 链权重

It is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the most storage provided.

分叉有可能在预期的共识中自然出现。EC依赖于加权链，以快速收敛于“一个真正的链”，每个区块增加链的权重。这意味着最重的链应该反映执行的工作量最多，或者在Filecoin的用例中，反映提供的存储空间最多。

In short, the weight at each block is equal to its `ParentWeight` plus that block's delta weight. Details of Filecoin's chain weighting function [are included here](https://observablehq.com/d/3812cd65c054082d).
简而言之，每个块的权值等于它的`父权重`加上该块的权值。Filecoin的链加权函数的详细信息[包含在这里](https://observablehq.com/d/3812cd65c054082d)。

Delta weight is a term composed of a few elements:

增量(Δ)权重是由以下几个元素组成的术语：

- wForkFactor: which seeks to cut the weight derived from rounds in which produced Tipsets do not correspond to what an honest chain is likely to have yielded (pointing to selfish mining or other non-collaborative miner behavior).
- wForkFactor: 这种方法旨在减少从发牌过程中产生的权重，而发牌过程中产生的权重与诚实的链条可能产生的权重不相符(指向自私的采矿或其他非合作的矿工行为)。
- wPowerFactor: which adds weight to the chain proportional to the total power backing the chain, i.e. accounted for in the chain's power table.
- wPowerFactor: 它使链的权重与支撑链的总能力成比例，即在链的能力表中占一定比例。
- wBlocksFactor: which adds weight to the chain proportional to the number of blocks mined in a given round. Like wForkFactor, it rewards miner cooperation (which will yield more blocks per round on expectation).
- wBlocksFactor: 这增加了权重的链在一个给定的回合成比例的块开采。就像wForkFactor一样，它也会奖励矿工合作(这将在预期的每轮中收获更多的块)。

The weight should be calculated using big integer arithmetic with order of operations defined above. We use brackets instead of parentheses below for legibility. We have:

权重应该使用上面定义的操作顺序的大整数算法来计算。为了便于阅读，我们在下面使用方括号而不是圆括号。我们有:

`w[r+1] = w[r] + (wPowerFactor[r+1] + wBlocksFactor[r+1]) * 2^8`

For a given tipset `ts` in round `r+1`, we define:

对于一个给定tipset `ts`的第`r+1`轮中，我们定义：

- `wPowerFactor[r+1]  = wFunction(totalPowerAtTipset(ts))`
- wBlocksFactor[r+1] =  `wPowerFactor[r+1] * wRatio * b / e`
  - with `b = |blocksInTipset(ts)|`
  - `e = expected number of blocks per round in the protocol`
  - and `wRatio in ]0, 1[`
Thus, for stability of weight across implementations, we take:
- wBlocksFactor[r+1] =  `(wPowerFactor[r+1] * b * wRatio_num) / (e * wRatio_den)`

We get:

我们得到：

- `w[r+1] = w[r] + wFunction(totalPowerAtTipset(ts)) * 2^8 + (wFunction(totalPowerAtTipset(ts)) * len(ts.blocks) * wRatio_num * 2^8) / (e * wRatio_den)`
 Using the 2^8 here to prevent precision loss ahead of the division in the wBlocksFactor.
 使用这里的2^8来防止wBlocksFactor中的除法之前的精度损失。

 The exact value for these parameters remain to be determined, but for testing purposes, you may use:
 这些参数的准确值仍有待确定，但出于测试目的，您可以使用:
 - `e = 5`
 - `wRatio = .5, or wRatio_num = 1, wRatio_den = 2`
- `wFunction = log2b` with
  - `log2b(X) = floor(log2(x)) = (binary length of X) - 1` and `log2b(0) = 0`. Note that that special case should never be used (given it would mean an empty power table.

```sh
Note that if your implementation does not allow for rounding to the fourth decimal, miners should apply the [tie-breaker below](#selecting-between-tipsets-with-equal-weight). Weight changes will be on the order of single digit numbers on expectation, so this should not have an outsized impact on chain consensus across implementations.
请注意，如果您的实现不允许四舍五入到第四个小数，则矿工应该应用[下面的tie-breaker](#selecting-between-tipset-with-equal-weight)。权重的变化将按照预期的个位数的顺序进行，因此这对跨实现链的共识不会有太大的影响。
```

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given Tipset should have
the same `ParentWeight` value) plus the delta weight of each parent. To make the
computation a bit easier, a block's `ParentWeight` is stored in the block itself (otherwise
potentially long chain scans would be required to compute a given block's weight).

`ParentWeight`是一个给定块的父块集合的总链权值。它的计算方法是:任意一个父块的`ParentWeight`(给定Tipset中的所有块都应该有相同的`ParentWeight`值)加上每个父块的增量权值。为了简化计算，块的`ParentWeight`存储在块本身中(否则可能需要长链扫描来计算给定块的权重)。

### Selecting between Tipsets with equal weight - 在具有相同权重的Tipset之间进行选择

When selecting between Tipsets of equal weight, a miner chooses the one with the smallest final ticket.

当在重量相等的Tipset之间进行选择时，矿工将选择最小的那个最终票据。

In the case where two Tipsets of equal weight have the same min ticket, the miner will compare the next smallest ticket (and select the Tipset with the next smaller ticket). This continues until one Tipset is selected.

如果两个相同权重的Tipset具有相同的最小票据，则矿工将比较下一个最小的票据(并在下一个更小的票据中选择Tipset)。这将一直进行下去，直到选择了一个Tipset。

The above case may happen in situations under certain block propagation conditions. Assume three blocks B, C, and D have been mined (by miners 1, 2, and 3 respectively) off of block A, with minTicket(B) < minTicket(C) < minTicket (D).

上述情况可能发生在某些块传播条件下。假设B、C和D三个区块(分别由矿工1、2和3开采)已经从A区块挖出，且minTicket(B) < minTicket(C) < minTicket(D)。

Miner 1 outputs their block B and shuts down. Miners 2 and 3 both receive B but not each others' blocks. We have miner 2 mining a Tipset made of B and C and miner 3 mining a Tipset made of B and D. If both succesfully mine blocks now, other miners in the network will receive new blocks built off of Tipsets with equal weight and the same smallest ticket (that of block B). They should select the block mined atop [B, C] since minTicket(C) < minTicket(D).

矿工1输出它们的区块B并关闭。矿工2和3都得到了B，但没有得到相互间的区块。我们有矿工2从B和C开挖一个Tipset，矿工3从B和D开挖一个Tipset .如果两者现在成功挖块,网络中的其他矿工将收到新的块以同样的最小重量和相同的最小的票据(块B)由Tipset建出，因为minTicket(C) < minTicket(D)，他们应该选择在[B, C]之上的块挖掘。

The probability that two Tipsets with different blocks would have all the same tickets can be considered negligible: this would amount to finding a collision between two 256-bit (or more) collision-resistant hashes.

具有不同块的两个tipset拥有所有相同票据的概率可以忽略不计:这相当于在两个256位(或更多)抗碰撞散列之间找到一个碰撞。

{{<label finality>}}
## Finality in EC - EC中的终结
EC enforces a version of soft finality whereby all miners at round N will reject all blocks that fork off prior to round N-F. For illustrative purposes, we can take F to be 500. While strictly speaking EC is a probabilistically final protocol, choosing such an F simplifies miner implementations and enforces a macroeconomically-enforced finality at no cost to liveness in the chain.

EC实施了一个软终结的版本，在第N轮的所有矿工将拒绝在第N-F轮之前分叉的所有区块。为了说明目的，我们可以将F设为500。严格地说，EC是一个可能的最终协议，选择这样的F简化了矿工实现，并在不牺牲链中的活动性的情况下实现了宏观经济强制的最终性。

{{<label consensus_faults>}}
## Consensus Faults - 共识故障

Due to the existence of potential forks in EC, a miner can try to unduly influence protocol fairness. This means they may choose to disregard the protocol in order to gain an advantage over the power they should normally get from their storage on the network. A miner should be slashed if they are provably deviating from the honest protocol.

由于EC中潜在分叉的存在，矿工可以试图过度影响协议公平性。这意味着他们可能会选择忽略协议，在网络上获得比正常从存储中获得的能力更大的优势。如果矿工确实偏离了诚实的原则，他们应该被砍掉。

This is detectable when a given miner submits two blocks that satisfy any of the following "consensus faults":

当给定的矿工提交满足以下任何“共识错误”的两个区块时，这是可以检测到的：

- (1) `double-fork mining fault`: two blocks mined at the same epoch.
- (1) `双分叉挖矿故障`: 两个区块在同一纪元开采。
{{< diagram src="diagrams/double_fork.dot.svg" title="Double-Fork Mining Fault" >}}

- (2) `time-offset mining fault`: two blocks mined off of the same Tipset at different epochs (i.e. with different `ElectionProof`s generated from the same input ticket).
- (2) `时间重置挖矿故障`: 在不同的纪元，从相同的Tipset中挖掘出两个区块(即从相同的输入票据中生成不同的`选举证明`)。
{{< diagram src="diagrams/time_offset.dot.svg" title="Time-Offset Mining Fault" >}}

- (3) `parent-grinding fault`: one block's parent is a Tipset that provably should have included a given block but does not. While it cannot be proven that a missing block was willfully omitted in general (i.e. network latency could simply mean the miner did not receive a particular block), it can when a miner has successfully mined a block two epochs in a row and omitted one. That is, this condition should be evoked when a miner omits their own prior block. When a miner's block at epoch e + 1 references a Tipset that does not include the block they mined at e both blocks can be submitted to prove this fault.
- (3) `父级研磨故障`:一个块的父级是一个可以证明应该包含给定块但没有包含的Tipset。虽然一般情况下不能证明丢失的块是故意缺失的(即网络延迟可能只是意味着矿工没有收到特定的块)，但是当矿工成功地连续两个纪元挖掘一个块并遗漏一个纪元时，它可以被证明。也就是说，当一个矿工遗失自己之前的块时，应该触发这个条件。当矿工在纪元 e + 1的块引用了一个不包含他们在e的块的Tipset时，可以提交来证明这个错误。
{{< diagram src="diagrams/parent_grinding.dot.svg" title="Parent-Grinding fault" >}}

Any node that detects any of the above events should submit both block headers to the `StoragePowerActor`'s `ReportConsensusFault` method. The "slasher" will receive a portion (TODO: define how much) of the offending miner's {{<sref pledge_collateral>}} as a reward for notifying the network of the fault.
(TODO: FIP of submitting commitments to block headers to prevent miners censoring slashers in order to gain rewards).
任何检测到上述任何事件的节点都应该将两个区块头提交到`存储能力角色`的`reportaccepsusfault`方法。`削减者`将收到违规矿工的{{<sref pledge_collateral "抵押担保">}}的一部分(TODO:定义多少)作为通知网络故障的奖励。
(TODO: FIP的提交对区块头的承诺来防止矿工删减削减者以获得奖励)。

It is important to note that there exists a third type of consensus fault directly reported by the `CronActor` on `StorageDeal` failures via the `ReportUncommittedPowerFault` method:

需要注意的是，`定时角色`通过`报告未承诺存力故障`方法直接报告的`存储密封`故障存在第三种类型的共识故障:

- (4) `uncommitted power fault` which occurs when a miner fails to submit their `PostProof` and is thus participating in leader election with undue power (see {{<sref storage_faults>}}).
- (4) `未承诺存力故障`当一个矿工没有提交他们的`时空证明的证明`，从而参与了不正当存力的领导人选举(见{{<sref storage_faults "存储故障">}})。

