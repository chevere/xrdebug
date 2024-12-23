test:
	gotestsum -- -v ./...

test-coverage:
	mkdir -p coverage
	gotestsum -- -v ./... -coverpkg=./... -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html

make open-coverage:
	open coverage/coverage.html

dev-deps:
	go install gotest.tools/gotestsum@latest

build:
	go build -o build/ ./...

build--all:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/xrdebug-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o build/xrdebug-linux-arm64
	GOOS=windows GOARCH=amd64 go build -o build/xrdebug-windows-amd64.exe
	GOOS=windows GOARCH=arm64 go build -o build/xrdebug-windows-arm64.exe
	GOOS=darwin GOARCH=amd64 go build -o build/xrdebug-macos-amd64
	GOOS=darwin GOARCH=arm64 go build -o build/xrdebug-macos-arm64
	GOOS=freebsd GOARCH=amd64 go build -o build/xrdebug-freebsd-amd64
	GOOS=freebsd GOARCH=arm64 go build -o build/xrdebug-freebsd-arm64

dist--all:
	mkdir -p dist
	cd build && for file in *; do \
		tar -czf ../dist/$$file.tar.gz $$file; \
	done
