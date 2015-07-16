#!/bin/bash

set -e

mkdir -p $GOPATH/src/github.com/influxdb
cd $GOPATH/src/github.com/influxdb
if [ -d influxdb ]
then
    cd influxdb
    git pull
else
    git clone https://github.com/nvcook42/influxdb.git
    cd influxdb
fi

./configure
make build
screen -d -m ./influxdb -stdout=true || cat influxdb.log

