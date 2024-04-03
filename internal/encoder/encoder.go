package encoder

import "strings"

func isSupportedContentType(contentType string) bool {
	supportedContentType :=
		[...]string{"application/json", "text/html; charset=utf-8"}

	for _, v := range supportedContentType {
		if strings.Contains(v, contentType) {
			return true
		}
	}

	return false
}
