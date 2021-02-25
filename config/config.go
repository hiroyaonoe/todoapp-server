/*
Package config は環境変数から設定を取得するパッケージ
*/
package config

import (
	"fmt"
	"os"
)

// Port はRouting用のポートを返す
func Port() string {
	return os.Getenv("ROUTING_PORT")
}

// DSN はデータベースに接続するためのData Source Nameを返す
func DSN(env string) string {
	dsn := map[string]string{}
	dsn["Production"] = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST_DEV"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	dsn["Test"] = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST_TEST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	return dsn[env]
}
