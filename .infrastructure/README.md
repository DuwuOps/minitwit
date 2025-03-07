To deploy

`docker compose down -v --rmi local` 

`docker compose up -d --build` 

`docker container ls`

`docker commit CONTAINERID tingariussorensen/minitwit:latest`

`docker push tingariussorensen/minitwit:latest`

`ssh root@159.223.8.210`

`docker container ls`

`docker stop CONTAINERID`

`docker rm CONTAINERID`

`docker image pull tingariussorensen/minitwit:latest`

`sudo docker run -d -p 0.0.0.0:80:8000 --restart=always -v /var/minitwit:/app/tmp tingariussorensen/minitwit:latest`
