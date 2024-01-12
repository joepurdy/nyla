package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mileusna/useragent"
)

type Event struct {
	SiteID      string
	CreatedAt   int32
	Type        string
	Event       string
	Referrer    string
	IsTouch     bool
	BrowserName string
	OSName      string
	DeviceType  string
	Country     string
	Region      string
	Timestamp   time.Time
}

type Events struct {
	DB *pgx.Conn
}

func (e *Events) Open() error {
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://postgres:password@localhost:5432/postgres",
	)
	if err != nil {
		return err
	} else if err := conn.Ping(context.Background()); err != nil {
		return err
	}

	e.DB = conn
	return nil
}

func (e *Events) Add(payload CollectorPayload, ua useragent.UserAgent, geo *GeoInfo) error {
	q := `
	INSERT INTO events
	(
		site_id, 
		created_at, 
		type, 
		event, 
		referrer,
		is_touch, 
		browser_name, 
		os_name,
		device_type, 
		country, 
		region
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	)
	`

	_, err := e.DB.Exec(
		context.Background(),
		q,
		payload.SiteID,
		nowToInt(),
		payload.Data.Type,
		payload.Data.Event,
		payload.Data.Referrer,
		"false",
		ua.Name,
		ua.OS,
		"not-implemented",
		geo.Country,
		geo.RegionName,
	)

	return err
}

func nowToInt() uint32 {
	now := time.Now().Format("20060102")
	i, err := strconv.ParseInt(now, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return uint32(i)
}
