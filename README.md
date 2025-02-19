# Introduction

This repository contains group c / "DuwuOps"'s official DevOps project.

# Run via Docker
## Using docker compose
Build image and start container: 
```
docker compose up
```
Stop and destroy container:
```
docker compose down
```

## Manual start and stop
start:
```
docker build --tag minitwit_image .
docker run --publish 8000:8000 --name minitwit_app minitwit_image 
```
Now go to http://localhost:8000/

stop:
```
docker stop minitwit_app
```

# Cleanup

Remove image. It is not always necessary to use -f it is only when you want to force deletion:
```
docker image rm minitwit_image -f
```
Remove container:
```
docker container rm minitwit_app
```
Clean local images during ``docker compose down``:
```
docker compose down -v --rmi local
```
# Additional information
See your images:
```
docker images
```
See containers:
```
docker ps --all
```
Run docker compose in detached mode:
```
docker compose up -d
```