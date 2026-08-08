package models

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShapeHAlignment string

const (
	Left    ShapeHAlignment = "LEFT"
	right   ShapeHAlignment = "RIGHT"
	CenterH ShapeHAlignment = "CENTER"
)

type ShapeVAlignment string

const (
	Top     ShapeVAlignment = "TOP"
	Bottom  ShapeVAlignment = "BOTTOM"
	CenterV ShapeVAlignment = "CENTER"
)

type ShapeType string

const (
	PointShapeType     ShapeType = "point"
	RectangleShapeType ShapeType = "rectangle"
	PolyShapeType      ShapeType = "poly"
	ContainerShapeType ShapeType = "container"
	TextShapeType      ShapeType = "text"
	IconShapeType      ShapeType = "icon"
	DoorShapeType      ShapeType = "door"
)

type Shape interface {
	GetType() ShapeType
}

var shapeRegistry = map[string]func() Shape{
	"point": func() Shape { return &PointShape{} },

	"rectangle": func() Shape { return &RectangleShape{} },

	"poly": func() Shape { return &PolyShape{} },

	"container": func() Shape { return &ContainerShape{} },

	"text": func() Shape { return &TextShape{} },

	"icon": func() Shape { return &IconShape{} },

	"door": func() Shape { return &DoorShape{} },
}

type AnyShape struct {
	Shape
}

func (a *AnyShape) decode(typeName string, decode func(any) error) error {
	ctor, ok := shapeRegistry[typeName]
	if !ok {
		return fmt.Errorf("unknown shape type %q", typeName)
	}

	shape := ctor()

	if err := decode(shape); err != nil {
		return err
	}

	a.Shape = shape

	return nil
}

//
// BSON
//

func (a *AnyShape) UnmarshalBSON(data []byte) error {

	var header struct {
		Type string `bson:"type"`
	}

	if err := bson.Unmarshal(data, &header); err != nil {
		return err
	}

	return a.decode(header.Type, func(v any) error {
		return bson.Unmarshal(data, v)
	})
}

func (a AnyShape) MarshalBSON() ([]byte, error) {
	if a.Shape == nil {
		return nil, fmt.Errorf("shape is nil")
	}

	return bson.Marshal(a.Shape)
}

//
// JSON
//

func (a *AnyShape) UnmarshalJSON(data []byte) error {

	var header struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &header); err != nil {
		return err
	}

	return a.decode(header.Type, func(v any) error {
		return json.Unmarshal(data, v)
	})
}

func (a AnyShape) MarshalJSON() ([]byte, error) {
	if a.Shape == nil {
		return []byte("null"), nil
	}

	return json.Marshal(a.Shape)
}

type BaseShape struct {
	Type ShapeType `bson:"type"`
}

func (b BaseShape) GetType() ShapeType {
	return b.Type
}

type PointShape struct {
	BaseShape `bson:",inline"`
	X         float64 `bson:"x"`
	Y         float64 `bson:"y"`
}

func (*PointShape) GetType() ShapeType {
	return PointShapeType
}

type RectangleShape struct {
	BaseShape `bson:",inline"`
	X         float64 `bson:"x"`
	Y         float64 `bson:"y"`
	Width     float64 `bson:"width"`
	Height    float64 `bson:"height"`
}

func (*RectangleShape) GetType() ShapeType {
	return RectangleShapeType
}

type PolyShape struct {
	BaseShape `bson:",inline"`
	X         float64      `bson:"x"`
	Y         float64      `bson:"y"`
	Points    []PointShape `bson:"points"`
}

func (*PolyShape) GetType() ShapeType {
	return PolyShapeType
}

type ContainerShape struct {
	BaseShape `bson:",inline"`
	X         float64         `bson:"x"`
	Y         float64         `bson:"y"`
	Width     float64         `bson:"width"`
	Height    float64         `bson:"height"`
	AlignX    ShapeHAlignment `bson:"alignX"`
	AlignY    ShapeVAlignment `bson:"alignY"`
	Children  []AnyShape      `bson:"children"`
}

func (*ContainerShape) GetType() ShapeType {
	return ContainerShapeType
}

type TextShape struct {
	BaseShape `bson:",inline"`
	X         float64         `bson:"x"`
	Y         float64         `bson:"y"`
	AlignX    ShapeHAlignment `bson:"alignX"`
	AlignY    ShapeVAlignment `bson:"alignY"`
	Text      string          `bson:"text"`
}

func (*TextShape) GetType() ShapeType {
	return TextShapeType
}

type IconShape struct {
	BaseShape `bson:",inline"`
	X         float64 `bson:"x"`
	Y         float64 `bson:"y"`
	Width     float64 `bson:"width"`
	Height    float64 `bson:"height"`
	Icon      string  `bson:"icon"`
}

func (*IconShape) GetType() ShapeType {
	return IconShapeType
}

type DoorShape struct {
	BaseShape `bson:",inline"`
	wallId    int32   `bson:"wallId"`
	length    float64 `bson:"length"`
	offset    float64 `bson:"offset"`
}

func (*DoorShape) GetType() ShapeType {
	return DoorShapeType
}
