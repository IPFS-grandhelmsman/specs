---
menuTitle: FileStore - 文件存储
title: "FileStore - Local Storage for Files - 文件存储 - 文件的本地存储"
---

The `FileStore` is an abstraction used to refer to any underlying system or device
that Filecoin will store its data to. It is based on Unix filesystem semantics, and
includes the notion of `Paths`. This abstraction is here in order to make sure Filecoin
implementations make it easy for end-users to replace the underlying storage system with
whatever suits their needs. The simplest version of `FileStore` is just the host operating
system's file system.

`文件存储`是一个抽象概念，用于引用任何Filecoin将其数据存储到的底层系统或设备。
它基于Unix文件系统语义，并且包括`路径`的概念。
这里的抽象是为了确保Filecoin实现使得最终用户可以很容易地替换成任何适合他们需要的底层存储系统。
`文件存储`最简单的版本就是主机操作系统的文件系统。

{{< readfile file="filestore.id" code="true" lang="go" >}}

# Varying user needs - 不同的用户需求

Filecoin user needs vary significantly, and many users -- especially miners -- will implement
complex storage architectures underneath and around Filecoin. The `FileStore` abstraction is here
to make it easy for these varying needs to be easy to satisfy. All file and sector local data
storage in the Filecoin Protocol is defined in terms of this `FileStore` interface, which makes
it easy for implementations to make swappable, and for end-users to swap out with their system
of choice.

Filecoin用户的需求差异很大，许多用户——尤其是矿工——将在Filecoin下面和周围实现复杂的存储体系结构。
这里的`文件存储`抽象使这些变化的需求易于满足。
Filecoin协议中的所有文件和扇区本地数据存储都是根据这个`文件存储`接口定义的，这使得实现可以很容易地进行可切换，并使最终用户能够与他们选择的系统进行交换。

# Implementation examples - 实现示例

The `FileStore` interface may be implemented by many kinds of backing data storage systems. For example:

`文件存储`接口可以由多种备份数据存储系统实现。例如:

- The host Operating System file system
- 主机操作系统文件系统
- Any Unix/Posix file system
- 任何Unix/Posix文件系统
- RAID-backed file systems
- RAID-backed文件系统
- Networked of distributed file systems (NFS, HDFS, etc)
- 分布式文件系统的网络(NFS、HDFS等)
- IPFS
- Databases
- NAS systems
- NAS系统
- Raw serial or block devices
- 原始串行或块设备
- Raw hard drives (hdd sectors, etc)
- 原始硬盘驱动器(硬盘驱动器扇区等)

Implementations SHOULD implement support for the host OS file system.

实现应该实现对主机OS文件系统的支持。

Implementations MAY implement support for other storage systems.

实现可以实现对其他存储系统的支持。


