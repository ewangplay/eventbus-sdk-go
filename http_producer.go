package sdk

import (
	"fmt"
)

// HttpProducer is a http-level type to publish to eventbus node.
type HttpProducer struct {
	config *Config
}

// NewHttpProducer returns an instance of HttpProducer for the specified config
//
// You can use a struct literal to set config value directly,
// such as: &Config{Nodes: []string{"192.168.100.1:4150", "192.168.100.2:4150"}}
func NewHttpProducer(config *Config) (*HttpProducer, error) {
	if len(config.Nodes) == 0 {
		return nil, fmt.Errorf("at lease one eventbus node must be set")
	}

	this := &HttpProducer{}
	this.config = config

	return this, nil
}

// Close initiates a graceful close of the HttpProducer
func (p *HttpProducer) Close() error {
	return nil
}

// Publish synchronously post the body to the specified eventbus node's http api,
// returning the http response body in result param if publish success, otherwise an error if publish failed
func (p *HttpProducer) Publish(body []byte) (result []byte, err error) {
	var url string
	for _, node := range p.config.Nodes {
		url = fmt.Sprintf("http://%s/v1/event", node)
		result, err = httpPost(url, body)
		if err != nil {
			Error("%v", err)
			continue
		}
		//if succ break it
		break
	}
	return
}
