package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

type redirect struct {
	from, to string
}

func findRedirects(filename string) (result []redirect, err error) {
	var (
		f        *os.File
		re       *regexp.Regexp
		line     []byte
		isPrefix bool
	)
	if f, err = os.Open(filename); err != nil {
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	if re, err = regexp.Compile("^::1 (.*) #redirects-to (.*)$"); err != nil {
		return
	}
	result = []redirect{}
	for {
		if line, isPrefix, err = r.ReadLine(); err == io.EOF {
			break
		}
		if isPrefix {
			err = fmt.Errorf("line too long: %s...", line)
			break
		}
		if match := re.FindSubmatch(line); match != nil {
			result = append(result, redirect{from: string(match[1]), to: string(match[2])})
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func main() {
	redirects, err := findRedirects("/etc/hosts")
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range redirects {
		http.Handle(r.from+"/", http.RedirectHandler(r.to, http.StatusSeeOther))
	}
	if err = http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
