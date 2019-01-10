#!/usr/bin/env bash

API_USER=notification-service

chown -R ${API_USER}:${API_USER} /opt/notification-service

if [ -x /etc/init.d/notification-service ]; then
    service notification-service start
fi
