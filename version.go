package main

import (
	_ "embed"
	"strings"
)

// versionFile is the single source of truth for the version, also read by build.sh.
//
//go:embed VERSION
var versionFile string

// Version is shown in the UI title and printed by --version.
var Version = "v" + strings.TrimSpace(versionFile)
