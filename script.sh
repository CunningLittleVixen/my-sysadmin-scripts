#!/usr/bin/env bash

INTERVAL=5


while true; do

    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')


    echo "--- $TIMESTAMP ---" >> monitor.log


    free -h >> monitor.log
    df -h >> monitor.log
    uptime >> monitor.log


    echo "" >> monitor.log

    sleep "$INTERVAL"
done