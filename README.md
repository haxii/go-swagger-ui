# Go Swagger UI

Cross-platform [Swagger UI](https://swagger.io/swagger-ui/) built with golang with only
a fat binary, which can help you serve your swagger documentation with only one command.

## Usage

Start a swagger documentation server on port `8000` for `/path/to/your/swagger.json`:

```
$ swaggerui -l "0.0.0.0:8000" -f "/path/to/your/swagger.json"
```

Everyone can read your documentation on `http://your.ip.address:8000`, you can also view
other online swagger file, such as swagger's [petstore](http://petstore.swagger.io/),
by opening 

```
http://your.ip.address:8000/?config=http://petstore.swagger.io/v2/swagger.json
```

To make it as a daemon:

```
$ swaggerui -s install -l "0.0.0.0:8000" -f "/path/to/your/swagger.json"
$ swaggerui -s start
```

which could run in background and autostart with your system.

More detailed usage:

```
Usage of swaggerui:

  -b    enable the topbar
  -f string
        swagger url or local file path (default "http://petstore.swagger.io/v2/swagger.json")
  -l string
        server's listening Address (default ":8080")
  -s string
        Send signal to a master process: install, remove, start, stop, status (default "status")
```

## Build

Source code is written in [go](https://golang.org/), [make](https://www.gnu.org/software/make/) and [xxd](https://www.systutorials.com/docs/linux/man/1-xxd/) is also needed to build the binary.

The Swagger UI's HTML pages in [dist](dist) folder is copied from the original [Swagger UI Source](https://github.com/swagger-api/swagger-ui/tree/master/dist), then converted into bytes array in file [static/static.go](static/static.go) using `xxd` command, to re-build it

```
$ make static
```

By typing the following command, you can get the cross platform distributions of this program


```
$ make release
```
