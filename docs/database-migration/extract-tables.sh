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
mkdir $TIMESTAMP && cd "$_"

# Copy database-file from prod-web-server


DATABASE_FILE=minitwit.db

echo "Copying $DATABASE_FILE to $ENV_DIR/$TIMESTAMP/$DATABASE_FILE"

scp root@134.209.137.191:/var/lib/docker/volumes/minitwit_sqliteDB/_data//minitwit.db $DATABASE_FILE


# Create data-dump files from local copy of database-file

if [[ ! -f $DATABASE_FILE ]] ; then
    echo "File '$DATABASE_FILE' is not here, aborting."
    exit
fi

# Create data-dump files from local copy of database-file
mkdir queries && cd "$_"

dump_table_data() {
    table_name=$1
    output_table=$table_name
    if [ "$table_name" == "user" ]; then
        output_table="${output_table}s"
    fi
    sqlite3 ../$DATABASE_FILE ".schema '$table_name'" > schema.$output_table.sql
    sqlite3 ../$DATABASE_FILE ".dump '$table_name'" > dump.$output_table.sql
    grep -vxF -f schema.$output_table.sql dump.$output_table.sql > data.$output_table.sql
    rm -f schema.$output_table.sql
    rm -f dump.$output_table.sql
    sed -i -E '/PRAGMA foreign_keys=OFF;/d' data.$output_table.sql
    sed -i -E '/BEGIN TRANSACTION;/d' data.$output_table.sql
    sed -i -E '/COMMIT;/d' data.$output_table.sql
    sed -i -E '/\/****** CORRUPTION ERROR *******\//d' data.$output_table.sql
    if [ "$table_name" == "user" ]; then
        sed -i -E "s/INSERT INTO $table_name VALUES/INSERT INTO $output_table VALUES/" data.$output_table.sql
        sed -i "1s/^/INSERT INTO users VALUES(0,'Uknown','unknown@unknown.com','\$2a\$10\$r1BVr0BDJVQd0f7NoN.p97jSBzyv67V97LrlbX9pRo.5mrYFFJG7T');\n/" data.$output_table.sql # To not lose messages and followings which use the user with user-id 0
    fi

    line_amount=$(wc -l < data.$output_table.sql)
    echo "Extracted $line_amount lines from $output_table in $DATABASE_FILE"
}

dump_table_data "user"
dump_table_data "message"
dump_table_data "follower"


# Remove previously added items
NEWEST_TIMESTAMP_DIR=$(ls "../../" | sed "/$TIMESTAMP/d" | grep -o "[0-9]\+" | tail -1)
NEWEST_QUERIES_PATH="../../$NEWEST_TIMESTAMP_DIR/queries"

filter() {
    table_name=$1
    dump_file=data.$table_name.sql
    if [ "$NEWEST_TIMESTAMP_DIR" == "" ]; then
        echo "No previous timestamped folders found for filtering."
        cat $dump_file > filtered_$dump_file
    else
        echo "  Filtering $dump_file from $NEWEST_QUERIES_PATH/$dump_file"
        
        # Sort the files
        sort "$dump_file" > $dump_file.new
        sort "$NEWEST_QUERIES_PATH/$dump_file" > $dump_file.old

        filtered_file=filtered_$dump_file
        comm -23 $dump_file.new $dump_file.old > $filtered_file # comm -23 means compare file 1 and 2 and only show lines unique to file 1.
        rm $dump_file.new
        rm $dump_file.old

        line_amount=$(wc -l < $filtered_file)
        echo "$line_amount lines occoured in $TIMESTAMP/data.$table_name.sql that did not in $NEWEST_TIMESTAMP_DIR/data.$table_name.sql"
    fi
}

filter "users"
filter "message"
filter "follower"


# Split sql-query-files into files of maximum 20000 lines each
mkdir split && cd "$_"


split_dump() {
    table_name=$1
    file_name=filtered_data.$table_name.sql
    mkdir $table_name
    split -dl 20000 --additional-suffix=.sql ../$file_name $table_name/
    echo "$file_name has been split into $(find $table_name -type f | wc -l) files"
}


split_dump "users"
split_dump "follower"
split_dump "message"