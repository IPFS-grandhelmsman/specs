---
title: System Actors - 系统角色
statusIcon: 🔁
entries:
- init_actor
- cron_actor
- account_actor
---

- There are two system actors required for VM processing:
- VM处理需要两个系统参与者：
  - [CronActor](#CronActor) - runs critical functions at every epoch
  - [定时器角色](#systems__filecoin_vm__sysactors__cron_actor) - 在每个时期运行关键函数
  - [InitActor](#InitActor) - initializes new actors
  - [初始化角色](#systems__filecoin_vm__sysactors__init_actor) - 初始化新的角色
- There is one more VM level actor:
- 还有一个VM级的角色：
  - [AccountActor](#AccountActor) - for user accounts.
  - [帐号角色](#systems__filecoin_vm__sysactors__account_actor) - 用于用户帐号(管理).
