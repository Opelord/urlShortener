package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

type PathsAndUrls struct{
	Path string `json:"path" yaml:"path"`
	Url string `json:"url" yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		url, ok := pathsToUrls[req.URL.Path]
		if ok {
			http.Redirect(writer, req, url, 301)
		} else {
			fallback.ServeHTTP(writer, req)
		}
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]PathsAndUrls, error){
	var pathsToUrls []PathsAndUrls
	err := yaml.Unmarshal(yml, &pathsToUrls)
	return pathsToUrls, err
}

// Below are two functions that are analogous to Yaml functions
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(jsn)
	if err != nil{
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsn []byte) ([]PathsAndUrls, error){
	var pathsToUrls []PathsAndUrls
	err := json.Unmarshal(jsn, &pathsToUrls)
	return pathsToUrls, err
}

// buildMap translates parsed json and yaml files
// to string map that can be handled by MapHandler function
func buildMap(parsed []PathsAndUrls) map[string]string {
	m := make(map[string]string)
	for _, entry := range parsed{
		m[entry.Path] = entry.Url
	}
	return m
}

