//go:build !js || !wasm

package server

import "embed"

//go:embed ui/*
var uiAssets embed.FS
