package cmd

import (
	cln "cf/client"
	pkg "cf/packages"

	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/k0kubun/go-ansi"
)

// RunWatch is called on running cf watch
func (opt Opts) RunWatch() {
	// check if contest id is present
	if opt.contest == "" {
		pkg.Log.Error("No contest id found")
		return
	}
	// header formatting for table
	headerfmt := color.New(color.FgBlue, color.Bold, color.Underline).SprintfFunc()

	if opt.SubCnt == 0 {
		// submissions aren't specified to be parsed
		// parse contest solve status instead

		// fetch contest solve status
		data, err := cln.WatchContest(opt.group, opt.contest, opt.contClass)
		pkg.PrintError(err, "Failed to extract contest solve status")

		// init table with header + color
		tbl := uitable.New()
		tbl.AddRow(headerfmt("#"), headerfmt("Name"),
			headerfmt("Status"), headerfmt("Count"))
		tbl.MaxColWidth = 40
		tbl.Separator = " | "

		// iterate over each row of data
		for _, prob := range data {
			if prob.ID != strings.ToUpper(opt.problem) && opt.problem != "" {
				// skip row, if doesn't match query
				continue
			}
			if strings.Index(prob.Count, "x") != -1 {
				//remove rune char from solve count text
				prob.Count = prob.Count[strings.Index(prob.Count, "x"):]
			} else {
				prob.Count = ""
			}
			// set color to be printed based on solved status
			clean := func(status string) string {
				switch status {
				case "accepted-problem":
					return color.New(color.BgGreen).Sprint("      ")
				case "rejected-problem":
					return color.New(color.BgRed).Sprint("      ")
				default:
					return color.New(color.BgWhite).Sprint("      ")
				}
			}
			// insert row to table
			tbl.AddRow(prob.ID, prob.Name, clean(prob.Status), prob.Count)
		}
		fmt.Println(tbl)
	} else {
		// infinite loop till verdicts declared
		for isFirst := true; true; isFirst = false {
			// timer to fetch data in interval of 1 second
			start := time.Now()
			// fetch contest submission status
			data, err := cln.WatchSubmissions(opt.group, opt.contest, opt.contClass, opt.problem)
			pkg.PrintError(err, "Failed to extract submissions in contest")

			// min function (since there golang lacks min/max uggh)
			min := func(a, b int) int {
				if a <= b {
					return a
				}
				return b
			}
			if isFirst == false {
				for i := 0; i <= min(opt.SubCnt, len(data)); i++ {
					ansi.CursorPreviousLine(1)
					ansi.EraseInLine(2)
				}
			}

			// create new table
			tbl := uitable.New()
			tbl.MaxColWidth = 20
			tbl.Separator = " | "
			tbl.AddRow(headerfmt("#"), headerfmt("When"), headerfmt("Name"), headerfmt("Lang"),
				headerfmt("Verdict"), headerfmt("Time"), headerfmt("Memory"))

			// do there exist submissions with pending verdicts?
			isPending := false
			for i, sub := range data {
				// break if exceeds reqd capacity
				if i >= opt.SubCnt {
					break
				}
				// insert row into table
				tbl.AddRow(sub.ID, sub.When, sub.Name, sub.Lang, sub.Verdict, sub.Time, sub.Memory)
				// update pending status
				if sub.Waiting == "true" {
					isPending = true
				}
			}
			fmt.Println(tbl)

			if isPending == false {
				break
			}

			time.Sleep(time.Second - time.Since(start))
		}
	}
	return
}
