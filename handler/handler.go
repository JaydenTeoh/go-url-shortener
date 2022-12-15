package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//if we can match a path, then redirect to it, else...
		path := r.URL.Path //get the path of the url of the http request

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// 1. Parse the yaml
	pathUrls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	// 2. Convert yaml array into map
	pathsToUrls := generateYamlMap(pathUrls)
	// 3. Return a map handler using the map
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls) //Unmarshal the yaml file and assign decoded values into the pointer to the pathUrls array
	if err != nil {
		return nil, err //unable to unmarshal yaml file and generate http.Handler
	}
	return pathUrls, nil
}

func generateYamlMap(data []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, values := range data {
		pathsToUrls[values.Path] = values.URL //map each path to their corresponding URL
	}
	return pathsToUrls
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
