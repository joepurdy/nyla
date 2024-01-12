package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	GEOIP_PROTO = getEnv("GEOIP_PROTO", "http")
	GEOIP_HOST  = getEnv("GEOIP_HOST", "localhost:8080")
)

// getEnv retrieves the value of the environment variable named by the key
// or returns the default value if the environment variable is not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type GeoInfo struct {
	IP         string `json:"ip"`
	Country    string `json:"country"`
	CountryISO string `json:"country_iso"`
	RegionName string `json:"region_name"`
	RegionCode string `json:"region_code"`
	City       string `json:"city"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

// ipFromRequest retrieves the client's IP address from the HTTP headers in a given request.
//
// The function takes two parameters: headers, a slice of strings representing the desired headers to check,
// and r, a pointer to an http.Request object containing the HTTP request information.
//
// The function returns a net.IP object, representing the client's IP address, and an error, if any.
func ipFromRequest(headers []string, r *http.Request) (net.IP, error) {
	remoteIP := ""
	for _, h := range headers {
		remoteIP = r.Header.Get(h)
		if http.CanonicalHeaderKey(h) == "X-Forwarded-For" {
			remoteIP = ipFromForwardedForHeader(remoteIP)
		}
		if remoteIP != "" {
			break
		}
	}

	if remoteIP == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, err
		}
		remoteIP = host
	}

	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP %s", remoteIP)
	}
	return ip, nil
}

// ipFromForwardedForHeader returns the first IP address from the "Forwarded-For" header value.
//
// It takes a string parameter `v` which represents the "Forwarded-For" header value.
// The function returns a string which is the first IP address from the header value.
func ipFromForwardedForHeader(v string) string {
	sep := strings.Index(v, ",")
	if sep == -1 {
		return v
	}
	return v[:sep]
}

// getGeoInfo retrieves geographical information for a given IP address.
//
// It takes a string parameter `ip` which represents the IP address for which
// the geographical information needs to be fetched.
//
// It returns a pointer to a `GeoInfo` struct and an error. The `GeoInfo`
// struct contains information like the country, region, city, latitude and
// longitude for the given IP address.
func getGeoInfo(ip string) (*GeoInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s/json?ip=%s", GEOIP_PROTO, GEOIP_HOST, ip), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info GeoInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	return &info, err
}
