package models

type FloorResponse struct {
	Institute string       `json:"institute"`
	Floor     int          `json:"floor"`
	Width     float32      `json:"width"`
	Height    float32      `json:"height"`
	Audiences []Auditorium `json:"audiences"`
	Service   []Service    `json:"service"`
}

type InstituteResponse struct {
	Name            string        `json:"name"`
	DisplayableName string        `json:"displayableName"`
	MinFloor        int           `json:"minFloor"`
	MaxFloor        int           `json:"maxFloor"`
	Url             string        `json:"url"`
	Latitude        float64       `json:"latitude"`
	Longitude       float64       `json:"longitude"`
	Icon            InstituteIcon `json:"icon"`
}
