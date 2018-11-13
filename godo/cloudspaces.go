package godo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// CloudSpaceConfig is used when creating a CloudSpace
type CloudSpaceConfig struct {
	CloudSpaceID           int     `json:"cloudspaceId,omitempty"`
	AccountID              int     `json:"accountId,omitempty"`
	Location               string  `json:"location,omitempty"`
	Name                   string  `json:"name,omitempty"`
	Access                 string  `json:"access,omitempty"`
	MaxMemoryCapacity      float64 `json:"maxMemoryCapacity,omitempty"`
	MaxCPUCapacity         int     `json:"maxCPUCapacity,omitempty"`
	MaxDiskCapacity        int     `json:"maxVDiskCapacity,omitempty"`
	MaxNetworkPeerTransfer int     `json:"maxNetworkPeerTransfer,omitempty"`
	MaxNumPublicIP         int     `json:"maxNumPublicIP,omitempty"`
	AllowedVMSizes         []int   `json:"allowedVMSizes,omitempty"`
}

// ResourceLimits contains all information related to resource limits
type ResourceLimits struct {
	CUM  float64 `json:"CU_M"`
	CUD  int     `json:"CU_D"`
	CUNP int     `json:"CU_NP"`
	CUI  int     `json:"CU_I"`
	CUC  int     `json:"CU_C"`
}

// CloudSpace contains all information related to a CloudSpace
type CloudSpace struct {
	Status            string         `json:"status"`
	UpdateTime        int            `json:"updateTime"`
	Externalnetworkip string         `json:"externalnetworkip"`
	Description       string         `json:"description"`
	ResourceLimits    ResourceLimits `json:"resourceLimits"`
	ID                int            `json:"id"`
	AccountID         int            `json:"accountId"`
	Name              string         `json:"name"`
	CreationTime      int            `json:"creationTime"`
	ACL               []struct {
		Status       string `json:"status"`
		CanBeDeleted bool   `json:"canBeDeleted"`
		Right        string `json:"right"`
		Type         string `json:"type"`
		UserGroupID  string `json:"userGroupId"`
	} `json:"acl"`
	Secret          string `json:"secret"`
	Gid             int    `json:"gid"`
	Location        string `json:"location"`
	Publicipaddress string `json:"publicipaddress"`
}

// CloudSpaceList returns a list of CloudSpaces
type CloudSpaceList []struct {
	Status            string `json:"status"`
	UpdateTime        int    `json:"updateTime"`
	Externalnetworkip string `json:"externalnetworkip"`
	Name              string `json:"name"`
	Descr             string `json:"descr"`
	CreationTime      int    `json:"creationTime"`
	ACL               []struct {
		Status       string `json:"status"`
		CanBeDeleted bool   `json:"canBeDeleted"`
		Right        string `json:"right"`
		Type         string `json:"type"`
		UserGroupID  string `json:"userGroupId"`
	} `json:"acl"`
	AccountACL struct {
		Status      string `json:"status"`
		Right       string `json:"right"`
		Explicit    bool   `json:"explicit"`
		UserGroupID string `json:"userGroupId"`
		GUID        string `json:"guid"`
		Type        string `json:"type"`
	} `json:"accountAcl"`
	Gid             int    `json:"gid"`
	Location        string `json:"location"`
	Publicipaddress string `json:"publicipaddress"`
	AccountName     string `json:"accountName"`
	ID              int    `json:"id"`
	AccountID       int    `json:"accountId"`
}

// CloudSpaceDeleteConfig used to delete a CloudSpace
type CloudSpaceDeleteConfig struct {
	CloudSpaceID int  `json:"cloudspaceId"`
	Permanently  bool `json:"permanently"`
}

// CloudSpaceService is an interface for interfacing with the CloudSpace
// endpoints of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/g8vdc/#/ApiDocs
type CloudSpaceService interface {
	Get(string) (*CloudSpace, error)
	Create(*CloudSpaceConfig) error
	Update(*CloudSpaceConfig) error
	Delete(*CloudSpaceConfig) error
}

// CloudSpaceServiceOp handles communication with the machine related methods of the
// OVC API
type CloudSpaceServiceOp struct {
	client *OvcClient
}

// Get individual CloudSpace
func (s *CloudSpaceServiceOp) Get(cloudSpaceID string) (*CloudSpace, error) {
	cloudSpaceIDMap := make(map[string]interface{})

	cloudSpaceIDInt, err := strconv.Atoi(cloudSpaceID)
	if err != nil {
		return nil, err
	}
	cloudSpaceIDMap["cloudspaceId"] = cloudSpaceIDInt
	cloudSpaceIDJson, err := json.Marshal(cloudSpaceIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/get", bytes.NewBuffer(cloudSpaceIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var cloudSpace = new(CloudSpace)
	err = json.Unmarshal(body, &cloudSpace)
	if err != nil {
		return nil, err
	}
	return cloudSpace, nil

}

// Create a new CloudSpace
func (s *CloudSpaceServiceOp) Create(cloudSpaceConfig *CloudSpaceConfig) error {
	cloudSpaceJSON, err := json.Marshal(*cloudSpaceConfig)
	log.Println(string(cloudSpaceJSON))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/create", bytes.NewBuffer(cloudSpaceJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Delete a CloudSpace
func (s *CloudSpaceServiceOp) Delete(cloudSpaceConfig *CloudSpaceDeleteConfig) error {
	cloudSpaceJSON, err := json.Marshal(*cloudSpaceConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/delete", bytes.NewBuffer(cloudSpaceJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Update an existing CloudSpace
func (s *CloudSpaceServiceOp) Update(cloudSpaceConfig *CloudSpaceConfig) error {
	cloudSpaceJSON, err := json.Marshal(*cloudSpaceConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/update", bytes.NewBuffer(cloudSpaceJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}