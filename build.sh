#!/bin/bash

DIR=$( cd $(dirname $BASH_SOURCE[0]) && pwd)
cd $DIR
cd mkdocs

mkdocs build --clean

cd $DIR

rm -rf docs/
mv mkdocs/site docs
