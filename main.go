package main

import (
	"github.com/thalassa-cloud/cli/cmd"
	buildversion "github.com/thalassa-cloud/cli/internal/version"
)

// these must be set by the compiler using LDFLAGS
// -X main.version= -X main.commit= -X main.date= -X main.builtBy=
var (
	version string
	commit  string
	date    string
	builtBy string
)

func main() {
	// execute Cobra root cmd
	cmd.Execute()
}

func init() {
	buildversion.Init(version, commit, date, builtBy)
}
