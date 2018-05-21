package cjson

import (
	"strings"
	"testing"
)

func TestNumberFormat(t *testing.T) {
	expect := func(f float64, val string) {
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

	expect(-150, "-150")
	expect(-150.1, "-1.501E2")
	expect(-15, "-15")
	expect(-15.1, "-1.51E1")
	expect(-1.5, "-1.5E0")
	expect(-1, "-1")
	expect(-0.15, "-1.5E-1")
	expect(-0.0, "0")
	expect(0, "0")
	expect(0.001, "1E-3")
	expect(0.01, "1E-2")
	expect(0.1, "1E-1")
	expect(0.11, "1.1E-1")
	expect(1, "1")
	expect(1.5, "1.5E0")
	expect(10, "10")
	expect(10.5, "1.05E1")
	expect(100, "100")
	expect(100.5, "1.005E2")
	expect(10000000000, "10000000000")
	expect(10000000000.1, "1.00000000001E10")
	expect(1000000000000000, "1000000000000000")
	expect(10000000000000000, "1E16")
	expect(1-(1<<53), "-9007199254740991")
	expect(-(1 << 53), "-9.007199254740992E15")
	expect((1<<53)-1, "9007199254740991")
	expect((1 << 53), "9.007199254740992E15")
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
	expect(`null null true false 1e1 1.5`, `null null true false 10 1.5E0`)
	expect(``, ``)
	expect(`
`, ``)
	expect(`99.9 [1,2] 3 [] null {} false true {} false`, `9.99E1[1,2]3[]null{}false true{}false`)
	expect(`-0.0 -0e0`, `0 0`)
}
