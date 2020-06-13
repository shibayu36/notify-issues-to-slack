VERSION = $(shell gobump show -r)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/shibayu36/notify-issues-to-slack.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

export GO111MODULE=on

devel-deps:
	go install ${u} golang.org/x/lint/golint \
	  github.com/mattn/goveralls              \
	  github.com/x-motemen/gobump/cmd/gobump    \
	  github.com/Songmu/goxz/cmd/goxz         \
	  github.com/Songmu/ghch/cmd/ghch         \
	  github.com/tcnksm/ghr

test:
	go test -v ./...

lint: devel-deps
	go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build:
	go build -ldflags=$(BUILD_LDFLAGS) .

bump: devel-deps
	_tools/releng

crossbuild:
	goxz -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) -arch=amd64,386 \
	  -d=./dist/v$(VERSION) .

upload:
	ghr v$(VERSION) dist/v$(VERSION)

release: bump crossbuild upload

.PHONY: test devel-deps lint cover build bump crossbuild upload release
