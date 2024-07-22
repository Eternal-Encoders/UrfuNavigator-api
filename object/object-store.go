package object

type ObjectStore interface {
	GetFile(fileName string) ([]byte, error)
}
