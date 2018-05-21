package cjson

import (
	"strings"
	"testing"
)

func TestNumberFormat(t *testing.T) {
	expectNum := func(f float64, val string) {
		var buf strings.Builder
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

func TestOkCanonicalize(t *testing.T) {
	expect := func(in, out string) {
		var buf strings.Builder
		_, err := Canonicalize(&buf, strings.NewReader(in))
		if err != nil {
			t.Error(err)
			return
		}
		if buf.String() != out {
			t.Errorf("%v became %v, not %v", in, buf.String(), out)
		}
	}
	expect("null", "null")
	expect(`[   null]`, `[null]`)
	expect(`{"b": {"hello": "world"}, "a": 10}`, `{"a":10,"b":{"hello":"world"}}`)
	expect(`[100.5, 1.5, 1.0, 1.0E0, 1.0e+0, 0.1e1]`, `[1.005E2,1.5E0,1,1,1,1]`)
	expect(`{"ab": [1, 2],"aba":1}
{"x":"y"}`, `{"ab":[1,2],"aba":1}{"x":"y"}`)
}
