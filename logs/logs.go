package logs

type Log interface {
	WriteLog(ip string, request string, args map[string]string) error
	InitDb() error
}
