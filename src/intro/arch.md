---
title: "Architecture Diagrams - 架构图"
---


# Filecoin Systems - Filecoin系统

<script type="text/javascript">

function statusIndicatorsShow() {
  var $uls = document.querySelectorAll('.statusIcon')
  $uls.forEach(function (el) {
    el.classList.remove('hidden')
  })
  return false; // stop click event
}

function statusIndicatorsHide() {
  var $uls = document.querySelectorAll('.statusIcon')
  $uls.forEach(function (el) {
    el.classList.add('hidden')
  })
  return false; // stop click event
}

</script>


Status Legend:

状态说明：

- 🛑 **Bare** - Very incomplete at this time.
- 🛑 **空的** - 这个时候非常不完整。
  - **Implementors:** This is far from ready for you.
  - **实现者:** 这还远没为你准备好。
- ⚠️ **Rough** -- work in progress, heavy changes coming, as we put in place key functionality.
- ⚠️ **粗糙的** - 工作正在进行，重大的变化即将到来，因为我们把关键的功能
  - **Implementors:** This will be ready for you soon.
  - **实现者:** 这即将为你准备好。
- 🔁 **Refining** - Key functionality is there, some small things expected to change. Some big things may change.
- 🔁 **精炼中** - 关键的功能就在那里，一些小事情预期会改变。一些大的事情可能会改变。
  - **Implementors:** Almost ready for you. You can start building these parts, but beware there may be changes still.
  - **实现者:** 几乎为你准备好了。您可以开始构建这些部分，但要注意可能仍有更改。
- ✅ **Stable** - Mostly complete, minor things expected to change, no major changes expected.
- ✅ **稳定的** - 大部分完成，小的事情预计会改变，没有大的变化预计。
  - **Implementors:** Ready for you. You can build these parts.
  - **实现者:** 为你准备好了。您可以构建这部分了。

[<a href="#" onclick="return statusIndicatorsShow();">Show</a> / <a href="#" onclick="return statusIndicatorsHide();">Hide</a> ] status indicators


{{< incTocMap "/docs/systems" 2 "colorful" >}}


# Overview Diagram

TODO:

- cleanup / reorganize
  - this diagram is accurate, and helps lots to navigate, but it's still a bit confusing
  - the arrows and lines make it a bit hard to follow. We should have a much cleaner version (maybe based on [C4](https://c4model.com))
- reflect addition of Token system
  - move data_transfers into Token

{{< diagram src="../diagrams/overview1/overview.dot.svg" title="Protocol Overview Diagram" >}}


# Protocol Flow Diagram -- deals off chain

{{< diagram src="../diagrams/sequence/full-deals-off-chain.mmd.svg" title="Protocol Sequence Diagram - Deals off Chain" >}}

# Protocol Flow Diagram -- deals on chain

{{< diagram src="../diagrams/sequence/full-deals-on-chain.mmd.svg" title="Protocol Sequence Diagram - Deals on Chain" >}}

# Parameter Calculation Dependency Graph

This is a diagram of the model for parameter calculation. This is made with [orient](https://github.com/filecoin-project/orient), our tool for modeling and solving for constraints.

{{< diagram src="../diagrams/orient/filecoin.dot.svg" title="Parameter Calculation Dependency Graph" >}}

