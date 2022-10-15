package service

import (
	guuid "github.com/google/uuid"
)

// NewUUID -
func NewUUID() (string, error) {

	// Use generate uuid of google
	return guuid.New().String(), nil
}
