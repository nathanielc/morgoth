#!/bin/bash

# Run the integration tests

source ./pythonpath.sh
python -m unittest discover -s src/morgoth -p 'int_test*.py'

