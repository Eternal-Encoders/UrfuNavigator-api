package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// EntityResponse is public-facing entity metadata (audit fields omitted).
type EntityResponse struct {
	Id              bson.ObjectID `json:"id"`
	DisplayableName string        `json:"displayableName"`
}

// AdminEntityResponse includes audit metadata for admin API responses.
type AdminEntityResponse struct {
	Id              bson.ObjectID `json:"id"`
	DisplayableName string        `json:"displayableName"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
	Author          bson.ObjectID `json:"author"`
	LastUpdatedBy   bson.ObjectID `json:"lastUpdatedBy"`
}

type BuildingGpsResponse struct {
	CentreAltitude float64       `json:"centreAltitude"`
	FloorId        bson.ObjectID `json:"floorId"`
}

type BuildingIconResponse struct {
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// BuildingResponse describes a campus building and its linked resources.
type BuildingResponse struct {
	EntityResponse
	Floors       []EntityResponse      `json:"floors"`
	Url          string                `json:"url"`
	Latitude     float64               `json:"latitude"`
	Longitude    float64               `json:"longitude"`
	Icon         BuildingIconResponse  `json:"icon"`
	ColorSchemes []bson.ObjectID       `json:"colorSchemes"`
	Gps          []BuildingGpsResponse `json:"gps"`
}

type LinearGpsResponse struct {
	B1 float64 `json:"b1"`
	B2 float64 `json:"b2"`
	A  float64 `json:"a"`
}

type PointGpsResponse struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type LinearMappingResponse struct {
	X LinearGpsResponse `json:"x"`
	Y LinearGpsResponse `json:"y"`
}

type StabilForceGpsResponse struct {
	Point PointGpsResponse `json:"point"`
	Force PointGpsResponse `json:"force"`
}

type FloorGpsResponse struct {
	Altitude float64                  `json:"altitude"`
	Linear   LinearMappingResponse    `json:"linear"`
	Forces   []StabilForceGpsResponse `json:"forces"`
}

type RoomResponse struct {
	EntityResponse
	Shape       AnyShape        `json:"shape"`
	PointId     *bson.ObjectID  `json:"pointId,omitempty"`
	Type        *GraphPointType `json:"type,omitempty"`
	Children    []AnyShape      `json:"children"`
	ColorSchema *bson.ObjectID  `json:"colorSchema,omitempty"`
	IsBorder    bool            `json:"isBorder"`
	IsFill      bool            `json:"isFill"`
}

type ServiceResponse struct {
	EntityResponse
	Shape       AnyShape       `json:"shape"`
	ColorSchema *bson.ObjectID `json:"colorSchema,omitempty"`
	IsBorder    bool           `json:"isBorder"`
	IsFill      bool           `json:"isFill"`
}

type PointNameResponse struct {
	Name         string                `json:"name"`
	Translations []TransaltionResponse `json:"translations"`
}

type TransaltionResponse struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

type WeekTimeResponse struct {
	Start    int32 `json:"start"`
	End      int32 `json:"end"`
	IsDayOff bool  `json:"isDayOff"`
}

type WeekResponse struct {
	Monday    WeekTimeResponse `json:"monday"`
	Tuesday   WeekTimeResponse `json:"tuesday"`
	Wednesday WeekTimeResponse `json:"wednesday"`
	Thursday  WeekTimeResponse `json:"thursday"`
	Friday    WeekTimeResponse `json:"friday"`
	Saturday  WeekTimeResponse `json:"saturday"`
	Sunday    WeekTimeResponse `json:"sunday"`
}

// GraphPointResponse is a navigation graph node linked to a building and floor.
type GraphPointResponse struct {
	EntityResponse
	BuildingId  bson.ObjectID       `json:"buildingId"`
	FloorId     bson.ObjectID       `json:"floorId"`
	X           float64             `json:"x"`
	Y           float64             `json:"y"`
	Links       []bson.ObjectID     `json:"links"`
	Types       []GraphPointType    `json:"types"`
	Names       []PointNameResponse `json:"names"`
	Time        WeekResponse        `json:"time,omitempty"`
	Description string              `json:"description,omitempty"`
	Info        string              `json:"info,omitempty"`
	IsPassFree  bool                `json:"isPassFree"`
}

// FloorDetailResponse is a floor with resolved rooms, services, and graph points.
type FloorDetailResponse struct {
	EntityResponse
	BuildingId bson.ObjectID        `json:"buildingId"`
	Elevation  float64              `json:"elevation"`
	Width      int32                `json:"width"`
	Height     int32                `json:"height"`
	Rooms      []RoomResponse       `json:"rooms"`
	Services   []ServiceResponse    `json:"services"`
	Graph      []GraphPointResponse `json:"graph"`
	Gps        *FloorGpsResponse    `json:"gps,omitempty"`
}

// PathResultResponse groups route segments by building and floor.
type PathResultResponse map[bson.ObjectID]map[bson.ObjectID][][]GraphPointResponse

type PathJSONResponse struct {
	Result PathResultResponse `json:"result"`
}

func ToEntityResponse(entity BaseDBSchema) EntityResponse {
	return EntityResponse{
		Id:              entity.Id,
		DisplayableName: entity.DisplayableName,
	}
}

func ToAdminEntityResponse(entity BaseDBSchema) AdminEntityResponse {
	return AdminEntityResponse{
		Id:              entity.Id,
		DisplayableName: entity.DisplayableName,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdateAt,
		Author:          entity.Author,
		LastUpdatedBy:   entity.LastUpdatedBy,
	}
}

func ToBuildingGpsResponse(gps BuildingGps) BuildingGpsResponse {
	return BuildingGpsResponse{
		CentreAltitude: gps.CentreAltitude,
		FloorId:        gps.FloorId,
	}
}

func ToFloorEntityResponses(floorIDs []bson.ObjectID, floors map[bson.ObjectID]*Floor) []EntityResponse {
	result := make([]EntityResponse, 0, len(floorIDs))
	for _, id := range floorIDs {
		if floor, ok := floors[id]; ok && floor != nil {
			result = append(result, ToEntityResponse(floor.BaseDBSchema))
			continue
		}
		result = append(result, EntityResponse{Id: id})
	}
	return result
}

func ToBuildingResponse(building Building, icon BuildingIconResponse, floors map[bson.ObjectID]*Floor) BuildingResponse {
	gps := make([]BuildingGpsResponse, len(building.Gps))
	for i, item := range building.Gps {
		gps[i] = ToBuildingGpsResponse(item)
	}

	return BuildingResponse{
		EntityResponse: ToEntityResponse(building.BaseDBSchema),
		Floors:         ToFloorEntityResponses(building.Floors, floors),
		Url:            building.Url,
		Latitude:       building.Latitude,
		Longitude:      building.Longitude,
		Icon:           icon,
		ColorSchemes:   building.ColorSchemes,
		Gps:            gps,
	}
}

func ToFloorGpsResponse(gps *FloorGps) *FloorGpsResponse {
	if gps == nil {
		return nil
	}

	forces := make([]StabilForceGpsResponse, len(gps.Forces))
	for i, force := range gps.Forces {
		forces[i] = StabilForceGpsResponse{
			Point: PointGpsResponse{X: force.Point.X, Y: force.Point.Y},
			Force: PointGpsResponse{X: force.Force.X, Y: force.Force.Y},
		}
	}

	return &FloorGpsResponse{
		Altitude: gps.Altitude,
		Linear: LinearMappingResponse{
			X: LinearGpsResponse{B1: gps.Linear.X.B1, B2: gps.Linear.X.B2, A: gps.Linear.X.A},
			Y: LinearGpsResponse{B1: gps.Linear.Y.B1, B2: gps.Linear.Y.B2, A: gps.Linear.Y.A},
		},
		Forces: forces,
	}
}

func ToRoomResponse(room Room) RoomResponse {
	return RoomResponse{
		EntityResponse: ToEntityResponse(room.BaseDBSchema),
		Shape:          room.Shape,
		PointId:        room.PointId,
		Type:           room.Type,
		Children:       room.Children,
		ColorSchema:    room.ColorSchema,
		IsBorder:       room.IsBorder,
		IsFill:         room.IsFill,
	}
}

func ToServiceResponse(service Service) ServiceResponse {
	return ServiceResponse{
		EntityResponse: ToEntityResponse(service.BaseDBSchema),
		Shape:          service.Shape,
		ColorSchema:    service.ColorSchema,
		IsBorder:       service.IsBorder,
		IsFill:         service.IsFill,
	}
}

func ToNameTranslationResponse(translation Transaltion) TransaltionResponse {
	return TransaltionResponse{
		Language: translation.Language,
		Value:    translation.Value,
	}
}

func ToPointNameResponse(name PointName) PointNameResponse {
	translations := make([]TransaltionResponse, len(name.Translations))
	for i, translation := range name.Translations {
		translations[i] = ToNameTranslationResponse(translation)
	}
	return PointNameResponse{
		Name:         name.Name,
		Translations: translations,
	}
}

func ToWeekTimeResponse(time WeekTime) WeekTimeResponse {
	return WeekTimeResponse{
		Start:    time.Start,
		End:      time.End,
		IsDayOff: time.IsDayOff,
	}
}

func ToWeekResponse(week Week) WeekResponse {
	return WeekResponse{
		Monday:    ToWeekTimeResponse(week.Monday),
		Tuesday:   ToWeekTimeResponse(week.Tuesday),
		Wednesday: ToWeekTimeResponse(week.Wednesday),
		Thursday:  ToWeekTimeResponse(week.Thursday),
		Friday:    ToWeekTimeResponse(week.Friday),
		Saturday:  ToWeekTimeResponse(week.Saturday),
		Sunday:    ToWeekTimeResponse(week.Sunday),
	}
}

func ToGraphPointResponse(point GraphPoint) GraphPointResponse {
	names := make([]PointNameResponse, len(point.Names))
	for i, name := range point.Names {
		names[i] = ToPointNameResponse(name)
	}

	return GraphPointResponse{
		EntityResponse: ToEntityResponse(point.BaseDBSchema),
		BuildingId:     point.BuildingId,
		FloorId:        point.FloorId,
		X:              point.X,
		Y:              point.Y,
		Links:          point.Links,
		Types:          point.Types,
		Names:          names,
		Time:           ToWeekResponse(point.Time),
		Description:    point.Description,
		Info:           point.Info,
		IsPassFree:     point.IsPassFree,
	}
}

func ToGraphPointResponses(points []GraphPoint) []GraphPointResponse {
	result := make([]GraphPointResponse, len(points))
	for i, point := range points {
		result[i] = ToGraphPointResponse(point)
	}
	return result
}

func ToRoomResponses(rooms []Room) []RoomResponse {
	result := make([]RoomResponse, len(rooms))
	for i, room := range rooms {
		result[i] = ToRoomResponse(room)
	}
	return result
}

func ToServiceResponses(services []Service) []ServiceResponse {
	result := make([]ServiceResponse, len(services))
	for i, service := range services {
		result[i] = ToServiceResponse(service)
	}
	return result
}

func ToFloorDetailResponse(floor Floor, rooms []Room, services []Service, graph []GraphPoint) FloorDetailResponse {
	return FloorDetailResponse{
		EntityResponse: ToEntityResponse(floor.BaseDBSchema),
		BuildingId:     floor.BuildingId,
		Elevation:      floor.Elevation,
		Width:          floor.Width,
		Height:         floor.Height,
		Rooms:          ToRoomResponses(rooms),
		Services:       ToServiceResponses(services),
		Graph:          ToGraphPointResponses(graph),
		Gps:            ToFloorGpsResponse(floor.Gps),
	}
}

func ToPathResultResponse(path map[bson.ObjectID]map[bson.ObjectID][][]GraphPoint) PathResultResponse {
	result := make(PathResultResponse, len(path))

	for buildingID, floors := range path {
		result[buildingID] = make(map[bson.ObjectID][][]GraphPointResponse, len(floors))
		for floorID, segments := range floors {
			convertedSegments := make([][]GraphPointResponse, len(segments))
			for i, segment := range segments {
				convertedSegments[i] = ToGraphPointResponses(segment)
			}
			result[buildingID][floorID] = convertedSegments
		}
	}

	return result
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Login     string    `json:"login"`
	Role      UserRole  `json:"role"`
}

type UserAdminResponse struct {
	AdminEntityResponse
	Login string   `json:"login"`
	Role  UserRole `json:"role"`
}

type CreatedResponse struct {
	Id bson.ObjectID `json:"id"`
}
