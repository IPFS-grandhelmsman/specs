---
title: Storage Miner Node(存储矿工节点)
---

```
type StorageMinerNode interface {
  FilecoinNode

  systems.Blockchain
  systems.Mining
  markets.StorageMarketProvider
  markets.MarketOrderBook
  markets.DataTransfers
}
```
