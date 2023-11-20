package models

import "time"

type Todo struct {
    ID     string  `json:"id"`
    Text  string  `json:"text"`
    Completed bool  `json:"completed" pg:",use_zero"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
