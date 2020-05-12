FROM golang AS builder

ENV GRIDURI http://selenium-hub:4444

WORKDIR /app
COPY selex.go go.mod go.sum ./
RUN go get -d
RUN go build

EXPOSE 8080

ENTRYPOINT ./selex --grid-uri=${GRIDURI}