REVISION = $(shell git describe --tags)
$(info    Make factom-cli $(REVISION))

LDFLAGS := "-s -w -X main.FactomcliVersion=$(REVISION)"

default: factom-cli
install: factom-cli-install
all: factom-cli-darwin-amd64 factom-cli-windows-amd64.exe factom-cli-windows-386.exe factom-cli-linux-amd64 factom-cli-linux-arm64 factom-cli-linux-arm7

BUILD_FOLDER := build

factom-cli:
	go build -trimpath -ldflags $(LDFLAGS)
factom-cli-install:
	go install -trimpath -ldflags $(LDFLAGS)

factom-cli-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-cli-darwin-amd64-$(REVISION)
factom-cli-windows-amd64.exe:
	env GOOS=windows GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-cli-windows-amd64-$(REVISION).exe
factom-cli-windows-386.exe:
	env GOOS=windows GOARCH=386 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-cli-windows-386-$(REVISION).exe
factom-cli-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-cli-linux-amd64-$(REVISION)
factom-cli-linux-arm64:
	env GOOS=linux GOARCH=arm64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-cli-linux-arm64-$(REVISION)
factom-cli-linux-arm7:
	env GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-cli-linux-arm7-$(REVISION)

.PHONY: clean

clean:
	rm -f factom-cli
	rm -rf build
