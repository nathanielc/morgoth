
The following is a walk through to install, configure, and run Morgoth. It will also demonstrate
how to detect anomalies with Morgoth. This tutorial assumes a Unix like environment but the concepts
will work fine on other OSes.

# Installation

Morgoth is go 'gettable' simply run:

```bash
$ go get github.com/nvcook42/morgoth
```

Morgoth should now be downloaded and installed on your system. Make sure to add
`$GOPATH/bin` to your `$PATH` in order for your system to find the newly installed morgoth binary.


Try it out:

```bash
$ morgoth
```

You should see an error about morgoth not finding its configuration file. Now let's install
and configure the basic depedencies of Morgoth.


# Basic Configuration

Let's create a sandbox for morgoth and copy in the example configuration file:

```bash
$ mkdir morgoth
$ cd morgoth
$ cp $GOPATH/src/github.com/nvcook42/morgoth/examples/morgoth.yaml.example ./morgoth.yaml
```

Take a look at the configuration file. All Morgoth configuration is contained in this [yaml](en.wikipedia.org/wiki/YAML) file.
There are a few basic sections:

* **engine** -- Defines the backend engine Morgoth should use. Currently only an [InfluxDB](influxdb.com) engine is supported.
* **schedule** -- Defines the schedule that Morgoth down samples the incoming data stream and searches for anomalies.
* **metrics** -- Defines which detection algorithms to apply to which metrics.
* **fittings** -- Defines which methods are enabled for input and output to Morgoth. For example the REST API is configured in this section.
* **morgoth** -- Defines properties of the Morgoth application itself (data directories, etc). This section is absent in the example since we just need the defaults.

More details of each of these sections can found in the [Configuration](/configuration/configuration) section.


## Dependencies

### InfluxDB

Morgoth uses [InfluxDB](influxdb.com) to store and process metric data. The default configuration we
just copied is looking for an InfluxDB instance to be running on `localhost` on port `8086`. It expects
to be able to login as the default root user and connect to a database called `morgoth`.

If you do not already have an InfluxDB instance running then I suggest using the provided Vagrantfile for InfluxDB. Just `cd $GOPATH/src/github.com/nvcook42/morgoth/influxdb/` and run
`vagrant up`. Installing [Vagrant](https://www.vagrantup.com/) is simple and can get your environment up and running quickly.

At this point either create a `morgoth` database on your running InfluxDB instance or change the configuration to point to a different endpoint.

Now that we have InfluxDB running Morgoth is ready to start detecting anomalies.


# Running Morgoth

Start Morgoth again from the sandbox directory.

```bash
$ morgoth
```

Morgoth algorithms can store metadata about what they have learned that persists from restart to restart.
This data is stored in the current directory under a `meta` directory. If you wish to change that directory
see [this](configuration/morgoth/)

# Detecting Anomalies

## Getting data into Morgoth

Now that Morgoth is running lets give it some data. The example configuration has configured Morgoth to
listen on port `2003` for graphite structured metrics.

Let's with start a simple example using the load on the host runing Morgoth. The [MGOF](#) algorithm, while simple,
requires that the data be bounded. Since load data isn't really bounded we will just use a good approximation.
Edit the configuration and set the `max` value in the `metrics` section to 1.5 times the number of CPUs on the 
host. Like this:

```yaml
metrics:
 - pattern: .*
   detectors:
     - mgof:
         min: 0
         max: 6 #NUM_CPUS * 1.5
```

Now let's start a background process to send the load metric to morgoth every second.

```bash
$ while true; do echo "load $(cat /proc/loadavg | awk '{print $1}') $(date +'%s')" | nc localhost 2003; sleep 1; done &
```

At this point Morgoth should be receiving the load of the host every second.

## Getting data out of Morgoth

Morgoth has a REST API listening on port `7000` by default. Simply curl the data to see the load metrics so far.

```bash
$ curl -X GET http://localhost:7000/data/load
```

You should see some JSON that lists several data points since we started reporting load to Morgoth.

Now we wait... Morgoth is detecting anomalies and does so by comparing one period of time to the next.

While we wait let's learn about [schedules](configuration/schedule/).

The default configuration schedule looks like this:

```yaml
schedule:
  rotations:
    - {period: 2m, resolution: 2s}
    - {period: 4m, resolution: 4s}
    - {period: 8m, resolution: 8s}
    - {period: 24m, resolution: 24s}
  delay: 15s
```

The smallest rotation is 2 minutes and the default `normal_count` for the MGOF algorithm is `3`. This means that
we need to wait 6 minutes before Morgoth will consider the load metric as not anomalous(aka normal).

This gives us a chance to query for detected anomalies since the first several rotations will be considered anomalous.

## Querying for anomalies

The Morgoth REST API also has an endpoint for querying anomalies. Run this curl command:

```bash
$ curl -X GET http://localhost:7000/anomalies/load
```

If it has been a few minutes than there should be a few anomalies recored so far.

Now we need to wait until Morgoth no longer marks any rotations as anomalous.

Once a few minutes have past without any 'anomalies' let's create a real anomaly for Morgoth to find.

Run this command to create some pointless be real load.

```bash
$ for i in {1..2}; do { while true; do true ; done & }; done
```

This command created two inifite while loops just spinning on the cpu. This should increase the load to at least 2 on the system.
Let this run for a few minutes.

Once you are statisfied that Morgoth was able to detect this anomaly you can kill both the while loops and the loops sending data to Morgoth
via this command:

```bash
$ kill $(jobs -p)
```


