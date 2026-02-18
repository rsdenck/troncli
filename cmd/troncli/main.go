package main

import (
"github.com/mascli/troncli/cmd/troncli/commands"
)

var (
version = "dev"
commit  = "none"
date    = "unknown"
)

func main() {
commands.Execute(version, commit, date)
}
