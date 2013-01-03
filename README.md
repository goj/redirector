Web server doing HTTP redirections based on /etc/hosts file
===========================================================

How it works
------------

Redirector starts a web server that listens at standard port `80`. It searches
`/etc/hosts` file for lines matching `::1 hostname #redirects-to url` pattern.
When you type `hostname` in your browser, browser hits your machine
(`::1` is ipv6's `127.0.0.1`) on which redirector is there to give you
HTTP redirect to `url`.

It may be useful when you don't have your own DNS server, or all you
have is OpenWRT's `dnsmasq` which doesn't support CNAME records.

If you have a choice, ditch the nasty hack redirector is and use DNS.

Usage (asumes you have `systemd`-based distro):
-----------------------------------------------

1. Edit your `/etc/hosts` file. Add entries like
```
::1 g #redirects-to http://gazeta.pl
::1 m #redirects-to https://mail.google.com
::1 cups #redirects-to http://localhost:631
::1 reddit.com www.reddit.com #redirects-to http://localhost:8080/stop-slacking-off
```

2. Install `redirector`
```
make
sudo make install
```

3. Run it
```
sudo systemctl start redirector.service
```

4. Make it start on system start
```
sudo systemctl enable redirector.service
```
