# Getting Started

To get started using morgoth follow the simple exercise below
or start by reading the configuration section of this documentation
to learn how to start using Morgoth for your own project.

## Install Morgoth

First install morgoth via `go get` set your `$GOPATH` and run:


```
$ go get github.com/nathanielc/morgoth/cmd/morgoth
```

The `morgoth` binary will be in `$GOPATH/bin`.

Next you will need to install and configure [Kapacitor](https://github.com/influxdata/kapacitor).

Generate a default Kapacitor config file:


```
$ kapacitord config > kapacitor.conf
```

Finally configure Kapacitor to use Morgoth by adding this section to the default configuration file you just created.


```
[udf]
[udf.functions]
   [udf.functions.morgoth]
       prog = "/path/to/bin/morgoth"
       timeout = "10s"
```


Start Kapacitor and if you do not get any errors you are good to go.


```
$ kapacitord -config ./kapacitor.conf
```

## Collecting Data

For this example we are going to use [Telegraf](https://github/com/influxdata/telegraf) to send data to Kapacitor.

[Install](https://github.com/influxdata/telegraf#installation) Telegraf and use this simplified configuration.


```
# Configuration for telegraf agent
[agent]
 ## Default data collection interval for all inputs
 interval = "1s"
 ## Rounds collection interval to 'interval'
 ## ie, if interval="10s" then always collect on :00, :10, :20, etc.
 round_interval = true

# Configure telegraf to send data to Kapacitor
# Kapacitor is write compatible with the InfluxDB database so
# we just configure an InfluxDB output and point it at Kapacitor.
[[outputs.influxdb]]
 urls = ["http://localhost:9092"]
 database = "telegraf" # required
 retention_policy = "default"
 precision = "s"

# Read metrics about cpu usage
[[inputs.cpu]]
 percpu = false
 totalcpu = true
 fielddrop = ["time_*"]
```

Place the above contents in a file `./telegraf.conf`.


Start Telegraf running:


```
   $ telegraf -config ./telegraf.conf
```

Confirm Kapacitor is receiving Telegraf data.

```
   $ kapacitor stats ingress
```

You should see an entry for the telegraf database and cpu measurement.

## Writing a Kapacitor task using Morgoth

Kapacitor processes data in "tasks".
To tell Kapacitor what to do for a task we need to write a TICKscript.
The script will define what data we are going to process and how to process it.
In our example we are going to use Morgoth with a one sigma fingerprinter.

Here is a generic template TICKscript for using Morgoth with a single sigma fingerprinter:
The script uses defaults for the cpu usage_idle data but can easily be modified to work on other datasets.

```go
// The measurement to analyze
var measurement = 'cpu'

// Optional group by dimensions
var groups = [*]

// Optional where filter
var whereFilter = lambda: TRUE

// The amount of data to window at once
var window = 1m

// The field to process
var field = 'usage_idle'

// The name for the anomaly score field
var scoreField = 'anomalyScore'

// The minimum support
var minSupport = 0.05

// The error tolerance
var errorTolerance = 0.01

// The consensus
var consensus = 0.5

// Number of sigmas allowed for normal window deviation
var sigmas = 3.0

stream
  // Select the data we want
  |from()
      .measurement(measurement)
      .groupBy(groups)
      .where(whereFilter)
  // Window the data for a certain amount of time
  |window()
     .period(window)
     .every(window)
     .align()
  // Send each window to Morgoth
  @morgoth()
     .field(field)
     .scoreField(scoreField)
     .minSupport(minSupport)
     .errorTolerance(errorTolerance)
     .consensus(consensus)
     // Configure a single Sigma fingerprinter
     .sigma(sigmas)
  // Morgoth returns any anomalous windows
  |alert()
     .details('')
     .crit(lamda: TRUE)
     .log('/tmp/cpu_alert.log')
```


Place the above contents into a file called `cpu_alert.tick`.
This script will take the incoming cpu data and window it into 1 minute buckets.
It will then pass each window to the Morgoth process.
Morgoth will return any window of data it thinks is anomalous.
We trigger an alert of all windows received from Morgoth and log the alert in the `/tmp/cpu_alert.log` file.

First we must define the task in Kapacitor:


```
$ # Define the task with the name cpu_alert
$ kapacitor define cpu_alert -type stream -dbrp telegraf.default -tick cpu_alert.tick
$ # Start the task
$ kapacitor enable cpu_alert
$ # Get info on the running task
$ kapacitor show cpu_alert
```

If you didn't get any error you are now successfully sending data to Morgoth via Kapacitor.


## Detecting your first anomaly

To detect our first anomaly we should let at least 1 minute pass so that Morgoth gets at least one window of good data.
After a minute or two has passed create some cpu activity.
Run this bash script to spawn two infinite bash loops that spin at 100% cpu.
We will kill these later.

```
$ for i in {1..2}; do { while true; do i=0 ; done & }; done
```

Again wait for a minute or so to pass and watch the alert log for alerts.

```
$ tail -F /tmp/cpu_alert.log
```

After a short wait you should see the critical alert triggered.

Kill the backgrounded jobs:

```
$ kill $(jobs -p)
```

After a short wait again you should see an OK alert in the log indicating the CPU usage has recovered.

## Interpreting the Anomaly Score

Morgoth returns an anomaly score when it detects an anomaly.
The anomaly score is defined as `1 - averageSuppport`, where `averageSuppport` is the average of the `support` values returned from each fingerprinter.
Remembering that `support = count / total`, where count is the number of times this event has been seen and total is the total number of events seen, we can interpret the support as a frequency percentage.
For example a support of 0.05 can be interpreted as: the event has been seen 5% of the time.
So the anomaly score can be interpreted as the percentage of time the event was not seen.

Specifically using the above script since we have only one fingerprinter `sigma` and assuming we got an anomaly score of 0.98 from Morgoth, we can interpret the score as:
Windows that are similar to the current window, as defined by the standard deviation and mean of the windows, have only been seen about 2% of the time.


## Next steps

At this point you should have a basic grasp of how to use Kapacitor and Morgoth.
Anomaly detection is not a simple task and requires that you experiment with different methods before you arrive at useful results.
At this point I would recommend playing around with some of the other fingerprinters and getting a feel for their settings.

Also Kapacitor is a capable tool for selecting specific sets of data and pre-processing the data as needed.
As in all machine learning algorithms garbage in equals garbage out, take the time to learn Kapacitor's TICKscript so that
you can send in clean useful data to the Morgoth algorithms.
