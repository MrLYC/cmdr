package define

import "embed"

//go:embed scripts/*
var EmbedFS embed.FS
