#!/bin/bash

set -e
trap "docker-compose -f ./deployments/docker-compose.int-tests.yml down" EXIT
docker-compose -f ./deployments/docker-compose.int-tests.yml up --build -d
test_status_code=0
docker-compose -f ./deployments/docker-compose.int-tests.yml run integration-tests  || test_status_code=$?
exit $test_status_code