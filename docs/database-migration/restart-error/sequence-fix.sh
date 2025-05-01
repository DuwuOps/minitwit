

# Get the current serial sequence of message at message_id
SQL_TEXT="SELECT pg_get_serial_sequence('message', 'message_id');"
echo $SQL_TEXT | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"

# Inspect the serial sequence of message at message_id
SQL_TEXT="SELECT s.* FROM message_message_id_seq s;"
echo $SQL_TEXT | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"

# Set the value of the current serial sequence of message at message_id to the current highest message_id +1
SQL_TEXT="SELECT setval(pg_get_serial_sequence('message', 'message_id'), MAX(message_id)+1) FROM message;"
echo $SQL_TEXT | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"