package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"urlShortener/urlshort"
)

var fileName = flag.String("f", "jsonFile.json", "Specify the name of YAML file")

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	fmt.Println(mapHandler)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
//	yaml := `
//- path: /urlshort
//  url: https://github.com/gophercises/urlshort
//- path: /urlshort-final
//  url: https://github.com/gophercises/urlshort/tree/solution
//`
	// Replacing above yaml variable with one read from a file
	f, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("Failed to read the file")
		return
	}
	defer f.Close()

	file, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Failed to read the file")
		return
	}
	extensionRunes := []rune(*fileName)
	extensionString := string(extensionRunes[len(extensionRunes)-4:])
	switch extensionString {
	case "yaml":
		yamlHandler, err := urlshort.YAMLHandler([]byte(file), mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	case "json":
		jsonHandler, err := urlshort.JSONHandler([]byte(file), mapHandler)
		if err != nil{
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
	default:
		fmt.Println("Wrong file")
		return
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
