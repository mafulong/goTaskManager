# 一个go的并发执行manager

用于并发执行有依赖关系的任务流, 优化延迟提高性能。

## 问题/业务场景
业务应用的API往往需要Load多个信息，并将这些信息进行组装成response进行返回。
但API延迟是比较重要的，需要尽可能的短。每个LOAD都对应了至少一个RPC/HTTP请求。

## 解决思路
因此，需要一个可并发执行多个task的框架来支撑，通过并发Load提高消息。
一个Load对应了一个及以上的domain的接口调用，不论是RPC还是HTTP。

抽象每个Load操作为一个Task, Task可以依赖其他Task, 不同Task可以并发执行，但任一Task只会在所有的依赖Task执行完成后才开始执行。
每个Task只会执行一次。 这就是一个任务流，这些Task组成了一个有向无环图。


考虑点：
- Goroutine不能太多，避免影响性能
- Task不能重复执行。否则会接口放大
- 依赖顺序需要保证

对于每个Task依赖和产出的数据通过DomainData统一管理，包括最初的Request和最终的Response也可以作为DomainData的一部分。
由于Task并发执行，因此DomainData里数据要拆分成不同domain并各自使用读写锁进行保护，避免并发问题。


## 使用方式

TODO

- Load失败可控制提前返回，停止后续Load执行[可选]
- 日志收集，请Fork后自行添加自己日志框架
- Error错误码 
