package models

import "time"

type AccessToken struct {
	ID string `json:"id"`
	// Email     string    `json:"email"`
	Token     string    `json:"token"`
	Roles     []string  `json:"roles"` // read, read-one, write, write-new-only, write-update-only
	Expiry    time.Time `json:"expiry"`
	UserID    string    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type AccessLog struct {
	ID             string    `json:"id"`
	Token          string    `json:"token"` // TODO: Change this to TokenID
	Path           string    `json:"path"`
	ClientIP       string    `json:"client_ip"`
	Method         string    `json:"method"`
	ResponseTime   int64     `json:"response_time"`
	ResponseSize   int       `json:"response_size"`
	StatusCode     int       `json:"status_code"`
	ServerHostname string    `json:"server_hostname"`
	BillID         string    `json:"bill_id"`
	Billed         bool      `json:"billed" pg:",use_zero"`
	CreatedAt      time.Time `json:"created_at"`
	// UpdatedAt      time.Time `json:"updated_at"` // TODO: Latter
}
