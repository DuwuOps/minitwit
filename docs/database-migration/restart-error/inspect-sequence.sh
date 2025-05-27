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

# Inspect the serial sequence of message at message_id
echo -e "\n #### INSPECTING MESSAGE ####"
echo -e "\n # Current serial sequence value for message:"
run_sql_query "SELECT s.* FROM message_message_id_seq s;"

echo -e "\n # Current max message_id in message:"
run_sql_query "SELECT MAX(message_id) FROM message;"


echo -e "\n #### INSPECTING USERS ####"
echo -e "\n # Current serial sequence value for users:"
run_sql_query "SELECT s.* FROM users_user_id_seq s;"

echo -e "\n # Current max user_id in users:"
run_sql_query "SELECT MAX(user_id) FROM users;"