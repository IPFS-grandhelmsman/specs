---
title: Sector Sealing - 扇区密封
---

{{< readfile file="sealing.id" code="true" lang="go" >}}

## Drawing randomness for sector commitments - 为扇区承诺绘制随机性

{{<sref ticket_chain "Tickets">}} are used as input to the SEAL above in order to tie Proofs-of-Replication to a given chain, thereby preventing long-range attacks (from another miner in the future trying to reuse SEALs).

{{<sref ticket_chain "票据">}} 被用作上述密封的输入以便将复制证明与给定的链连接起来，从而防止远程攻击(将来可能会有另一个矿工试图重新密封)

The ticket has to be drawn from a finalized block in order to prevent the miner from potential losing storage (in case of a chain reorg) even though their storage is intact.

即使他们的存储是完整的，票据必须从最后确定的块中绘制以防止矿工可能丢失存储(在区块重组的情况下)。

Verification should ensure that the ticket was drawn no farther back than necessary by the miner. We note that tickets can uniquely be associated to a given round in the protocol (lest a hash collision be found), but that the round number is explicited by the miner in `commitSector`.

校验应确保矿工不会把票据绘制得比必要的更远。我们注意到，在协议中，票据可以唯一地与给定的轮关联(以免发现散列冲突)，但是轮数是由`承诺扇区`中的矿工指定的。

We present precisely how ticket selection and verification should work. In the below, we use the following notation:

我们精确地展示了票据选择和验证将如何工作。在下面，我们使用以下符号:

- `F`-- Finality (number of rounds)
- `F`-- 终结 (轮数的号码)
- `X`-- round in which SEALing starts
- `X`-- 开始密封中的回合
- `Z`-- round in which the SEAL appears (in a block)
- `Z`-- 出现密封中的回合(在一个区块里)
- `Y`-- round announced in the SEAL `commitSector` (should be X, but a miner could use any Y <= X), denoted by the ticket selection
- `Y`-- 密封`承诺密封`中的回合声明 (应该X, 但矿工可以使用任意的Y <= X), 票据选择的标志
 - `T`-- estimated time for SEAL, dependent on sector size
 - `T`-- 密封估计时间，取决于扇区大小
 - `G = T + variance`-- necessary flexibility to account for network delay and SEAL-time variance.
 - `G = T + variance`-- 必要的灵活性来解释网络延迟和密封时间的差异。

We expect Filecoin will be able to produce estimates for sector commitment time based on sector sizes, e.g.:
`(estimate, variance) <--- SEALTime(sectors)`
G and T will be selected using these.

我们期望Filecoin能够根据部门的大小来估算部门的承诺时间，例如:
`(estimate, variance) <--- SEALTime(sectors)`
G和T将被选择使用这些

#### Picking a Ticket to Seal - 取一张票据进行密封

When starting to prepare a SEAL in round X, the miner should draw a ticket from X-F with which to compute the SEAL.

当开始准备在X轮的密封时，矿工应该从X-F中抽取一张票据来计算密封。

#### Verifying a Seal's ticket - 校验一张密封的票据

When verifying a SEAL in round Z, a verifier should ensure that the ticket used to generate the SEAL is found in the range of rounds [Z-T-F-G, Z-T-F+G].

当在Z轮验证一个密封时，一个验证者应确保用于生成密封的票据在轮次的这个范围[Z-T-F-G, Z-T-F+G]被找到。

#### In Detail - 细节

```
                               Prover
           ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
          │

          ▼
         X-F ◀───────F────────▶ X ◀──────────T─────────▶ Z
     -G   .  +G                 .                        .
  ───(┌───────┐)───────────────( )──────────────────────( )────────▶
      └───────┘                 '                        '        time
 [Z-T-F-G, Z-T-F+G]
          ▲

          └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
                              Verifier
```

Note that the prover here is submitting a message on chain (i.e. the SEAL). Using an older ticket than necessary to generate the SEAL is something the miner may do to gain more confidence about finality (since we are in a probabilistically final system). However it has a cost in terms of securing the chain in the face of long-range attacks (specifically, by mixing in chain randomness here, we ensure that an attacker going back a month in time to try and create their own chain would have to completely regenerate any and all sectors drawing randomness since to use for their fork's power).

注意，这里的验证者是在链上提交消息(即密封)。使用比必要更旧的票据来生成密封是矿工可能会做的一些事情以获得更多关于最终结果的信心(因为我们是在一个概率性的最终系统中)。虽然它在链的安全条款中有一条有代价的条款来面对远程范围的攻击(具体地说,这里通过混合链的随机性,我们确保一个攻击者最终回退一个月来尝试与创建他们自己的链也必须完全重新生成任何和任何所有扇区的绘制随机性以用于他们分叉的能力)。

We break this down as follows:

我们将其分解如下:

- The miner should draw from `X-F`.
- 这矿工将从`X-F`绘制
- The verifier wants to find what `X-F` should have been (to ensure the miner is not drawing from farther back) even though Y (i.e. the round of the ticket actually used) is an unverifiable value.
- 验证者想要找到`X-F`应该是什么(以确保矿工没有从更远的地方提取)，即使Y(即实际使用该回合的票据d)是一个无法验证的值。
- Thus, the verifier will need to make an inference about what `X-F` is likely to have been based on:
- 因此，验证者需要对`X-F`可能是基于什么做出推断:
  - (known) round in which the message is received (Z)
  - (已知)已接收消息的轮数(Z)
  - (known) finality value (F)
  - (已知)最终值(F)
  - (approximate) SEAL time (T)
  - (近似)密封时间(T)
- Because T is an approximate value, and to account for network delay and variance in SEAL time across miners, the verifier allows for G offset from the assumed value of `X-F`: `Z-T-F`, hence verifying that the ticket is drawn from the range `[Z-T-F-G, Z-T-F+G]`.
- 因为T是一个近近值，为了考虑到矿工之间的网络延迟和密封时间的差异，验证器允许G偏移假定的值`X-F`:`Z-T-F`，从而验证票据是从`[Z-T-F, Z-T-F+G]`范围内绘制的。
