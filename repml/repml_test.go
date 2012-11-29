package repml

import (
	"testing"
	"os"
	"io"
)

func TestProduce(t *testing.T) {
	f, e := os.OpenFile("./test.html", os.O_CREATE | os.O_TRUNC | os.O_WRONLY,0)
	if e!=nil {
		t.Fatal(e)
	}
	defer f.Close()

	r := New(f,"My Test")

	r.Heading("Head 1")
	r.Paragraph(`In an expression switch, the switch expression is evaluated and the case expressions, which need not be constants, are evaluated left-to-right and top-to-bottom; the first one that equals the switch expression triggers execution of the statements of the associated case; the other cases are skipped. If no case matches and there is a "default" case, its statements are executed. There can be at most one default case and it may appear anywhere in the "switch" statement. A missing switch expression is equivalent to the expression true. `)
	r.Section()
	r.Heading("Sub 1")
	r.Paragraph("ADSFADSF sdf asd fasd ");
	r.Heading("Sub 2")
	r.Paragraph("ADSFADSF sdf asd fasd ");
	r.Figure()
	io.WriteString(r.Writer(),`
<svg style="border: 1px solid black; overflow: hidden;height: 6cm" viewbox="-1 -1 2 2" preserveAspectRatio="xMidYMid meet">
<circle cx="0" cy="0" r="1" />
</svg>

`);
io.WriteString(r.Writer(),`
<svg style="border: 1px solid black; overflow: hidden; height:6cm;" viewbox="0 0 33 33" preserveAspectRatio="xMidYMid meet">
<g transform="translate(-33,178) scale(1,-1)">
	<path d="M33,145 L66,178 L70,178"/>
</g>
</svg>

`);
	r.Caption("This is a figure")
	r.End()
	r.End()
	es := r.Finish()
	if es != nil {
		for _, v := range es {
			t.Error(v)
		}
	}

}
