
Please contribute your anomaly detection alogorithms to morgoth.


Morgoth Structure
=================

Morgoth is composed of several core components:

* Metric Configuration
* Input plugins
* Anomaly Detectors
* Notification plugins

Metrics Configuration
---------------------
Every metric will be tagged with several configurations.

Each metrics will be associated with one or more anomaly detection alogrithim.

Input plugins
-------------

Morgoth can accept data from many different sources. And more plugins can be
written to expand its compatibility.

Anomaly Detectors
-----------------

Morgoth allows several anomaly detection alogrithms to be used across all the
metrics. This allows each alogrithm to be tailored to the specific
characteristics of the metrics data.


Notification plugins
--------------------

The notification system is also plugable.

