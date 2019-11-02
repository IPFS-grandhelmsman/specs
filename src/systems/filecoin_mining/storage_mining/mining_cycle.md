---
title: Storage Mining Cycle - 存储挖掘周期
---

Block miners should constantly be performing Proofs of SpaceTime, and also checking if they have a winning `ticket` to propose a block at each height/in each round. Rounds are currently set to take around 30 seconds, in order to account for network propagation around the world. The details of both processes are defined here.

区块矿工应该不断地证明时空，并检查他们是否有一张抽中的`票据`，在每个高度/在每个回合中提出区块。为了考虑到全球范围内的网络传播，当前的轮数设置为30秒左右。这两个过程的详细信息在这里定义

# The Miner Actor - 矿工角色

After successfully calling `CreateStorageMiner`, a miner actor will be created on-chain, and registered in the storage market. This miner, like all other Filecoin State Machine actors, has a fixed set of methods that can be used to interact with or control it.
在成功调用`创建存储矿工`之后，一个矿工角色将被创建在链上，并注册在存储市场中。与所有其他Filecoin状态机角色一样，这个miner有一组固定的方法，可以用来与之交互或控制它。

{{< readfile file="storage_miner_actor.id" code="true" lang="go" >}}
{{< readfile file="storage_miner_actor.go" code="true" lang="go" >}}

## Owner Worker distinction - 拥有者与工人区别

The miner actor has two distinct 'controller' addresses. One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key.

矿工角色有两个不同的“控制器”地址。一个是工人，这个地址将负责做所有的工作，提交证明，承诺新的扇区，和所有其他日常活动。拥有者地址是创建矿工，支付了抵押品，并有块奖励支付给它的地址。这种区别的原因是允许不同的角色发挥不同的作用。一个例子是，所有者是一个多签名钱包，或一个冷存储密钥，而工人密钥是一个“热钱包”密钥。

## Storage Mining Cycle - 存储挖掘周期

Storage miners must continually produce Proofs of SpaceTime over their storage to convince the network that they are actually storing the sectors that they have committed to. Each PoSt covers a miner's entire storage.

存储矿工必须不断地为他们的存储产生时空证明，以使网络相信他们承诺的实际存储扇区。每个存储证明都覆盖矿工的全部存储空间。

### Step 0: Registration - 第0步：注册

To initially become a miner, a miner first register a new miner actor on-chain. This is done through the storage power actor's `CreateStorageMiner` method. The call will then create a new miner actor instance and return its address.

