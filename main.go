package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func printHelp(help bool) {

	if (help == true) {
		fmt.Println("Options: ")
		fmt.Println("--help, Provides information on how to use this")
		fmt.Println("--url, Provide a URL to perform a POST Request")
		fmt.Println("--profile, Provides the number of times to perform the POST Request \n")
		return 
	}

	return 

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
	
	printHelp(*helpPtr)

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

	// Profile Variables
	var min int64 = 0
	var max int64 = 0
	var mean int64 = 0
	var median []int64
	var successes float64 = 0
	var errorCodes []string
	var minSize int = 0
	var maxSize int = 0

	// Send requests however many times specified (default 1)
	for i := 0; i < *profilePtr; i++ {
		// Connects to the address provided  
		conn, err := net.Dial("tcp", host + ":" + port)	
		checkError(err)

		var start = time.Now()
		_, err = conn.Write([]byte("POST " + endpoint + " HTTP/1.0\r\nHost:" + host + headers ))
		checkError(err)

		var time int64 = time.Since(start).Nanoseconds()
		mean += time
		median = append(median, time)

		// Check for minimum and maximum times
		if (min == 0 || time < min) {
			min = time
		}
		if (max == 0 || time > max) {
			max = time
		}

		result, err := ioutil.ReadAll(conn)
		checkError(err)

		// String processing to get error codes and size
		var resultArray []string
		resultArray = strings.Split(strings.TrimSpace(string(result)), "\n")
		for i := range resultArray { 

			// Retrieving error codes
			if strings.Contains(resultArray[i], "200 OK") {
				successes+= 1
			} else if (strings.Contains(resultArray[i], "HTTP/1.1") && !strings.Contains(resultArray[i], "200 OK")) {
				var temp string = strings.Split(resultArray[i], "HTTP/1.1")[1]
				errorCodes = append(errorCodes, strings.TrimSpace(temp))
			}

			if strings.Contains(resultArray[i], "Content-Length:") {
				var temp, err = strconv.Atoi(strings.TrimSpace(strings.Split(resultArray[i], "Content-Length:")[1]))
				checkError(err)

				// Checking for minimum and maximum size responses
				if temp < minSize || minSize == 0 {
					minSize = temp
				}
				if temp > maxSize {
					maxSize = temp
				}
			}
		}
	}

	medianInt := make([]int, len(median))
	for i, val := range median {
		medianInt[i] = int(val)
	  }
	sort.Ints(medianInt)

	// Calculate median
	var medianValue int
	if (*profilePtr % 2 == 1) {
		medianValue = medianInt[*profilePtr / 2]
	} else {
		medianValue = ( medianInt[*profilePtr / 2] + medianInt[(*profilePtr - 1) / 2 ] ) / 2
	}


	// Print Results
	fmt.Printf("Request Number: %d\n", *profilePtr)
	fmt.Printf("Fastest Request (µs): %d\n", min)
	fmt.Printf("Slowest Request (µs): %d\n", max)
	fmt.Printf("Mean Request Time (µs): %d\n", mean / int64(*profilePtr))
	fmt.Printf("Median Request Time (µs): %d\n", medianValue)
	fmt.Printf("Percent Successful Requests: %.3f\n", successes / float64(*profilePtr))
	fmt.Printf("Error Codes: %v\n", errorCodes)
	fmt.Printf("Size of Largest Response: %d\n", maxSize)
	fmt.Printf("Size of Smallest Response: %d\n", minSize)

    os.Exit(0)

}