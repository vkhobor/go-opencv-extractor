package web

import "embed"

// Relative paths are not supported in embed.FS.
// Meaning that accessing a file in a will be interpreted as acessing from where embed.go is located.
// So don't do this: go:embed all:dist/extractor/browser

//go:embed all:dist
var StaticFiles embed.FS

const PrefixForClientFiles = "dist/extractor/browser/"
const IndexHtml = "dist/extractor/browser/index.html"
