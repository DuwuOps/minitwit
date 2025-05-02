#!/bin/bash
##
# Script to drop and re-create all tables in prod-database's postgresql server.
# Author: David Martin Sørensen
# Date: 11/04/2025
##

drop_table() {
    table_name=$1
    
    echo # New line
    echo "Dropping table $table_name..."

    # DROP
    DROP_STR="DROP TABLE IF EXISTS $table_name;"
    echo -e $DROP_STR | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"
}

drop_table "follower"
drop_table "message"
drop_table "users"
drop_table "latest_processed"

ENV_DIR=$(dirname "$0")
QUERIES_DIR="${ENV_DIR}/../../src/queries"

setup_table() {
    table_name=$1

    echo # New line
    echo "Setting up table $table_name..."

    # CREATE
    cat "$QUERIES_DIR/schema.$table_name.sql" | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"
}

setup_table "users"
setup_table "follower"
setup_table "message"
setup_table "latest_processed"