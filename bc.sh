#!/bin/bash

# check if builds properly
docker-compose up -d
if [[ $? -ne 0 ]]; then
    echo "Could not build all containers"
    exit 1
fi

# check all tests
go test ./...
if [[ $? -ne 0 ]]; then
    echo "Failed test"
    exit 1
fi

# Cleanup
docker-compose down
docker rmi chefcurrgobbles_ccgobbles_server
docker rmi chefcurrgobbles_metrics_server

# push to git
git push -u origin head:master

echo "Successfully pushed to codebase"