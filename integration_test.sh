#!/bin/bash

# Run the integration tests

# Start mongo ...
mongod --fork --logpath /dev/null

# Start influxdb ...



sleep 1

source ./pythonpath.sh
python -m unittest discover -s src/morgoth -p 'int_test*.py'
ret=$?

# Stop mongod
pkill mongod

# Stop influxdb


exit $ret
