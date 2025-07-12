#!/bin/bash

# Get max open files system-wide and per-process (static, so fetch once)
MAX_FILES_SYS=$(sysctl -n kern.maxfiles)
MAX_FILES_PROC=$(sysctl -n kern.maxfilesperproc)

# Get your app's PID (update this line as needed)
APP_PID=$(pgrep -f redis-document-api | head -n1)

while true; do
  echo "------ $(date) ------"
  # System-wide open files
  CUR_FILES_SYS=$(sudo lsof | wc -l)
  echo "Open files (system-wide): $CUR_FILES_SYS / $MAX_FILES_SYS"

  # Per-process open files
  if [ -n "$APP_PID" ]; then
    CUR_FILES_PROC=$(lsof -p $APP_PID 2>/dev/null | wc -l)
    echo "Open files (process $APP_PID): $CUR_FILES_PROC / $MAX_FILES_PROC"
  else
    echo "Open files (process): PID not found"
  fi

  # Sockets to Redis (default port 6379)
  CUR_SOCKETS_REDIS=$(netstat -an | grep 6379 | wc -l)
  echo "Open sockets to Redis: $CUR_SOCKETS_REDIS"

  # Sockets in TIME_WAIT
  CUR_TIMEWAIT=$(netstat -an | grep TIME_WAIT | wc -l)
  echo "Sockets in TIME_WAIT: $CUR_TIMEWAIT"

  echo "-----------------------------"
  sleep 1
done
