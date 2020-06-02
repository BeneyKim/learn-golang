package main

import (
	"fmt"
	"time"

	"github.com/BeneyKim/learn-golang/morestrings"
	"github.com/google/go-cmp/cmp"
)

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

var weekday string

//var test string

func init() {
	weekday = time.Now().Weekday().String()
	fmt.Printf("Today is %s\n", weekday)
}

func init() {
	test := time.Now().Weekday().String()
	fmt.Printf("Today Test is %v\n", test)
}

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

func main() {
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
	fmt.Println(cmp.Diff("Hello World", "Hello Go"))

	for pos, char := range "日本\x80語" { // \x80 is an illegal UTF-8 encoding
		fmt.Printf("character %#U starts at byte position %d\n", char, pos)
		//fmt.Printf("Type %T %v %v\n", pos, len(char), 0)
	}

	// var test interface{}
	// test = 3
	// fmt.Println("%v", test.(type))

	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}

	a := [...]string{1: "no error", 2: "Eio", 5: "invalid argument"}
	s := []string{1: "no error", 2: "Eio", 5: "invalid argument"}
	m := map[int]string{1: "no error", 2: "Eio", 5: "invalid argument"}
	fmt.Printf("%T %v\n", a, a)
	fmt.Printf("%T %v\n", s, s)
	fmt.Printf("%T %v\n", m, m)

	fmt.Printf("Today is %s\n", weekday)
	//fmt.Printf("Today is %s %v\n", weekday, test)

}
