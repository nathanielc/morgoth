#!/bin/bash

# Run the integration tests

source ./pythonpath.sh
python -m unittest discover -s src/morgoth -p 'int_test*.py'

if [ $? -ne 0 ]
then
    cat $GOPATH/src/github.com/influxdb/influxdb/influxdb.log
fi
