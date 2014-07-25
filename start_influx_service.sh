#!/bin/bash

#git clone https://github.com/nvcook42/influxdb.git
git clone https://github.com/influxdb/influxdb.git

cd influxdb
./configure
make build
ls -l .
./daemon

