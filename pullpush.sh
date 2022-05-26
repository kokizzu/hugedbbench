#!/usr/bin/env bash

set -e # exit on error

if [ $# -eq 0 ] ; then
  echo "Usage: 
  ./pullpush.sh 'the commit message'"
  exit
fi

# format indentation
goimports -w **/*.go
echo "codes formatted.."

# update deps
go get -u -v github.com/kokizzu/gotro@latest
go mod tidy

# add and commit all files
git add .
git status
read -p "Press Ctrl+C to exit, press any enter key to check the diff..
"

# recheck again
git diff --staged
echo 'Going to commit with message: '\"$*\"
read -p "Press Ctrl+C to exit, press any enter key to really commit..
"

git commit -m "$*" && git pull && git push origin master
