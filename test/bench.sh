#!/bin/bash
echo "generating example data"
python3 gen.py

echo "building executable"
rm ./test
go build ./test.go

hyperfine "./test ./1MB.json" "./test -libjson=false ./1MB.json"
hyperfine "./test ./5MB.json" "./test -libjson=false ./5MB.json"
hyperfine "./test ./10MB.json" "./test -libjson=false ./10MB.json"
