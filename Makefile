

clean:
	sudo du -hs $(wildcard 2021/cockroachdb/*/)
	sudo rm -rf $(wildcard 2021/cockroachdb/*/)
	sudo du -hs $(wildcard 2021/tidb/ti*/)
	sudo rm -rf $(wildcard 2021/tidb/ti*/)
	sudo du -hs $(wildcard 2021/yugabytedb/yb*/) 
	sudo rm -rf $(wildcard 2021/yugabytedb/yb*/)
