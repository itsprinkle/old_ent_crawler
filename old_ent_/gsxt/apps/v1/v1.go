package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"gsxt/credit"
	"gsxt/gsxt/creditd"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	wfile, err := os.Create(os.Args[2])
	r := bufio.NewReader(file)
	w := bufio.NewWriterSize(wfile, 1024000)

	i := 0
	for line, err := creditd.Readln(r); err == nil; line, err = creditd.Readln(r) {
		if line == "" {
			continue
		}
		i += 1
		if i%2000 == 0 {
			println(i)
			w.Flush()
		}
		var v2 credit.InfoV2
		if err := json.NewDecoder(strings.NewReader(line)).Decode(&v2); err != nil {
			println(i, err.Error())
			continue
		} else {
			v1 := creditd.ToV1(v2)
			body, err := json.Marshal(v1)
			if err == nil {
				w.Write(body)
				w.WriteString("\n")
			}
		}
	}
	w.Flush()
	wfile.Close()
	file.Close()
}
