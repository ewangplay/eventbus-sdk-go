package sdk

import (
	"encoding/json"
	"fmt"
)

// NodesInfo structure in cluster
type NodesInfo struct {
	Producers []NodeInfo `json:"producers"`
}

// Eventbus node info structure
type NodeInfo struct {
	RemoteAddress    string   `json:"remote_address"`
	HostName         string   `json:"hostname"`
	BroadcastAddress string   `json:"broadcast_address"`
	TcpPort          uint16   `json:"tcp_port"`
	HttpPort         uint16   `json:"http_port"`
	Version          string   `json:"version"`
	Topics           []string `json:"topics"`
}

// GetNodesInfo gets the node info list of specified seeker
// Seeker address format: "192.168.10.100:4161"
func GetNodesInfo(seekerAddr string) (nodesInfo *NodesInfo, err error) {
	url := fmt.Sprintf("http://%s/nodes", seekerAddr)
	data, err := httpGet(url)
	if err != nil {
		Error("%v", err)
		return
	}

	nodesInfo = &NodesInfo{}
	err = json.Unmarshal(data, nodesInfo)
	if err != nil {
		Error("%v", err)
		return
	}
	return
}
