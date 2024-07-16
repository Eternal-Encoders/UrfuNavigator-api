package models

type InstituteIcon struct {
	Url string `bson:"filename" json:"url"`
	Alt string `json:"alt"`
}

type Institute struct {
	Name            string  `json:"name"`
	DisplayableName string  `json:"displayableName"`
	MinFloor        int     `json:"minFloor"`
	MaxFloor        int     `json:"maxFloor"`
	Url             string  `json:"url"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Icon            string  `json:"icon"`
}
