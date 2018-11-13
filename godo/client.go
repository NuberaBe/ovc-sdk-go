package godo

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Hostname     string
	ClientId     string
	ClientSecret string
}

// Credentials used to authenticate
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// OvcClient struct for all Terraform provider methods
type OvcClient struct {
	JWT       string
	ServerURL string
	Access    string

	Machines MachineService
}

func (c *OvcClient) Do(req *http.Request) ([]byte, error) {
	var body []byte
	client := &http.Client{}
	req.Header.Set("Authorization", "bearer "+c.JWT)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	body, err = ioutil.ReadAll(resp.Body)
	log.Println("Status code: " + resp.Status)
	log.Println("Body: " + string(body))
	if resp.StatusCode > 202 {
		return body, errors.New(string(body))
	}
	if err != nil {
		return body, errors.New(string(body))
	}

	if err != nil {
		return body, errors.New(string(body))
	}
	return body, nil
}

func NewLogin(c *Config) string {
	authForm := url.Values{}
	authForm.Add("grant_type", "client_credentials")
	authForm.Add("client_id", c.ClientId)
	authForm.Add("client_secret", c.ClientSecret)
	authForm.Add("response_type", "id_token")
	req, _ := http.NewRequest("POST", "https://itsyou.online/v1/oauth/access_token", strings.NewReader(authForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing login request")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body")
	}
	jwt := string(bodyBytes)
	defer resp.Body.Close()
	return jwt
}
