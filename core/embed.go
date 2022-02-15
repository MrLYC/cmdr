package core

import "embed"

//go:generate mockgen -destination=mock/fs.go -package=mock io/fs FS
//go:embed scripts/*
var EmbedFS embed.FS
