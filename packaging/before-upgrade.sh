#!/usr/bin/env bash

API_USER=notification-service
getent passwd ${API_USER} > /dev/null

if [ $? -ne 0 ]; then
    useradd ${API_USER}
fi

if [ -x /etc/init.d/notification-service ]; then
    service notification-service stop
fi