### D7024E Mobile and Distributed Systems

Team 8
Members: Chen Xinyu , Brandon Chong

## sprint0 commands:

# To run the docker containers

docker-compose up -d

# This command executes 'ping -c 5 node42' inside the node1 container

docker-compose exec node1 ping -c 5 node42

# To clean up and remove all the containers and the network, run the command:

docker-compose down
