#!/bin/sh
set -e

mkdir -p /home/app/output
cd /home/app/output

git -c core.compression=0 clone --progress --depth 1 --single-branch "$GIT_REPOSITORY_URL" .

/bs
