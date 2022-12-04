package common

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"asong.cloud/go-algorithm/leaf/config"
)

func NewMysqlClient(conf *config.Server) *sql.DB {
	connInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.Mysql.Username, conf.Mysql.Password, conf.Mysql.Host, conf.Mysql.Db)
	var err error
	db, err := sql.Open("mysql", connInfo)
	if err != nil {
		fmt.Printf("init mysql err %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("ping mysql err: %v", err)
	}
	db.SetMaxIdleConns(conf.Mysql.Conn.MaxIdle)
	db.SetMaxOpenConns(conf.Mysql.Conn.Maxopen)
	db.SetConnMaxLifetime(5 * time.Minute)
	fmt.Println("init mysql successc")
	return db
}
