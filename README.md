# Introduction

This repository contains group c / "DuwuOps"'s official DevOps project.

# Run via Docker
## Not using docker compose
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

## Using docker compose
TODO

# Cleanup

remove image: 

It is not always necessary to use -f it is only when you want to force deletion.
```
docker image rm minitwit_image -f
```
remove container:
```
docker container rm minitwit_app
```

## Get an overview
See your images:
```
docker images
```
See containers:
```
docker ps --all
```