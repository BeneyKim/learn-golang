package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/BeneyKim/learn-golang/morestrings"
)

const (
	SERVER_ADDRESS = "cloudsim.ntel.ml"
	//SERVER_ADDRESS = "127.0.0.1"
)

func main() {

	quit := make(chan int)

	go makeConsumer(quit)
	go makeProvider()

	<-quit

}

func makeProvider() {

	providerAddr := fmt.Sprintf("%s:%d", SERVER_ADDRESS, 23001)
	conn, err := net.Dial("tcp", providerAddr)
	checkError(err)

	data := make([]byte, 512)

	for {
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Server send : " + string(data[:n]))
		time.Sleep(time.Duration(3) * time.Second)

		fmt.Println("Returing back : " + morestrings.ReverseRunes(string(data[:n])))
		conn.Write([]byte(morestrings.ReverseRunes(string(data[0:n]))))
	}

}

func makeConsumer(quit chan int) {

	consumerAddr := fmt.Sprintf("%s:%d", SERVER_ADDRESS, 23002)
	conn, err := net.Dial("tcp", consumerAddr)
	checkError(err)

	data := make([]byte, 512)

	for {
		var s string
		fmt.Scanln(&s)
		conn.Write([]byte(s))
		time.Sleep(time.Duration(3) * time.Second)

		n, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Returned : " + string(data[:n]))
	}

	quit <- 0
}

func checkError(err error) {

	_, file, line, _ := runtime.Caller(1)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error (%s:%d): %s", file, line, err.Error())
		os.Exit(1)
	}
}
