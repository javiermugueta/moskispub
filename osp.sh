#!/bin/sh
# jmu, may 2021
#
# packages
#
go get github.com/oracle/oci-go-sdk/common github.com/oracle/oci-go-sdk/streaming github.com/oracle/oci-go-sdk/common github.com/eclipse/paho.mqtt.golang github.com/google/uuid io/ioutil github.com/oracle/oci-go-sdk/identity context fmt io/ioutil os github.com/oracle/oci-go-sdk/common time
#
# build
#
go build osp.go
#
# run
#
./osp 
#