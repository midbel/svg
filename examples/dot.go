package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/midbel/svg/dot"
)

func main() {
	scan := flag.Bool("s", false, "scan only")
	flag.Parse()

	r, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	switch {
	case *scan:
		scanOnly(r)
	default:
		parseOnly(r)
	}
}

func parseOnly(r io.Reader) {
	err := dot.Parse(r)
	if err != nil {
		fmt.Println(">> parsing fail:", err)
	}
}

func scanOnly(r io.Reader) {
	s, _ := dot.Scan(r)
	for i := 0; ; i++ {
		tok := s.Scan()
		fmt.Println(i, tok)
		if tok.IsEOF() || tok.IsInvalid() {
			break
		}
	}
}
