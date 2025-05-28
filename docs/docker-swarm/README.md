This is a guide of how to set up and update a docker swarm.

## Preparation-steps

1. Have 3 Digital Ocean droplets with Docker installed.
2. SSH into the droplet which should become the Manager-node, and execute the following command. (Note: Replace `$MANAGER_IP` with the droplet's IP).
```bash
docker swarm init --advertise-addr $MANAGER_IP
```
4. Copy the token-command outputted by the previous step. It should be something akin to the following:
```bash
docker swarm join --token LONG_TOKEN_TEXT $MANAGER_IP:2377
```
5. SSH into the other droplets, which should become the Worker-nodes, and execute the copied command.
6. Close the SSH-connection to the worker droplets. You do not ever have to do anything there again.
8. Copy and paste the docker-compose files into the Manager-droplet:
```bash
ssh root@$MANAGER_IP "mkdir ~/.deploy"
scp docker-compose.yml root@$MANAGER_IP:~/.deploy/docker-compose.yml
scp docker-compose.deploy.yml root@$MANAGER_IP:~/.deploy/docker-compose.deploy.yml
```
7. In a terminal with SSH-acces to Manager-droplet, with the correct env-variables saved, run the following command:
```shell
docker stack deploy -c ~/.deploy/docker-compose.yml -c ~/.deploy/docker-compose.deploy.yml --prune minitwit
```
Done! You now have this project running on a Swarm!


## How to remove the previous containers, if they are running.

If you need to remove the previous containers made through docker-compose, run the following command in a temrinal with SSH-connection,  and with the correct env-variables saved:
```shell
docker compose -f docker-compose.yml -f docker-compose.deploy.yml down --rmi local
```