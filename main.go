package main

import (
	"fmt"
	"flag"
)

func validation(help bool, url string) (int) {

	if (help == true) {
		fmt.Println("Hello")
		return 0
	}

	return 1

}

func main() {
	urlPtr := flag.String("text", "", "URL to query")
	helpPtr := flag.Bool("help", false, "Help Information")
	flag.Parse()

	number := validation(*helpPtr, *urlPtr)



	
	fmt.Printf("url: %s, help: %t, exit: %v\n", *urlPtr, *helpPtr, number )
}