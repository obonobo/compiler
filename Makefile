.PHONY: default build build-static clean install grammar tool grammar-stdout
default: build


SHELL		=	bash
out			=	esacc
codegen 	=	./resources/a2/workspace/tool.go
# gram		=	./resources/a2/workspace/COMP442.grammar.BNF.grm.noebnf.noambiguity.pure
gram		= 	./resources/a3/workspace/grammar-mangled-attribute.grm
codegen_out	=	./core/token/gen.go


download:
	go get -d -v

build: download
	export GOOS=linux
	export GO111MODULE=on
	go build -o $(out)

install:
	install -o root -g root -m 0755 $(out) /bin/$(out)

# Adds some flags for building the app statically linked
build-static: download
	export GOOS=linux
	export GO111MODULE=on
	export CGO_ENABLED=0
	go build \
		-ldflags="-extldflags=-static" \
		-tags osusergo,netgo \
		-o $(out)

clean:
	rm -rf ./$(out) ./vendor

test:
	go clean --testcache && go test -v -timeout 30s ./...

# Generates the parser table from grammar production rules.
grammar:
	$(codegen) --compile $(gram) > $(codegen_out)
	gofmt -w $(codegen_out)

grammar-stdout:
	$(codegen) --compile $(gram)

tool:
	$(codegen) $(gram)
