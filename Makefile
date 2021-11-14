

clean:
	sudo rm -rf $(wildcard 2021/cockroachdb/*1/)
	sudo rm -rf $(wildcard 2021/tidb/t*/) 2021/tidb/logs
	sudo rm -rf $(wildcard 2021/yugabytedb/yb*1/)
