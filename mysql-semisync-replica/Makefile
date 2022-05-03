create-network:
	@docker network create mysql-semisync-net

run-master:
	@docker run \
	--detach \
	--platform "linux/x86_64" \
	--name=semisync-master \
	--env="MYSQL_ROOT_PASSWORD=Passw0rd" \
	--publish 6610:3306 \
	--volume=`pwd`/master/conf.d/my-custom.cnf:/etc/mysql/conf.d/mysql.conf.cnf \
	--volume=`pwd`/master/data:/var/lib/mysql/ \
	--volume=`pwd`/master/dump:/root/dump/ \
	--network=mysql-semisync-net \
	mysql:5.7

start-master:
	@docker start semisync-master

kill-master:
	@docker kill --signal=9 semisync-master

stop-master:
	@docker stop semisync-master

console-master:
	@docker exec -it semisync-master bash

start-slave1:
	@docker run \
	--detach \
	--platform "linux/x86_64" \
	--name=semisync-slave1 \
	--env="MYSQL_ROOT_PASSWORD=Passw0rd" \
	--publish 6611:3306 \
	--volume=`pwd`/slave1/conf.d/my-custom.cnf:/etc/mysql/conf.d/mysql.conf.cnf \
	--volume=`pwd`/slave1/data:/var/lib/mysql/ \
	--volume=`pwd`/slave1/dump:/root/dump/ \
	--network=mysql-semisync-net \
	mysql:5.7

stop-slave1:
	@docker stop semisync-slave1

console-slave1:
	@docker exec -it semisync-slave1 bash

start-slave2:
	@docker run \
	--detach \
	--platform "linux/x86_64" \
	--name=semisync-slave2 \
	--env="MYSQL_ROOT_PASSWORD=Passw0rd" \
	--publish 6612:3306 \
	--volume=`pwd`/slave2/conf.d/my-custom.cnf:/etc/mysql/conf.d/mysql.conf.cnf \
	--volume=`pwd`/slave2/data:/var/lib/mysql/ \
	--volume=`pwd`/slave2/dump:/root/dump/ \
	--network=mysql-semisync-net \
	mysql:5.7

stop-slave2:
	@docker stop semisync-slave2

console-slave2:
	@docker exec -it semisync-slave2 bash
