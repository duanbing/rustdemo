#!/bin/bash

go build -buildmode=c-archive -o libmath.a math.go
