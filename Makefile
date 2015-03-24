NAME=kui
VERSION=$(shell cat VERSION)

local: *.go
	godep go build -ldflags "-X main.Version $(VERSION)-dev" -o build/kui

release:
	rm -rf release && mkdir release
	go get github.com/progrium/gh-release/...
	cp build/* release
	gh-release create jimmidyson/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

clean:
	rm -f build

.PHONY: release build
