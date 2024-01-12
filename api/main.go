// nyla-api - GDPR compliant privacy focused web analytics
// Copyright (C) 2024 Joe Purdy
// mailto:nyla AT purdy DOT dev
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mileusna/useragent"
)

// Version is provided at compile time
var Version = "devel"

var (
	events *Events = &Events{}
)

type CollectorData struct {
	Type      string `json:"type"`
	Event     string `json:"event"`
	UserAgent string `json:"ua"`
	Hostname  string `json:"hostname"`
	Referrer  string `json:"referrer"`
}

type CollectorPayload struct {
	SiteID string        `json:"site_id"`
	Data   CollectorData `json:"data"`
}

// decodePayload decodes a base64 encoded string and unmarshals it into a CollectorPayload struct.
//
// It takes a single parameter:
//   - s: a string representing the base64 encoded data.
//
// It returns two values:
//   - payload: a CollectorPayload struct representing the decoded and unmarshaled data.
//   - err: an error if there was any issue during the decoding or unmarshaling process.
func decodePayload(s string) (payload CollectorPayload, err error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &payload)
	return
}

func main() {
	fmt.Println("nyla version:", Version)

	if err := events.Open(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/collect", collect)

	fmt.Println("listening on :9876")
	http.ListenAndServe(":9876", nil)
}

func collect(w http.ResponseWriter, r *http.Request) {
	defer w.WriteHeader(http.StatusOK)

	data := r.URL.Query().Get("data")
	payload, err := decodePayload(data)
	if err != nil {
		fmt.Print(err)
	}

	ua := useragent.Parse(payload.Data.UserAgent)

	for k, v := range r.Header {
		fmt.Println(k, v)
	}

	ip, err := ipFromRequest([]string{"X-Forwarded-For", "X-Real-IP"}, r)
	if err != nil {
		fmt.Println("error getting IP:", err)
	}

	geoInfo, err := getGeoInfo(ip.String())
	if err != nil {
		fmt.Println("error getting geo info:", err)
	}

	hash, err := generatePrivateIDHash(ip.String(), trk.Data.UserAgent, trk.Data.Hostname, trk.SiteID)
	if err != nil {
		fmt.Println("error generating private ID hash:", err)
	}

	if err := events.Add(payload, hash, ua, geoInfo); err != nil {
		fmt.Println(err)
	}
}
