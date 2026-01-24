#!/bin/sh

git clone "$GIT_REPOSITORY_URL" /home/app/output

exec go run src/main.go
