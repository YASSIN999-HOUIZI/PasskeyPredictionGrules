package main

import (
	// "fmt"
	"os"
	s "strings"
)

func main(){
	args := os.Args[1:]
	if s.ToLower(args[0]) == "server" {
		runServer()
	} else {
		test()
	}
	
	// fmt.Println("Device information saved to device_info.csv")
}
