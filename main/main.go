package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"urlshort"
)

func main() {
	var jsonFile, yamlFile string
	flag.StringVar(&jsonFile, "json", "", "path to json file.")
	flag.StringVar(&yamlFile, "yaml", "", "path to yaml file.")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	Handler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	if jsonFile != "" {
		fmt.Println(jsonFile)
		jsonData, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			panic(err)
		}
		// Build the JSONHandler using the mapHandler as the
		// fallback
		Handler, err = urlshort.JsonHandler([]byte(jsonData), Handler)
		if err != nil {
			panic(err)
		}
	} else if yamlFile != "" {
		yamlData, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			panic(err)
		}
		Handler, err = urlshort.YAMLHandler([]byte(yamlData), Handler)
		if err != nil {
			panic(err)
		}

	}
	fmt.Println("Starting the server on :8080")
	err := http.ListenAndServe(":8080", Handler)
	if err != nil {
		panic(err)
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
