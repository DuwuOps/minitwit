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

# Use Golang linter locally
You will need to do set up according to your prefrences and eitor. See documentation here: [golangci-lint](https://golangci-lint.run/welcome/integrations/)

## For VSCode
- Install the Go extension (you might need to switch to pre-release version)
- in VSCode `ctrl+shift+p` and enter open user settings (JSON)
- In this insert: 
```JSON
"go.lintTool": "golangci-lint-v2",
"go.lintFlags": [
  "--path-mode=abs",
  "--fast-only"
],
"go.formatTool": "custom",
"go.alternateTools": {
  "customFormatter": "golangci-lint-v2"
},
"go.formatFlags": [
  "fmt",
  "--stdin"
]
```
- Follow one of the many ways to install golangci-lint according to your system. [install](https://golangci-lint.run/welcome/install/)
    - I used `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest` but it is not recomended as mistakes tend to happen.
- Now you should be able to run `golangci-lint run` and see all the things it has found that should be cleaned. 
    - If this doesn't work, but you do have golangci-lint in your go path try: `~/go/bin/golangci-lint-v2 run` 