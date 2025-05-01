# To see the amount of current connections:
SQL_TEXT="SELECT COUNT(*) FROM pg_stat_activity;"
echo $SQL_TEXT | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"

# To see the details of the current connections:
SQL_TEXT="SELECT sa.* FROM pg_stat_activity sa;"
echo $SQL_TEXT | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"

# To see the current maximum connections:
SQL_TEXT="SHOW max_connections;"
echo $SQL_TEXT | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"

# To see the maximum allowed value for the "current maximum connections"-setting:
SQL_TEXT="SELECT min_val, max_val FROM pg_settings WHERE NAME='max_connections';"
echo $SQL_TEXT | ssh root@164.90.227.119 "docker exec -i database psql -U admin -d minitwit"