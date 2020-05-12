package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(writer, request, url, http.StatusFound)
		}
		fallback.ServeHTTP(writer, request)
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
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func parseYaml(yamlBytes []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yamlBytes, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMap(pathUrls []pathUrl) map[string]string{
	pathToUrls := make(map[string]string)

	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.Url
	}
	return pathToUrls
}


type pathUrl struct {
	Path string `yaml:"path"`
	Url string `yaml:"url"`
}
