#!/bin/bash
##
# Script to drop and re-create all tables in prod-database's postgresql server.
# Author: David Martin Sørensen
# Date: 11/04/2025
##

DROP_STR="DROP TABLE IF EXISTS message;\n
DROP TABLE IF EXISTS follower;\n
DROP TABLE IF EXISTS users;"

echo -e $DROP_STR | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"


ENV_DIR=$(dirname "$0")
QUERIES_PATH="${ENV_DIR}/../../src/queries/schema.sql"
echo $QUERIES_PATH

cat $QUERIES_PATH | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"