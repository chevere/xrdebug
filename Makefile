vet:
	go vet -v ./...

test:
	go test -v ./...

test-coverage:
	go test -v ./... -coverprofile=coverage.out

open-coverage:
	go tool cover -html=coverage.out

.PHONY: build
build:
	go build -o build/os/arch/

build-all:
	@mkdir -p build
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/linux/amd64/
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o build/linux/arm64/
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/windows/amd64/
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o build/windows/arm64/
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/macos/amd64/
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/macos/arm64/
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o build/freebsd/amd64/
	GOOS=freebsd GOARCH=arm64 go build -ldflags="-s -w" -o build/freebsd/arm64/

dist-all:
ifndef VERSION
	@echo "Usage: make dist-all VERSION=<version>"
	@exit 1
endif
	@mkdir -p dist
	tar -C build/linux/amd64 -czvf dist/xrdebug-$(VERSION)-linux-amd64.tar.gz .
	tar -C build/linux/arm64 -czvf dist/xrdebug-$(VERSION)-linux-arm64.tar.gz .
	tar -C build/windows/amd64 -czvf dist/xrdebug-$(VERSION)-windows-amd64.tar.gz .
	tar -C build/windows/arm64 -czvf dist/xrdebug-$(VERSION)-windows-arm64.tar.gz .
	tar -C build/macos/amd64 -czvf dist/xrdebug-$(VERSION)-macos-amd64.tar.gz .
	tar -C build/macos/arm64 -czvf dist/xrdebug-$(VERSION)-macos-arm64.tar.gz .
	tar -C build/freebsd/amd64 -czvf dist/xrdebug-$(VERSION)-freebsd-amd64.tar.gz .
	tar -C build/freebsd/arm64 -czvf dist/xrdebug-$(VERSION)-freebsd-arm64.tar.gz .
