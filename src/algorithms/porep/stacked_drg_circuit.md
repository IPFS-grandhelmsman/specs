## Stacked DRG: Offline PoRep Circuit Spec - 推叠DRG：离线复制证明电路规范

### Stacked DRG Overview

Stacked DRG PoRep is based on layering DRG graphs `LAYERS` times. The data represented in each DRG layer is a labeling based on previously labeled nodes. The final labeled layer is the SDR key, and the 'final layer' of replication the replica, an encoding of the original data using the generated key.

堆叠DRG复制证明是基于DRG图形的`层`时间分层。每个DRG层中表示的数据是基于先前标记的节点的标记。最后一个标记层是SDR密钥，复制副本的“最后一层”使用生成的密钥对原始数据进行编码。

- `ReplicaId` is a unique replica identifier (see the Filecoin Proofs spec for details).
- `副本Id` 是唯一的副本标识符(详细信息请参阅Filecoin证明规范)。
- `CommD` is the Merkle tree root hash of the input data to the first layer.
- `CommD` 是向第一层输入数据的Merkle树根散列。
- `CommC` is the Merkle tree root hash of the SDR column commitments.
- `CommC` 是SDR列承诺的Merkle树根散列。
- `CommRLast` is the Merkle tree root hash of the replica.
- `CommRLast` 是Merkle树的根哈希的副本。
- `CommR` is the on-chain commitment to the replica, dervied as the hash of the concatenation of `CommC` and `CommRLast`.
- `CommR` 是对副本的链上承诺，是`CommC`和`CommRLast`连接的散列。

The (offline) proof size in SDR is too large for blockchain usage (~3MB). We use SNARKs to generate a proof of knowledge of a correct SDR proof. In other words, we implement the SDR proof verification algorithm in an arithmetic circuit and use SNARKs to prove that it was evaluated correctly.

SDR中(脱机)的证明大小对于区块链使用来说太大了(~3MB)。我们使用SNARK生成一个关于正确SDR证明的知识证明。换句话说，我们在一个算术电路中实现了SDR证明验证算法，并使用SNARK来证明其被正确评估。

This circuit proves that given a Merkle root `CommD`, `CommRLast`, and `commRStar`, that the prover knew the correct replicated data at each layer.

该电路证明，给定一个Merkle根`CommD`，`CommRLast`，`commRStar`，证明者知道正确的复制数据在每一层。

### Spec notation - 规范符号

- **Fr**: Field element of BLS12-381
- **Fr**: BLS12-381的字段元素
- **UInt**: Unsigned integer
- **UInt**: 无符号整型
- **{0..x}**: From `0` (included) to `x` (not included) (e.g. `[0,x)` )
- **{0..x}**: 从`0`(含)到`x`(不含)(即`[0,x)`)
- **Check**:
- **检查**:
  - If there is an equality, create a constraint
  - 如果存在等式，则创建约束
  - otherwise, execute the function
  - 否则，执行函数
- **Inclusion path**: Binary representation of the Merkle tree path that must be proven packed into a single `Fr` element.
- **包含路径**: 二进制表示的Merkle树路径，必须证明压缩成一个单一的`Fr`元素。

# Offline PoRep circuit 离线复制证明电路

## Public Parameters - 公开的参数

*Parameters that are embeded in the circuits or used to generate the circuit*

*嵌入电路或用于产生电路的参数*

