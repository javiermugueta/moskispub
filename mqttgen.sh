#!/bin/sh
# 
# randomly generates json messages in test.mosquito.org
#
topic="jmu-signals"
while :
do 
    MSG="{\"t\":${RANDOM},\"d\":${RANDOM},\"v\":${RANDOM}}"
    echo $MSG
    mosquitto_pub -h test.mosquitto.org -t $topic -m $MSG 
    sleep 1
done


