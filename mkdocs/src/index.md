# Morgoth

Welcome to the Morgoth documentation. This site is a working document.
Basic tutorials and examples have been documented but there is still more to do.

## State of Morgoth
Currently Morgoth is pre-production. I (the author) have a Morgoth instance running along side
a production metrics stream in a non-critical path. A single 4 core machine running Morgoth backed by a
single 4 core InfluxDB instance easily handles several thousand metrics. There are many factors that
contribute to load and I do not have solid performance metrics at this time.


A few back end features are either not implemented or not completely functional at this time. InfluxDB the
main engine backing Morgoth is going through significant changes for its version 0.9.0. There are many
features promised that will make several of the tasks currently not implemented easier to implement. InfluxDB
projects a release of version 0.9.0 in March of 2015. Once their release has stabilized a I will refactor Morgoth
to use the newer version of InfluxDB.




