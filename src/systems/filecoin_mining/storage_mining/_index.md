---
menuTitle: Storage Miner - 存储矿工
statusIcon: 🔁
title: Storage Miner - 存储矿工
entries:
  - mining_cycle
  - storage_miner_actor
  - mining_scheduler
---

{{<label storage_mining_subsystem>}}

TODO:

- rename "Storage Mining Worker" ?
- 重命名为 "存储挖掘工人" ?

Filecoin Storage Mining Subsystem

Filecoin存储挖掘子系统

{{< readfile file="storage_mining_subsystem.id" code="true" lang="go" >}}

# Sector in StorageMiner State Machine (new one) - 存储矿工状态机中的扇区(新建)

{{< diagram src="diagrams/sector_state_fsm.dot.svg" title="Sector State (new one)" >}}

{{< diagram src="diagrams/sector_state_legend.dot.svg" title="Sector State Legend (new one)" >}}

# Sector in StorageMiner State Machine (both) - 存储矿工状态机中的扇区(两者)

{{< diagram src="diagrams/sector_fsm.dot.svg" title="Sector State Machine (both)" >}}
