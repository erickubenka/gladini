# Selex

`Selex`, a simple abbreviation for "Selenium exporter".  
Selex will use the Selenium Grid Hub API, provided under http://selenium-hub:NODE_Port/status, and is implemented in a simple `Golang` HTTP Listener. To parse the JSON into objects, Selex provide some structs. 

> !!! - IMPORTANT! Selex must be build on your local machine, unless the container is `NOT` deployed into a public registry.

## Build

Run these commands to build the Docker container locallly, if you want to:
````bash
docker build --tag selex:1.0 .
docker run --publish 8080:8080 --detach --name selex selex:1.0 
# Verify it works
docker stop selex
docker rm selex
```` 