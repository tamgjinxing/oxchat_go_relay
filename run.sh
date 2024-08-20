#!/bin/bash
PIDS=`ps -ef | grep "oxchat_go_relay" | grep -v grep |awk '{print $2}'`
for PID in $PIDS ; do
        kill $PID > /dev/null 2>&1
done

go build -o oxchat_go_relay

nohup ./oxchat_go_relay >> ./log/oxchat_relay.log 2>&1&