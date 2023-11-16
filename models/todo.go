package models

import "time"

type Todo struct {
    ID     string  `json:"id"`
    Text  string  `json:"text"`
    Completed bool  `json:"completed"`
	// Date  float64 `json:"date"` // TODO: implement latter
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

