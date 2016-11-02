#!/bin/bash

PID_LIST=""

consul agent -dev -advertise=127.0.0.1 & pid=$! 
PID_LIST+=" $pid"


DIR=./db/dummy-db
go run $DIR/main.go $DIR/todo_list.go $DIR/user_profile.go & pid=$!
PID_LIST+=" $pid"

DIR=./auth/jwt
go run $DIR/main.go $DIR/token.go $DIR/user.go & pid=$!
PID_LIST+=" $pid"

DIR=./todo/simple
go run $DIR/main.go $DIR/todo.go & pid=$!
PID_LIST+=" $pid"

DIR=./router
go run $DIR/main.go $DIR/login.go $DIR/profile.go $DIR/todo.go & pid=$!
PID_LIST+=" $pid"

trap "kill $PID_LIST" SIGINT

wait $PID_LIST