- `LAYERS : UInt`: Number of DRG layers.
- `LAYER_CHALLENGES : [LAYERS]UInt`: Number of challenges per layer.
- `EXPANSION_DEGREE: UInt`: Degree of each bipartite expander graph to extend dependencies between layers.
- `BASE_DEGREE: UInt`: Degree of each Depth Robust Graph.
- `TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is (log_2(Size of original data in bytes/32 bytes per leaf)).
- `PARENT_COUNT : UInt`: Defined as `EXPANSION_DEGREE+BASE_DEGREE`.

## Public Inputs - 公开的输入

*Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

*验证者用于生成SNARK证明和验证者用于验证它的输入*

- `ReplicaId : Fr`: A unique identifier for the replica.
- `CommD : Fr`: the Merkle tree root hash of the original data (input to the first layer).
- `CommR : Fr`: The Merkle tree root hash of the final replica (output of the last layer).
- `InclusionPath : [LAYERS][]Fr`: Inclusion path for the challenged data and replica leaf.
- `ParentInclusionPath : [LAYERS][][PARENT_COUNT]Fr`:  Inclusion path for the parents of the corresponding `InclusionPath[l][c]`.

Design notes:

- `CommRLast` is a private input used during during Proof-of-Spacetime.
   To enable this, the prover must store `CommC` and use it to prove that `CommRLast` is included in `CommR` [TODO: define 'included' language.]
- `InclusionPath` and `ParentInclusionPath`: Each layer `l` has `LAYER_CHALLENGES[l]` inclusion paths.

## Private Inputs - 私有输入

*Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

*验证者用于生成SNARK证明的输入，验证者不需要这些输入来验证证明*

- `CommR : [LAYERS-1]Fr`: Commitment of the the encoded data at each layer.

  Note: Size is `LAYERS-1` since the commitment to the last layer is `CommRLast`

- `DataProof : [LAYERS][][TREE_DEPTH]Fr`: Merkle tree inclusion proof for the current layer unencoded challenged leaf.

- `ReplicaProof : [LAYERS][][TREE_DEPTH]Fr`: Merkle tree inclusion proof for the current layer encoded challenged leaves.

- `ParentProof : [LAYERS][][PARENT_COUNT][TREE_DEPTH]Fr`: Pedersen hashes of the Merkle inclusion proofs of the parent leaves for each challenged leaf at layer `l`.

- `DataValue : [LAYERS][]Fr`: Value of the unencoded challenged leaves at layer `l`.

- `ReplicaValue : [LAYERS][]Fr`: Value of the encoded leaves for each challenged leaf at layer `l`.

- `ParentValue : [LAYERS][][PARENT_COUNT]Fr`: Value of the parent leaves for each challenged leaf at layer `l`.

## Circuit - 电路

In high level, we do 4 checks:

1. **ReplicaId Check**: Check the binary representation of the ReplicaId
2. **Inclusion Proofs Checks**: Check the inclusion proofs
3. **Encoding Checks**: Check that the data has been correctly encoding into a replica
4. **CommRStar Check**: Check that CommRStar has been generated correctly

Detailed

```go
// 1: ReplicaId Check - Check ReplicaId is equal to its bit representation
let ReplicaIdBits : [255]Fr = Fr_to_bits(ReplicaId)
assert(Packed(replica_id_bits) == ReplicaId)

let DataRoot, ReplicaRoot Fr

for l in range LAYERS {

  if l == 0 {
    DataRoot = CommD
  } else {
    DataRoot = CommR[l-1]
  }

  if l == LAYERS-1 {
    ReplicaRoot = CommRLast
  } else {
    ReplicaRoot = CommR[l]
  }

  for c in range LAYERS_CHALLENGES[l] {
    // 2: Inclusion Proofs Checks
    // 2.1: Check inclusion proofs for data leaves are correct
    assert(MerkleTreeVerify(DataRoot, InclusionPath[l][c], DataProof[l][c], DataValue[l][c]))
    // 2.2: Check inclusion proofs for replica leaves are correct
    assert(MerkleTreeVerify(ReplicaRoot, InclusionPath[l][c], ReplicaProof[l][c], ReplicaValue[l][c]))
    // 2.3: Check inclusion proofs for parent leaves are correct
    for p in range PARENT_COUNT {
      assert(MerkleTreeVerify(ReplicaRoot, ParentInclusionPath[l][c][p], ParentProof[l][c][p]))
    }

    // 3: Encoding checks - Check that replica leaves have been correctly encoded
    let ParentBits [PARENT_COUNT][255]Fr
    for p in range PARENT_COUNT {
      // 3.1: Check that each ParentValue is equal to its bit representation
      let parent = ParentValue[l][c][p]
      ParentBits[p] = Fr_to_bits(parent)
      assert(Packed(ParentBits[p]) == parent)
    }

    // 3.2: KDF check - Check that each key has generated correctly
    // PreImage = ReplicaIdBits || ParentBits[1] .. ParentBits[PARENT_NODES]
    let PreImage = ReplicaIdBits
    for parentbits in ParentBits {
      PreImage.Append(parentbits)
    }
    let key Fr = SHA256(PreImage)
    assert(SHA256(PreImage) == key)

    // 3.3: Check that the data has been encoded to a replica with the right key
    assert(ReplicaValue[l][c] == DataValue[l][c] + key)
  }

  // 4: CommRStar check - Check that the CommRStar constructed correctly
  let hash = ReplicaId
  for l in range LAYERS-1 {
    hash.Append(CommR[l])
  }
  hash.Append(CommRLast)

  assert(CommRStar == PedersenHash(hash))
  // TODO check if we need to do packing/unpacking
}
```



## Verification of offline porep proof - 离线存储证明的证明校验

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
- Parent checks: For each `leaf = InclusionPath[l][c]`:
  - **Check** that all `ParentsInclusionPaths_[l][c][0..PARENT_COUNT}` are the correct parent leaves of `leaf` in the DRG graph, if a leaf has less than `PARENT_COUNT`, repeat the leaf with the highest label in the graph.
  - **Check** that the parent leaves are in ascending numerical order.

