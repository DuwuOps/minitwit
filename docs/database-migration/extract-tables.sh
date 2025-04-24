#!/bin/bash
##
# SQlite-database extraction script for DuwuOps/minitwit
# Author: David Martin Sørensen
# Date: 04/04/2025
##


# Make a new timestamped directory to work from

ENV_DIR=$(dirname "$0")
cd $ENV_DIR

TIMESTAMP=$(date +"%s")
mkdir $TIMESTAMP
cd $TIMESTAMP

# Copy database-file from prod-web-server


DATABASE_FILE=minitwit.db

echo "Copying $DATABASE_FILE to $ENV_DIR/$TIMESTAMP/$DATABASE_FILE"

scp root@134.209.137.191:/var/lib/docker/volumes/minitwit_sqliteDB/_data//minitwit.db $DATABASE_FILE