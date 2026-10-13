//go:build !js || !wasm

package web

import "embed"

// Files contains the website used by the native server.
//
//go:embed *.html manifest.json sw.js assets data
var Files embed.FS
