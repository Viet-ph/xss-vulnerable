package utils

import (
	"strings"
)

func InputValidater(original string) string {

	blackList := []string{"<script", "script>", "eval(",
						  "eval&#40;", "javascript:", "javascript&#58;",
						  "fromCharCode", "&#62;", "&#60;", "&lt;", "&rt;"}
	for _, word := range blackList {
		original = strings.Replace(original, word, "", -1)
	}

	return original
}
