.PHONY: all install reload clean

all: redirector

redirector: main.go
	go build github.com/goj/redirector

install: redirector
	cp redirector /usr/bin/
	cp systemd/redirector.service /etc/systemd/system

clean:
	rm redirector
