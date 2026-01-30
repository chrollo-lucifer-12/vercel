#!/bin/bash

if [ "$#" -ne 2 ]; then
  echo "Usage: $(basename "$0") <git-repo-url> <branch>"
  exit 1
fi

REPO_URL="$1"
BRANCH="$2"


REPO_NAME=$(basename "$REPO_URL" .git)
WORKDIR="/tmp/git-checks"
REPO_DIR="$WORKDIR/$REPO_NAME"

mkdir -p "$WORKDIR"

if [ ! -d "$REPO_DIR/.git" ]; then
  git clone -b "$BRANCH" "$REPO_URL" "$REPO_DIR" > /dev/null 2>&1
fi

cd "$REPO_DIR" || exit 1


git fetch origin "$BRANCH" > /dev/null 2>&1

LOCAL_HASH=$(git rev-parse "$BRANCH")
REMOTE_HASH=$(git rev-parse "origin/$BRANCH")

if [ "$LOCAL_HASH" = "$REMOTE_HASH" ]; then
  echo "Branch unchanged"
else
  echo "Branch changed"
fi
