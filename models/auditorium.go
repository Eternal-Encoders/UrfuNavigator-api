package models

type AuditorumDoor struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
	Fill   string  `json:"fill"`
}

type AuditoriumChild struct {
	Type       string  `json:"type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Identifier string  `json:"identifier"`
	AlignX     string  `json:"alignX"`
	AlignY     string  `json:"alignY"`
}

type Auditorium struct {
	Id       string            `json:"id"`
	X        float64           `json:"x"`
	Y        float64           `json:"y"`
	Width    float32           `json:"width"`
	Height   float32           `json:"height"`
	Fill     *string           `json:"fill"`
	Stroke   *string           `json:"stroke"`
	PointId  *string           `json:"pointId"`
	Children []AuditoriumChild `json:"children"`
	Doors    []AuditorumDoor   `json:"doors"`
}
