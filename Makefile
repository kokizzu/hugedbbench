
clean-docker:
	docker rm $(docker ps -a -f status=exited -q)
	docker rm -f $(docker ps -a -q)

clean-2021:
	sudo du -hs $(wildcard 2021/cockroachdb/*/)
	sudo rm -rf $(wildcard 2021/cockroachdb/*/)
	sudo du -hs $(wildcard 2021/tidb/ti*/)
	sudo rm -rf $(wildcard 2021/tidb/ti*/)
	sudo du -hs $(wildcard 2021/yugabytedb/yb*/) 
	sudo rm -rf $(wildcard 2021/yugabytedb/yb*/)

clean-2021mq:
	sudo du -hs $(wildcard 2021mq/kafka/*/)
	sudo rm -rf $(wildcard 2021mq/kafka/*/)
	sudo du -hs $(wildcard 2021mq/nats-jetstream/*/)
	sudo rm -rf $(wildcard 2021mq/nats-jetstream/*/)
	sudo du -hs $(wildcard 2021mq/redpanda/*/)
	sudo rm -rf $(wildcard 2021mq/redpanda/*/)
	sudo du -hs $(wildcard 2021mq/tidbAsQ/ti*/)
	sudo rm -rf $(wildcard 2021mq/tidbAsQ/ti*/)
	sudo du -hs $(wildcard 2021mq/tarantoolAsQ/t*/)
	sudo rm -rf $(wildcard 2021mq/tarantoolAsQ/t*/)
	sudo du -hs $(wildcard 2021mq/clickhouseAsQ/c*/)
	sudo rm -rf $(wildcard 2021mq/clickhouseAsQ/c*/)
	sudo du -hs $(wildcard 2021mq/jetstream/*/)
	sudo rm -rf $(wildcard 2021mq/jetstream/*/)
