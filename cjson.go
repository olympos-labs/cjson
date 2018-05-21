package cjson

import (
	"io"
	"math"
	"strconv"
)

func writeNumber(dst io.Writer, f float64) (int, error) {
	if -(1<<53)+1 < f && f < (1<<53)-1 {
		_, frac := math.Modf(f)
		if frac == 0.0 {
			return io.WriteString(dst, strconv.FormatInt(int64(f), 10))
		}
	}
	bs := strconv.AppendFloat([]byte{}, f, 'E', -1, 64)
	// this is kind of stupid, but oh well: We need to strip away pluses and
	// leading zeroes in exponents.
	var e int
	for i := range bs {
		if bs[i] == 'E' {
			e = i
		}
	}
	writeFrom := e + 1
	offset := 0
	hasPlus := bs[e+1] == '+'
	if hasPlus {
		offset++
		writeFrom = e + 2
	}
	hasLeadingZero := bs[e+2] == '0'
	if hasLeadingZero {
		offset++
		writeFrom = e + 3
	}
	for writeFrom < len(bs) {
		bs[writeFrom-offset] = bs[writeFrom]
		writeFrom++
	}
	bs = bs[:len(bs)-offset]
	return dst.Write(bs)
}
