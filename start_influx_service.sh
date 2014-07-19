#!/bin/bash

git clone https://github.com/nvcook42/influxdb.git

cd influxdb
./configure
make
./daemon

