# goCrudServer

To run code tou need to run DB on a container 

run those commands:

docker pull postgres
docker run -d --name postgresCont -p 5433:5432 -e POSTGRES_PASSWORD=pass123 postgres
