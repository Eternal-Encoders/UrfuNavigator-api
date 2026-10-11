package models

type BuildingColorTheme struct {
	BuildingBorder     *string                   `bson:"buildingBorder" json:"buildingBorder,omitempty"`
	BuildingFill       *string                   `bson:"buildingFill" json:"buildingFill,omitempty"`
	BuildingBackground *string                   `bson:"buildingBackground" json:"buildingBackground,omitempty"`
	RoomBorder         *string                   `bson:"roomBorder" json:"roomBorder,omitempty"`
	RoomFill           *string                   `bson:"roomFill" json:"roomFill,omitempty"`
	RoomText           *string                   `bson:"roomText" json:"roomText,omitempty"`
	RoomTypeBorder     map[GraphPointType]string `bson:"roomTypeBorder" json:"roomTypeBorder,omitempty"`
	RoomTypeFill       map[GraphPointType]string `bson:"roomTypeFill" json:"roomTypeFill,omitempty"`
	RoomTypeText       map[GraphPointType]string `bson:"roomTypeText" json:"roomTypeText,omitempty"`
}

type BuildingColorSchema struct {
	BaseDBSchema    `bson:",inline" json:",inline"`
	AccentColor     string             `bson:"accentColor" json:"accentColor"`
	WhiteColorTheme BuildingColorTheme `bson:"whiteColorTheme" json:"whiteColorTheme"`
	DarkColorTheme  BuildingColorTheme `bson:"darkColorTheme" json:"darkColorTheme"`
}
