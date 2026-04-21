package audit

import "embed"

//go:embed content
var contentFS embed.FS

func mustReadContent(path string) string {
	b, err := contentFS.ReadFile("content/" + path)
	if err != nil {
		panic("missing embedded audit content: " + path)
	}
	return string(b)
}
