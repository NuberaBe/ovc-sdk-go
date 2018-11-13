package godo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// MachineList is a list of machines
// Returned when using the List method
type MachineList []struct {
	Status      string `json:"status"`
	StackID     int    `json:"stackId"`
	UpdateTime  int    `json:"updateTime"`
	ReferenceID string `json:"referenceId"`
	Name        string `json:"name"`
	Nics        []struct {
		Status      string `json:"status"`
		MacAddress  string `json:"macAddress"`
		ReferenceID string `json:"referenceId"`
		DeviceName  string `json:"deviceName"`
		Type        string `json:"type"`
		Params      string `json:"params"`
		NetworkID   int    `json:"networkId"`
		GUID        string `json:"guid"`
		IPAddress   string `json:"ipAddress"`
	} `json:"nics"`
	SizeID       int   `json:"sizeId"`
	Disks        []int `json:"disks"`
	CreationTime int   `json:"creationTime"`
	ImageID      int   `json:"imageId"`
	Storage      int   `json:"storage"`
	Vcpus        int   `json:"vcpus"`
	Memory       int   `json:"memory"`
	ID           int   `json:"id"`
}

// MachineConfig is used when creating a machine
type MachineConfig struct {
	MachineID    string        `json:"machineId"`
	CloudspaceID int           `json:"cloudspaceId"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	SizeID       int           `json:"sizeId"`
	ImageID      int           `json:"imageId"`
	Disksize     int           `json:"disksize"`
	DataDisks    []interface{} `json:"datadisks"`
	Permanently  string        `json:"permanently"`
}

// MachineInfo contains all information related to a cloudspace
type MachineInfo struct {
	Cloudspaceid int    `json:"cloudspaceid"`
	Status       string `json:"status"`
	UpdateTime   int    `json:"updateTime"`
	Hostname     string `json:"hostname"`
	Locked       bool   `json:"locked"`
	Name         string `json:"name"`
	CreationTime int    `json:"creationTime"`
	Sizeid       int    `json:"sizeid"`
	Disks        []struct {
		Status  string `json:"status,omitempty"`
		SizeMax int    `json:"sizeMax,omitempty"`
		Name    string `json:"name,omitempty"`
		Descr   string `json:"descr,omitempty"`
		ACL     struct {
		} `json:"acl"`
		Type string `json:"type"`
		ID   int    `json:"id"`
	} `json:"disks"`
	Storage int `json:"storage"`
	ACL     []struct {
		Status       string `json:"status"`
		CanBeDeleted bool   `json:"canBeDeleted"`
		Right        string `json:"right"`
		Type         string `json:"type"`
		UserGroupID  string `json:"userGroupId"`
	} `json:"acl"`
	OsImage  string `json:"osImage"`
	Accounts []struct {
		GUID     string `json:"guid"`
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"accounts"`
	Interfaces []struct {
		Status      string `json:"status"`
		MacAddress  string `json:"macAddress"`
		ReferenceID string `json:"referenceId"`
		DeviceName  string `json:"deviceName"`
		IPAddress   string `json:"ipAddress"`
		Params      string `json:"params"`
		NetworkID   int    `json:"networkId"`
		GUID        string `json:"guid"`
		Type        string `json:"type"`
	} `json:"interfaces"`
	Imageid     int         `json:"imageid"`
	ID          int         `json:"id"`
	Description interface{} `json:"description"`
}

// MachineService is an interface for interfacing with the Machine
// endpoints of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/g8vdc/#/ApiDocs
type MachineService interface {
	List() (*MachineList, error)
	Get(id string) (*MachineInfo, error)
	Create(*MachineConfig) error
	Update(*MachineConfig) error
	Delete(*MachineConfig) error
	Template(string, *MachineConfig) error
}

// MachineServiceOp handles communication with the machine related methods of the
// OVC API
type MachineServiceOp struct {
	client *OvcClient
}

// List all machines
func (s *MachineServiceOp) List(cloudSpaceID int) (*MachineList, error) {
	cloudSpaceIDMap := make(map[string]interface{})
	cloudSpaceIDMap["cloudspaceId"] = cloudSpaceID
	cloudSpaceIDJSON, err := json.Marshal(cloudSpaceIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/list", bytes.NewBuffer(cloudSpaceIDJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var machines = new(MachineList)
	err = json.Unmarshal(body, &machines)
	if err != nil {
		return nil, err
	}
	return machines, nil
}

// Get individual machine
func (s *MachineServiceOp) Get(id string) (*MachineInfo, error) {
	machineIDMap := make(map[string]interface{})
	machineIDMap["machineId"], _ = strconv.Atoi(id)
	machineIDJson, err := json.Marshal(machineIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/get", bytes.NewBuffer(machineIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var machineInfo = new(MachineInfo)
	err = json.Unmarshal(body, &machineInfo)
	if err != nil {
		return nil, err
	}
	return machineInfo, nil
}

// Create a new machine
func (s *MachineServiceOp) Create(machineConfig *MachineConfig) error {
	machineJSON, err := json.Marshal(*machineConfig)
	log.Println(string(machineJSON))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/create", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Update an existing machine
func (s *MachineServiceOp) Update(machineConfig *MachineConfig) error {
	machineJSON, err := json.Marshal(*machineConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/update", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return err
}

// Delete an existing machine
func (s *MachineServiceOp) Delete(machineConfig *MachineConfig) error {
	machineJSON, err := json.Marshal(*machineConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/delete", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Template creates an image of the existing machine
func (s *MachineServiceOp) Template(machineConfig *MachineConfig, templateName string) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = machineConfig.MachineID
	machineMap["templateName"] = templateName
	machineJSON, err := json.Marshal(machineMap)
	if err != nil {
		return nil
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/delete", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}