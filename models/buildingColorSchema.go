package models

type BuildingColorTheme struct {
	BuildingBorder     *string                   `bson:"buildingBorder"`
	BuildingFill       *string                   `bson:"buildingFill"`
	BuildingBackground *string                   `bson:"buildingBackground"`
	RoomBorder         *string                   `bson:"roomBorder"`
	RoomFill           *string                   `bson:"roomFill"`
	RoomText           *string                   `bson:"roomText"`
	RoomTypeBorder     map[GraphPointType]string `bson:"roomTypeBorder"`
	RoomTypeFill       map[GraphPointType]string `bson:"roomTypeFill"`
	RoomTypeText       map[GraphPointType]string `bson:"roomTypeText"`
}

type BuildingColorSchema struct {
	BaseDBSchema    `bson:",inline" json:",inline"`
	AccentColor     string             `bson:"accentColor"`
	WhiteColorTheme BuildingColorTheme `bson:"whiteColorTheme"`
	DarkColorTheme  BuildingColorTheme `bson:"darkColorTheme"`
}
