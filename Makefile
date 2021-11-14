

clean:
	sudo du -hs --max-depth 0 $(wildcard 2021/cockroachdb/*/)
	sudo rm -rf $(wildcard 2021/cockroachdb/*/)
	sudo du -hs --max-depth 0 $(wildcard 2021/tidb/ti*/)
	sudo rm -rf $(wildcard 2021/tidb/ti*/)
	sudo du -hs --max-depth 0 $(wildcard 2021/yugabytedb/yb*/) 
	sudo rm -rf $(wildcard 2021/yugabytedb/yb*/)
