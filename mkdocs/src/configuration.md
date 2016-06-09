# Configuration

Morgoth can run as either a standalone daemon that listens on a socket or be invoked as a child process of Kapacitor.

## Socket UDF

To use Morgoth as a socket UDF start the morgoth process with the `-socket` option.

```
$ morgoth -socket /path/to/morgoth/socket
```

Next you will need to configure Kapacitor to use the morgoth socket.

```
[udf]
[udf.functions]
   [udf.functions.morgoth]
       socket = "/path/to/morgoth/socket"
       timeout = "10s"
```


## Process UDF

To use Morgoth as a child process of Kapacitor all you need to do is configure Kapacitor.

```
[udf]
[udf.functions]
   [udf.functions.morgoth]
       prog = "/path/to/bin/morgoth"
       timeout = "10s"
```


## Logging

Morgoth allows different logging levels DEBUG, INFO, WARN, ERROR or OFF.
You can set the default logging level via the flag `-log-level`

```
$ morgoth -log-level warn
```

