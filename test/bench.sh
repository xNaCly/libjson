#!/bin/bash
echo "fetching example data"
curl -C https://microsoftedge.github.io/Demos/json-dummy-data/64KB.json -o 64KB.json 2>> /dev/null
curl -C https://microsoftedge.github.io/Demos/json-dummy-data/128KB.json -o 128KB.json 2>> /dev/null
curl -C https://microsoftedge.github.io/Demos/json-dummy-data/256KB.json -o 256KB.json 2>> /dev/null
curl -C https://microsoftedge.github.io/Demos/json-dummy-data/512KB.json -o 512KB.json 2>> /dev/null
# curl https://microsoftedge.github.io/Demos/json-dummy-data/1MB.json -o 1MB.json 2>> /dev/null

echo "[libjson] building executable"
go build libjson.go

echo -n "[libjson] 64KB: "
./libjson ./64KB.json
echo -n "[libjson] 128KB: "
./libjson ./128KB.json
echo -n "[libjson] 256KB: "
./libjson ./256KB.json
echo -n "[libjson] 512KB: "
./libjson ./512KB.json

echo "[gojson] building executable"
go build gojson.go

echo -n "[gojson] 64KB: "
./gojson ./64KB.json
echo -n "[gojson] 128KB: "
./gojson ./128KB.json
echo -n "[gojson] 256KB: "
./gojson ./256KB.json
echo -n "[gojson] 512KB: "
./gojson ./512KB.json
