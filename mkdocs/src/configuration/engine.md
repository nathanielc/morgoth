#  Engine

The only supported engine currently is [InfluxDB](https://influxdb.com)

The configuration consists of these values:

```yaml
---
engine:
  influxdb:
    host: localhost
    port: 8086
    user: root
    password: root
    database: metrics
    # The name of the measurement to use when anomalies are written back to InfluxDB.
    anomaly_measurement: anomaly
    # The name of the tag to use for the measurement name.
    measurement_tag: msrmnt # Note this cannot be 'measurement' since InfluxDB seems to have bugs associated with using key words.
```
