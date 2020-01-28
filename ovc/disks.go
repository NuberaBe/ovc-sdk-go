package ovc

import (
	"encoding/json"
	"errors"
	"strconv"
)

// DiskConfig is used when creating a disk
type DiskConfig struct {
	AccountID   int    `json:"accountId,omitempty"`
	GridID      int    `json:"gid,omitempty"`
	MachineID   int    `json:"machineId,omitempty"`
	DiskName    string `json:"diskName,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Size        int    `json:"size,omitempty"`
	Type        string `json:"type,omitempty"`
	SSDSize     int    `json:"ssdSize,omitempty"`
	IOPS        int    `json:"iops,omitempty"`
	DiskID      int    `json:"diskId,omitempty"`
	Detach      bool   `json:"detach,omitempty"`
	Permanently string `json:"permanently,omitempty"`
}

// DiskDeleteConfig is used when deleting a disk
type DiskDeleteConfig struct {
	DiskID      int  `json:"diskId"`
	Detach      bool `json:"detach"`
	Permanently bool `json:"permanently"`
}

// DiskAttachConfig is used when attatching a disk to a machine
type DiskAttachConfig struct {
	DiskID    int `json:"diskId"`
	MachineID int `json:"machineId"`
}

// DiskInfo contains all information related to a disk
type DiskInfo struct {
	ReferenceID         string        `json:"referenceId"`
	DiskPath            string        `json:"diskPath"`
	Images              []interface{} `json:"images"`
	GUID                int           `json:"guid"`
	ID                  int           `json:"id"`
	PCIBus              int           `json:"pci_bus"`
	PCISlot             int           `json:"pci_slot"`
	AccountID           int           `json:"accountId"`
	SizeUsed            int           `json:"sizeUsed"`
	Descr               string        `json:"descr"`
	GridID              int           `json:"gid"`
	Role                string        `json:"role"`
	Params              string        `json:"params"`
	Type                string        `json:"type"`
	Status              string        `json:"status"`
	RealityDeviceNumber int           `json:"realityDeviceNumber"`
	Passwd              string        `json:"passwd"`
	Iotune              struct {
		TotalIopsSec int `json:"total_iops_sec"`
	} `json:"iotune"`
	Name    string        `json:"name"`
	SizeMax int           `json:"sizeMax"`
	Meta    []interface{} `json:"_meta"`
	ACL     struct {
	} `json:"acl"`
	Iqn           string `json:"iqn"`
	BootPartition int    `json:"bootPartition"`
	Login         string `json:"login"`
	Order         int    `json:"order"`
	Ckey          string `json:"_ckey"`
}

// DiskList is a list of disks
// Returned when using the List method
type DiskList []struct {
	Username    interface{} `json:"username"`
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Size        int         `json:"sizeMax"`
	Type        string      `json:"type"`
	ID          int         `json:"id"`
	AccountID   int         `json:"accountId"`
}

// DiskService is an interface for interfacing with the Disk
// endpoints of the OVC API
type DiskService interface {
	Resize(*DiskConfig) error
	List(int, string) (*DiskList, error)
	Get(string) (*DiskInfo, error)
	GetByName(string, int, string) (*DiskInfo, error)
	Create(*DiskConfig) (string, error)
	CreateAndAttach(*DiskConfig) (string, error)
	Attach(*DiskAttachConfig) error
	Detach(*DiskAttachConfig) error
	Update(*DiskConfig) error
	Delete(*DiskDeleteConfig) error
}

// DiskServiceOp handles communication with the disk related methods of the
// OVC API
type DiskServiceOp struct {
	client *Client
}

// List all disks
func (s *DiskServiceOp) List(accountID int, diskType string) (*DiskList, error) {
	diskMap := make(map[string]interface{})
	diskMap["accountId"] = accountID
	if len(diskType) != 0 {
		diskMap["type"] = diskType
	}

	body, err := s.client.Post("/cloudapi/disks/list", diskMap)
	if err != nil {
		return nil, err
	}

	disks := new(DiskList)
	err = json.Unmarshal(body, &disks)
	if err != nil {
		return nil, err
	}

	return disks, nil
}

// CreateAndAttach a new Disk and attaches it to a machine
func (s *DiskServiceOp) CreateAndAttach(diskConfig *DiskConfig) (string, error) {
	body, err := s.client.Post("/cloudapi/machines/addDisk", *diskConfig)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Create a new Disk
func (s *DiskServiceOp) Create(diskConfig *DiskConfig) (string, error) {
	body, err := s.client.Post("/cloudapi/disks/create", *diskConfig)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Attach attaches an existing disk to a machine
func (s *DiskServiceOp) Attach(diskAttachConfig *DiskAttachConfig) error {
	_, err := s.client.Post("/cloudapi/machines/attachDisk", *diskAttachConfig)
	return err
}

// Detach detaches an existing disk from a machine
func (s *DiskServiceOp) Detach(diskAttachConfig *DiskAttachConfig) error {
	_, err := s.client.Post("/cloudapi/machines/detachDisk", *diskAttachConfig)
	return err
}

// Update updates an existing disk
func (s *DiskServiceOp) Update(diskConfig *DiskConfig) error {
	switch {
	case diskConfig.Size != 0:
		_, err := s.client.Post("/cloudapi/disks/resize", *diskConfig)
		if err != nil {
			return err
		}

		fallthrough

	case diskConfig.IOPS != 0:
		_, err := s.client.Post("/cloudapi/disks/limitIO", *diskConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete an existing Disk
func (s *DiskServiceOp) Delete(diskConfig *DiskDeleteConfig) error {
	_, err := s.client.Post("/cloudapi/disks/delete", *diskConfig)
	return err
}

// Get individual Disk
func (s *DiskServiceOp) Get(diskID string) (*DiskInfo, error) {
	diskIDMap := make(map[string]interface{})
	diskIDInt, err := strconv.Atoi(diskID)
	if err != nil {
		return nil, err
	}
	diskIDMap["diskId"] = diskIDInt

	body, err := s.client.Post("/cloudapi/disks/get", diskIDMap)
	if err != nil {
		return nil, err
	}
	diskInfo := new(DiskInfo)
	err = json.Unmarshal(body, &diskInfo)
	if err != nil {
		return nil, err
	}

	return diskInfo, nil
}

// GetByName gets a disk by its name
func (s *DiskServiceOp) GetByName(name string, accountID int, diskType string) (*DiskInfo, error) {
	disks, err := s.client.Disks.List(accountID, diskType)
	if err != nil {
		return nil, err
	}
	for _, dk := range *disks {
		if dk.Name == name {
			did := strconv.Itoa(dk.ID)
			return s.client.Disks.Get(did)
		}
	}

	return nil, errors.New("Could not find disk based on name")
}

// Resize resizes a disk. Can only increase the size of a disk
func (s *DiskServiceOp) Resize(diskConfig *DiskConfig) error {
	_, err := s.client.Post("/cloudapi/disks/resize", *diskConfig)
	return err
}
