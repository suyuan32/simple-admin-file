#!/bin/sh

nginx -c  /etc/nginx/nginx.conf &

/app/file_api -f /app/etc/file.yaml