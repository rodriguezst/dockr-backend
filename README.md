# dockr-backend
Simple backend to start/stop docker containers

docker run -v /var/run/docker.sock:/var/run/docker.sock -p {dest_port}:8080 {image}
