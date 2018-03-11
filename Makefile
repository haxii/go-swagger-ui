VERSION=$(shell git describe --abbrev=0 --tags)
BUILD=$(shell git rev-parse --short HEAD)

# Inject the build version (commit hash) into the executable.
LDFLAGS := -ldflags "-X main.Build=$(BUILD) -X main.Version=$(VERSION)"

.PHONY: static
static:
	echo 'package st
	atic' > ./static/static.go
	echo '// Files static in-memory cache of swagger UI dist static files' >> ./static/static.go
	echo '// change by makefile, do NOT edit directly' >> ./static/static.go
	echo 'var Files = map[string][]byte{' >> ./static/static.go
	ls ./dist/ | while read staticFileName; do \
		echo '"'$$staticFileName'":{' >> ./static/static.go; \
		xxd -i ./dist/$$staticFileName | sed '1d' | sed '$$d' | sed '$$d' | sed '$$s/$$/},/'>> ./static/static.go; \
	done
	echo '}' >> ./static/static.go
	gofmt -w ./static/static.go


# `make clean` cleans up everything
.PHONY: clean
clean:
	rm -rf bin release swaggerui-*.tar.gz

################################################################################
## Below are commands for shipping distributable binaries for each platfomr.  ##
################################################################################
PLATFORMS := windows/386  windows/amd64 linux/386 linux/amd64 darwin/amd64
release: $(PLATFORMS)
.PHONY: release $(PLATFORMS)

# Handy variables to pull OS and arch from $PLATFORMS.
temp = $(subst /, ,$@)
os   = $(word 1, $(temp))
arch = $(word 2, $(temp))

$(PLATFORMS):
	mkdir -p release/swaggerui-$(VERSION)-$(os)-$(arch)
	cp -r README.md release/swaggerui-$(VERSION)-$(os)-$(arch)/
	GOOS=$(os) GOARCH=$(arch) \
		go build $(LDFLAGS) -v -i -o bin/swaggerui ./swagger.go
	cp bin/swaggerui* release/swaggerui-$(VERSION)-$(os)-$(arch)/
	cd release; tar -czvf ../swaggerui-$(VERSION)-$(os)-$(arch).tar.gz swaggerui-$(VERSION)-$(os)-$(arch)
