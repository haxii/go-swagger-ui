.PHONY: static
static:
	echo 'package static' > ./static/static.go
	echo '// static in-memory cache of swagger UI dist static files' >> ./static/static.go
	echo '// change by makefile, do NOT edit directly' >> ./static/static.go
	echo 'var staticFiles = map[string][]byte{' >> ./static/static.go
	ls ./dist/ | while read staticFileName; do \
		echo '"'$$staticFileName'":{' >> ./static/static.go; \
		xxd -i ./dist/$$staticFileName | sed '1d' | sed '$$d' | sed '$$d' | sed '$$s/$$/},/'>> ./static/static.go; \
	done
	echo '}' >> ./static/static.go
	gofmt -w ./static/static.go
