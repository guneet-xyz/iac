#!/bin/sh

cd app

gofmt -l .
go vet .
