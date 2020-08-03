package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ewangplay/eventbus-sdk-go"
)

func main() {
	var pAddr = flag.String("addr", "127.0.0.1:8091", "eventbus server http address")
	var pFile = flag.String("file", "", "path to data file which contains the request body")
	flag.Parse()

	//Data filename must not be empty
	if *pFile == "" {
		flag.Usage()
		return
	}

	// Read data file content
	body, err := ioutil.ReadFile(*pFile)
	if err != nil {
		fmt.Printf("Read data file(%s) error: %v\n", *pFile, err)
		return
	}

	fmt.Printf("Request body:\n%s\n", body)

	c, err := sdk.NewClient([]string{*pAddr})
	if err != nil {
		fmt.Printf("New http producer error: %v\n", err)
		return
	}

	result, err := c.Publish(body)
	if err != nil {
		fmt.Printf("Publish body error: %v\n", err)
		return
	}

	fmt.Printf("Publish body succ.\nResponse result:\n%s\n", result)

	c.Close()
}
