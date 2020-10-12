package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"time"
)

func validation(help bool, url string) (int) {

	if (help == true) {
		fmt.Println("Hello")
		return 0
	}

	return 1

}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}


func main() {
	// CLI Arguments processing 
	helpPtr := flag.Bool("help", false, "Help Information")
	urlPtr := flag.String("url", "", "URL to fetch")
	profilePtr := flag.Int("profile", 1, "Number of requests to be made")
	flag.Parse()
	validation(*helpPtr, *urlPtr)

	// URL String Parsing
	u, err := url.Parse(*urlPtr)
	checkError(err)
	var host string = u.Hostname()
	var port string = u.Port()

	// Get endpoint
	var endpoint string = u.EscapedPath()
	if endpoint == "" {
		endpoint = "/"
	} 

	// Headers
	var headers = "\r\nAccept: application/json\r\nConnection: close\r\n\r\n"

	var min int64 = 0
	var max int64 = 0

	var mean int64 = 0

	// Send requests however many times specified (default 1)
	for i := 0; i < *profilePtr; i ++ {
		// Connects to the address provided  
		conn, err := net.Dial("tcp", host + ":" + port)	
		checkError(err)

		var start = time.Now()
		_, err = conn.Write([]byte("POST " + endpoint + " HTTP/1.0\r\nHost: " + host + headers ))
		checkError(err)
		var end = time.Now()

		var time int64 = time.Since(start).Nanoseconds()
		mean += time

		// Check for minimum and maximum times
		if (min == 0 || time < min) {
			min = time
		}
		if (max == 0 || time > max) {
			max = time
		}
		

		result, err := ioutil.ReadAll(conn)
		checkError(err)

		fmt.Println(string(result[1]))
		fmt.Println(end.Sub(start))

	}

	// Print Results
	fmt.Printf("Request Number: %d\n", *profilePtr)
	fmt.Printf("Fastest Request (µs): %d\n", min)
	fmt.Printf("Slowest Request (µs): %d\n", max)
	fmt.Printf("Mean Request Time (µs): %d\n", mean / int64(*profilePtr))




    os.Exit(0)

}