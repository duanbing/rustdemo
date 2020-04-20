#!/bin/bash

(
cd server
go build server.go
)

(
cd client
go build client.go
)
