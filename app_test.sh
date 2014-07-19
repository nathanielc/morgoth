#!/bin/bash

#Run full end to end application tests

source ./pythonpath.sh
python -m unittest discover -s src/morgoth -p 'app_test*.py'
