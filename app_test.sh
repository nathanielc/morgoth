#!/bin/bash

#Run full end to end app test

python -m unittest discover -s src/morgoth -p 'app_test*.py'
