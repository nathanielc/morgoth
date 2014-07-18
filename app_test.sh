#!/bin/bash

#Run full end to end application tests

python -m unittest discover -s src/morgoth -p 'app_test*.py'
