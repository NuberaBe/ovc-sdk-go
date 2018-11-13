package godo

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Account contains
type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AccountService is an interface for interfacing with the Account
// endpoints of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/g8vdc/#/ApiDocs
type AccountService interface {
	GetIDByName(string) (int, error)
}

// AccountServiceOp handles communication with the account related methods of the
// OVC API
type AccountServiceOp struct {
	client *OvcClient
}

// GetIDByName returns the account ID based on the account name
func (s *AccountServiceOp) GetIDByName(account string) (int, error) {
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/accounts/list", nil)
	if err != nil {
		return 0, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	var accounts = new([]Account)
	err = json.Unmarshal(body, &accounts)
	for _, acc := range *accounts {
		if acc.Name == account {
			return acc.ID, nil
		}
	}
	return -1, errors.New("Account not found")
}
