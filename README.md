# UStack

> Hi, welcome to here, UStack is comming soon ...

`UStack` is a communication management component embedded in user space applications. It provides functions including but not limited to package statistics, forwarding, discarding, load balancing, encoding and decoding, and supports remote dynamic behavior adjustment. It is currently being implemented as planned.


```
                  +----------+   +----------+
                  | endpoint |   | endpoint |
+--------------+ +-+----------+---+----------+----+
|              | |         upper deck             |
|  controller  | +--------------------------------+
|              |
|              |
|  +--------+  | +--------------------------------+
|  | event  |  | |         data processor         |
|  +--------+  | +--------------------------------+
|              |
|  +--------+  | +--------------------------------+
|  | stat   |  | |         data processor         |
|  +--------+  | +--------------------------------+
|              |
|  +--------+  | +--------------------------------+
|  | log    |  | |         ... ... ... ...        |
|  +--------+  | +--------------------------------+
|              |
|              | +--------------------------------+
|  +--------+  | |         data processor         |
|  |   ...  |  | +--------------------------------+
|  +--------+  |
|              |
|              | +--------------------------------+
|              | |         lower deck             |
+--------------+ +---------+-------------+--------+
                           |  transport  |
                           +-------------+
```