#!/bin/bash
echo "generating example data"
python3 gen.py

echo "building executable"
rm ./test
go build -o ./test ../cmd/lj.go

hyperfine "./test -s ./1MB.json" "./test -s -libjson=false ./1MB.json"
hyperfine "./test -s ./5MB.json" "./test -s -libjson=false ./5MB.json"
hyperfine "./test -s ./10MB.json" "./test -s -libjson=false ./10MB.json"
