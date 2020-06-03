/* ThreadedEchoServer
 */

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
)

// RelayConn represents
type RelayConn struct {
	provider net.Conn
	consumer net.Conn
}

func (thisRelayConn RelayConn) String() string {
	provider := "nil"
	if thisRelayConn.provider != nil {
		provider = thisRelayConn.provider.RemoteAddr().String()
	}
	consumer := "nil"
	if thisRelayConn.consumer != nil {
		consumer = thisRelayConn.consumer.RemoteAddr().String()
	}
	return fmt.Sprintf("(%v --> %v)", provider, consumer)
}

// RelayConnPool represents
type RelayConnPool struct {
	pool map[string]*RelayConn
	mux  sync.Mutex
	ch   chan string
}

// RegisterProvider stores service provider informaion
func (thisRelayConnPool *RelayConnPool) registerProvider(newProviderConn net.Conn) (string, *RelayConn) {
	remoteAddr := newProviderConn.RemoteAddr().String()

	thisRelayConnPool.mux.Lock()
	defer thisRelayConnPool.mux.Unlock()
	// Lock so only one goroutine at a time can access the map c.v.
	newRelayConn := new(RelayConn)
	newRelayConn.provider = newProviderConn
	newRelayConn.consumer = nil
	thisRelayConnPool.pool[remoteAddr] = newRelayConn
	log.Println("New Provider is connected from", remoteAddr)
	log.Println("Current size of pool is", len(thisRelayConnPool.pool))
	thisRelayConnPool.ch <- remoteAddr

	return remoteAddr, newRelayConn
}

// UnregisterProvider stores service provider informaion
func (thisRelayConnPool *RelayConnPool) unregisterProvider(key string) {
	<-thisRelayConnPool.ch
	thisRelayConnPool.mux.Lock()
	defer thisRelayConnPool.mux.Unlock()

	// Lock so only one goroutine at a time can access the map c.v.
	delete(thisRelayConnPool.pool, key)
	log.Printf("A Provider %s is Removed from Pool\n", key)
}

// RegisterConsumer retrieve service provider informaion
func (thisRelayConnPool *RelayConnPool) registerConsumer(newConsumerConn net.Conn) *RelayConn {

	<-thisRelayConnPool.ch
	thisRelayConnPool.mux.Lock()
	defer thisRelayConnPool.mux.Unlock()

	for key, realyConn := range thisRelayConnPool.pool {
		// provider.conn.
		if realyConn.consumer == nil {
			delete(thisRelayConnPool.pool, key)

			log.Printf("A Provider %s is Removed from Pool\n", key)
			log.Println("Current size of pool is", len(thisRelayConnPool.pool))
			realyConn.consumer = newConsumerConn
			return realyConn
		}
	}

	return nil
}

func main() {

	// Log Setting
	logFile, err := os.OpenFile("csimlog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// Business Logic
	relayConns := new(RelayConnPool)
	relayConns.pool = make(map[string]*RelayConn)
	relayConns.ch = make(chan string, 10)
	quitProvider := make(chan int)
	quitConsumer := make(chan int)

	go handleProviders(quitProvider, relayConns)
	go handleConsumers(quitConsumer, relayConns)

	<-quitProvider
	<-quitConsumer
}

func handleProviders(quit chan int, relayConnPool *RelayConnPool) {

	providerAddr := ":23001"
	listener, err := net.Listen("tcp", providerAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		key, relayConn := relayConnPool.registerProvider(conn)
		go handleProvider(key, relayConn, relayConnPool)
	}

	quit <- 0
}

func handleProvider(key string, relayConn *RelayConn, relayConnPool *RelayConnPool) {

	in := relayConn.provider
	// relay.consumer is still nil here.

	defer in.Close()

	var buf [512]byte
	for {
		// read upto 512 bytes
		n, err := in.Read(buf[0:])
		out := relayConn.consumer // could be nil

		if err != nil {
			log.Printf("Error from connection %s : %v", in.RemoteAddr().String(), err)

			if out == nil {
				relayConnPool.unregisterProvider(key)
			} else {
				out.Close()
			}
			return
		}

		if out != nil {
			// write the n bytes read
			_, err2 := out.Write(buf[0:n])
			if err2 != nil {
				return
			}
		}
	}
}

func handleConsumers(quit chan int, relayConnPool *RelayConnPool) {

	consumerAddr := ":23002"
	listener, err := net.Listen("tcp", consumerAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		relayConn := relayConnPool.registerConsumer(conn)

		log.Printf("Start Reply between %s\n", relayConn)
		go handleConnection(relayConn)
	}

	quit <- 0
}

func handleConnection(relayConn *RelayConn) {

	in := relayConn.consumer
	out := relayConn.provider

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
		log.Printf("Fatal error (%s:%d): %s", file, line, err.Error())
		os.Exit(1)
	}
}
