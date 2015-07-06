package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type redirect struct {
	from, shortTo, longTo string
}

func findRedirects(filename string) (result []*redirect, err error) {
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
	result = []*redirect{}
	for {
		if line, isPrefix, err = r.ReadLine(); err == io.EOF {
			break
		}
		if isPrefix {
			err = fmt.Errorf("line too long: %s...", line)
			break
		}
		if match := re.FindSubmatch(line); match != nil {
			var shortTo, longTo string
			allRedirects := strings.Split(string(match[2]), " ")
			shortTo = allRedirects[0]
			if len(allRedirects) > 1 {
				longTo = allRedirects[1]
			}
			result = append(result, &redirect{
				from:    string(match[1]),
				shortTo: shortTo,
				longTo:  longTo})
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func addRedirect(r *redirect) {
	http.HandleFunc(r.from+"/", func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path[1:]
		if r.longTo != "" && len(req.URL.Path) > 1 {
			http.Redirect(w, req, r.longTo+path, http.StatusSeeOther)
		} else {
			http.Redirect(w, req, r.shortTo, http.StatusSeeOther)
		}
	})
}

func main() {
	redirects, err := findRedirects("/etc/hosts")
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range redirects {
		addRedirect(r)
	}
	if err = http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
