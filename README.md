
<div align="center">
<a href="https://gopherize.me">
<img src="assets/gopher.png" height="120" alt="gopher with moustache">
</a>

### ROTXY
Rotating Proxy Server

[![golangci-lint](https://github.com/MrMarble/rotxy/actions/workflows/golangci.yml/badge.svg)](https://github.com/MrMarble/rotxy/actions/workflows/golangci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrmarble/rotxy)](https://goreportcard.com/report/github.com/mrmarble/rotxy)
![Lines of code](https://img.shields.io/tokei/lines/github/mrmarble/rotxy)
</div>

---

Rotxy is a simple proxy server that rotates the proxy used for each request.

## Why?

I needed a proxy server that rotates the proxy used for each request. I couldn't find any that did this, so I made one.

## How?

Rotxy has 12 proxy providers built in. You can add your own by using the `--from-url` or `--from-file` flags. Rotxy will then rotate through all the available proxies for each connection.

## Usage

```bash
$ rotxy --help
Usage: rotxy

A simple proxy rotator

Flags:
  -h, --help                 Show context-sensitive help.
  -v, --verbose=INT          Enable verbose logging
      --version              Print version information and quit
  -p, --port=8080            Port to listen on.
  -h, --host="0.0.0.0"       Host to listen on.
  -s, --strategy="random"    Proxy strategy to use.

Download
  -d, --download                   Download proxies from internal list.
  -f, --from-file=FROM-FILE,...    File containing proxies.
  -u, --from-url=FROM-URL,...      URL to download proxies from.
  -U, --update-delay=1h            Update delay.

Check
  -c, --prune          Prune non-working proxies
  -y, --cycle=1        Number of cycles to prune proxies.
  -t, --timeout=30s    Timeout for proxy checks.
      --tls            Enable TLS for proxy checks.
```
