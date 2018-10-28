package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/haxii/daemon"
	"github.com/haxii/go-swagger-ui/static"
)

var (
	// Build of git, got by LDFLAGS on build
	Build = "-unknown-"
	// Version of git, got by LDFLAGS on build
	Version = "-unknown-"
)

var (
	_ = flag.String("s", daemon.UsageDefaultName, daemon.UsageMessage)

	serverAddr  = flag.String("l", ":8080", "server's listening Address")
	swaggerFile = flag.String("f",
		"http://petstore.swagger.io/v2/swagger.json",
		"swagger url or local file path")
	localSwaggerDir = flag.String("d", "/swagger", "swagger files vhost dir")
	enableTopbar    = flag.Bool("b", false, "enable the topbar")

	isNativeSwaggerFile   = false
	nativeSwaggerFileName = ""
)

const (
	querySwaggerURLKey  string = "url"
	querySwaggerFileKey string = "file"
)

func main() {
	daemon.Make("-s",
		"swaggerui",
		"Swagger UI service").Run(serve)
}

func serve() {
	flag.Parse()
	fmt.Printf("Server listening on %s\n", *serverAddr)

	// test if swagger file is a local one
	if fileStat, err := os.Stat(*swaggerFile); err == nil &&
		fileStat.Mode().IsRegular() {
		isNativeSwaggerFile = true
		nativeSwaggerFileName = filepath.Base(*swaggerFile)
	}
	if isNativeSwaggerFile {
		fmt.Printf("Using default local swagger file %s\n", *swaggerFile)
	} else {
		fmt.Printf("Using default online swagger file %s\n", *swaggerFile)
	}
	if *enableTopbar {
		fmt.Println("Topbar enabled")
	} else {
		fmt.Println("Topbar disabled")
	}
	fmt.Println("Swagger UI version", Version, ", build", Build)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*serverAddr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Path[1:]
	if len(source) == 0 {
		source = "index.html"
	}

	// serve the local file
	if isNativeSwaggerFile && source == nativeSwaggerFileName {
		http.ServeFile(w, r, *swaggerFile)
		return
	} else if strings.HasPrefix(source, "swagger/") {
		// we treat path started with swagger as a direct request of a local swagger file
		http.ServeFile(w, r, filepath.Join(*localSwaggerDir,
			source[len("swagger/"):]))
		return
	}

	// server the swagger UI
	//
	// find the in-memory static files
	staticFile, exists := static.Files[source]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// set up the content type
	switch filepath.Ext(source) {
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// return back the non-index files
	if source != "index.html" {
		w.Header().Set("Content-Length", strconv.Itoa(len(staticFile)))
		w.Write(staticFile)
		return
	}

	// set up the index page
	targetSwagger := *swaggerFile
	if f := r.URL.Query().Get(querySwaggerFileKey); len(f) > 0 {
		// requesting a local file, join it with a `swagger/` prefix
		targetSwagger = filepath.Join("swagger", f)
	} else if url := r.URL.Query().Get(querySwaggerURLKey); len(url) > 0 {
		// deal with the query swagger firstly
		targetSwagger = url
	} else if isNativeSwaggerFile {
		// for a native swagger file, use the filename directly
		targetSwagger = nativeSwaggerFileName
	}
	// replace the target swagger file in index
	indexHTML := string(staticFile)
	indexHTML = strings.Replace(indexHTML,
		"https://petstore.swagger.io/v2/swagger.json",
		targetSwagger, -1)
	if *enableTopbar {
		indexHTML = strings.Replace(indexHTML,
			"SwaggerUIBundle.plugins.DownloadUrl, HideTopbarPlugin",
			"SwaggerUIBundle.plugins.DownloadUrl", -1)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(indexHTML)))
	fmt.Fprint(w, indexHTML)
}
