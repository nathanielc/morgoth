#!/bin/bash


#yum install golang-1.3.3 -y
#yum install -y git mercurial bzr protobuf-compiler flex bison valgrind \
#  gcc-c++ libstdc++-static make autoconf libtool zlib-devel bzip2-libs \
#  bzlib2-devel
#
## Fix missing link for bz2 libs
#cd /usr/lib64
#ln -s libbz2.so.1 libbz2.so
#
#echo "export GOPATH=/gopath" >> /home/vagrant/.bashrc
#export GOPATH=/gopath
#mkdir -p $GOPATH/src/github.com/influxdb
#cd $GOPATH/src/github.com/influxdb
#git clone https://github.com/influxdb/influxdb.git
#
#cd $GOPATH/src/github.com/influxdb/influxdb
#
#./configure
#make

#$GOPATH/src/github.com/influxdb/influxdb/influxdb -config $GOPATH/src/github.com/influxdb/influxdb/config.sample.toml &

yum install -y http://s3.amazonaws.com/influxdb/influxdb-latest-1.x86_64.rpm

/etc/rc.d/init.d/influxdb start

