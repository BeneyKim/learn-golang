#!/bin/bash

echo "Starting SSH Daemon..."
service ssh start

echo "Starting Service..."
${WORK_DIR}${SERVICE_NAME} &

echo "STARTUP COMPLETED"

while /bin/true; do
  sleep 3600
done
