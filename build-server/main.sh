#!/bin/sh
set -e


git clone "$GIT_REPOSITORY_URL" /home/app/output

exec /home/app/app
