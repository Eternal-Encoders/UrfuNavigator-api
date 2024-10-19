package utils

import (
	"fmt"
	"strings"
)

func CreateTable(tableName string, columns map[string]string, order string) string {
	columnsStr := ""
	for column, columnType := range columns {
		columnsStr += fmt.Sprintf("%s %s, ", column, columnType)
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) ORDER BY %s", tableName, columnsStr, order)
	return query
}

func CreateEnum(values []string) string {
	query := fmt.Sprintf("ENUM ('%s')", strings.Join(values, "', '"))
	return query
}

func CreateArray(values []string) string {
	query := fmt.Sprintf("Array('%s')", strings.Join(values, "', '"))
	return query
}
