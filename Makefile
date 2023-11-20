# Build target: compiles the Go code and generates an executable named "bankApp" in the "bin" directory.
build:
	@go build -o bin/bankApp

# Run target: depends on the "build" target, executes the compiled "bankApp" executable.
run: build
	@./bin/bankApp

# Test target: runs tests for the entire Go project in verbose mode.
test:
	@go test -v ./...
