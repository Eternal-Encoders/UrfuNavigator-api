package models

type Service struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Data   string  `json:"data"`
	Stroke *string `json:"stroke"`
	Fill   *string `json:"fill"`
}
