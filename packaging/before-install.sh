#!/usr/bin/env bash

API_USER=notification-service

getent passwd ${API_USER} > /dev/null

if [ $? -ne 0 ]; then
    echo "Creating user ${API_USER}"
    useradd -U ${API_USER} -d /opt/notification-service
fi