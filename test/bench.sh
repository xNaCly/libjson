#!/bin/bash
echo "fetching example data"
curl https://microsoftedge.github.io/Demos/json-dummy-data/1MB.json -C - -o 1MB.json 2>> /dev/null
curl https://microsoftedge.github.io/Demos/json-dummy-data/5MB.json -C - -o 5MB.json 2>> /dev/null

echo "building executable"
go build test.go

hyperfine "./test ./1MB.json" "./test -libjson=false ./1MB.json"
hyperfine "./test ./5MB.json" "./test -libjson=false ./5MB.json"
