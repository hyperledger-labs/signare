#!/usr/bin/env bash
set -e

CONTAINER_NAME="postgres_10_12"

# Container name expected as an env variable
if [[ "${SIGNARE_POSTGRES_CONTAINER_NAME}" != "" ]]
then
	CONTAINER_NAME="${SIGNARE_POSTGRES_CONTAINER_NAME}"
fi

docker cp create-databases.sql "${CONTAINER_NAME}":/
docker exec -e PGPASSWORD=postgres -it "${CONTAINER_NAME}" psql -U postgres -a -f /create-databases.sql
