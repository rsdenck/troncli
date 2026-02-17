package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	commands := []string{
		"audit", "bash", "completion", "container", "disk", "doctor",
		"firewall", "group", "network", "pkg", "plugin", "process",
		"remote", "service", "system", "user",
	}

	f, err := os.Create("COMMAND.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "# TRONCLI Command Reference\n\nGenerated on %s\n\n", time.Now().Format(time.RFC1123))

	fmt.Fprintf(f, "## Main Help\n```\n")
	out, _ := exec.Command("go", "run", "cmd/troncli/main.go", "--help").CombinedOutput()
	f.Write(out)
	fmt.Fprintf(f, "```\n\n")

	fmt.Fprintf(f, "## Subcommands\n")
	for _, cmd := range commands {
		fmt.Fprintf(f, "### %s\n```\n", cmd)
		out, _ := exec.Command("go", "run", "cmd/troncli/main.go", cmd, "--help").CombinedOutput()
		f.Write(out)
		fmt.Fprintf(f, "```\n\n")
	}
	fmt.Println("COMMAND.md generated successfully.")
}
