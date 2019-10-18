---
title: Relayer Node - 中继节点
---

```
type RelayerNode interface {
  FilecoinNode

  blockchain.MessagePool
  markets.MarketOrderBook
}
```
