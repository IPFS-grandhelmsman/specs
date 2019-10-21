---
menuTitle: IpldStore - Ipld存储
title: "IpldStore - Local Storage for hash-linked data - IpldStore——哈希链接数据的本地存储"
---

{{< readfile file="../../../../libraries/ipld/ipld.id" code="true" lang="go" >}}

TODO:

- What is IPLD
- 什么是IPLD
  - hash linked data
  - 哈希关联数据
  - from IPFS
  - 来自IPFS
- Why is it relevant to filecoin
- 为什么与filecoin有关
  - all network datastructures are definitively IPLD
  - 所有网络数据结构都是绝对的IPLD
  - all local datastructures can be IPLD
  - 所有的本地数据结构都可以是IPLD
- What is an IpldStore
- 什么是IpldStore
  - local storage of dags
  - dags的本地储存
- How to use IpldStores in filecoin
- 如何在filecoin中使用IpldStores
  - pass it around
  - 传递下去
- One ipldstore or many
- 一个或多个ipldstore
  - temporary caches
  - 临时缓存
  - intermediately computed state
  - 中间计算状态
- Garbage Collection
- 垃圾收集
