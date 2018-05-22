package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"olympos.io/encoding/cjson"
)

var version = flag.Bool("version", false, "print json_canonicalize version/help")
var help = flag.Bool("help", false, "print json_canonicalize version/help")
var newlines = flag.Bool("n", false, "print a newline after each standalone JSON element")

func main() {
	flag.Parse()
	if *version || *help {
		flag.Usage()
		os.Exit(2)
	}
	enc := cjson.NewEncoder(os.Stdout)
	dec := json.NewDecoder(os.Stdin)
	enc.SetStreamSpace(!*newlines)
	var err error
	for {
		var val interface{}
		err = dec.Decode(&val)
		if err != nil {
			break
		}
		err = enc.Encode(val)
		if err != nil {
			break
		}
		if *newlines {
			fmt.Fprintln(os.Stdout)
		}
	}
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Reader error: %s\n", err.Error())
		os.Exit(1)
	}
}
