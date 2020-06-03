package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/BeneyKim/learn-golang/morestrings"
)

func main() {

	quit := make(chan int)
	go makeProviderOnly(quit)

	<-quit

}

func makeProviderOnly(quit chan int) {

	providerAddr := "127.0.0.1:23001"
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

	quit <- 0

}

func checkError(err error) {

	_, file, line, _ := runtime.Caller(1)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error (%s:%d): %s", file, line, err.Error())
		os.Exit(1)
	}
}
