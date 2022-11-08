package main

import (
	"log"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mrmarble/rotxy/pkg/proxy"
	"github.com/mrmarble/rotxy/pkg/server"
)

var (
	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong) error {
	log.Printf("rotxy has version %s built from %s on %s\n", version, commit, date)
	app.Exit(0)

	return nil
}

type CLI struct {
	Verbose int         `help:"Enable verbose logging" short:"v" type:"counter"`
	Version VersionFlag `name:"version" help:"Print version information and quit"`

	Port     int    `help:"Port to listen on." default:"8080" short:"p"`
	Host     string `help:"Host to listen on." default:"0.0.0.0" short:"h"`
	Strategy string `help:"Proxy strategy to use." default:"random" enum:"random,round-robin" short:"s"`

	// Download
	Download bool          `help:"Download proxies from internal list." short:"d" group:"Download"`
	FromFile []string      `help:"File containing proxies." type:"existingfile" short:"f" name:"from-file" group:"Download"`
	FromURL  []string      `help:"URL to download proxies from." short:"u" name:"from-url" group:"Download"`
	Updelay  time.Duration `help:"Update delay." default:"1h" short:"U" name:"update-delay" group:"Download"`

	// Prune proxies
	Prune   bool          `help:"Prune non-working proxies" default:"false" short:"c" group:"Check"`
	Cycle   uint          `help:"Number of cycles to prune proxies." default:"1" group:"Check" short:"y"`
	Timeout time.Duration `help:"Timeout for proxy checks." default:"30s" group:"Check" short:"t"`
	TLS     bool          `help:"Enable TLS for proxy checks." default:"false" group:"Check"`
}

func (c *CLI) Run() error {
	if !c.Download && c.FromFile == nil && len(c.FromURL) == 0 {
		log.Println("No proxies specified. Use -d, -f or -u to specify proxies.")
		return nil
	}
	sources := []string{}
	if c.Download {
		sources = append(sources, proxy.InternalProxyList...)
	}
	sources = append(sources, c.FromURL...)
	sources = append(sources, c.FromFile...)

	proxyList := proxy.NewProxyList(sources)
	log.Println("Parsing proxies from", len(sources), "sources...")
	proxyList.Update()
	log.Println("Found", len(proxyList.Proxies), "proxies")

	if c.Prune {
		log.Println("Pruning non-working proxies (tls:", c.TLS, c.Cycle, ", cycle(s), timeout:", c.Timeout, ")")
		proxyList.Prune(c.Cycle, c.Timeout, c.TLS)
		log.Println("Found", len(proxyList.Proxies), "working proxies")
	}

	if len(proxyList.Proxies) == 0 {
		log.Println("No proxies available")
		return nil
	}

	if c.Updelay != 0 {
		go update(proxyList, c.Updelay, c.Prune, c.TLS, c.Timeout, c.Cycle)
	}

	server.Listen(c.Port, c.Host, proxy.Provider(c.Strategy, proxyList), c.Verbose > 0)
	return nil
}

func main() {
	// If running without any extra arguments, default to the --help flag
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	cli := &CLI{}
	ctx := kong.Parse(cli, kong.UsageOnError(), kong.Name("rotxy"), kong.Description("A simple proxy rotator"))
	err := ctx.Run()

	ctx.FatalIfErrorf(err)
}

func update(proxies *proxy.ProxyList, delay time.Duration, check, tls bool, timeout time.Duration, cycle uint) {
	ticker := time.NewTicker(delay)
	for {
		<-ticker.C
		log.Println("Updating proxies from sources...")
		proxies.Update()
		if check {
			proxies.Prune(cycle, timeout, tls)
		}
	}
}
