#!/usr/bin/env sh

echo "run linter and test..." \
  && golangci-lint run --timeout=5m --fix ./... \
  && go test -coverprofile=coverage.out -v -p 1 ./...