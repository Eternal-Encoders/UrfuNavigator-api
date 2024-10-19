package logs

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"
	"urfunavigator/index/utils"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Clickhouse struct {
	conn clickhouse.Conn
	db   string
}

func NewClickhouse(
	uri string,
	db string,
	user string,
	password string,
) *Clickhouse {
	ctx := context.Background()
	dialCount := 0

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{uri},
		Auth: clickhouse.Auth{
			Username: user,
			Password: password,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:      time.Duration(10) * time.Second,
		MaxOpenConns:     5,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		BlockBufferSize:  10,
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		log.Fatal(err)
	}

	return &Clickhouse{
		conn: conn,
		db:   db,
	}
}

func (c *Clickhouse) InitDb() error {
	ctx := context.Background()

	requests := []string{
		"MainHandler",
		"FloorHandler",
		"InstituteHandler",
		"InstitutesHandler",
		"PointsHandler",
		"PointIdHandler",
		"PathHandler",
		"ObjectHandler",
		"SearchHandler",
	}
	apiCallsTypes := map[string]string{
		"ip":        "IPv4",
		"request":   utils.CreateEnum(requests),
		"args":      "Map(String, String)",
		"timestamp": "DateTime",
	}

	if err := c.conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", c.db)); err != nil {
		return err
	}
	if err := c.conn.Exec(ctx, utils.CreateTable(fmt.Sprintf("%s.api_calls", c.db), apiCallsTypes, "timestamp")); err != nil {
		return err
	}

	return nil
}

func (c *Clickhouse) WriteLog(ip string, request string, args map[string]string) error {
	ctx := context.Background()
	if err := c.conn.AsyncInsert(
		ctx,
		fmt.Sprintf("INSERT INTO %s.%s (ip, request, args, timestamp) VALUES (toIPv4(?), ?, ?, now())", c.db, "api_calls"),
		true,
		ip,
		request,
		args,
	); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
