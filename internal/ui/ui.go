package ui

import (
	"fmt"
	"os"
)

type UI struct {
	Quiet   bool
	NoColor bool
	Verbose  bool
}

func NewUI() *UI {
	return &UI{}
}

func (u *UI) Info(msg string) {
	if u.Quiet {
		return
	}
	fmt.Println("⚠ ", msg)
}

func (u *UI) Success(msg string) {
	if u.Quiet {
		return
	}
	fmt.Println("✓ ", msg)
}

func (u *UI) Error(msg string) {
	if u.Quiet {
		return
	}
	fmt.Println("✗ ", msg)
}

func (u *UI) Running(msg string) {
	if u.Quiet {
		return
	}
	fmt.Println("◇ ", msg)
}

func (u *UI) Fatal(msg string) {
	fmt.Println("✗ ", msg)
	os.Exit(1)
}
