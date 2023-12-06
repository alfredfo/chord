#!/bin/bash

# Run the first command in the background
go run ./cmd/chord/. &

# Store the process ID (PID) of the first command
pid1=$!
# Run the second command in the background
go run ./cmd/chord/. -ja 0.0.0.0 -jp 1234 -p 1235 &

# Store the PID of the second command
pid2=$!

sleep 10 || wait

kill -9 $pid1 $pid2 > /dev/null 2>&1

echo "Both commands have completed or were terminated."