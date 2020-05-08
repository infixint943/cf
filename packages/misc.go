package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

// Log is struct holding functions to print colored to stderr
// (lightweight replacement for Logger)
var (
	writer = os.Stderr
	Green  = color.New(color.FgGreen)
	Blue   = color.New(color.FgBlue)
	Red    = color.New(color.FgRed)
	Yellow = color.New(color.FgYellow)

	Log struct {
		Success, Notice, Info, Error,
		Warning func(text ...interface{})
	}
)

func init() {
	Log.Success = func(text ...interface{}) { Green.Fprintln(writer, text...) }
	Log.Notice = func(text ...interface{}) { fmt.Fprintln(writer, text...) }
	Log.Info = func(text ...interface{}) { Blue.Fprintln(writer, text...) }
	Log.Error = func(text ...interface{}) { Red.Fprintln(writer, text...) }
	Log.Warning = func(text ...interface{}) { Yellow.Fprintln(writer, text...) }
}

// PrintError outputs error (with custom message)
// and exits the program execution (if err != nil)
func PrintError(err error, desc string) {
	if err != nil {
		if desc != "" {
			Log.Error(desc)
		}
		Log.Error(err.Error())
		os.Exit(0)
	}
}

// CreateFile copies data to dst (create if not exists)
// Returns absolute path to destination file
func CreateFile(data, dst string) string {
	out, err := os.Create(dst)
	PrintError(err, "File "+dst+" couldn't be created!")
	defer out.Close()

	out.WriteString(data)
	return dst
}

// GetText extracts text from particular html data
func GetText(sel *goquery.Selection, query string) string {
	str := sel.Find(query).Text()
	return strings.TrimSpace(str)
}

// GetAttr extracts attribute valur of particular html data
func GetAttr(sel *goquery.Selection, query, attr string) string {
	str := sel.Find(query).AttrOr(attr, "")
	return strings.TrimSpace(str)
}
