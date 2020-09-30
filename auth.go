package auth

import (
	"fmt"
	"net/http"
	"encoding/json"
)

// Validator defines the interface for all the possible validation processes
type Validator interface {
	IsValid(subject string) bool
}

// NewCredentialsValidator creates a validator for a given credentials pair
func NewCredentialsValidator(credentials Credentials) Validator {
	url := credentials.Url
	return authHeader{url}
}

type authHeader struct {
	url string
}

type authInfo struct {
	session_id string
	user_id int
}

// IsValid implements the Validator interface
func (a authHeader) IsValid(value string) bool {
	url := a.url + "?cookie=" + value
	fmt.Println("Make a request to:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("can't make get request to ", url)
		return false
	}
	if resp.StatusCode <= 200 || resp.StatusCode >= 299 {
		fmt.Println("HTTP Status is not in the 2xx range ", resp.StatusCode)
		return false
	}
	
	defer resp.Body.Close()
	info := new(authInfo)
	err = json.NewDecoder(resp.Body).Decode(info)
    if err != nil {
		fmt.Println("can't read body from", url)
		return false
	}
	fmt.Println("session_id:", info.session_id)
	fmt.Println("user_id:", info.user_id)
	return true
}
