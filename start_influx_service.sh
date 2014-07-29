#!/bin/bash


mkdir -p $GOPATH/src/github.com/influxdb
cd $GOPATH/src/github.com/influxdb
git clone https://github.com/nvcook42/influxdb.git
cd influxdb

./configure
make build
screen ./influxdb

