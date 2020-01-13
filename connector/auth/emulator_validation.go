package auth

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// IsTokenFromEmulator checks if the request is from emulator by probing the bearer token in request header
func IsTokenFromEmulator(authHeader string) bool {
	if authHeader == "" {
		return false
	}
	// The Auth Header format
	// "Bearer eyJ0e[...Big Long String...]XAiO"
	parts := strings.Split(authHeader, " ")
	// parts[0] = Bearer
	// parts[1] = eyJ0e[...Big Long String...]XAiO
	if len(parts) != 2 {
		return false
	}

	authScheme := parts[0]
	token := parts[1]

	if authScheme != "Bearer" {
		return false
	}

	// Parse token
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, nil)
	// Check if token is getting parsed correctly
	fmt.Printf("token %#v", parsedToken, err)
	if err != nil {
		return false
	}

	if _, ok := claims["iss"]; !ok {
		return false
	}
	return true
}

// EmulatorTokenValidator provides functionality to check if a request is from an emulator
type EmulatorTokenValidator struct {
}
