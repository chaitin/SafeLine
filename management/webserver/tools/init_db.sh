docker rm -f mgt-postgres-dev

#docker run -it -d -e POSTGRES_DB=safeline-ce -e POSTGRES_USER=safeline-ce -e POSTGRES_PASSWORD=safeline-ce -p 127.0.0.1:5432:5432 --name mgt-postgres-dev chaitin.cn/library/postgres:15.2
docker run -it -d -e POSTGRES_DB=safeline-ce -e POSTGRES_USER=safeline-ce -e POSTGRES_PASSWORD=safeline-ce -p 127.0.0.1:5432:5432 --name mgt-postgres-dev postgres:15.2