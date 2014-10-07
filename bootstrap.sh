#!/bin/bash


echo "export GOPATH=/gopath" >> /etc/profile.d/gopath.sh
echo "export PATH=/gopath/bin:\$PATH" >> /etc/profile.d/gopath.sh
source /etc/profile.d/gopath.sh

mkdir -p /gopath/{bin,pkg,src}

#rpm -Uvh http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm


yum install golang git mercurial -y

go get github.com/tools/godep

cd /gopath/src/github.com/nvcook42/morgoth/
godep restore

chown -R vagrant:vagrant /gopath

godep go test ./...
