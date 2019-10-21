---
menuTitle: PieceStore 碎片存储
title: PieceStore - storing and indexing pieces - 碎片存储 - 存储和索引碎片
---


A `PieceStore` is an object that can store and retrieve pieces
from some local storage. The `PieceStore` additionally keeps
an index of pieces.

`碎片存储`是一个可以从一些本地存储中进行存储和检索片段的对象。`PieceStore`还保存了块的索引。

{{< readfile file="piece_store.id" code="true" lang="go" >}}
