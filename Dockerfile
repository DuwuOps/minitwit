FROM golang:1.23.6-alpine

# Install build dependencies for virtual machine
# Claude AI helped us with this line.
RUN apk add --no-cache git sqlite gcc musl-dev 

# Create the proper module structure
WORKDIR /minitwit

# Download module dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy source code & build binary
COPY ./src/*.go ./src/
COPY ./src/datalayer ./src/datalayer
COPY ./src/handlers ./src/handlers
COPY ./src/models ./src/models
COPY ./src/routes ./src/routes
COPY ./src/template_rendering ./src/template_rendering

RUN go build -o minitwit ./src/main.go

# Delete source code from container
RUN rm -rf ./src/
RUN rm go.mod
RUN rm go.sum

# Copy non-source-code files 
COPY ./src/templates ./templates
COPY ./src/static ./static
COPY ./src/queries ./queries


# Expose port and run binary-file
EXPOSE 8000
CMD ["./minitwit"]
