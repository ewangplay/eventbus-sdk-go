package sdk

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

// Producer is a tcp-level type to publish to eventbus node.
//
// A Producer instance is 1:1 with a destination `eventbus node`
// and will lazily connect to that instance (and re-connect)
// when Publish commands are executed.
type Producer struct {
	nsqProducers []*nsq.Producer
	currIndex    int
}

// NewProducer returns an instance of Producer for the specified config
//
// You can use a struct literal to set config value directly,
// such as: &Config{Nodes: []string{"192.168.100.1:4150", "192.168.100.2:4150"}}
//
// NewProducer will use the first valid eventbus node to create the Producer instance
func NewProducer(config *Config) (*Producer, error) {
	if len(config.Nodes) == 0 {
		return nil, fmt.Errorf("at lease one eventbus node must be set")
	}

	cfg := nsq.NewConfig()

	var err error
	var p *nsq.Producer
	var producers []*nsq.Producer

	for _, node := range config.Nodes {
		p, err = nsq.NewProducer(node, cfg)
		if err != nil {
			Error("New producer for node(%s) fail: %v", node, err)
			continue
		}
		err = p.Ping()
		if err != nil {
			Error("Ping node(%s) fail: %v", node, err)
			continue
		}
		producers = append(producers, p)
	}

	if len(producers) == 0 {
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("no any producer instance created successfully")
		}
	}

	// Set global log options
	if gLogger != nil {
		p.SetLogger(gLogger, nsq.LogLevel(gLogLevel))
	}

	this := &Producer{}
	this.nsqProducers = producers
	this.currIndex = 0

	return this, nil
}

// Close initiates a graceful close of the Producer
//
// NOTE: this blocks until completion
func (p *Producer) Close() error {
	for _, producer := range p.nsqProducers {
		if producer != nil {
			producer.Stop()
		}
	}
	return nil
}

// Publish synchronously publishes a message body to the specified topic, returning
// an error if publish failed
func (p *Producer) Publish(topic string, body []byte) (err error) {
	if len(p.nsqProducers) == 0 {
		return fmt.Errorf("no any producer instance to publish message")
	}

	retryCount := len(p.nsqProducers)
	for {
		producer := p.nsqProducers[p.currIndex]
		err = producer.Publish(topic, body)
		if err != nil {
			retryCount--
			if retryCount == 0 {
				//publish error and reached retry max count
				return
			} else {
				//publish error and not reached retry max count
				p.currIndex++
				if p.currIndex >= len(p.nsqProducers) {
					p.currIndex = 0
				}
				continue
			}
		}

		//if publish succ, break out
		break
	}

	// Roll the current index of producers
	p.currIndex++
	if p.currIndex >= len(p.nsqProducers) {
		p.currIndex = 0
	}

	return
}

// SetLogger sets producer instance scope log options: logger and log-level
func (p *Producer) SetLogger(l ILogger, lvl LogLevel) {
	for _, producer := range p.nsqProducers {
		if producer != nil {
			producer.SetLogger(l, nsq.LogLevel(lvl))
		}
	}
}
