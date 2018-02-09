all: exec
	
exec:
	@go build -o bin/exec -v cmd/exec/*.go
	@chmod +x bin/exec
	
deps:
	@dep ensure || \
		go get -u github.com/golang/dep/cmd/dep && \
		dep ensure
