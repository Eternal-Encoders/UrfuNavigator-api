package models

// Floors
type FloorQuery struct {
	Id string `query:"id,required"`
}

// Buildings
type BuildingQuery struct {
	Id string `query:"id,required"`
}

// Graph Points
type PointsQuery struct {
	BuildingId *string `query:"buildingId"`
	FloorId    *string `query:"floorId"`
	Type       *string `query:"type"`
	Name       *string `query:"name"`
	Length     *int    `query:"length"`
}

type PointIdQuery struct {
	Id string `query:"id,required"`
}

type PathQuery struct {
	From string `query:"from,required"`
	To   string `query:"to,required"`
}

// Media
type MediaDeleteBody struct {
	Filename string `json:"filename"`
}

type MediaUploadBody struct {
	Filename string `json:"filename"`
	Image    []byte `json:"image"`
}

// Points Find
type SearchQuery struct {
	Name   string `query:"name,required"`
	Length *int   `query:"length"`
}

// Auth
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Admin Users
type CreateUserRequest struct {
	DisplayableName string   `json:"displayableName"`
	Login           string   `json:"login"`
	Password        string   `json:"password"`
	Role            UserRole `json:"role"`
}

type UpdateUserRequest struct {
	DisplayableName *string   `json:"displayableName,omitempty"`
	Login           *string   `json:"login,omitempty"`
	Password        *string   `json:"password,omitempty"`
	Role            *UserRole `json:"role,omitempty"`
}
