package report

import (
	"io"
	"github.com/vron/report/repml"
)

// A Report object is the object that things to include in the 
// report is added to. If an error occurse when writing the
// status is changed so a error may be retrieved
// All reporting functions return the Reporter itself so calls
// may be chained.
type Report interface {
	// Start a new section (i.e. increase heading level by one)
	Section() Report
	// Start a new figure
	Figure() Report
	// End the last started thing
	End() Report

	// Add a new heading at the current level
	Heading(string) Report
	// Add a paragraph with text
	Paragraph(string) Report
	// Add a caption to the current context
	Caption(string) Report

	// Return nil if all is ok, otherwise the first generated error
	IsError() error

	// Called at the end to add any mandatory ending, maybe defer this?
	Finish() error
}

// Create a new report writing to specified writer, defaults
// to a html5 reporter.
func New(w io.Writer, title string) Report {
	return repml.New(w,title)
}
