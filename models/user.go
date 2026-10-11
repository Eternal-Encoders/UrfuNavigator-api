package models

type UserRole string

const (
	ReaderRole UserRole = "reader"
	WriterRole UserRole = "writer"
	AdmineRole UserRole = "admine"
)

type User struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	Login        string   `bson:"login"`
	PasswordHash string   `bson:"passwordHash"`
	Role         UserRole `bson:"role"`
}
