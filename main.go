package main

import (
	"fmt"
	"flag"
	"net"
	"net/url"
	"os"
	"io/ioutil"
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
	var host = u.Hostname()
	var port = u.Port()

	// Get endpoint
	var endpoint = u.EscapedPath()
	if endpoint == "" {
		endpoint = "/"
	} 

	// Headers
	var headers = "\r\nAccept: application/json\r\nConnection: close\r\n\r\n"

	// Send requests however many times specified (default 1)
	for i := 0; i < *profilePtr; i ++ {
		// Connects to the address provided  
		conn, err := net.Dial("tcp", host + ":" + port)	
		checkError(err)

		_, err = conn.Write([]byte("POST " + endpoint + " HTTP/1.0\r\nHost: " + host + headers ))
		checkError(err)

		result, err := ioutil.ReadAll(conn)
		checkError(err)

		fmt.Println(string(result))
	}

	


    os.Exit(0)

}