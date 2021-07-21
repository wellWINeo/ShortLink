##
# Project ShortLink
#
# @file
# @version 1.0

ENTRYPOINT=cmd/main.go
BINARY=ShortLink

all: build

build:
	go build ${ENTRYPOINT}
	mv main ${BINARY}

run: build
	./${BINARY}

clean:
	go clean
	rm ${BINARY}

generate:
	go generate ./...

test:
	go test ./...

# end
