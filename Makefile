.PHONY: static
static:
	echo 'package main' > ./static.go
	echo '// static in-memory cache of swagger UI dist static files' >> ./static.go
	echo '// change by makefile, do NOT edit directly' >> ./static.go
	echo 'var staticFiles = map[string][]byte{' >> ./static.go
	ls ./dist/ | while read staticFileName; do \
		echo '"'$$staticFileName'":{' >> ./static.go; \
		xxd -i ./dist/$$staticFileName | sed '1d' | sed '$$d' | sed '$$d' | sed '$$s/$$/},/'>> ./static.go; \
	done
	echo '}' >> ./static.go
	gofmt -w ./static.go
