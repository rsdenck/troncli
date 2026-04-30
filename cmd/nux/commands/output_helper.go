package commands

import (
	"github.com/rsdenck/nux/internal/output"
)

func printSuccess(data interface{}, message string) {
	output.NewSuccess(data).WithMessage(message).Print()
}

func printError(err string, code string) {
	output.NewError(err, code).Print()
}

func printList(items interface{}, total int, message string) {
	output.NewList(items, total).WithMessage(message).Print()
}
