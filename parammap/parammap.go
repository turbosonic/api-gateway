package parammap

import (
	"net/http"
	"strings"

	"goji.io/pat"
)

func GetParams(originalPattern string, r *http.Request) map[string]string {
	params := make(map[string]string)

	// slice up the url
	elements := strings.Split(originalPattern, "/")

	// check each one for the : prefix
	for _, e := range elements {
		if strings.HasPrefix(e, ":") {
			// if it has the prefix add it to the map
			params[e] = pat.Param(r, strings.Replace(e, ":", "", 1))
		}
	}

	return params
}
