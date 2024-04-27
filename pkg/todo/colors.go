package todo

import (
	"github.com/fatih/color"
)

func red(s string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(s)
}

func green(s string) string {
	green := color.New(color.FgGreen).SprintFunc()
	return green(s)
}

func blue(s string) string {
	blue := color.New(color.FgBlue).SprintFunc()
	return blue(s)
}
