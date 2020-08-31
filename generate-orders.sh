#!/bin/bash

for ((c=0; c<$(jot -r 1 5 10); c++))
do
	go run ccgobbles_client/ccgobbles/ccgobbles_client.go &
done
wait
