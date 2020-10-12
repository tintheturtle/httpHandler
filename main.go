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
	helpPtr := flag.Bool("help", false, "Help Information")
	urlPtr := flag.String("url", "", "URL to fetch")
	flag.Parse()

	validation(*helpPtr, *urlPtr)

	
	u, err := url.Parse(*urlPtr)
	checkError(err)

	conn, err := net.Dial("tcp", u.Hostname()+":"+u.Port())	
	checkError(err)

	var endpoint string
	if u.EscapedPath() == "" {
		endpoint = "/"
	} else {
		endpoint = u.EscapedPath()
	}

	_, err = conn.Write([]byte("POST " + endpoint + " HTTP/1.0\r\nHost: " + u.Hostname() + "\r\nAccept: application/json\r\nConnection: close\r\n\r\n"))
	checkError(err)

	result, err := ioutil.ReadAll(conn)
	checkError(err)

    fmt.Println(string(result))

    os.Exit(0)

}