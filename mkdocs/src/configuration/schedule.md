# Schedule

All the activity within Morgoth is controlled by a schedule. The schedule consists
of a list of rotations and a global delay:


## Rotations


Rotations define how data is down sampled and when anomaly detection is triggered.

Rotations have two properties:

* Period
* Resolution

A rotation with a period of 5 minutes and a resolution of 10 seconds will down sample the raw metric data
so that there is only a point every 10 seconds and will keep at most 5 minutes of data. This rotation will
then always contain 30 data points. Every 5 minutes those 30 data points will be passed into the configured
anomaly detection algorithms.

In this way each rotation defines a down sampled buffer of data to consider when searching for anomalies.
Most algorithms will perform best with at least 30 data points per rotation.

### Example configuration

```yaml
schedule:
  rotations:
    -  {period: 5m, resolution: 10s}
    -  {period: 10m, resolution: 20s}
    -  {period: 20m, resolution: 40s}
```

This configuration defines three rotations each with 30 data points.


## Delay

Because of some internal limitations data needs to arrive at Morgoth within a reasonable time of the timestamp
on the data points. Configuring the delay defines the maximum time that Morgoth will tolerate data arriving late.

### Example configuration

```yaml
schedule:
  delay: 15s
```


# Global Behavior

Currently the schedule is global to entire Morgoth application. This has been sufficient for current use
cases. Defining a different schedule for different sets of metrics just adds unnecessary complexity. If
different schedules are needed just run multiple instances of Morgoth.

