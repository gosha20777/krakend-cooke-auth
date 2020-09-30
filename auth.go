package auth

import (
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
)

// Validator defines the interface for all the possible validation processes
type Validator interface {
	IsValid(subject string) (*authInfo, error)
}

// NewCredentialsValidator creates a validator for a given credentials pair
func NewCredentialsValidator(credentials Credentials) Validator {
	url := credentials.Url
	return authHeader{url}
}

type authHeader struct {
	Url string
}

type authInfo struct {
	SessionId string `json:"session_id"`
	UserId int `json:"user_id"`
}

// IsValid implements the Validator interface
func (a authHeader) IsValid(value string) (*authInfo, error) {
	url := a.Url + "?cookie=" + value

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("can't make get request to ", url)
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Println("HTTP Status is not in the 2xx range ", resp.StatusCode)
		return nil, errors.New("HTTP Status is not in the 2xx range")
	}
	
	defer resp.Body.Close()
	info := new(authInfo)
	err = json.NewDecoder(resp.Body).Decode(info)
    if err != nil {
		fmt.Println("can't read body from", url)
		return nil, err
	}
	return info, nil
}
