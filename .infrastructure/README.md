To deploy

`docker compose down -v --rmi local` 

`docker compose up -d --build` 

`docker container ls`

`docker commit ${CONTAINER_ID} tingariussorensen/minitwit:latest`

`docker push tingariussorensen/minitwit:latest`

`ssh root@159.223.8.210`

`docker container ls`

`docker stop ${CONTAINER_ID}`

`docker rm ${CONTAINER_ID}`

`docker image pull tingariussorensen/minitwit:latest`

`sudo docker run -d -p 0.0.0.0:80:8000 --restart=always -v sqliteDB:/minitwit/tmp tingariussorensen/minitwit:latest`

`export DATABASE_FILE_PATH="/tmp/minitwit.go`

`export LATESTPROCESSED_PATH="/tmp/latest_processed_sim_action_id.txt`

`docker cp ${DATABASE_FILE_PATH} ${CONTAINER_ID}:/minitwit/tmp/`

`docker cp ${LATESTPROCESSED_PATH} ${CONTAINER_ID}:/minitwit/`

