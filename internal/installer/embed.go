package installer

import (
	"embed"
	"strings"
	"time"
)

//go:embed content_skills
var docsFS embed.FS

// mustReadDoc reads an embedded markdown file from docs/.
// Panics if the file is missing (compile-time guarantee via embed).
func mustReadDoc(path string) string {
	b, err := docsFS.ReadFile("content_skills/" + path)
	if err != nil {
		panic("missing embedded doc: " + path)
	}
	return string(b)
}

// mustReadAsset reads an embedded asset file from docs/assets/.
func mustReadAsset(name string) string {
	return mustReadDoc("assets/" + name)
}

// mustReadAssetWithDate reads an asset and replaces {{date}} with today's date.
func mustReadAssetWithDate(name string) string {
	return strings.ReplaceAll(
		mustReadAsset(name),
		"{{date}}",
		time.Now().Format("2006-01-02"),
	)
}
