#!/bin/sh
cd "$(dirname "$0")" || exit 1
exec ./bin/cinebite -server -port 127.0.0.1:8080
