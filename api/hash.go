package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// generateSalt generates a salt string based on the given IP and site ID.
//
// Parameters:
// - ip: the IP address used to generate the salt string.
// - siteID: the site ID used to generate the salt string.
//
// Returns:
// - string: the generated salt string.
func generateSalt(ip, siteID string) string {
	// HACK: Use the current date to change the salt every day
	// Replace this with a system that cannot be easily reverse engineered
	currentDate := time.Now().Format("20060102")
	return ip + "_" + siteID + "_" + currentDate
}

// generatePrivateIDHash generates a private ID hash based on the given inputs.
//
// Parameters:
// - ip: the IP address of the user
// - userAgent: the user agent string of the user's browser
// - hostname: the hostname of the server
// - siteID: the ID of the site
//
// Returns:
// - string: the generated private ID hash
// - error: an error if the hash generation fails
func generatePrivateIDHash(ip, userAgent, hostname, siteID string) (string, error) {
	salt := generateSalt(ip, siteID)
	data := salt + userAgent + hostname + siteID

	hasher := sha256.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}
