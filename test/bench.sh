#!/bin/bash
echo "fetching example data"
curl https://microsoftedge.github.io/Demos/json-dummy-data/64KB.json -C - -o 64KB.json 2>> /dev/null
curl https://microsoftedge.github.io/Demos/json-dummy-data/128KB.json -C - -o 128KB.json 2>> /dev/null
curl https://microsoftedge.github.io/Demos/json-dummy-data/256KB.json -C - -o 256KB.json 2>> /dev/null
curl https://microsoftedge.github.io/Demos/json-dummy-data/512KB.json -C - -o 512KB.json 2>> /dev/null
curl https://microsoftedge.github.io/Demos/json-dummy-data/1MB.json -C - -o 1MB.json 2>> /dev/null
curl https://microsoftedge.github.io/Demos/json-dummy-data/5MB.json -C - -o 5MB.json 2>> /dev/null

echo "building executable"
go build test.go

echo -n "[libjson] 64KB: "
./test ./64KB.json
echo -n "[libjson] 128KB: "
./test ./128KB.json
echo -n "[libjson] 256KB: "
./test ./256KB.json
echo -n "[libjson] 512KB: "
./test ./512KB.json
echo -n "[libjson] 1MB: "
./test ./1MB.json
echo -n "[libjson] 5MB: "
./test ./5MB.json

echo -n "[gojson] 64KB: "
./test -libjson=false ./64KB.json
echo -n "[gojson] 128KB: "
./test -libjson=false ./128KB.json
echo -n "[gojson] 256KB: "
./test -libjson=false ./256KB.json
echo -n "[gojson] 512KB: "
./test -libjson=false ./512KB.json
echo -n "[gojson] 1MB: "
./test -libjson=false ./1MB.json
echo -n "[gojson] 5MB: "
./test -libjson=false ./5MB.json
