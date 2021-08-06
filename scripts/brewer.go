package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

const colorbrewer = "https://colorbrewer2.org/export/colorbrewer.json"

func main() {
	name := flag.String("n", "colors", "package name")
	flag.Parse()
	if err := os.MkdirAll(filepath.Dir(flag.Arg(0)), 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	w, err := os.Create(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer w.Close()

	if err := writeColours(w, *name); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func writeColours(w io.Writer, name string) error {
	res, err := http.Get(colorbrewer)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var (
		dat = make(map[string]interface{})
		ws  = bufio.NewWriter(w)
	)
	defer ws.Flush()

	fmt.Fprintf(ws, "package %s", name)
	fmt.Fprintln(ws)

	if err := json.NewDecoder(res.Body).Decode(&dat); err != nil {
		return err
	}
	keys := sortKeys(dat)
	for _, k := range keys {
		xs, ok := dat[k].(map[string]interface{})
		if !ok {
			continue
		}
		keys2 := sortKeys(xs)
		for _, x := range keys2 {
			xs, ok := xs[x].([]interface{})
			if !ok {
				continue
			}
			fmt.Fprintf(ws, "var %s%s = []string{", k, x)
			fmt.Fprintln(ws)
			for _, x := range xs {
				fmt.Fprintf(ws, "\t\"%s\",", x)
				fmt.Fprintln(ws)
			}
			fmt.Fprintln(ws, "}")
			fmt.Fprintln(ws)
		}
	}
	fmt.Fprintln(ws)
	return nil
}

func sortKeys(data map[string]interface{}) []string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
