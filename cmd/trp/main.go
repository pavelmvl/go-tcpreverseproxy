package main

import (
	"flag"
	"log"
	"trp/internal/tcpreverseproxy"
)

func main() {
	var serverUri string
	var clientUri string
	var bufferSize int
	flag.StringVar(&clientUri, "connect", "127.0.0.1:9999", "set URI like <ip>:<port>")
	flag.StringVar(&serverUri, "listen", "0.0.0.0:20000", "set URI like <ip>:<port>")
	flag.IntVar(&bufferSize, "bufferSize", 64*1024, "set recv buffer size")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	tcpReverseProxy := tcpreverseproxy.NewTcpReversePproxy(clientUri, serverUri, int(bufferSize))
	tcpReverseProxy.Start()
}
