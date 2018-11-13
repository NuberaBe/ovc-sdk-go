package godo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// PortForwardingConfig is used when creating a portforward
type PortForwardingConfig struct {
	CloudspaceID     int    `json:"cloudspaceId"`
	SourcePublicIP   string `json:"sourcePublicIp"`
	SourcePublicPort int    `json:"sourcePublicPort"`
	SourceProtocol   string `json:"sourceProtocol"`
	PublicIP         string `json:"publicIp"`
	PublicPort       int    `json:"publicPort"`
	MachineID        int    `json:"machineId"`
	LocalPort        int    `json:"localPort"`
	Protocol         string `json:"protocol"`
}

// PortForwardingList is a list of portforwards
// Returned when using the List method
type PortForwardingList []struct {
	Protocol    string `json:"protocol"`
	LocalPort   string `json:"localPort"`
	MachineName string `json:"machineName"`
	PublicIP    string `json:"publicIp"`
	LocalIP     string `json:"localIp"`
	MachineID   int    `json:"machineId"`
	PublicPort  string `json:"publicPort"`
	ID          int    `json:"id"`
}

// ForwardingService is an interface for interfacing with the portforwards
// endpoints of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/g8vdc/#/ApiDocs
type ForwardingService interface {
	Create(*PortForwardingConfig) error
	List(*PortForwardingConfig) (*PortForwardingList, error)
	Delete(*PortForwardingConfig) error
	Update(*PortForwardingConfig) error
}

// ForwardingServiceOp handles communication with the machine related methods of the
// OVC API
type ForwardingServiceOp struct {
	client *OvcClient
}

// Create a new portforward
func (c *OvcClient) Create(portForwardingConfig *PortForwardingConfig) error {
	portForwardingJSON, err := json.Marshal(*portForwardingConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.ServerURL+"/cloudapi/portforwarding/create", bytes.NewBuffer(portForwardingJSON))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Update an existing portforward
func (c *OvcClient) Update(portForwardingConfig *PortForwardingConfig) error {
	portForwardingJSON, err := json.Marshal(*portForwardingConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.ServerURL+"/cloudapi/portforwarding/updateByPort", bytes.NewBuffer(portForwardingJSON))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Delete an existing portforward
func (c *OvcClient) Delete(portForwardingConfig *PortForwardingConfig) error {
	portForwardingJSON, err := json.Marshal(*portForwardingConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.ServerURL+"/cloudapi/portforwarding/deleteByPort", bytes.NewBuffer(portForwardingJSON))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// List all portforwards
func (c *OvcClient) List(portForwardingConfig *PortForwardingConfig) (*PortForwardingList, error) {
	portForwardingJSON, err := json.Marshal(*portForwardingConfig)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.ServerURL+"/cloudapi/portforwarding/list", bytes.NewBuffer(portForwardingJSON))
	body, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var portForwardingList = new(PortForwardingList)
	err = json.Unmarshal(body, &portForwardingList)
	if err != nil {
		return nil, err
	}
	return portForwardingList, nil
}
