---
title: Clock(时钟)
statusIcon: ✅
---

{{< readfile file="clock_subsystem.id" code="true" lang="go" >}}
{{< readfile file="clock_subsystem.go" code="true" lang="go" >}}

Filecoin assumes weak clock synchrony amongst participants in the system. That is, the system relies on participants having access to a globally synchronized clock, tolerating bounded delay in honest clock lower than epoch time (more on this in a forthcoming paper).

Filecoin假定系统中的参与者之间存在弱时钟同步。也就是说，系统依赖于能够访问全局同步时钟的参与者，容忍低于epoch时间的有限延迟(在即将发表的论文中有更多相关内容)。

Filecoin relies on this system clock in order to secure consensus, specifically ensuring that participants are only running leader elections once per epoch and enabling miners to catch such deviations from the protocol. Given a system start and epoch time by the genesis block, the system clock allows miners to associate epoch and wall clock time, thereby enabling them to reason about block validity and give the protocol liveness.

Filecoin依靠这个系统时钟来确保共识，特别是确保参与者每纪元只进行一次领导人选举，并使矿工能够捕捉到这种偏离协议的情况。当创世块给出了一个系统的开始时间和历元时间，系统时钟允许矿工将历元和挂钟时间联系起来，从而使他们能够推断出块的有效性并赋予协议活性。

## Clock uses - 时钟使用
Specifically, the Filecoin system clock is used:

具体来说，Filecoin系统时钟被用于：

- to validate incoming blocks and ensure they were mined in the appropriate round, looking at the wall clock time in conjunction with the block's `ElectionProof` (which contains the epoch number) (see {{<sref leader_election>}} and {{<sref block_validation>}}).
- 要验证传入的块并确保它们在适当的轮中被挖掘，请结合块的`选举证明`(其中包含纪元号)查看时钟时间 (参考{{<sref leader_election "领导人选举">}} 和 {{<sref block_validation "块校验">}})。
- to help protocol convergence by giving miners a specific cutoff after which to reject incoming blocks in this round (see {{<sref chain_sync>}}).
- 为了帮助协议收敛，给矿工一个特定的截止时间，在此之后拒绝这一轮的传入块(参见{{<sref chain_sync "区块同步">}})。
- to maintain protocol liveness by allowing participants to try leader election in the next round if no one has produced a block in this round (see {{<sref storage_power_consensus>}}).
- 为了保持协议的活性，允许参与者们在下一轮中尝试领导人选举，如果在这一轮中没有人产生一个块(参见{{<sref storage_power_consensus "存储能力共识">}})。

In order to allow miners to do the above, the system clock must:

为了让矿工做到以上几点，系统时钟必须:

1. have low clock drift: at most on the order of 1s (i.e. markedly lower than epoch time) at any given time.
1. 有低的时钟漂移:在任何给定的时间，最多在1的数量级(即明显低于epoch时间)。
2. maintain accurate network time over many epochs: resyncing and enforcing accurate network time.
2. 在多个纪元上保持准确的网络时间：重新同步和强制执行准确的网络时间。
3. set epoch number on client initialization equal to `epoch ~= (current_time - genesis_time) / epoch_time`
3. 在客户端初始化时将纪元号设置为`纪元号约~=(当前时间-创建时间) / 纪元时间`

It is expected that other subsystems will register to a NewRound() event from the clock subsystem.

预期其他子系统将注册到来自时钟子系统的新一轮事件。

## Clock Requirements - 时钟要求

Computer-grade clock crystals can be expected to have drift rates on the order of [1ppm](https://www.hindawi.com/journals/jcnc/2008/583162/) (i.e. 1 microsecond every second or .6 seconds a week), therefore, in order to respect the first above-requirement,

计算机级时钟晶体的漂移率可达[1ppm](https://www.hindawi.com/journals/jcnc/2008/583162/)(即1微秒/秒或每周0.6秒)，因此，为了满足上述第一个要求，

- clients SHOULD query an NTP server (`pool.ntp.org` is recommended) on an hourly basis to adjust clock skew.
- 客户端应该每小时查询一个NTP服务器(建议使用“pool.ntp.org”)来调整时钟倾斜
  - We recommend one of the following:
  - 我们建议如下:
    - `pool.ntp.org` (can be catered to a [specific zone](https://www.ntppool.org/zone))
    - `pool.ntp.org` (可以被满足一个[特定时区](https://www.ntppool.org/zone))
    - `time.cloudflare.com:1234` (more on [Cloudflare time services](https://www.cloudflare.com/time/))
    - `time.cloudflare.com:1234` (更多[Cloudflare时间服务](https://www.cloudflare.com/time/))
    - `time.google.com` (more on [Google Public NTP](https://developers.google.com/time))
    - `time.google.com` (更多[Google公开的NTP](https://developers.google.com/time))
    - `ntp-b.nist.gov` ([NIST](https://tf.nist.gov/tf-cgi/servers.cgi) servers require registration)
    - `ntp-b.nist.gov` ([NIST](https://tf.nist.gov/tf-cgi/servers.cgi)服务器需要注册)
  - We further recommend making 3 measurements in order to drop by using the network to drop outliers
  - 我们还建议进行3次测量，以便通过使用网络来剔除离群值
  - See how [go-ethereum does this](https://github.com/ethereum/go-ethereum/blob/master/p2p/discv5/ntp.go) for inspiration
  - 参考[go-ethereum如何做](https://github.com/ethereum/go-ethereum/blob/master/p2p/discv5/ntp.go)来获取灵感
- clients CAN consider using cesium clocks instead for accurate synchrony within larger mining operations
- 客户端可以考虑在大型挖矿操作中使用铯钟来精确同步

Assuming a majority of rational participants, the above should lead to relatively low skew over time, with seldom more than 10-20% clock skew that should be rectified periodically by the network, as is the case in other networks. This assumption can be tested over time by ensuring that:

假设大多数合理的参与者都是这样的，那么随着时间的推移，上述情况应该会导致相对较低的歪斜，很少会有超过10-20%的时钟偏差应该由网络周期性地进行纠正，就像其他网络中的情况一样。确保:

- (real-time) epoch time is as dictated by the protocol
- (实时时间)历元时间由协议规定
- (historical) the current epoch number is as expected
- (历史的)当前的纪元数正如所预料的那样

## Future work - 未来的工作

If either of the above metrics show significant network skew over time, future versions of Filecoin may include potential timestamp/epoch correction periods at regular intervals.

如果上述任何一个指标随时间的推移显示出显著的网络偏差，那么Filecoin的未来版本可能会包括定期间隔的潜在时间戳/历元校正周期。

More generally, future versions of the Filecoin protocol will use Verifiable Delay Functions (VDFs) to strongly enforce block time and fulfill this leader election requirement; we choose to explicitly assume clock synchrony until hardware VDF security has been proven more extensively.

更一般地说，未来版本的Filecoin协议将使用可验证延迟函数(VDFs)来强制执行块时间并满足领袖选举的要求;我们选择显式地假定时钟同步，直到硬件VDF安全性得到更广泛的验证。

