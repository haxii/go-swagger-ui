package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	jsoniter "github.com/json-iterator/go"

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
	querySwaggerHost    string = "host"
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

func serveLocalFile(localFilePath string, w http.ResponseWriter, r *http.Request) {
	newHost := r.URL.Query().Get("host")
	if len(newHost) == 0 {
		http.ServeFile(w, r, localFilePath)
		return
	}
	isJSON := false
	switch filepath.Ext(localFilePath) {
	case ".json":
		isJSON = true
	case ".yaml":
		fallthrough
	case ".yml":
		isJSON = false
	default:
		http.Error(w, "unknown swagger file: "+localFilePath, http.StatusBadRequest)
		return
	}

	// open file
	file, err := os.Open(localFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not exists", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	swg := new(map[string]interface{})
	if isJSON {
		dec := jsoniter.NewDecoder(file)
		err = dec.Decode(swg)
	} else {
		dec := yaml.NewDecoder(file)
		err = dec.Decode(swg)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	(*swg)["host"] = newHost
	var resp []byte
	if isJSON {
		resp, err = jsoniter.Marshal(swg)
	} else {
		resp, err = yaml.Marshal(swg)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Write(resp)
}

func handler(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Path[1:]
	if len(source) == 0 {
		source = "index.html"
	}

	// serve the local file
	localFile := ""
	if isNativeSwaggerFile && source == nativeSwaggerFileName {
		localFile = *swaggerFile
	} else if strings.HasPrefix(source, "swagger/") {
		// we treat path started with swagger as a direct request of a local swagger file
		localFile = filepath.Join(*localSwaggerDir, source[len("swagger/"):])
	}
	if len(localFile) > 0 {
		serveLocalFile(localFile, w, r)
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
		base, err := url.Parse("swagger/")
		if err != nil {
			return
		}
		target, err := url.Parse(f)
		if err != nil {
			return
		}
		targetSwagger = base.ResolveReference(target).String()
		if h := r.URL.Query().Get(querySwaggerHost); len(h) > 0 {
			targetSwagger += "?host=" + h
		}
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
