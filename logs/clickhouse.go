package logs

import (
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseData struct {
	Ip        net.IP
	Request   string
	Args      []map[string]string
	Timestamp time.Time
}

type Clickhouse struct {
	batchSize int
	// batchData driver.Batch
	conn clickhouse.Conn
	db   string
}

func NewClickhouse(
	uri string,
	db string,
	user string,
	password string,
	batchSize int,
) *Clickhouse {
	// ctx := context.Background()
	// dialCount := 0

	// conn, err := clickhouse.Open(&clickhouse.Options{
	// 	Addr: []string{uri},
	// 	Auth: clickhouse.Auth{
	// 		Username: user,
	// 		Password: password,
	// 	},
	// 	DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
	// 		dialCount++
	// 		var d net.Dialer
	// 		return d.DialContext(ctx, "tcp", addr)
	// 	},
	// 	Settings: clickhouse.Settings{
	// 		"max_execution_time": 60,
	// 	},
	// 	Compression: &clickhouse.Compression{
	// 		Method: clickhouse.CompressionLZ4,
	// 	},
	// 	DialTimeout:      time.Duration(10) * time.Second,
	// 	MaxOpenConns:     5,
	// 	MaxIdleConns:     5,
	// 	ConnMaxLifetime:  time.Duration(10) * time.Minute,
	// 	ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	// 	BlockBufferSize:  10,
	// 	TLS: &tls.Config{
	// 		InsecureSkipVerify: true,
	// 	},
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := conn.Ping(ctx); err != nil {
	// 	if exception, ok := err.(*clickhouse.Exception); ok {
	// 		log.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
	// 	}
	// 	log.Fatal(err)
	// }

	// return &Clickhouse{
	// 	conn:      conn,
	// 	db:        db,
	// 	batchSize: batchSize,
	// }
	return &Clickhouse{
		conn:      nil,
		db:        db,
		batchSize: batchSize,
	}
}

func (c *Clickhouse) InitDb() error {
	// ctx := context.Background()

	// requests := []string{
	// 	"MainHandler",
	// 	"FloorHandler",
	// 	"InstituteHandler",
	// 	"InstitutesHandler",
	// 	"PointsHandler",
	// 	"PointIdHandler",
	// 	"PathHandler",
	// 	"ObjectHandler",
	// 	"SearchHandler",
	// }
	// apiCallsTypes := map[string]string{
	// 	"Ip":        "IPv4",
	// 	"Request":   utils.CreateEnum(requests),
	// 	"Args":      "Array(Tuple(Name String, Value String))",
	// 	"Timestamp": "DateTime",
	// }

	// if err := c.conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", c.db)); err != nil {
	// 	return err
	// }
	// if err := c.conn.Exec(ctx, utils.CreateTable(fmt.Sprintf("%s.api_calls", c.db), apiCallsTypes, "Timestamp")); err != nil {
	// 	return err
	// }

	// c.CreateBatchInsert()

	return nil
}

func (c *Clickhouse) CreateBatchInsert() {
	// ctx := context.Background()
	// batch, err := c.conn.PrepareBatch(
	// 	ctx,
	// 	fmt.Sprintf("INSERT INTO %s.%s", c.db, "api_calls"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.batchData = batch
}

func (c *Clickhouse) WriteLog(ip string, request string, args map[string]string) error {
	// tempArgs := []map[string]string{}
	// for k, v := range args {
	// 	tempArgs = append(tempArgs, map[string]string{
	// 		"Name":  strings.Clone(k),
	// 		"Value": strings.Clone(v),
	// 	})
	// }
	// data := ClickHouseData{
	// 	Ip:        net.ParseIP(ip),
	// 	Request:   strings.Clone(request),
	// 	Args:      tempArgs,
	// 	Timestamp: time.Now(),
	// }

	// if batchErr := c.batchData.AppendStruct(&data); batchErr != nil {
	// 	log.Println(batchErr)
	// 	return batchErr
	// }
	// if c.batchData.Rows() < c.batchSize {
	// 	return nil
	// }

	// if sendErr := c.batchData.Send(); sendErr != nil {
	// 	log.Println(sendErr)
	// 	return sendErr
	// }

	// c.CreateBatchInsert()
	return nil
}
