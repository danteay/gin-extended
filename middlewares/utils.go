package middlewares

import "regexp"

func skipRegexpPath(list []string, path string) bool {
	if length := len(list); length > 0 {
		for _, reg := range list {
			exp := regexp.MustCompile(reg)

			if exp.MatchString(path) {
				return true
			}
		}

	}

	return false
}
