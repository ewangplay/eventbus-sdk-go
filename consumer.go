package sdk

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

// Message is the fundamental data type containing
// the id, body, and other metadata
type Message struct {
	ID        string
	Body      []byte
	Timestamp int64
	Attempts  uint16
}

// Consumer is a high-level type to consume from EventBus.
//
// If configured, it will poll eventbus seeker instances and handle connection (and
// reconnection) to any discovered eventbus nodes.
type Consumer struct {
	nsqConsumer *nsq.Consumer
	handler     *Handler
}

// NewConsumer creates a new instance of Consumer for the specified topic/channel
//
// You can use a struct literal to set config value directly,
// such as: &Config{Seekers: []string{"192.168.110.10:4160"}}
//
// If the eventbus nodes and seeker set up at the same time, NewConsumer will connect to
// eventbus seeker preferentially.
func NewConsumer(topic string, channel string, config *Config) (*Consumer, error) {
	// Check the params
	if topic == "" {
		Error("Topic must not be empty")
		return nil, fmt.Errorf("topic must not be empty")
	}
	if channel == "" {
		channel = fmt.Sprintf("%s#ephemeral", GenerateUUID())
		Info("Input channel name is empty, using the auto-generated name: %s", channel)
	}
	if len(config.Seekers) == 0 && len(config.Nodes) == 0 {
		Error("Eventbus seekers and nodes can not be empty at the same time")
		return nil, fmt.Errorf("eventbus seekers and nodes can not be empty at the same time")
	}

	maxInFlight := evaluateMaxInFlight(config)
	Info("Evaluate max in flight: %d", maxInFlight)

	// Build the nsq config
	cfg := nsq.NewConfig()
	if maxInFlight > 0 {
		cfg.MaxInFlight = maxInFlight
	}
	if config.ClientID != "" {
		cfg.ClientID = config.ClientID
	}
	if config.MaxAttempts > 0 {
		cfg.MaxAttempts = config.MaxAttempts
	}

	// Create nsq consumer instance
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		Error("%v", err)
		return nil, err
	}

	// Set global log options
	if gLogger != nil {
		c.SetLogger(gLogger, nsq.LogLevel(gLogLevel))
	}

	// Add message handler to nsq consumer instance
	var handler *Handler
	if maxInFlight > 0 {
		handler = NewHandler(maxInFlight)
	} else {
		handler = NewHandler(1)
	}
	c.AddHandler(handler)

	// Connect to eventbus seekers or nodes instance
	if len(config.Seekers) > 0 {
		if len(config.Seekers) == 1 {
			err = c.ConnectToNSQLookupd(config.Seekers[0])
			if err != nil {
				Error("%v", err)
				return nil, fmt.Errorf("connect to eventbus seeker[%s] fail: %v", config.Seekers[0], err)
			}
		} else {
			err = c.ConnectToNSQLookupds(config.Seekers)
			if err != nil {
				Error("%v", err)
				return nil, fmt.Errorf("connect to eventbus seeker[%s] fail: %v", config.Seekers, err)
			}
		}
	} else {
		if len(config.Nodes) == 1 {
			err = c.ConnectToNSQD(config.Nodes[0])
			if err != nil {
				Error("%v", err)
				return nil, fmt.Errorf("connect to eventbus node[%s] fail: %v", config.Nodes[0], err)
			}
		} else {
			err = c.ConnectToNSQDs(config.Nodes)
			if err != nil {
				Error("%v", err)
				return nil, fmt.Errorf("connect to eventbus nodes[%s] fail: %v", config.Nodes, err)
			}
		}
	}

	this := &Consumer{}
	this.nsqConsumer = c
	this.handler = handler

	return this, nil
}

// Close will initiate a graceful stop of the Consumer
func (c *Consumer) Close() error {
	if c.nsqConsumer != nil {
		c.nsqConsumer.Stop()
		<-c.nsqConsumer.StopChan
	}
	return nil
}

//Consume will return a message channel which you can get each message from it
func (c *Consumer) Consume() (msgs <-chan *Message, err error) {
	if c.handler != nil {
		return c.handler.GetMessage(), nil
	}
	return nil, fmt.Errorf("consumer instance invalid")
}

// SetLogger sets consumer instance scope log options: logger and log-level
func (c *Consumer) SetLogger(l ILogger, lvl LogLevel) {
	if c.nsqConsumer != nil {
		c.nsqConsumer.SetLogger(l, nsq.LogLevel(lvl))
	}
}

func evaluateMaxInFlight(config *Config) int {
	var maxInFlight int

	if len(config.Seekers) > 0 {
		for _, seeker := range config.Seekers {
			nodesInfo, err := GetNodesInfo(seeker)
			if err == nil {
				maxInFlight += len(nodesInfo.Producers)
			} else {
				//Assume eventbus cluster contains 5 nodes
				maxInFlight += 5
			}
		}
	} else {
		maxInFlight = len(config.Nodes)
	}

	if maxInFlight < config.MaxInFlight {
		maxInFlight = config.MaxInFlight
	}

	return maxInFlight
}
