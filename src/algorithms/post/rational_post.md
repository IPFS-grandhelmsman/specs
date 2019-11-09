---
title: Rational-PoSt - 时空证明原理
---

This document describes Rational-PoSt, the Proof-of-Spacetime used in Filecoin.

本文档描述了在Filecoin中使用的证明时间的方法。

# High Level API - 高级API

## Fault Detection - 故障检测

Fault detection happens over the course of the life time of a sector. When the sector is for some reason unavailable, the miner is responsible to submit the known `faults`, before the PoSt challenge begins. (Using the `AddFaults` message to the chain).
Only faults which have been reported at challenge time, will be accounted for. If any other faults have occured the miner can not submit a valid PoSt for this proving period.

故障检测发生在扇区的整个生命周期中。当该扇区是由于某种原因不可用时，矿工负责提交已知的“故障”，在时空证明考查开始之前。(使用`增加故障`消息到链)。
只有在考查时报告的错误才会被考虑。如果出现任何其他故障，矿工不能为这个证明时期提交有效的时空证明的。

The PoSt generation then takes the latest available `faults` of the miner to generate a PoSt matching the committed sectors and faults.

然后，这个时空证明生成获取矿工最新可用的“故障”，以生成与提交扇区和故障匹配的时空证明。

When a PoSt is successfully submitted all faults are reset and assumed to be recovered. A miner must either (1) resolve a faulty sector and accept challenges against it in the next proof submission, (2) report a sector faulty again if it persists but is eventually recoverable, (3) report a sector faulty *and done* if the fault cannot be recovered.

当一个时空证明被成功提交时，所有的错误都被重置并被认为已经恢复。矿商必须(1)解决一个故障扇区，并在下一次提交证明时接受对它的质疑;(2)如果扇区故障持续存在，但最终可以恢复，则必须再次报告扇区故障;(3)如果不能恢复，则必须报告扇区故障*并完成*。

If the miner knows that the sectors are permanently lost, they can submit them as part of the `doneSet`, to ensure they are removed from the proving set.

如果矿工知道这些扇区永久性丢失，他们可以将其作为`完成集`的一部分提交，以确保这些扇区被从证明集中移除。

{{% notice note %}}
**Note**: It is important that all faults are known (i.e submitted to the chain) prior to challenge generation, because otherwise it would be possible to know the challenge set, before the actual challenge time. This would allow a miner to report only faults on challenged sectors, with a gurantee that other faulty sectors would not be detected.

**注**:重要的是要知道(即已提交到链)所有的错误。在考查生成之前，因为否则就有可能在实际考查时间之前知道考查集。这将允许矿工只报告受到考查的扇区的故障，并保证其它有故障的扇区不会被发现。
{{% /notice %}}


{{% notice todo %}}
**TODO**: The penalization for faults is not clear yet.

**TODO**: 对故障的处罚还不明确。
{{% /notice %}}

## Fault Penalization - 故障惩罚

Each reported fault carries a penality with it.

每个报告了错误的都有相应的惩罚。

{{% notice todo %}}
**TODO**: Define the exact penality structure for this.

**TODO**: 定义精确的惩罚结构。
{{% /notice %}}

## Generation - 生成

`GeneratePoSt` generates a __*Proof of Spacetime*__ over all  __*sealed sectors*__ of a single miner— identified by their `commR` commitments. This is accomplished by performing a series of merkle inclusion proofs (__*Proofs of Retrievability*__). Each proof is of a challenged node in a challenged sector. The challenges are generated pseudo-randomly, based on the provided `seed`. At each time step, a number of __*Proofs of Retrievability*__ are performed.

`GeneratePoSt`生成了一个单一矿商所有 __*密封扇区*__ 的 __*时空证明*__ ——由他们的“commR”承诺来确定。这是通过执行一系列merkle包含证明来实现的(__*可收回性证明*__)。每一个证明都是被考查扇区中的一个被考查节点的证明。这些考查是伪随机生成的，基于提供的`种子`。在每个时间步骤，若干个__*可收回性证明*__被执行。

```go
// Generate a new PoSt.
func GeneratePoSt(sectorSize BytesAmount, sectors SectorSet, seed Seed, faults FaultSet) PoStProof {
    // Generate the Merkle Inclusion Proofs + Faults

    challenges := DerivePoStChallenges(seed, faults, sectorSize, SortAsc(GetSectorIds(sectors)))
    challengedSectors := []
    inclusionProofs := []

    for i := 0; i < len(challenges); i++ {
        challenge := challenges[i]

        // Leaf index of the selected sector
        inclusionProof, isFault := GenerateMerkleInclusionProof(challenge.Sector, challenge.Leaf)
        if isFault {
            // faulty sector, need to post a fault to the chain and try to recover from it
            return Fatal("Detected late fault")
        }

        inclusionProofs[n] = inclusionProof
        challengedSectors[i] = sectors[challenge.Sector]
    }

    // Generate the snark
    snarkProof := GeneratePoStSnark(sectorSize, challenges, challengedSectors, inclusionProofs)

    return snarkProof
}
```

## Verification - 校验

`VerifyPoSt` is the functional counterpart to `GeneratePoSt`. It takes all of `GeneratePoSt`'s output, along with those of `GeneratePost`'s inputs required to identify the claimed proof. All inputs are required because verification requires sufficient context to determine not only that a proof is valid but also that the proof indeed corresponds to what it purports to prove.

