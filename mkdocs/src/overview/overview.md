# Overview


Morgoth is a flexible anomaly detection application. It is designed to integrate
with many existing technologies and provide powerful anomaly detection for those systems.
Currently the only supported backend is an [InfluxDB](http://influxdb.com) 0.9.x database.

## Architecture

Morgoth consist of an anomaly detection framework and an alerting framework.

## Anomaly Detection

The anomaly detection framework follows a basic workflow:

 1. Morgoth selects windows of data from its backend (InfluxDB) on configured intervals.
 2. The windows of data are processed through the anomaly detection algorithms (more on this later).
 3. Any detected anomalies are written back to the backend so that alerts can be triggered.


## Alerting

The alerting framework is very light weight. It consists of scheduled queries that check against a configured threshold.
If the threshold is crossed and alert is fired.


## Getting help

There are several places you can look for help with Morgoth.

* [Mailing List](https://groups.google.com/forum/#!forum/morgoth) -- Search or post a question
* IRC on freednode in #morgoth -- Come chat
* [Github](https://github.com/nathanielc/morgoth) -- File an issue or submit a PR.

