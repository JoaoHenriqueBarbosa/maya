.PHONY: demo simple clean

demo:
	GOOS=js GOARCH=wasm go build -o examples/demo/app.wasm examples/demo/main.go

simple:
	GOOS=js GOARCH=wasm go build -o examples/simple/app.wasm examples/simple/main.go

clean:
	rm -f examples/*/app.wasm

all: demo simple
