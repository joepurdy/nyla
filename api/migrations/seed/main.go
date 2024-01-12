package main

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	TOTALROWS  = 1500000
	TOTALUSERS = 125321
	TOTALPAGES = 500
)

var (
	SITEID    = []string{"site1", "site2", "site3"}
	MONTHS    = []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	TYPES     = []string{"page", "event"}
	PAGES     = make([]string, TOTALPAGES)
	REFERRERS = []string{"", "google", "twitter", "reddit", "siteabc.com"}
	DEVICES   = []string{"desktop", "tablet", "phone"}
	BROWSERS  = []string{"chrome", "firefox", "edge"}
	OSNAME    = []string{"linux", "windows", "macos"}
	COUNTRIES = []string{"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8"}
)

const (
	collisionProbability  = 0.1 // 10% chance of collision
	numPreGeneratedHashes = 100 // Number of unique pre-generated hashes
)

var preGeneratedHashes []string

func init() {
	// Generate a set of pre-defined hashes
	preGeneratedHashes = make([]string, numPreGeneratedHashes)
	for i := range preGeneratedHashes {
		preGeneratedHashes[i] = generateNewMockHash()
	}
}

func generateNewMockHash() string {
	bytes := make([]byte, 32) // SHA256 hash is 32 bytes
	_, err := cryptorand.Read(bytes)
	if err != nil {
		// Handle error here
		fmt.Println(err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

func generateMockHash() string {
	// Randomly decide whether to generate a new hash or use a pre-generated one
	if rand.Float64() < collisionProbability {
		// Return a pre-generated hash to simulate a collision
		return preGeneratedHashes[rand.Intn(len(preGeneratedHashes))]
	}
	// Generate a new unique hash
	return generateNewMockHash()
}

func main() {
	genPages()

	const numGoroutines = 10 // Number of goroutines
	chunkSize := TOTALROWS / numGoroutines

	var wg sync.WaitGroup
	tempFiles := make([]string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tempFile, err := os.CreateTemp("", "data_chunk_*.tmp")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer tempFile.Close()

			tempFiles[i] = tempFile.Name()

			for j := i * chunkSize; j < (i+1)*chunkSize; j++ {
				row := genInsert()
				if _, err := tempFile.WriteString(row); err != nil {
					fmt.Println(err)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// create data directory if it doesn't exist
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err := os.Mkdir("data", 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Concatenate temp files
	finalFile, err := os.OpenFile("data/dump.data", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer finalFile.Close()

	for _, tempFileName := range tempFiles {
		tempFile, err := os.Open(tempFileName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if _, err := io.Copy(finalFile, tempFile); err != nil {
			fmt.Println(err)
		}

		tempFile.Close()
		os.Remove(tempFileName) // Clean up temp file
	}
}

func genInsert() string {
	device := DEVICES[rand.Intn(len(DEVICES))]
	is_touch := "true"
	if device == "desktop" {
		is_touch = "false"
	}

	anonID := generateMockHash()

	qry := "%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n"

	return fmt.Sprintf(qry,
		anonID,
		SITEID[rand.Intn(len(SITEID))],
		genCreatedAt(),
		TYPES[rand.Intn(len(TYPES))],
		PAGES[rand.Intn(len(PAGES))],
		REFERRERS[rand.Intn(len(REFERRERS))],
		is_touch,
		BROWSERS[rand.Intn(len(BROWSERS))],
		OSNAME[rand.Intn(len(OSNAME))],
		device,
		COUNTRIES[rand.Intn(len(COUNTRIES))],
		"no need",
		time.Now().Format(time.RFC3339),
	)
}

func genCreatedAt() uint32 {
	year := time.Now().Year() - rand.Intn(3)
	month := MONTHS[rand.Intn(len(MONTHS))]
	day := rand.Intn(27) + 1 // let's play safe...

	d := fmt.Sprintf("%d%s%d", year, month, day)
	i, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		return 20231205
	}
	return uint32(i)
}

func genPages() {
	for i := 0; i < TOTALPAGES; i++ {
		sep := rand.Intn(3)
		if sep == 0 {
			PAGES[i] = "/"
		} else {
			p := "/"
			for j := 0; j < sep; j++ {
				p += pageName() + "/"
			}

			PAGES[i] = p
		}
	}
}

func pageName() string {
	return fmt.Sprint(
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
	)
}
