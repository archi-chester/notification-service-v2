#!/usr/bin/env bash

if [ -x /etc/init.d/notification-service ]; then
    service notification-service stop
fi