FROM golang:1.23.6-alpine

#Claude ai helped us with this line
RUN apk add --no-cache git sqlite gcc musl-dev 

WORKDIR /app
COPY go.mod go.sum ./
COPY *.go ./
COPY ./templates /app/templates
COPY ./static /app/static
COPY ./queries /app/queries

RUN go mod download

RUN go build -o minitwit ./minitwit.go

EXPOSE 8000

CMD ["./minitwit"]