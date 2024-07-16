package models

type Floor struct {
	Institute string       `json:"institute"`
	Floor     int          `json:"floor"`
	Width     float32      `json:"width"`
	Height    float32      `json:"height"`
	Audiences []Auditorium `json:"audiences"`
	Service   []Service    `json:"service"`
	Graph     []string     `json:"graph"`
}
