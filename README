
Morgoth
=======

Morgoth is a metric collection system with two main focuses:

1. Detect anomalies in the data in real time.
2. Allow storage of arbitrary metrics with little overhead.

It is built on top of MongoDB.


Detect Anomalies in Real Time
-----------------------------

Morgoth is watching the metrics in real time. It will provide notification of
any anomalies in the metrics given to it.

Store Arbitrary Metrics
-----------------------

Morgoth stores only the data points it is given. This provides two key
advantages.

First an application can log metrics of dynamic aspects. This way
anomalies on per user or per client can be detected.

Second the application can aggresively log a metric and then decide later how
important the metric is.

TODO: Implement aggregation so metrics can be reduced in size at a later date.


Anomaly Detection
=================

Since morgoth is primarily focused on anomaly detection we should explain
its methods...
