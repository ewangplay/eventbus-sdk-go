package sdk

//Config is a struct of EventBus client options
type Config struct {
	// Mutiple EventBus Nodes Addresses
	// Producer will use the first node to publish message if sets multiple nodes
	// Must be set when using Producer side
	Nodes []string

	// Mutiple EventBus Seekers addresses
	// Consumer will use these seekers to discover the specified topic messages
	// It is recommended to set seekers adddress so that topics are discovered automatically when using Consumer side,
	// unless you want to connect to a single EventBus node instance
	Seekers []string

	// Maximum number of times this consumer will attempt to process a message before giving up
	// min:"0" max:"65535" default:"5"`
	MaxAttempts uint16

	// Identifiers sent to eventbus representing this client
	// (defaults: short hostname)
	ClientID string

	// Maximum number of messages to allow in flight (concurrency knob)
	// min:"0" default:"1"`
	MaxInFlight int
}
