#!/bin/bash

set -e

GV="$1"

./hack/generate_group.sh "client,lister,informer" github.com/linuxsuren/kde/pkg/client github.com/linuxsuren/kde/api "${GV}" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"
rm -rf pkg/client
mv github.com/linuxsuren/kde/pkg/client pkg/
rm -rf github.com
