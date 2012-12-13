package repml

import (
	//"report"
	"errors"
	"html"
	"io"
	"strconv"
)

// TODO: Implement checking of current status before accepting input

const (
	s_body = iota
	s_head1
	s_head2
	s_head3
	s_head4
	s_head5
	s_fig
	s_tbl
	s_tblrow
)

type state int
type stack []state

func (s stack) peek() state {
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return -1
}
func (s *stack) pop() state {
	if len(*s) > 0 {
		r := (*s)[len(*s)-1]
		*s = (*s)[0 : len(*s)-1]
		return r
	}
	return -1
}
func (s *stack) push(st state) {
	*s = append(*s, st)
}

type Report struct {
	e []error
	w io.Writer
	s stack
}

func (r *Report) seterr(e error) {
	if r.e == nil {
		r.e = []error{}
	}
	if e != nil {
		r.e = append(r.e, e)
	}
}
func (r *Report) ws(s string) {
	_, e := io.WriteString(r.w, s)
	r.seterr(e)
}

func New(w io.Writer, title string) *Report {
	r := &Report{e: nil, w: w}
	r.ws(`<!DOCTYPE html><html><head><link rel="stylesheet" type="text/css" href="style.css" /><title>`)
	r.ws(html.EscapeString(title))
	r.ws("</title></head><body>\n")
	// We start in body state!
	r.s = stack{s_body}
	return r
}

// TODO: When we know where we are we should escape everything
func (r *Report) Finish() []error {
	// Close all the states except the last one
	for i := len(r.s); i > 1; i-- {
		r.End()
	}
	if len(r.s) != 1 || r.s[0] != s_body {
		r.seterr(errors.New("Expect to be at body state at end"))
	}
	r.ws(`</body></html>`)
	return r.e
}

func (r *Report) Section() *Report {
	switch r.s.peek() {
	case s_body, s_head1, s_head2, s_head3, s_head4:
		r.s.push(r.s.peek() + 1)
		r.ws("<span>\n")
	case s_head5:
		r.s.push(r.s.peek())
		r.seterr(errors.New("Cannot have more nested heading levels"))
	default:
		r.seterr(errors.New("Not a legal positon for starting Section"))
	}
	return r
}
func (r *Report) Figure() *Report {
	switch r.s.peek() {
	case s_body, s_head1, s_head2, s_head3, s_head4:
		r.ws("<figure>\n")
		r.s.push(s_fig)
	default:
		r.seterr(errors.New("Must be in text mode to start figure"))
	}
	return r
}
func (r *Report) Table() *Report {
	switch r.s.peek() {
	case s_body, s_head1, s_head2, s_head3, s_head4, s_head5, s_fig:
		r.ws("<table>\n")
		r.s.push(s_tbl)
	default:
		r.seterr(errors.New("Must be in text or fig mode to start table"))
	}
	return r
}
func (r *Report) Row() *Report {
	switch r.s.peek() {
	case s_tbl:
		r.ws("<tr>\n")
		r.s.push(s_tblrow)
	default:
		r.seterr(errors.New("Must be in table mode to start row"))
	}
	return r
}
func (r *Report) Cell(str string) *Report {
	switch r.s.peek() {
	case s_tblrow:
		r.ws("<tc>" + html.EscapeString(str) + "</tc>")
	default:
		r.seterr(errors.New("Must be in table row mode to start cell"))
	}
	return r
}

func (r *Report) End() *Report {
	switch r.s.peek() {
	case s_body:
		r.seterr(errors.New("Allready in body mode"))
	case s_head1, s_head2, s_head3, s_head4, s_head5:
		r.ws("</span>\n")
		r.s.pop()
	case s_fig:
		r.ws("</figure>\n")
		r.s.pop()
	case s_tbl:
		r.ws("</table>\n")
		r.s.pop()
	case s_tblrow:
		r.ws("</tr>\n")
		r.s.pop()
	default:
		r.seterr(errors.New("Not recognized state, popping anyway"))
		r.s.pop()
	}
	return r
}
func (r *Report) Heading(str string) *Report {
	switch r.s.peek() {
	case s_body, s_head1, s_head2, s_head3, s_head4, s_head5:
		r.ws("<h" + strconv.Itoa(int(r.s.peek()+1)) + ">")
		r.ws(html.EscapeString(str))
		r.ws("</h" + strconv.Itoa(int(r.s.peek()+1)) + ">\n")
	default:
		r.seterr(errors.New("Can only create heading in text mode " + strconv.Itoa(int(r.s.peek()))))
	}
	return r
}
func (r *Report) Paragraph(str string) *Report {
	switch r.s.peek() {
	case s_body, s_head1, s_head2, s_head3, s_head4, s_head5:
		r.ws("<p>")
		r.ws(html.EscapeString(str))
		r.ws("</p>\n")
	default:
		r.seterr(errors.New("Can only create paragraph in text mode"))
	}
	return r
}
func (r *Report) Caption(str string) *Report {
	switch r.s.peek() {
	case s_fig:
		r.ws("<figcaption>")
		r.ws(html.EscapeString(str))
		r.ws("</figcaption>\n")
	default:
		r.seterr(errors.New("Caption is only supported inside image"))
	}
	return r
}

// Returns the underlying writer that may be used to add any input,
// note that if you do something here that changes the state of the Report
// your computer might explode...
func (r *Report) Writer() io.Writer {
	return r.w
}

// Let us also implement the io.Writer interface!
func (r *Report) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

func (r *Report) IsError() []error {
	return r.e
}
