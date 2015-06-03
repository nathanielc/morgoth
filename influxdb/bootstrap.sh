#!/bin/bash

yum install -y http://get.influxdb.org/influxdb-0.9.0_rc30-1.x86_64.rpm

/etc/rc.d/init.d/influxdb start

