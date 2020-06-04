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

	"github.com/google/uuid"
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
	return fmt.Sprintf("Provider: %v ---> Consumer: %v", provider, consumer)
}

// RelayConnPool represents
type RelayConnPool struct {
	providerPool map[string]*RelayConn
	consumerPool map[string]*RelayConn
	mutex        sync.Mutex
}

func (thisRelayConnPool RelayConnPool) String() string {
	return fmt.Sprintf("\n%s\nProvider Pool Length: %v\n%v\nConsumer Pool Length: %v\n%v\n%s\n",
		"----- Relay Connection Pool Status -----",
		len(thisRelayConnPool.providerPool),
		thisRelayConnPool.providerPool,
		len(thisRelayConnPool.consumerPool),
		thisRelayConnPool.consumerPool,
		"----------------------------------------")
}

// registerProvider stores service provider informaion
func (thisRelayConnPool *RelayConnPool) registerProvider(newProviderConn net.Conn) (key string, newRelayConn *RelayConn) {

	thisRelayConnPool.mutex.Lock()
	defer thisRelayConnPool.mutex.Unlock()
	defer log.Println(thisRelayConnPool)

	// key = newProviderConn.RemoteAddr().String()
	key = uuid.New().String()

	for k, conn := range thisRelayConnPool.consumerPool {
		if conn.provider == nil {
			delete(thisRelayConnPool.consumerPool, k)
			newRelayConn = conn
			newRelayConn.provider = newProviderConn
			log.Printf("Starting Reply between %s\n", newRelayConn)
			return
		}
	}

	newRelayConn = new(RelayConn)
	newRelayConn.provider = newProviderConn
	newRelayConn.consumer = nil
	thisRelayConnPool.providerPool[key] = newRelayConn
	return
}

// unregisterProvider stores service provider informaion
func (thisRelayConnPool *RelayConnPool) unregisterProvider(key string) {

	thisRelayConnPool.mutex.Lock()
	defer thisRelayConnPool.mutex.Unlock()
	defer log.Println(thisRelayConnPool)

	delete(thisRelayConnPool.providerPool, key)
}

// registerConsumer retrieve service provider informaion
func (thisRelayConnPool *RelayConnPool) registerConsumer(newConsumerConn net.Conn) (key string, newRelayConn *RelayConn) {

	thisRelayConnPool.mutex.Lock()
	defer thisRelayConnPool.mutex.Unlock()
	defer log.Println(thisRelayConnPool)

	// key = newConsumerConn.RemoteAddr().String()
	key = uuid.New().String()

	for k, conn := range thisRelayConnPool.providerPool {
		if conn.consumer == nil {
			delete(thisRelayConnPool.providerPool, k)
			newRelayConn = conn
			newRelayConn.consumer = newConsumerConn
			log.Printf("Starting Reply between %s\n", newRelayConn)
			return
		}
	}

	newRelayConn = new(RelayConn)
	newRelayConn.provider = nil
	newRelayConn.consumer = newConsumerConn
	thisRelayConnPool.consumerPool[key] = newRelayConn
	return
}

// unregisterProvider stores service provider informaion
func (thisRelayConnPool *RelayConnPool) unregisterConsumer(key string) {

	thisRelayConnPool.mutex.Lock()
	defer thisRelayConnPool.mutex.Unlock()
	defer log.Println(thisRelayConnPool)

	delete(thisRelayConnPool.consumerPool, key)
}

func main() {

	// Log Setting
	logFile, err := os.OpenFile("cloud-sim-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// Business Logic
	relayConns := new(RelayConnPool)
	relayConns.providerPool = make(map[string]*RelayConn)
	relayConns.consumerPool = make(map[string]*RelayConn)
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

		log.Printf("New Provider is Connected from %s\n", conn.RemoteAddr().String())

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
			log.Printf("Error from Provider(%s) : %v", key, err)

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

		log.Printf("New Customer is Connected from %s\n", conn.RemoteAddr().String())

		key, relayConn := relayConnPool.registerConsumer(conn)
		go handleConsumer(key, relayConn, relayConnPool)
	}

	quit <- 0
}

func handleConsumer(key string, relayConn *RelayConn, relayConnPool *RelayConnPool) {

	in := relayConn.consumer
	// relay.consumer is still nil here.

	defer in.Close()

	var buf [512]byte
	for {
		// read upto 512 bytes
		n, err := in.Read(buf[0:])
		out := relayConn.provider // could be nil

		if err != nil {
			log.Printf("Error from Consumer(%s) : %v", key, err)

			if out == nil {
				relayConnPool.unregisterConsumer(key)
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

func checkError(err error) {

	_, file, line, _ := runtime.Caller(1)

	if err != nil {
		log.Printf("Fatal error (%s:%d): %s", file, line, err.Error())
		os.Exit(1)
	}
}
