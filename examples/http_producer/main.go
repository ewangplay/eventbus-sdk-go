package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ewangplay/eventbus-sdk-go"
)

func main() {
	var pNodeHttpAddr = flag.String("node-http-address", "127.0.0.1:8091", "eventbus node http address")
	var pDataFile = flag.String("data-file", "", "path to data file which contains the request body")
	flag.Parse()

	//Data filename must not be empty
	if *pDataFile == "" {
		flag.Usage()
		return
	}

	// Read data file content
	body, err := ioutil.ReadFile(*pDataFile)
	if err != nil {
		fmt.Printf("Read data file(%s) error: %v\n", *pDataFile, err)
		return
	}

	fmt.Printf("HTTP Request body:\n%s\n", body)

	cfg := &sdk.Config{
		Nodes: []string{*pNodeHttpAddr},
	}

	p, err := sdk.NewHttpProducer(cfg)
	if err != nil {
		fmt.Printf("New http producer error: %v\n", err)
		return
	}

	result, err := p.Publish(body)
	if err != nil {
		fmt.Printf("Publish body error: %v\n", err)
		return
	}

	fmt.Printf("Publish body succ.\nHTTP Response result:\n%s\n", result)

	p.Close()
}
