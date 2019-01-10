#!/usr/bin/env bash

API_USER=notification-service

chown -R ${API_USER}:${API_USER} /opt/notification-service
chkconfig notification-service on