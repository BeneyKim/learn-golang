package main1

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	fmt.Println("Time=", time.Now())
	// GET 호출
	resp, err := http.Get("http://google.com")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// for k, v := range resp.Header {
	// 	fmt.Print(k)
	// 	fmt.Print(":")
	// 	fmt.Println(v)
	// }

	date := resp.Header["Date"][0]
	fmt.Println("Date:", date)
	now, _ := time.Parse(time.RFC1123, date)
	fmt.Println("Now:", now)
	rand.Seed(now.Unix())
	fmt.Println("My favorite number is", rand.Intn(10))
}
