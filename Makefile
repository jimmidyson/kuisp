# KUISP - A utility to serve static content & reverse proxy to RESTful services
#
# Copyright 2015 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

NAME=kuisp
VERSION=$(shell cat VERSION)
GO=GO15VENDOREXPERIMENT=1 go
pkgs = $(shell $(GO) list ./... | grep -v /vendor/)

local: *.go
	$(GO) build -ldflags "-X main.Version=$(VERSION)-dev" -o build/kuisp

arm:
	GOOS=linux GOARCH=arm $(GO) build -ldflags "-X main.Version=$(VERSION)" -o build/kuisp-linux-arm

bump:
	$(GO) get -u github.com/fabric8io/gobump
	gobump patch

release: bump
	$(GO) get -u github.com/progrium/gh-release
	rm -rf build release && mkdir build release
	for os in linux freebsd darwin ; do \
	GOOS=$$os GOARCH=amd64 $(GO) build -ldflags "-X main.Version=$(VERSION)" -o build/kuisp-$$os-amd64 ; \
	tar --transform 's|^build/||' --transform 's|-.*||' -czvf release/kuisp-$(VERSION)-$$os-amd64.tar.gz build/kuisp-$$os-amd64 README.md LICENSE ; \
	done
	GOOS=linux GOARCH=arm $(GO) build -ldflags "-X main.Version=$(VERSION)" -o build/kuisp-linux-arm
	tar --transform 's|^build/||' --transform 's|-.*||' -czvf release/kuisp-$(VERSION)-linux-arm.tar.gz build/kuisp-linux-arm README.md LICENSE ; \
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "-X main.Version=$(VERSION)" -o build/kuisp-$(VERSION)-windows-amd64.exe
	zip release/kuisp-$(VERSION)-windows-amd64.zip build/kuisp-$(VERSION)-windows-amd64.exe README.md LICENSE && \
		echo -e "@ build/kuisp-$(VERSION)-windows-amd64.exe\n@=kuisp.exe"  | zipnote -w release/kuisp-$(VERSION)-windows-amd64.zip
	go get github.com/progrium/gh-release/...
	gh-release create jimmidyson/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

test:
	$(GO) get -u github.com/jstemmer/go-junit-report
	OUTPUT=`$(GO) test -short -race -v $(pkgs)` && echo "$${OUTPUT}" | tee /dev/tty | go-junit-report -set-exit-code > $${CIRCLE_TEST_REPORTS:-.}/junit.xml

clean:
	rm -rf build release

.PHONY: release clean test bump
