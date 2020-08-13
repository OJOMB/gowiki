package wikiUtils

import (
	"os"
	"strings"
)

const PATH_SEPARATOR = string(os.PathSeparator)

func ConstructPath(fields []string, fileExtension string) string {
	return strings.Join(fields, PATH_SEPARATOR) + fileExtension
}
