package cli

import (
	"fmt"

	"github.com/fatih/color"
)

// BoolSetting provides a human-readable serialization for bool values.
func BoolSetting(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// Printf prints the given text using fmt.Printf
// in a way where colors work on Windows.
func Printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(color.Output, format, a...)
	if err != nil {
		fmt.Printf(format, a...)
	}
}

// Println prints the given text using fmt.Println
// in a way where colors work on Windows.
func Println(a ...interface{}) {
	_, err := fmt.Fprintln(color.Output, a...)
	if err != nil {
		fmt.Println(a...)
	}
}

// PrintlnColor prints using the given color function.
// If that doesn't work, it falls back to printing without color.
func PrintlnColor(color *color.Color, messages ...interface{}) {
	_, err := color.Println(messages...)
	if err != nil {
		fmt.Println(messages...)
	}
}

func PrintEntry(label, value string) {
	Printf("  %s: %s\n", label, value)
}

// PrintError prints the given error message to the console.
func PrintError(err error) {
	fmt.Println()
	PrintlnColor(color.New(color.Bold).Add(color.FgRed), "Error:", err.Error())
	fmt.Println()
}

func PrintHeader(text string) {
	boldUnderline := color.New(color.Bold).Add(color.Underline)
	PrintlnColor(boldUnderline, text+":")
}

// PrintLabelAndValue prints the label bolded and underlined
// the value indented on the next line
// followed by an empty line.
func PrintLabelAndValue(label, value string) {
	PrintHeader(label)
	Println(Indent(value))
	fmt.Println()
}

// PrintDriverAction logs activities from a code hosting driver on the CLI.
func PrintDriverAction(template string, messages ...interface{}) {
	fmt.Println()
	_, err := color.New(color.Bold).Printf(template, messages...)
	if err != nil {
		fmt.Printf(template, messages...)
	}
}

func StringSetting(text string) string {
	if text != "" {
		return text
	}
	return "(not set)"
}