要初始化成为一个矿工，矿工首先在链上注册一个新的矿工角色。这是通过存储市场角色的[`创建存储矿工`](actors.md#createstorageminer)方法完成的。这个调用然后将创建一个新的矿工角色实例并返回它的地址。

The next step is to place one or more storage market asks on the market. This is done off-chain as part of storagee market functions. A miner may create a single ask for their entire storage, or partition their storage up in some way with multiple asks (at potentially different prices).

下一步是在市场上放置一个或多个存储市场请求。这是作为存储者市场功能的一部分进行的。一个矿工可能会为他们的整个存储创建一个单独的请求，或者以某种方式用多个询价(价格可能不同)来分割他们的存储。

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the {{<sref storage_market>}}.

之后，他们需要与客户进行交易，并开始用数据填充扇区。有关交易的更多信息，请参见{{<sref storage_market "存储市场">}}。

When they have a full sector, they should seal it. This is done by invoking the {{<sref sector_sealer>}}.

当他们有一个完整的扇区，他们应该密封它。这是通过在扇区上调用{{<sref sector_sealer "扇区密封器">}}来完成的。

#### Changing Worker Addresses - 变更矿工地址

Note that any change to worker keys after registration (TODO: spec how this works) must be appropriately delayed in relation to randomness lookback for SEALing data (see [this issue](https://github.com/filecoin-project/specs/issues/415)).

请注意，在注册后对工作密钥的任何更改(TODO: 规范这个如何工作)都必须适当延迟，以便对密封数据进行随机性回退(参见[这个讨论](https://github.com/filecoin-project/specs/issues/415))。

### Step 1: Commit - 第1步：提交

When the miner has completed their first seal, they should post it on-chain using the {{<sref storage_miner_actor>}}'s `ProveCommitSector` function. If the miner had zero committed sectors prior to this call, this begins their proving period.

当矿工完成他们的第一个密封时，他们应该使用[提交扇区](actors.md#commitsector)把它到发送链上。如果矿工在此之前没有提交的扇区，这将开始了它们的证明周期

The proving period is a fixed amount of time in which the miner must submit a Proof of Space Time to the network.

证明周期是一个固定的时间，在这个时间内，矿工必须向网络提交一份时空证明。

During this period, the miner may also commit to new sectors, but they will not be included in proofs of space time until the next proving period starts.
For example, if a miner currently PoSts for 10 sectors, and commits to 20 more sectors. The next PoSt they submit (i.e. the one they're currently proving) will be for 10 sectors again, the subsequent one will be for 30.

在此期间，矿工也可以承诺新的扇区，但在下一个验证周期开始之前，这些扇区不会被包括在时空证明中。
例如，如果一个矿工当前的为10个扇区运行时空证明，并承诺20多个扇区。他们提交的下一个时空证明(即他们目前正在证明那个)将再次涉及10个扇区，随后的下一个将是30个。

TODO: sectors need to be globally unique. This can be done either by having the seal proof prove the sector is unique to this miner in some way, or by having a giant global map on-chain is checked against on each submission. As the system moves towards sector aggregation, the latter option will become unworkable, so more thought needs to go into how that proof statement could work.

待做:扇区需要全局唯一。这可以通过以下两种方式来实现:用密封证明验证该扇区在某种方式上是唯一的，或者是在链上有一个区大的全局表对每一个提交进行检查。随着系统向扇区聚合方向发展，后一种选择将变得不可行，因此需要更多地考虑证明表述如何工作。

### Step 2: Proving Storage (PoSt creation) - 第2步：验证存储 (时空证明创建)

```go
func ProveStorage(sectorSize BytesAmount, sectors []commR) PoStProof {
    challengeBlockHeight := miner.ProvingPeriodEnd - POST_CHALLENGE_TIME

    // Faults to be used are the currentFaultSet for the miner.
    faults := miner.currentFaultSet
    seed := GetRandFromBlock(challengeBlockHeight)
    return GeneratePoSt(sectorSize, sectors, seed, faults)
}
```

Note: See ['Proof of Space Time'](proof-of-spacetime.md) for more details.

注：有关详细信息，请参见[“时空证明”](proof-of-spacetime.md)。

The proving set remains consistent during the proving period. Any sectors added in the meantime will be included in the next proving set, at the beginning of the next proving period.

在证明期间，证明集保持一致。在此期间增加的扇区将被包括在下一个证明集中，在下一个证明期开始时。

### Step 3: PoSt Submission - 第3步：时空证明提交

When the miner has completed their PoSt, they must submit it to the network by calling [SubmitPoSt](actors.md#submitpost). There are two different times that this *could* be done.

当矿工完成他们的时空证明后，他们必须通过调用[提交时空证明](actors.md#submitpost)将其提交给网络。有两种不同的时间*可以*这样做。

1. **Standard Submission**: A standard submission is one that makes it on-chain before the end of the proving period. The length of time it takes to compute the PoSts is set such that there is a grace period between then and the actual end of the proving period, so that the effects of network congestion on typical miner actions is minimized.
1. **标准提交**：一个标准的提交是一个使它在证明期结束前上链。计算这些时空证明所花费的时间是这样设置的:从那时到证明周期的实际结束之间有一段宽限期，从而使网络拥塞对典型矿工操作的影响最小化。
2. **Penalized Submission**: A penalized submission is one that makes it on-chain after the end of the proving period, but before the generation attack threshold. These submissions count as valid PoSt submissions, but the miner must pay a penalty for their late submission. (See '[Faults](faults.md)' for more information)
2. **惩罚提交**：一个惩罚的提交是在证明周期结束后上链，但要在生成攻击阈值之前。这些提交都是有效的时空证明，但是矿工必须为他们的逾期提交支付罚款。(参考'[故障](faults.md)'获得更多信息)
   - Note: In this case, the next PoSt should still be started at the beginning of the proving period, even if the current one is not yet complete. Miners must submit one PoSt per proving period.
   - 注意:在这种情况下，下一个时空证明仍应在证明周期开始时开始，即使当前的这个尚未完成。矿商必须为每一个证明周期提交一个时空证明。

Along with the PoSt submission, miners may also submit a set of sectors that they wish to remove from their proving set. This is done by selecting the sectors in the 'done' bitfield passed to `SubmitPoSt`.

在时空证明提交的同时，矿工也可以提交一组他们希望从他们的证明集中移除的扇区。这是通过在“完成”位字段中传递给`提交时空证明`中选择扇区来完成的。

# Stop Mining - 停止挖掘

In order to stop mining, a miner must complete all of its storage contracts, and remove them from their proving set during a PoSt submission. A miner may then call [`DePledge()`](actors.md#depledge) to retrieve their collateral. `DePledge` must be called twice, once to start the cooldown, and once again after the cooldown to reclaim the funds. The cooldown period is to allow clients whose files have been dropped by a miner to slash them before they get their money back and get away with it.

为了停止采矿，采矿者必须完成所有的存储契约，并在时空证明提交后将它们从证明集中移除。然后，矿工可以调用[`DePledge()`](actors.md#depledge)来取回他们的抵押品。`DePledge`必须调用两次，一次启动冷却，一次在冷却后收回资金。冷却期是在矿工拿回他们的钱并带它离开之前允许已被矿工丢弃的文件的客户去削减他们的文件。

# Faults 故障

Faults are described in the [faults document](faults.md).

故障被描述在[故障文档](faults.md)中.

## On Being Slashed (WIP, needs discussion) - 关于被削减(WIP，需要讨论)

If a miner is slashed for failing to submit their PoSt on time, they currently lose all their pledge collateral. They do not necessarily lose their storage collateral. Storage collateral is lost when a miner's clients slash them for no longer having the data. Missing a PoSt does not necessarily imply that a miner no longer has the data. There should be an additional timeout here where the miner can submit a PoSt, along with 'refilling' their pledge collateral. If a miner does this, they can continue mining, their mining power will be reinstated, and clients can be assured that their data is still there.

如果一名矿工因在按时提交他们的时空证明失败时而被削减，他们目前将失去所有的抵押品。它们不一定会失去它们的存储抵押品。存储抵押品在矿工的客户不再要数据而削减了它们时被遗失。缺失一个时空证明并不一定意味着矿工不再拥有数据。这里应该有一个额外的超时时间，矿工可以提交一个带“再填充”他们的抵押品一起的时空证明。如果矿工这样做，他们可以继续挖掘，他们的挖掘能力将恢复，客户可以确保他们的数据仍然在那里。

TODO: disambiguate the two collaterals across the entire spec

待做:在整个规范中消除两个抵押品的歧义

Review Discussion Note: Taking all of a miners collateral for going over the deadline for PoSt submission is really really painful, and is likely to dissuade people from even mining filecoin in the first place (If my internet going out could cause me to lose a very large amount of money, that leads to some pretty hard decisions around profitability). One potential strategy could be to only penalize miners for the amount of sectors they could have generated in that timeframe.

回顾讨论记录:接管所有的矿工抵押品后提交的最后期限是真的痛苦,甚至可能会阻止人们挖掘filecoin首先(如果我的互联网出去可能会导致我失去了大量的钱,导致一些非常艰难的决定盈利能力)。一种可能的策略是，仅根据矿业公司在这段时间内所能创造的行业数量来惩罚它们。
回顾讨论记录:接管所有的矿工抵押品的超过最后期限的存储证明提交的是真的痛苦,甚至可能首先会阻止人们挖掘Filecoin(如果我的互联网出局可能会导致我失去了大量的钱,引导围绕盈利能力的一些决定非常艰难)。一个可能的策略是只惩罚那些本可以在那个时间段内完成的扇区的数量矿商。

# Future Work - 未来的工作

There are many ideas for improving upon the storage miner, here are ideas that may be potentially implemented in the future.

有许多想法可以改进存储矿工，以下是将来可能实现的一些想法。

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **扇区重密封**：矿工应该能够“重新密封”扇区，去允许他们获取一组已几乎释放了的数据片的扇区，并将尚未过期的数据片合并成一个(或多个)扇区。
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, and will not be implemented in the initial release of Filecoin, but could provide interesting capabilities down the road.
- **扇区转移**：矿工应该能够将存储数据的责任重新委托给另一个矿工。由于许多原因，这很棘手，并且不会在Filecoin的初始版本中实现，但可为未来提供有趣的功能埋下了伏笔。
