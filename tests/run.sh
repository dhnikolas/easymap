#!/bin/bash
echo 'package dynamic_test' > tests/dynamic_test.go
echo 'import "github.com/dhnikolas/easymap/tests/testdata/in"' >> tests/dynamic_test.go
echo 'import "testing"' >> tests/dynamic_test.go
echo 'import "reflect"' >> tests/dynamic_test.go
go run cmd/app/main.go copygen ./tests/testdata/in/in.go:Source Out >> tests/dynamic_test.go
cat ./tests/test_copygen.go.test >> tests/dynamic_test.go
go test ./tests/dynamic_test.go
rm ./tests/dynamic_test.go
