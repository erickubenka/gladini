# Selvidere

`Selvidere` is a simple watch service, that will remove old and unhealthy Selenium nodes automatically from Selenium Hub.
There are currently two mechanisms.

1. Remove unhealthy Nodes
2. Remove duplicates - Sometimes Kubernetes will restart a Selenium node port. When this happens, the node registers with same UUID on Hub.  
Currently Selenium do not handle these duplicates, so `Selvidere` will do.

## Build

Run these commands to build the Docker container locallly, if you want to:
````bash
docker build --tag selvidere:1.0 .
docker run --detach --name selvidere selvidere:1.0 
# Verify it works
docker stop selvidere
docker rm selvidere
```` 