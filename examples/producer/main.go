package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ewangplay/eventbus-sdk-go"
)

func main() {
	var pNodeAddr = flag.String("node-address", "127.0.0.1:4150", "eventbus node address")
	var pTopic = flag.String("topic", "", "topic to be published")
	flag.Parse()

	//Topic must not be empty
	if *pTopic == "" {
		flag.Usage()
		return
	}

	// Create new Producer instance
	cfg := &sdk.Config{
		Nodes: []string{*pNodeAddr},
	}
	p, err := sdk.NewProducer(cfg)
	if err != nil {
		fmt.Printf("New producer error: %v\n", err)
		return
	}

	// Set producer instance-level log options
	logger := log.New(os.Stdout, "[Producer-Test] ", log.Ldate|log.Ltime)
	p.SetLogger(logger, sdk.LogLevelInfo)

	// Publish messages to eventbus
	bio := bufio.NewReader(os.Stdin)
	for {
		line, err := bio.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "quit" || line == "exit" {
			break
		}

		err = p.Publish(*pTopic, []byte(line))
		if err != nil {
			fmt.Printf("Publish message[%s: %s] error: %v\n", *pTopic, line, err)
			continue
		}

		fmt.Printf("Publish message[%s: %s] succ\n", *pTopic, line)
	}

	p.Close()
}
