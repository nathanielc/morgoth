# Schedules

Since Morgoth does not ingest data itself it needs to be told were and how to get data.
This is were schedules come in, they define a query and a period that the query should be run.
Each time the query is run Morgoth adds a time clause resticting it to just the window of time since the last time the query ran.

The basic components of a scheduled query are:

1. query -- The query to run without any time clause.
2. period -- How often to run the query. Can be any duration string (s,m,h,w,d, etc.)
3. delay -- How long to wait from real time before executing the query. Usefull if data takes time to arrive at InfluxDB.
4. tags -- Static set of tags to set on all windows returned from the query


As the queries execute data receives data and processes by sending it to the mapper.
Since the mapper uses tags to make mapping decisions its possible to define a set of static tags to be set on the results of a given query.
This allows the mapper to make decisions based on which retention policy, etc. the data came from.

Here is an example schedules config section:

```yaml
schedules:
  - query: SELECT value FROM "day".cpu_idle GROUP BY *
    period: 30s
    delay: 10s
    tags:
      ret: day
  - query: SELECT value FROM "week".cpu_idle GROUP BY *
    period: 5m
    delay: 1m
    tags:
      ret: week
  - query: SELECT value FROM "month".cpu_idle GROUP BY *
    period: 1h
    delay: 1m
    tags:
      ret: month
```

Notice that these queries are selecting data from different retention policies and adding a tag as such.
This way Morgoth can map to a new detector for each set of data for different retention policies allowing for detection of anomalies of widely different time scales.

Also note the  `GROUP BY *` InfluxDB will not return tags in the result unless the are requested, so either add a `GROUP BY *` or explicity name each tag used in the mappings section.


