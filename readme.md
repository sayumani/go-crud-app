[![Build Status](https://travis-ci.com/sayumani/go-crud-app.svg?branch=master)](https://travis-ci.com/sayumani/go-crud-app)

# Items App

# Steps To Follow to Configure The Project

- Install go lang in your machine https://golang.org/dl/
- Set environment variable GOPATH to your go workspace folder
- Install vs code as IDE
- Install docker

# Step to run in docker

- Set DB_HOST=products-postgres in .env
- Run docker-compose up

# Step to configure postgres admin and creating table

- Run docker container ls
- Note down postgres container id
- Run docker inspect <postgres container id> | grep IPAddress and notedown IP address displayed.This is needed to configure postgres admin
- Open localhost:5050
- Login with credentails in .env file
- Click add server
- Provide the server name
- Provide host name as the IP address in step 3
- Provide port as 5432
- Click add
- Create database and table as in query.sql

# Step to stop and remove containers in docker

- docker-compose down

# Step to rebuild and start in docker

- docker-compose up --build

# Steps To Run locally

- Set DB_HOST=localhost in .env
- Run go run main.go

# Steps to run db migration

- download [goose](https://github.com/letsencrypt/goose)
- add $GOPATH/bin to path variable
- run goose -env=<envirmonent> up
