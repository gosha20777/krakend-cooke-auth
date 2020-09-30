package auth

import (
	"fmt"
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

// IsValid implements the Validator interface
func (a authHeader) IsValid(value string) bool {
	url := a.url + "?value=" + value
	fmt.Println("Make a request to:", url)
	return true
}
