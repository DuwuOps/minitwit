#################### BUILD BINARY ####################

FROM golang:1.23.6-alpine AS builder

# Install build-dependencies for virtual machine
RUN apk add --no-cache git sqlite gcc musl-dev

WORKDIR /build

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



####################  RUN BINARY  ####################

FROM alpine:3.18

# Install run-dependencies for virtual machine
RUN apk add --no-cache sqlite

WORKDIR /minitwit

# Copy binary from build-phase
COPY --from=builder /build/minitwit .

# Copy non-source-code files
COPY ./src/templates ./templates
COPY ./src/static ./static
COPY ./src/queries ./queries

# Expose port and run binary-file
EXPOSE 8000
CMD ["./minitwit"]
