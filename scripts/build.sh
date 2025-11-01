#!/bin/bash
GOOS=wasip1 GOARCH=wasm go build -o ../dist/plugin-morphe-py-types-v1.0.0.wasm ../cmd/plugin/main.go
