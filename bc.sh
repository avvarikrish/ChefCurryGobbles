#!/bin/bash

# rebuild all docker images
docker-compose build --no-cache
if [[ $? -ne 0 ]]; then
    echo "Could not build all images"
    exit 1
fi

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

# push to git
git push -u origin master

echo "Successfully pushed to codebase"