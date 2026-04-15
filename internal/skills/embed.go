package skills

import (
	"embed"
	"strings"
	"time"
)

//go:embed content
var docsFS embed.FS

func mustReadDoc(path string) string {
	b, err := docsFS.ReadFile("content/" + path)
	if err != nil {
		panic("missing embedded doc: " + path)
	}
	return string(b)
}

func mustReadAsset(name string) string {
	return mustReadDoc("assets/" + name)
}

func mustReadAssetWithDate(name string) string {
	return strings.ReplaceAll(
		mustReadAsset(name),
		"{{date}}",
		time.Now().Format("2006-01-02"),
	)
}
