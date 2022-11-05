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

func Listen(port int, iterator proxy.ProxyIterator) {
	middleProxy := goproxy.NewProxyHttpServer()
	middleProxy.Verbose = false

	middleProxy.ConnectDial = func(network, addr string) (net.Conn, error) {
		proxy := iterator()
		return middleProxy.NewConnectDialToProxy("http://"+proxy)(network, addr)
	}

	middleProxy.Tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		proxy := iterator()
		return middleProxy.NewConnectDialToProxy("http://"+proxy)(network, addr)
	}

	log.Println("Serving at port", port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), middleProxy)
	if err != nil {
		panic(err)
	}
}
