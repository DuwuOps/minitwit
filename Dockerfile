FROM golang:1.23.6-alpine

# Claude AI helped us with this line.
RUN apk add --no-cache git sqlite gcc musl-dev 

# Create the proper module structure
WORKDIR /minitwit

# Copy everything from src to the root of the module
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./src/*.go ./src/

COPY ./src/datalayer ./src/datalayer
COPY ./src/handlers ./src/handlers
COPY ./src/models ./src/models
COPY ./src/routes ./src/routes
COPY ./src/template_rendering ./src/template_rendering

RUN go build -o minitwit ./src/main.go

RUN rm -rf ./src/
RUN rm go.mod
RUN rm go.sum

COPY ./src/templates ./templates
COPY ./src/static ./static
COPY ./src/queries ./queries


EXPOSE 8000
CMD ["./minitwit"]
