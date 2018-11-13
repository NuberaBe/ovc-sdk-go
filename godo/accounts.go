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

// getAccountId returns the account ID based on the account name
func (c *OvcClient) getAccountID(account string) (int, error) {
	req, err := http.NewRequest("POST", c.ServerURL+"/cloudapi/accounts/list", nil)
	if err != nil {
		return 0, err
	}
	body, err := c.Do(req)
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
