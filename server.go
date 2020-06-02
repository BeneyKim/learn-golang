/* ThreadedEchoServer
 */

package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
)

// Provider represents
type Provider struct {
	conn  net.Conn
	inUse bool
}

// Providers represents
type Providers struct {
	pool map[string]Provider
	mux  sync.Mutex
	ch   chan string
}

// Push stores service provider informaion
func (thisProviders *Providers) Push(newConn net.Conn) {
	remoteAddr := newConn.RemoteAddr().String()

	thisProviders.mux.Lock()
	defer thisProviders.mux.Unlock()
	// Lock so only one goroutine at a time can access the map c.v.
	newProvider := Provider{conn: newConn, inUse: false}
	thisProviders.pool[remoteAddr] = newProvider
	fmt.Println("New Provider is connected from %s", remoteAddr)
	fmt.Println("Current # of Providers is %d", len(thisProviders.pool))
	thisProviders.ch <- remoteAddr
}

// GetUnused retrieve service provider informaion
func (thisProviders *Providers) GetUnused() net.Conn {

	<-thisProviders.ch
	thisProviders.mux.Lock()
	defer thisProviders.mux.Unlock()

	for key, provider := range thisProviders.pool {
		if provider.inUse == false {
			delete(thisProviders.pool, key)
			fmt.Println("A Provider %s is Removed from Pool", key)
			fmt.Println("Current # of Providers is %d", len(thisProviders.pool))
			return provider.conn
		}
	}

	return nil
}

func main() {

	providerPool := new(Providers)
	providerPool.pool = make(map[string]Provider)
	providerPool.ch = make(chan string, 10)
	quitProvider := make(chan int)
	quitConsumer := make(chan int)

	go handleProviders(quitProvider, providerPool)
	go handleConsumers(quitConsumer, providerPool)

	<-quitProvider
	<-quitConsumer
}

func handleProviders(quit chan int, providerPool *Providers) {

	providerAddr := ":23001"
	listener, err := net.Listen("tcp", providerAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		providerPool.Push(conn)
	}

	quit <- 0
}

func handleConsumers(quit chan int, providerPool *Providers) {

	consumerAddr := ":23002"
	listener, err := net.Listen("tcp", consumerAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		provider := providerPool.GetUnused()
		fmt.Println("Start Reply between %s and %s", conn.RemoteAddr(), provider.RemoteAddr())

		go handleConnection(provider, conn)
		go handleConnection(conn, provider)
	}

	quit <- 0
}

func handleConnection(in net.Conn, out net.Conn) {

	defer in.Close()

	var buf [512]byte
	for {
		// read upto 512 bytes
		n, err := in.Read(buf[0:])
		if err != nil {
			out.Close()
			return
		}

		// write the n bytes read
		_, err2 := out.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {

	_, file, line, _ := runtime.Caller(1)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error (%s:%d): %s", file, line, err.Error())
		os.Exit(1)
	}
}
