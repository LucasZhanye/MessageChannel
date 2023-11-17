package web

import "embed"

//go:embed static
var Web embed.FS

var FileName = "static"
