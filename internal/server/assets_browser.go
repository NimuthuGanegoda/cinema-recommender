//go:build js && wasm

package server

import "embed"

// Browser builds use the same Go API handlers without embedding website photos.
var uiAssets embed.FS
