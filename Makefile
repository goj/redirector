.PHONY: all install clean

all: redirector

redirector:
	go build github.com/goj/redirector

install: redirector
	cp redirector /usr/bin/
	cp systemd/redirector.service /etc/systemd/system

clean:
	rm redirector
