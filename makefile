
.PHONY: all clean

test:
	GO111MODULE=on go test -tags=native github.com/bokuweb/gopher-boy/...

reg:
	reg-cli ./test/actual ./test/expect ./test/diff

reg-update:
	reg-cli ./test/actual ./test/expect ./test/diff -U

build:
	GO111MODULE=on go build -tags="native" -o "gopher-boy" "cmd/gopher-boy/main.go"
