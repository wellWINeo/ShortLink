##
# Project ShortLink
#
# @file
# @version 1.0

ENTRYPOINT=cmd/main.go
BINARY=ShortLink

all: build

build:
	go buld ${ENTRYPOINT} -o {BINARY}

run: build
	./{BINARY}

clean:
	go clean
	rm ${BINARY}

# end
