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

local: *.go
	go build -ldflags "-X main.Version=$(VERSION)-dev" -o build/kuisp

release:
	rm -rf build release && mkdir build release
	for os in linux freebsd darwin ; do \
		GOOS=$$os ARCH=amd64 godep go build -ldflags "-X main.Version=$(VERSION)" -o build/kuisp-$$os-amd64 ; \
		tar --transform 's|^build/||' --transform 's|-.*||' -czvf release/kuisp-$(VERSION)-$$os-amd64.tar.gz build/kuisp-$$os-amd64 README.md LICENSE ; \
	done
	GOOS=windows ARCH=amd64 godep go build -ldflags "-X main.Version=$(VERSION)" -o build/kuisp-$(VERSION)-windows-amd64.exe
	zip release/kuisp-$(VERSION)-windows-amd64.zip build/kuisp-$(VERSION)-windows-amd64.exe README.md LICENSE && \
		echo -e "@ build/kuisp-$(VERSION)-windows-amd64.exe\n@=kuisp.exe"  | zipnote -w release/kuisp-$(VERSION)-windows-amd64.zip
	go get github.com/progrium/gh-release/...
	gh-release create jimmidyson/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

clean:
	rm -rf build release

.PHONY: release clean
