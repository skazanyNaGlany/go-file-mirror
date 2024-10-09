#!/bin/bash

go test -coverprofile=coverage.out -coverpkg=.. && go tool cover -html=coverage.out

