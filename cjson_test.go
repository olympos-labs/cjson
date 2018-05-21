package cjson

import (
	"bytes"
	"testing"
)

func TestNumberFormat(t *testing.T) {
	expectNum := func(f float64, val string) {
		var buf bytes.Buffer
		_, err := writeNumber(&buf, f)
		if err != nil {
			t.Error(err)
			return
		}
		if buf.String() != val {
			t.Errorf("%f became %s, not %s", f, buf.String(), val)
		}
	}

	expectNum(-150, "-150")
	expectNum(-150.1, "-1.501E2")
	expectNum(-15, "-15")
	expectNum(-15.1, "-1.51E1")
	expectNum(-1.5, "-1.5E0")
	expectNum(-1, "-1")
	expectNum(-0.15, "-1.5E-1")
	expectNum(0, "0")
	expectNum(0.001, "1E-3")
	expectNum(0.01, "1E-2")
	expectNum(0.1, "1E-1")
	expectNum(0.11, "1.1E-1")
	expectNum(1, "1")
	expectNum(1.5, "1.5E0")
	expectNum(10, "10")
	expectNum(10.5, "1.05E1")
	expectNum(100, "100")
	expectNum(100.5, "1.005E2")
	expectNum(10000000000, "10000000000")
	expectNum(10000000000.1, "1.00000000001E10")
	expectNum(1000000000000000, "1000000000000000")
	expectNum(10000000000000000, "1E16")

}
