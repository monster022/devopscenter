package helper

import (
	"database/sql"
	"devopscenter/configuration"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var SqlContext *sql.DB

func MysqlConnect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		configuration.Configs.MysqlUsername,
		configuration.Configs.MysqlPassword,
		configuration.Configs.MysqlAddress,
		configuration.Configs.MysqlPort,
		configuration.Configs.MysqlDatabases,
		configuration.Configs.MysqlCharset)
	MysqlEngine, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Connect Databases Failed: ", err)
		panic(err)
	}
	MysqlEngine.SetMaxOpenConns(100)
	MysqlEngine.SetMaxIdleConns(50)
	SqlContext = MysqlEngine
}
