# GO Swagger UI

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
