#!/usr/bin/env bash

mkdir build/
go build -o build/streaming-management main.go
cp -r obs-assets/ build/obs-assets
cd build/
zip -r ../streaming-management-$RUNNER_OS ./*
