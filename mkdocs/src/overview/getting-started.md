
The following is a walk through to install, configure, and run Morgoth. It will also demonstrate
how to detect anomalies with Morgoth. This tutorial assumes a Unix like environment but the concepts
will work fine on other OSes.

# Installation

Morgoth is go 'gettable' simply run:

```bash
$ go get github.com/nathanielc/morgoth/cmd/morgothd
```

Morgothd should now be downloaded and installed on your system. Make sure to add
`$GOPATH/bin` to your `$PATH` in order for your system to find the newly installed morgoth binary.


Try it out:

```bash
$ morgothd
```

You should see an error about morgoth not finding its configuration file. Now let's install
and configure the basic dependencies of Morgoth.


# Basic Configuration

Let's create a sandbox for morgoth and copy in the example configuration file:

```bash
$ mkdir morgoth
$ cd morgoth
$ cp $GOPATH/src/github.com/nathanielc/morgoth/morgoth.yaml.example ./morgoth.yaml
```

Take a look at the configuration file. All Morgoth configuration is contained in this [yaml](http://en.wikipedia.org/wiki/YAML) file.
There are a few basic sections:

* **engine** -- Defines the back end engine Morgoth should use. Currently only an [InfluxDB](http://influxdb.com) engine is supported.
* **schedules** -- Defines the schedule for selecting windows of data from the backend engine.
* **mappings** -- Defines which detection algorithms to apply to which metrics.
* **alerts** -- Defines a set of queries to execute in order to do simple threshold based alerting.

More details of each of these sections can found in the [Configuration](/configuration/configuration) section.


## Dependencies

### InfluxDB

Morgoth uses [InfluxDB](http://influxdb.com) (specifically 0.9.x) as a source of metric data.
Note, Morgoth does not ingest any data, it simple expects to query alreadt existing data from InfluxDB and process it.
The default configuration we just copied is looking for an InfluxDB instance to be running on `localhost` on port `8086`.
It expects to be able to login as the default root user and connect to a database called `morgoth`.

At this point either create a `morgoth` database on your running InfluxDB instance or change the configuration to point to a different endpoint.

Now that we have InfluxDB running Morgoth is ready to start detecting anomalies.

# Running Morgoth

Start Morgoth again from the sandbox directory.

```bash
$ morgothd
```

Morgoth algorithms can store metadata about what they have learned that persists from restart to restart.
This data is stored in a [BoltDB](https://github.com/boltdb/bolt) database.

# Detecting Anomalies

## Getting data for Morgoth

Now that Morgoth is running lets tell it how to find data.
In the [schedules](/configuration/schedules) section we define several queries to run an on what frequency.
We do not need to difine a `time` clause as Morgoth will do this for you, based on the current time.
The example configuration is selecting data from the `cpu_idle` measurement and grouping by all tags.
The result of this query is that serveral series are returned each with a unique tag set for the `cpu_idle` measurement.
Each of these series or windows as Morgoth calls them are passed to the mappings in order to determine what to do with them.

But before we can see this in action we need to get data into InfluxDB.
If you already have data change the query to select data from a measurement in your data set.
If you do not have data in InfluxDB install [Telegraf](https://github.com/influxdb/telegraf) and start it up with simple config that only has the cpu and mem plugins enabled.

Something like this:

```toml
[influxdb]
url = "http://localhost:8086"
database = "metrics"

[agent]
interval = "5s"
debug = false
hostname = "morgoth1"

[cpu]
  # no configuration
[mem]
  # no configuration
```

Now that telegraf is running confirm that you are getting data in InfluxDB.

## Mapping Data

Now that InfluxDB has data, Morgoth can query it.
Each window or series of data returned from InfluxDB will be sent to the mapper inside of Morgoth.
The mapper is configured via the [mappings](/configuration/mappings) section.
Each mapping is a contains a `name` regex that will match against the measurement name, `cpu_idle` in this case, and a set of tag regexs.
Each tag regex must also match for the window to match the mapping.
Once a window with its name and tags matches a mapping it will be processed by an instance of a detector with the configured settings.
The default configuration references some [fingerprinters](/concepts/fingerprints) etc.
Basically we have told the mapping to use a basic Sigma approach to finding anomalies in the `cpu_idle` data.


The way the anomaly detection algorithms work is that they will always mark the first couple windows as anomalous.
This gives us a chance to see an anomaly get recorded.

After a few seconds run:

```sql
SELECT start FROM anomaly GROUP BY *
```

Notice that `cpu_idle` has been converted into a tag now and the measurement we queried was the `anomaly` measurement.
The measurement tag and the anomaly measurement name can be configured as well(see [here](/configuration/engine)


## Alerting on Anomalies

Now that we have a few anomalies Morgoth example configuration should have also queried for the number of anomalies and written a notification to a log (alerts.log based on the example config).
Check this log now to see if the alert was fired.
If so everything has worked as expected, you can now play around with the configuration to get alerts for anomalies you care about.
Just like the scheduled queries above the alert queries will have an time where clause automatically append to them based on their period.

