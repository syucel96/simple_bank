package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks the validity of a token and returns its payload if valid
	VerifyToken(token string) (*Payload, error)
}
