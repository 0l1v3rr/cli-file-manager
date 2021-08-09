build:
	go build -o ./bin/cli-file-manager cmd/cli-file-manager/main.go
run: build
	./bin/cli-file-manager
clean:
	rm -rf bin/