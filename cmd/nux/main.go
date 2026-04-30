package main

import (
	"github.com/rsdenck/nux/cmd/nux/commands"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	commands.Execute(version, commit, date)
}
