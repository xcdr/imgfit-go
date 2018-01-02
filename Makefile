BUILD_VERSION 	= 1.0.0
BUILD_BRANCH 	= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE		= $(shell date +%Y%m%d%H%M)

LDFLAGS 		= "-X main.version=$(BUILD_VERSION) -X main.build=$(BUILD_DATE).$(BUILD_BRANCH)"

build: prepare imgfit perm

prepare:
	mkdir -p build/imgfit/etc
	cp config.yml build/imgfit/etc

	mkdir -p build/imgfit/bin

imgfit: lib_config.go lib_image.go lib_server.go main.go
	go build -ldflags ${LDFLAGS} -o build/imgfit/bin/imgfit lib_config.go lib_image.go lib_server.go main.go

test:
	rm -rf /tmp/imgfit
	go test -v
	rm -rf /tmp/imgfit

bench:
	rm -rf /tmp/imgfit
	go test -v -run=bench -bench=.
	rm -rf /tmp/imgfit

perm:
	chmod -R og+rw build

clean:
	go clean
	rm -rf build
