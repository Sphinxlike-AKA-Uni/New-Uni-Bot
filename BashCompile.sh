#!/bin/bash
#go build -gccgoflags "-L /lib64 -l pthread" Uni.go
go build Uni.go
if [ $? == 0 ]; then
 ./Uni
fi
