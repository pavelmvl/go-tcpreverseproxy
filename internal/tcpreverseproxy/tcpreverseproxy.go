package tcpreverseproxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type TcpReverseProxy struct {
	clientUri     string
	serverUri     string
	bufferSize    int
	clientsConn   []net.Conn
	clientsConnMu sync.Mutex
}

func NewTcpReversePproxy(connectUri, listenUri string, bSize int) *TcpReverseProxy {
	return &TcpReverseProxy{
		clientUri:   connectUri,
		serverUri:   listenUri,
		bufferSize:  bSize,
		clientsConn: make([]net.Conn, 0, 2),
	}
}

func (trp *TcpReverseProxy) Write(p []byte) (int, error) {
	var wg sync.WaitGroup
	//log.Println("Write ", len(p), " bytes")
	trp.clientsConnMu.Lock()
	defer trp.clientsConnMu.Unlock()
	for _, c := range trp.clientsConn {
		go func(conn net.Conn) {
			wg.Add(1)
			n, e := conn.Write(p)
			_ = n
			if e != nil {
				// log.Printf("connection %s:%v got %v bytes and error %v",
				// 	conn.RemoteAddr().Network(), conn.RemoteAddr().String(), n, e)
				trp.RemoveAndCloseConnection(conn, e)
			}
			wg.Done()
		}(c)
	}
	wg.Wait()
	return len(p), nil
}

func (trp *TcpReverseProxy) AddConnection(conn net.Conn) {
	trp.clientsConnMu.Lock()
	defer trp.clientsConnMu.Unlock()
	trp.clientsConn = append(trp.clientsConn, conn)
}

func (trp *TcpReverseProxy) RemoveAndCloseConnection(conn net.Conn, err error) {
	defer conn.Close()
	// trp.clientsConnMu.Lock()
	// defer trp.clientsConnMu.Unlock()
	for n, c := range trp.clientsConn {
		if c == conn {
			trp.clientsConn = append(trp.clientsConn[:n], trp.clientsConn[n+1:]...)
			log.Println("disconnected", conn.RemoteAddr().String(), "error: ", err)
		}
	}
}

func (trp *TcpReverseProxy) Start() error {
	var err error
	var listener net.Listener
	var conn net.Conn
	listener, err = net.Listen("tcp", trp.serverUri)
	if err != nil {
		return err
	}
	timeoutValue := 3 * time.Second
	t := time.Now()
	for err = net.ErrClosed; err != nil; conn, err = net.Dial("tcp", trp.clientUri) {
		if t.Sub(time.Now()) > timeoutValue {
			return fmt.Errorf("timeout %v connection to %v", timeoutValue, trp.clientUri)
		}
		log.Println(err)
		time.Sleep(100 * time.Millisecond)
	}
	go func() {
		log.Printf("Run server %v:%v", listener.Addr().Network(), listener.Addr().String())
		for {
			c, e := listener.Accept()
			if e != nil {
				// handle error
				log.Println(e)
			}
			trp.AddConnection(c)
			log.Println("New connection", c.RemoteAddr())
		}
	}()
	buffer := make([]byte, trp.bufferSize)
	for {
		written, err := io.CopyBuffer(trp, conn, buffer)
		log.Printf("copy %v bytes, error: %v\n", written, err)
	}
	//return nil
}

func (trp *TcpReverseProxy) SimpleHandlerConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("client %v connected", conn.RemoteAddr().String())
}
