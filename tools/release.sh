#!/usr/bin/env bash

mkdir build/
go build -o build/streaming-management main.go
cp -r obs-assets/ build/obs-assets
zip --junk-paths streaming-management-$RUNNER_OS build/
