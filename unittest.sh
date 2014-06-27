#!/bin/bash

# Run the basic unittests

source ./pythonpath.sh
python -m unittest discover -s src/morgoth
