build:
	go build -o ./bin/cfm cmd/cli-file-manager/main.go
run: build
	./bin/cfm
clean:
	rm -rf bin/
update:
	git pull https://github.com/0l1v3rr/cli-file-manager.git
	build
	@echo "Successful update!"