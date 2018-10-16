/*
Copyright 2018 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Parse URL given as command line arguments or line-by-line on standard input, and output one JSON object per line.
//
// Properties:
//   scheme
//   username
//   [password]
//   hostname
//   path
//   [query]
//   [fragment]
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func usage() {
	os.Stderr.WriteString("usage: url2jsonl <url>...\n")
	os.Exit(2)
}

func main() {
	if len(os.Args) > 1 {
		if arg := os.Args[1]; len(arg) == 0 || arg[0] == '-' {
			usage()
		}
		for _, raw := range os.Args[1:] {
			process(raw)
		}
	} else {
		// TODO check if os.Stdin is a TTY
		var buf [4096]byte
		r := bufio.NewScanner(os.Stdin)
		r.Buffer(buf[:0], 4096)
		for r.Scan() {
			process(r.Text())
		}
		if err := r.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func printError(err error) {
	errJSON, _ := json.Marshal(err.Error())
	fmt.Printf("{\"error\":%s}\n", errJSON)
}

func process(rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		printError(err)
		return
	}
	type urlJSON struct {
		Scheme   string     `json:"scheme"`
		Username string     `json:"username"`
		Password string     `json:"password,omitempty"`
		Host     string     `json:"host"`
		Path     string     `json:"path"`
		Query    url.Values `json:"query,omitempty"`
		Fragment string     `json:"fragment,omitempty"`
	}
	uj := urlJSON{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     u.EscapedPath(),
		Query:    u.Query(),
		Fragment: u.Fragment,
	}
	if ui := u.User; u.User != nil {
		uj.Username = ui.Username()
		uj.Password, _ = ui.Password()
	}
	buf, _ := json.Marshal(uj)
	os.Stdout.Write(buf)
	os.Stdout.WriteString("\n")
}
