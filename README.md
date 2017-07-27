
Morgoth [![Build Status](https://travis-ci.org/nathanielc/morgoth.svg?branch=master)](https://travis-ci.org/nathanielc/morgoth)
=======

Morgoth is a framework for flexible anomaly detection algorithms packaged to be used with [Kapacitor](https://github.com/influxdata/kapacitor/)

Morgoth provides a framework for implementing the smaller pieces of an anomaly detection problem.

The basic framework is that Morgoth maintains a dictionary of normal behaviors and compares new windows of data to the normal dictionary.
If the new window of data is not found in the dictionary then it is considered anomalous.

Morgoth uses algorithms, called fingerprinters, to compare windows of data to determine if they are similar.
The [Lossy Counting Algorithm](http://www.vldb.org/conf/2002/S10P03.pdf)(LCA) is used to maintain the dictionary of normal windows.
The LCA is a space efficient algorithm that can account for drift in the normal dictionary, more on LCA below.

Morgoth uses a consensus model where each fingerprinter votes for whether it thinks the current window is anomalous.
If the total votes percentage is greater than a consensus threshold then the window is considered anomalous.

## Getting started

### Install

Morgoth can be installed via go:

```sh
go get github.com/nathanielc/morgoth/cmd/morgoth
```

### Configuring

Morgoth can run as either a child process of Kapacitor or as a standalone daemon that listens on a socket.

#### Child Process

Morgoth is a UDF for [Kapacitor](https://github.com/influxdata/kapacitor).
Add this configuration to Kapacitor in order to enable using Morgoth.

```
[udf]
  [udf.functions]
    [udf.functions.morgoth]
      prog = "/path/to/bin/morgoth"
      timeout = "10s"
```

Restart Kapacitor and you are ready to start using Morgoth within Kapacitor.

#### Socket

To use Morgoth as a socket UDF start the morgoth process with the `-socket` option.

```
   morgoth -socket /path/to/morgoth/socket
```

Next you will need to configure Kapacitor to use the morgoth socket.

```
[udf]
  [udf.functions]
    [udf.functions.morgoth]
      socket = "/path/to/morgoth/socket"
      timeout = "10s"
```

Restart Kapacitor and you are ready to start using Morgoth within Kapacitor.


### TICKscript

Here is an example TICKscript for detecting anomalies in cpu data from [Telegraf](https://github.com/influxdata/telegraf).

```javascript
stream
    |from()
        .measurement('cpu')
        .where(lambda: "cpu" == 'cpu-total')
        .groupBy(*)
    |window()
        .period(1m)
        .every(1m)
    @morgoth()
        // track the 'usage_idle' field
        .field('usage_idle')
        .errorTolerance(0.01)
        // The window is anomalous if it occurs less the 5% of the time.
        .minSupport(0.05)
        // Use the sigma fingerprinter
        .sigma(3.0)
        // Multiple fingerprinters can be defined...
```


## Fingerprinters

A fingerprinter is a method that can determine if a window of data is similar to a previous window of data.
In effect the fingerprinters take fingerprints of the incoming data and can compare fingerprints of new data to see if they match.
These fingerprinting algorithms provide the core of Morgoth as they are the means by which Morgoth determines if a new window of data is new or something already observed.

An example fingerprinting algorithm is a *sigma* algorithm that computes the mean and standard deviation of a window and store them as the fingerprint for the window.
When a new window arrives it compares the fingerprint (mean, stddev) of the new window to the previous window.
If the windows are too far apart then they are not considered at match.

By defining several fingerprinting algorithms Morgoth can decide whether new data is anomalous or normal.

## Lossy Counting Algorithm

The LCA counts frequent items in a stream of data.
It is *lossy* because to conserve space it will drop less frequent items.
The result is that the algorithm will find frequent items but may loose track of less frequent items.
More on the specific mathematical properties of the algorithm can be found below.

There are two parameters to the algorithm, error tolerance (e) and minimum support (m).
First e is in the range [0, 1] and is an error bound, interpreted as a percentage value.
For example given and e = 0.01 (1%), items less the 1% frequent in the data set can be dropped.
Decreasing e will require more space but will keep track of less frequent items.
Increasing e will require less space but will loose track of less frequent items.
Second m is in the range [0, 1] and is a minimum support such that items that are considered frequent have at least m% frequency.
For example if m = 0.05 (5%) then if an item has a support less than 5% it is not considered frequent, aka normal.
The minimum support becomes the threshold for when items are considered anomalous.

Notice that m > e, this is so that we reduce the number of false positives.
For example say we set e = 5% and m = 5%.
If a *normal* behavior X, has a true frequency of 6% than based on variations in the true frequency, X might fall below 5% for a small interval and be dropped.
This will cause X's frequency to be underestimated, which will cause it to be flagged as an anomaly, triggering a false positive.
By setting e < m we have a buffer to help mitigate creating false positives.


### Properties

The Lossy Counting algorithm has three properties:

1. there are no false negatives,
2. false positives are guaranteed to have a frequency of at least (m - e)*N,
3. the frequency of an item can underestimated by at most e*N,

where N is the number of items encountered.

The space requirements for the algorithm are at most (1 / e) * log(e*N).
It has also been show that if the item with low frequency are uniformly random than the space requirements are no more than 7 / e.
This means that as Morgoth continues to processes windows of data its memory usage will grow as the log of the number of windows and can reach a stable upper bound.



## Metrics

Morgoth exposes metrics about each detector and fingerprinter.
The metrics are exposed as a promethues `/metrics` endpoint over HTTP.
By default the metrics HTTP endpoint binds to `:6767`.

>NOTE: Using the metrics HTTP endpoint only makes sense if you are using Morgoth in socket mode as otherwise each new process would collide on the bind port.

Metrics will have some or all of these labels:

* task - the Kapacitor task ID.
* node - the ID of the morgoth node within the Kapacitor task.
* group - the Kapacitor group ID.
* fingerprinter - the unique name for the specific fingerprinter, i.e. `sigma-0`.


The most useful metric for debugging why Morgoth is not behaving as expected is likely to be the `morgoth_unique_fingerprints` gauge.
The metric reports the number of unique fingerprints each fingerprinter is tracking.
This is useful because if the number is large or growing with each new window its likely that the fingerprinter is erroneously marking every window as anomalous.
By providing visibility into each fingerprinter, Morgoth can be tuned as needed.

Using Kapacitor's scraping service you can scrape the Morgoth UDF process for these metrics and consume them within Kapacitor.
See this [tutorial](https://docs.influxdata.com/kapacitor/latest/pull_metrics/scraping-and-discovery/) for more information.

