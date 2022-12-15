package main

import (
	"fmt"
	"net/http"

	urlshort "github.com/JaydenTeoh/url-shortener/handler"
)

func main() {
	mux := defaultMux() //fallback http.Handler

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /github
  url: https://github.com/JaydenTeoh
- path: /linkedin
  url: https://www.linkedin.com/in/jayden-teoh/
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux() //ServeMux type has a ServeHTTP() method, meaning that it too satisfies the http. Handler interface.
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my URL shortener.")
	fmt.Fprintln(w, "  Add the following /:path to the URL to access different pages.")
	fmt.Fprintln(w, "     /github: My Github Page")
	fmt.Fprintln(w, "     /linkedin: My Linkedin Page")
}
