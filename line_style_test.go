package output

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestLineStyleIsValid(t *testing.T) {
	t.Parallel()

	testhelpers.TestEnumIsValid(t, []LineStyle{
		LineStyleSolid,
		LineStyleDashed,
		LineStyleDotted,
		"invalid",
		"",
	}, []bool{
		true,
		true,
		true,
		false,
		false,
	})
}

func TestLineStyleString(t *testing.T) {
	t.Parallel()

	tests := []testhelpers.StringEnumTestCase[LineStyle]{
		{Value: LineStyleSolid, Want: "solid"},
		{Value: LineStyleDashed, Want: "dashed"},
		{Value: LineStyleDotted, Want: "dotted"},
	}

	testhelpers.TestEnumString(t, "LineStyle.String", tests, func(l LineStyle) string { return l.String() })
}
