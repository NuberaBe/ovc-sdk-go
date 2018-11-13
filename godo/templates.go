package godo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TemplateList []struct {
	Username    interface{} `json:"username"`
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Size        int         `json:"size"`
	Type        string      `json:"type"`
	ID          int         `json:"id"`
	AccountID   int         `json:"accountId"`
}

type TemplateService interface {
	List(int) (*TemplateList, error)
}

type TemplateServiceOp struct {
	client *OvcClient
}

func (s *TemplateServiceOp) List(accountID int) (*TemplateList, error) {
	templateMap := make(map[string]interface{})
	templateMap["accountId"] = 4
	templateJSON, err := json.Marshal(templateMap)
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/images/list", bytes.NewBuffer(templateJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var templates = new(TemplateList)
	err = json.Unmarshal(body, &templates)
	if err != nil {
		return nil, err
	}
	return templates, nil
}
