GO=go
NPM=cd files/app && npm
BIN=Brucheion
NODE_MODULES=files/app/node_modules

.PHONY: all dev build release app test clean deps

all: deps test build

build: app brucheion

brucheion:
	$(GO) build -o $(BIN) -v

release: deps app test
	./scripts/release.sh

app:
	$(NPM) run build

test: app
	$(GO) test -v ./...
	cd app && npm test

clean:
	$(GO) clean
	rm -f $(BIN)
	rm -r $(NODE_MODULES)

dev:
	$(NPM) run dev

deps: $(NODE_MODULES)

$(NODE_MODULES): files/app/package.json files/app/package-lock.json
	$(NPM) install
