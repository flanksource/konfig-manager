package ui

import "embed"

//nolint
//go:embed out/*
//go:embed out/_next/static/*/_*
//go:embed out/_next/static/chunks/*/_*
var StaticContent embed.FS
