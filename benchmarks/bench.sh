#!/bin/bash
echo "generating example data"
python3 gen.py

echo "building executable"
rm ./test
go build -o ./test ../cmd/lj.go

for SIZE in 1MB 5MB 10MB 100MB; do
    hyperfine \
        --warmup 1 \
        --runs 10 \
        "./test -s ./${SIZE}.json" \
        "./test -s -libjson=false ./${SIZE}.json"
done
