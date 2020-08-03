package sdk

import (
	"fmt"
)

// Client is a http-level type to publish event to eventbus services.
type Client struct {
	addrs []string
}

// NewClient returns an instance of Client for the specified eventbus services
//
// You can use a struct literal to set eventbus services directly,
// such as: []string{"192.168.100.1:8091", "192.168.100.2:8091"}
func NewClient(addrs []string) (*Client, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("at lease one eventbus server must be set")
	}

	this := &Client{}
	this.addrs = addrs

	return this, nil
}

// Close initiates a graceful close of the Client
func (p *Client) Close() error {
	return nil
}

// Publish synchronously post the event body to the specified eventbus service
func (p *Client) Publish(body []byte) (result []byte, err error) {
	var url string
	for _, addr := range p.addrs {
		url = fmt.Sprintf("http://%s/v1/event", addr)
		result, err = httpPost(url, body)
		if err != nil {
			continue
		}
		//if succ break it
		break
	}
	return
}
