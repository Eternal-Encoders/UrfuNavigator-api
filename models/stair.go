package models

type Stair struct {
	Id    string   `bson:"_id" json:"id"`
	Links []string `json:"links"`
}

type StairParsed struct {
	Id    string       `json:"id"`
	Links []GraphPoint `json:"links"`
}
