#!/bin/bash
##
# 
# Author: David Martin Sørensen
# Date: 04/04/2025
##

ENV_DIR=$(dirname "$0")
NEWEST_TIMESTAMP_DIR=$(ls $ENV_DIR |  grep -o "[0-9]\+" | tail -1)
cd $ENV_DIR/$NEWEST_TIMESTAMP_DIR
echo "Executing queries in newest timestamped folder: '$NEWEST_TIMESTAMP_DIR'"

execute_all() {
    table_name=$1
    dir_path=$OUTPUT_DIR/$table_name
    file_amount=$(find $dir_path -type f | wc -l)

    echo "Executing the $file_amount files in $dir_path"
    for  (( i=0; i < $file_amount; ++i ))
    do
        # The split command does not simply split the files into 0,1,2,3... Instead it has a weird naming scheme were it goes 9001 after 89
        if [ $((i<10)) -eq 1 ] 
        then
            file_name=0${i}.sql
        elif [ $((99<i)) -eq 1 ] # If i-90 >= 10
        then
            new_i="$(($i-90))"
            file_name=90$new_i.sql
        elif [ $((89<i)) -eq 1 ] # If i-90 < 10
        then
            new_i="$(($i-90))"
            file_name=900$new_i.sql
        else
            file_name=${i}.sql
        fi
        echo "Executing $file_name"

        echo -e "SET client_min_messages TO WARNING; \n$(cat $dir_path/$file_name)" | ssh root@164.90.227.119 "docker exec -i database psql -q -U admin -d minitwit"
    done

    if [ "$table_name" == "users" ]; then
        sqlite3 minitwit.db "select count(*) from user;"
    else
        sqlite3 minitwit.db "select count(*) from $table_name;"
    fi
  
    echo $(echo "select count(*) from $table_name;" | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit --pset=\"footer=off\"" | grep -o "[0-9]\+" )
}


execute_all "users"
execute_all "message"
execute_all "follower"