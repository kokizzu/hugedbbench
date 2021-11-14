

clean:
	sudo rm -rf $(wildcard 2021/cockroachdb/*/)
	sudo rm -rf $(wildcard 2021/tidb/t*/) 2021/tidb/logs
	sudo rm -rf $(wildcard 2021/yugabytedb/yb*/)
