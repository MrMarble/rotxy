package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/mrmarble/rotxy/pkg/proxy"
)

func Listen(port int, host string, iterator proxy.ProxyIterator, verbose bool) {
	middleProxy := goproxy.NewProxyHttpServer()
	middleProxy.Verbose = verbose

	middleProxy.ConnectDial = func(network, addr string) (net.Conn, error) {
		proxy := iterator()
		return middleProxy.NewConnectDialToProxy("http://"+proxy)(network, addr)
	}

	middleProxy.Tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		proxy := iterator()
		return middleProxy.NewConnectDialToProxy("http://"+proxy)(network, addr)
	}

	log.Printf("Listening at http://%s:%d", host, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), middleProxy)
	if err != nil {
		panic(err)
	}
}
