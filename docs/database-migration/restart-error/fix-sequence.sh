#!/bin/bash
##
# 
# Author: David Martin Sørensen
# Date: 26/04/2025
##

# Helper function(s)
run_sql_query() {
    echo $1 | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"
}

echo -e "\n #### FIXING MESSAGE ####"
# Set the value of the current serial sequence of message at message_id to the current highest message_id +1
run_sql_query "SELECT setval(pg_get_serial_sequence('message', 'message_id'), MAX(message_id)+1) FROM message;"


echo -e "\n #### FIXING USERS ####"
# Set the value of the current serial sequence of users at user_id to the current highest user_id +1
run_sql_query "SELECT setval(pg_get_serial_sequence('users', 'user_id'), MAX(user_id)+1) FROM users;"
