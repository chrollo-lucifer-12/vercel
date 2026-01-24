#!/bin/sh

git clone "$GIT_REPOSITORY_URL" /home/app/output

exec node index.js