`校验时空证明`是`生成时空证明`的功能对等物。它需要`生成时空证明`的所有输出，以及`生成时空证明`的输入，以识别所要求的证明。所有的输入都是必须的，因为验证不仅需要足够的上下文来确定一个证明是有效的，而且还要确定这个证明确实符合它所要证明的内容。

```go
// Verify a PoSt.
func VerifyPoSt(sectorSize BytesAmount, sectors SectorSet, seed Seed, proof PoStProof, faults FaultSet) bool {
    challenges := DerivePoStChallenges(seed, faults, sectorSize, SortAsc(GetSectorIds(sectors)))
    challengedSectors := []

    // Match up commitments with challenges
    for i := 0; i < len(challenges); i++ {
        challengedSectors[i] = sectors[challenges[i].Sector]
    }

    // Verify snark
    return VerifyPoStSnark(sectorSize, challenges, challengedSectors)
}
```

## Types - 类型

```go
// The random challenge seed, provided by the chain.
Seed [32]byte
```

```go
type Challenge struct {
    Sector SectorID
    Leaf Uint
}
```

## Challenge Derivation - 考查推导

```go
// Derive the full set of challenges for PoSt.
func DerivePoStChallenges(seed Seed, faults FaultSet, sectorSize Uint, sortedSectors []SectorID) [POST_CHALLENGES_COUNT]Challenge {
    challenges := []

    for n := 0; n < POST_CHALLENGES_COUNT; n++ {
        attemptedSectors := {SectorID:bool}
        while challenges[n] == nil {
            challenge := DerivePoStChallenge(seed, n, attempt, sectorSize, sortedSectors)

            // check if we landed in a faulty sector
            if !faults.Contains(challenge.Sector) {
                // Valid challenge
                challenges[n] = challenge
            }

            // invalid challenge, regenerate
            attemptedSectors[challenge.Sector] = true

            if len(attemptedSectors) >= len(sortedSectors) {
                Fatal("All sectors are faulty")
            }
        }
    }

    return challenges
}

// Derive a single challenge for PoSt.
func DerivePoStChallenge(seed Seed, n Uint, attempt Uint, sectorSize Uint, sortedSectors []SectorID) Challenge {
    nBytes := WriteUintToLittleEndian(n)
    data := concat(seed, nBytes, WriteUintToLittleEndian(attempt))
    challengeBytes := blake2b(data)

    sectorChallenge := ReadUintLittleEndian(challengeBytes[0..8])
    leafChallenge := ReadUintLittleEndian(challengeBytes[8..16])

    sectorIdx := sectorChallenge % sectorCount

    return Challenge {
        Sector: sortedSectors[sectorIdx],
        Leaf: leafChallenge % (sectorSize / NODE_SIZE),
    }
}
```


# PoSt Circuit - 时空证明电路

## Public Parameters - 公共参数

*Parameters that are embeded in the circuits or used to generate the circuit*

*嵌入电路或用于生成电路的参数*

- `POST_CHALLENGES_COUNT: UInt`: Number of challenges.
- `时空证明_考查_计数：无符号整型`：考查的数量
- `POST_TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is `(log_2(Size of original data in bytes/32 bytes per leaf))`.
- `时空证明_树_深度:无符号整型`: Merkle树的深度。注意，这是`(log_2(原始数据的大小，以字节为单位/每个叶子32字节)')。
- `SECTOR_SIZE: UInt`: The size of a single sector in bytes.
- `扇区_大小:无符号整型`：单个扇区的大小(以字节为单位)。

## Public Inputs - 公共输入

*Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

*验证者用于生成SNARK证明和验证者用于验证的输入*

- `CommRs: [POST_CHALLENGES_COUNT]Fr`: The Merkle tree root hashes of all replicas, ordered to match the inclusion paths and challenge order.
- `CommRs: [时空证明_考查_计数]Fr`:所有副本的Merkle树根哈希值，为匹配包含路径和考查顺序排序。
- `InclusionPaths: [POST_CHALLENGES_COUNT]Fr`: Inclusion paths for the replica leafs, ordered to match the `CommRs` and challenge order. (Binary packed bools)
- `包容性路径: [时空证明_考查_计数]Fr`:复制叶子的包含路径，为匹配`CommRs`和考查的顺序排序。(二进制bool包装)

## Private Inputs - 私有输入

*Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

*验证者用于生成SNARK证明的输入，验证者不需要这些来验证证明*

- `InclusionProofs: [POST_CHALLENGES_COUNT][TREE_DEPTH]Fr`: Merkle tree inclusion proofs, ordered to match the challenge order.
- 包含证明: [时空证明_考查_计数][树_深度]Fr ': Merkle树包含证明，为匹配考查的顺序排序。
- `InclusionValues: [POST_CHALLENGES_COUNT]Fr`: Value of the encoded leaves for each challenge, ordered to match challenge order.
- `包含值：[时空证明_考查_计数]Fr`: 为每个考查编码的叶子的值，为匹配考查的顺序排序。

## Circuit - 电路

### High Level - 高级

In high level, we do 1 check:
在高级别，我们做1个检查：

1. **Inclusion Proofs Checks**: Check the inclusion proofs
1. **包含性证明检查**：检查包含证明

### Details - 详情

```go
for c in range POST_CHALLENGES_COUNT {
  // Inclusion Proofs Checks
  assert(MerkleTreeVerify(CommRs[c], InclusionPath[c], InclusionProof[c], InclusionValue[c]))
}
```

## Verification of PoSt proof - 校验时空证明的证明

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
- SNARK 证明检查：**检查** 给定SNARK证明和公共输入，SNARK验证输出为true
