---
menuTitle: Piece - 数据片
statusIcon: 🔁
title: Piece - a part of a file - 数据片 - 文件的一部分
entries:
- piece_store
---


A `Piece` is an object that represents a whole or part of a `File`,
and is used by `Clients` and `Miners` in `Deals`. `Clients` hire `Miners`
to store `Pieces`.

一个`数据片`是一个对象，代表一个`文件`的整体或部分，
`客户`和`矿工`在`交易`中使用。`客户`雇佣`矿工`
存储`数据片`。

{{< readfile file="piece.id" code="true" lang="go" >}}
