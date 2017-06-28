#!/bin/bash
ROOT_PATH="github.com/qinguoan/prometheus/"
mkdir -p $1/../src/$ROOT_PATH

rsync -avz $1/*  $1/../src/$ROOT_PATH  --exclude=.svn --exclude=rpm/ --exclude=.git

export GOPATH=$1/../
export GOROOT=/usr/local/golang
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT/bin:/Users/Tinker/go/bin:$GOBIN

cmd="prometheus"
bin="prometheus"
cd $1/../src/$ROOT_PATH/cmd/$cmd  && pwd &&  go build -o ${bin} "main.go" "config.go"

cd $1/rpm
rpm_create $2.spec -v $3 -r $4 -p /home/a/