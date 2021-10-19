# D7024E Kademlia_lab

# Setup
## Install dependencies
The following dependencies is required to run this program:
- Docker
- Golang

## Step 1: Build the docker image
To build the docker image, run the following commands from the root directory of this project:
```
docker build -t kadlab .
```

## Step 2: Run the containers with Docker Compose
Modify the number of nodes in the docker-compose file.
```
replicas: 50
```
Run this command from the root directory of the project to start the containers:
```
docker-compose up
```

To terminate and remove the containers, run command:
```
docker-compose down
```

## Step 3: Start another command prompt from the root directory
Run this command to attach a containers ongoing input and output:
```
docker attach <container_name>
```  

## Step 4: Execute CLI commands
The CLI has three main functions. Available commands:
```
- put \<string> (stores data on k closest nodes to hash)	
- get \<string> (fetches data object if it is stored in the network)
- exit (terminates this node)

```