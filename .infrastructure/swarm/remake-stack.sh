
docker stack remove minitwit
docker rmi -f $(docker images -aq)
docker stack deploy -c docker-compose.yml -c docker-compose.deploy.yml --prune minitwit