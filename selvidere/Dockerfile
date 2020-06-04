FROM golang AS builder

ENV GRIDURI http://selenium-hub:4444
ENV INTERVAL 5

WORKDIR /app
COPY selvidere.go go.mod go.sum ./
RUN go get -d
RUN go build

ENTRYPOINT ./selvidere --grid-uri=${GRIDURI} --interval=${INTERVAL}