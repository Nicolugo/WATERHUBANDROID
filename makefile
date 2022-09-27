
.PHONY: all clean

test:
	GO111MODULE=on go test -tags=native github.com/bokuweb/gopher-boy/...

reg:
	reg-cli ./test/actual ./test/expect ./test/diff

reg-update:
	reg-cli ./test/actual ./test/expect ./test/diff -U

build:
	GO111MODULE=on go build -tags="native" -o "gopher-boy" "cmd/gopher-boy/main.go"

build-wasm:
	GOOS=js GOARCH=wasm go build -tags=wasm -o "docs/main.wasm" "cmd/gopher-boy/wasm_main.go"

serve:
	xdg-open 'http://localhost:5000'
	serve -r docs -a :5000 || (go get -v github.com/mattn/serve && serve  -r docs -a :5000)

clean:
	rm -f *.wasm	