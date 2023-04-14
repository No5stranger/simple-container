BIN_FILE="t-snapshot"

build-test:
	@git checkout snapshot
	@git pull origin
	@go build -o $(BIN_FILE)
