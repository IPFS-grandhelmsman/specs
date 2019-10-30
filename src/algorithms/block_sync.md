---
title: "BlockSync - 块同步"
---

- **Name**: Block Sync
- **名称**: 块同步
- **Protocol ID**: `/fil/sync/blk/0.0.1`
- **协议ID**: `/fil/sync/blk/0.0.1`

The blocksync protocol is a small protocol that allows Filecoin nodes to request ranges of blocks from each other. It is a simple request/response protocol.

块同步协议是一个小协议，它允许Filecoin节点相互请求块的范围。它是一个简单的请求/响应协议。

The request requests a chain of a given length by the hash of its highest block. The `Options` allow the requester to specify whether or not blocks and messages to be included.

请求通过其最高块的散列请求给定长度的链。“选项”允许请求者指定是否包含块和消息。

The response contains the requested chain in reverse iteration order. Each item in the `Chain` array contains the blocks for that tipset if the `Blocks` option bit in the request was set, and if the `Messages` bit was set, the messages across all blocks in that tipset. The `MsgIncludes` array contains one array of integers for each block in the `Blocks` array. Each of the arrays in `MsgIncludes` contains a list of indexes of messages from the `Messages` array that are in each `Block` in the blocks array.

响应以反迭代顺序包含请求的链。如果设置了请求中的`区块`选项位，则`链`数组中的每一项都包含该tipset的块;如果设置了`消息`位，则该tipset中所有块的消息都包含在该`消息`数组中。`MsgIncludes`数组为`区块`数组中的每个块包含一个整数数组。`MsgIncludes`中的每个数组都包含一个来自`消息`数组的消息索引列表，这些`消息`数组位于区块数组中的每个`区块`中。

```sh
type BlockSyncRequest struct {
    ## The TipSet being synced from
	start [&Block]
    ## How many tipsets to sync
	requestLength UInt
    ## Query options
    options Options
}
```

```sh
type Options enum {
    # Include only blocks
    | Blocks 0
    # Include only messages
    | Messages 1
    # Include messages and blocks
    | BlocksAndMessages 2
}

type BlockSyncResponse struct {
	chain [TipSetBundle]
	status Status
}

type TipSetBundle struct {
  blocks [Blocks]
  secpMsgs [SignedMessage]
  secpMsgIncludes [[UInt]]

  blsMsgs [Message]
  blsMsgIncludes [[Uint]]
}

type Status enum {
    ## All is well.
    | Success 0
    ## Sent back fewer blocks than requested.
    | PartialResponse 101
    ## Request.Start not found.
    | BlockNotFound 201
    ## Requester is making too many requests.
    | GoAway 202
    ## Internal error occured.
    | InternalError 203
    ## Request was bad
    | BadRequest 204
}
```

## Example - 例子

The TipSetBundle
TipSet包

```
Blocks: [b0, b1]
secpMsgs: [mA, mB, mC, mD]
secpMsgIncludes: [[0, 1, 3], [1, 2, 0]]
```

corresponds to:
对应于：

```
Block 'b0': [mA, mB, mD]
Block 'b1': [mB, mC, mA]
```
