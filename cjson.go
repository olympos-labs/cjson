package cjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
)

// Canonicalize canonicalizes src and puts the result into dst.
func Canonicalize(dst io.Writer, src io.Reader) (int64, error) {
	c := &canonicalizer{
		dec: json.NewDecoder(src),
	}
	var written int64
	var err error
	for {
		var w int64
		w, err = c.value(dst)
		if err != nil {
			break
		}
		written += w
	}
	if err == io.EOF {
		err = nil
	}
	return written, err
}

type canonicalizer struct {
	scratch    bytes.Buffer
	dec        *json.Decoder
	needsSpace bool
}

func (c *canonicalizer) value(dst io.Writer) (int64, error) {
	tok, err := c.dec.Token()
	if err != nil {
		return 0, err
	}
	switch tok := tok.(type) {
	case string:
		c.needsSpace = false
		return c.writeString(dst, tok)
	case float64:
		var written int64
		if c.needsSpace {
			w, err := dst.Write([]byte{' '})
			written += int64(w)
			if err != nil {
				return written, err
			}
		}
		w, err := writeNumber(dst, tok)
		written += int64(w)
		c.needsSpace = true
		return written, err
	case json.Delim:
		switch tok {
		case '[':
			c.needsSpace = false
			return c.array(dst)
		case '{':
			c.needsSpace = false
			return c.object(dst)
		}
	case bool:
		var written int64
		if c.needsSpace {
			w, err := dst.Write([]byte{' '})
			written += int64(w)
			if err != nil {
				return written, err
			}
		}
		var w int
		if tok {
			w, err = dst.Write([]byte("true"))
		} else {
			w, err = dst.Write([]byte("false"))
		}
		written += int64(w)
		c.needsSpace = true
		return written, err
	default:
		if tok == nil {
			var written int64
			if c.needsSpace {
				w, err := dst.Write([]byte{' '})
				written += int64(w)
				if err != nil {
					return written, err
				}
			}
			c.needsSpace = true
			w, err := dst.Write([]byte("null"))
			written += int64(w)
			return written, err
		}
	}
	panic(fmt.Sprintf("unknown/unexpected JSON token for value %v", tok))
}

func (c *canonicalizer) array(dst io.Writer) (int64, error) {
	var written int64
	w, err := dst.Write([]byte{'['})
	written += int64(w)
	if err != nil {
		return written, err
	}
	first := true
	for {
		if !c.dec.More() {
			_, err := c.dec.Token()
			if err != nil {
				return written, err
			}
			w, err := dst.Write([]byte{']'})
			written += int64(w)
			return written, err
		}
		if !first {
			w, err = dst.Write([]byte{','})
			written += int64(w)
			if err != nil {
				return written, err
			}
		}
		first = false
		w64, err := c.value(dst)
		c.needsSpace = false
		written += w64
		if err != nil {
			return written, err
		}
	}
}

func (c *canonicalizer) object(dst io.Writer) (int64, error) {
	var values tuples
	for {
		if !c.dec.More() {
			_, err := c.dec.Token()
			if err != nil {
				return 0, err
			}
			return c.writeObject(dst, values)
		}
		var key string
		tok, err := c.dec.Token()
		if err != nil {
			return 0, err
		}
		switch tok := tok.(type) {
		case string:
			key = tok
		default:
			return 0, fmt.Errorf("Unexpected type %T (%v) reading JSON object, expected string key", tok, tok)
		}
		buf := new(bytes.Buffer)
		_, err = c.value(buf)
		c.needsSpace = false
		if err != nil {
			return 0, err
		}
		values = append(values, tuple{key: key, val: buf.Bytes()})
	}
}

type tuple struct {
	key string
	val []byte
}
type tuples []tuple

func (t tuples) Len() int           { return len(t) }
func (t tuples) Less(i, j int) bool { return t[i].key < t[j].key }
func (t tuples) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func (c *canonicalizer) writeObject(dst io.Writer, values tuples) (int64, error) {
	var written int64
	w, err := dst.Write([]byte{'{'})
	written += int64(w)
	if err != nil {
		return written, err
	}
	sort.Sort(values)
	first := true
	for _, value := range values {
		if !first {
			w, err = dst.Write([]byte{','})
			written += int64(w)
			if err != nil {
				return written, err
			}
		}
		first = false
		w64, err := c.writeString(dst, value.key)
		written += w64
		if err != nil {
			return written, err
		}
		dst.Write([]byte{':'})
		w, err = dst.Write(value.val)
		written += int64(w)
		if err != nil {
			return written, err
		}
	}
	w, err = dst.Write([]byte{'}'})
	written += int64(w)
	return written, err
}

func (c *canonicalizer) writeString(dst io.Writer, s string) (int64, error) {
	c.scratch.Reset()
	c.scratch.WriteByte('"')
	for _, r := range s {
		if r == '\\' || r == '"' {
			c.scratch.WriteByte('\\')
		}
		c.scratch.WriteRune(r)
	}
	c.scratch.WriteByte('"')
	return c.scratch.WriteTo(dst)
}

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
