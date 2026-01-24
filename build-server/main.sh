#!/bin/sh
set -e

mkdir -p /home/app/output
cd /home/app/output

git clone "$GIT_REPOSITORY_URL" .

/bs
