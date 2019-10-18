---
title: Retrieval Miner Node - 检索矿工节点
---

```
type RetrievalMinerNode interface {
  FilecoinNode

  blockchain.Blockchain
  markets.RetrievalMarketProvider
  markets.MarketOrderBook
  markets.DataTransfers
}
```
