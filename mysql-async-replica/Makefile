create-network:
	@docker network create mysql-async-net

start-master:
	@docker run \
	--detach \
	--platform "linux/x86_64" \
	--name=async-master \
	--env="MYSQL_ROOT_PASSWORD=Passw0rd" \
	--publish 6603:3306 \
	--volume=`pwd`/master/conf.d/my-custom.cnf:/etc/mysql/conf.d/mysql.conf.cnf \
	--volume=`pwd`/master/data:/var/lib/mysql/ \
	--volume=`pwd`/master/dump:/root/dump/ \
	--network=mysql-async-net \
	mysql:5.7

stop-master:
	@docker stop async-master
	@docker rm async-master

console-master:
	@docker exec -it async-master bash

start-slave:
	@docker run \
	--detach \
	--platform "linux/x86_64" \
	--name=async-slave \
	--env="MYSQL_ROOT_PASSWORD=Passw0rd" \
	--publish 6604:3306 \
	--volume=`pwd`/slave/conf.d/my-custom.cnf:/etc/mysql/conf.d/mysql.conf.cnf \
	--volume=`pwd`/slave/data:/var/lib/mysql/ \
	--volume=`pwd`/slave/dump:/root/dump/ \
	--network=mysql-async-net \
	mysql:5.7

stop-slave:
	@docker stop async-slave
	@docker rm async-slave

console-slave:
	@docker exec -it async-slave bash
