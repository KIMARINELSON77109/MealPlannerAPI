#!/bin/bash
# export GOPATH=$HOME/git/MealPlannerAPI
export GOPATH=$PWD
mkdir bin
go get  github.com/go-sql-driver/mysql
go get  github.com/gorilla/mux
go install github.com/go-sql-driver/mysql
go install github.com/gorilla/mux
rm bin/application

go clean main

go build -o bin/application main