package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

var MysqlDb *sql.DB
var MysqlDbErr error

// 初始化链接
func init() {
	var cfg cfg
	if _, err := toml.DecodeFile("./config/config.toml", &cfg); err != nil {
		fmt.Println(err)
	}

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", cfg.Db.Username, cfg.Db.Password, cfg.Db.Endpoint, cfg.Db.Port, cfg.Db.Database, cfg.Db.Charset)

	// 打开连接失败
	MysqlDb, MysqlDbErr = sql.Open("mysql", dbDSN)
	//defer MysqlDb.Close();
	if MysqlDbErr != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlDbErr.Error())
	}

	// 最大连接数
	MysqlDb.SetMaxOpenConns(100)
	// 闲置连接数
	MysqlDb.SetMaxIdleConns(20)
	// 最大连接周期
	MysqlDb.SetConnMaxLifetime(100 * time.Second)

	if MysqlDbErr = MysqlDb.Ping(); nil != MysqlDbErr {
		panic("数据库链接失败: " + MysqlDbErr.Error())
	}

}

type Webhook struct {
	Name     string `db:"name"`
	Author   string `db:"author"`
	CommitID string `db:"commitID"`
}

// func GetRule(Q string) {
// 	//获取告警规则信息
// 	wk := new(Webhook)
// 	row := MysqlDb.QueryRow("select  name,author,commitid from wkorder where Name=? ORDER BY id desc LIMIT 1;", Q)
// 	if err := row.Scan(&wk.Name, &wk.Author, &wk.CommitID); err != nil {
// 		fmt.Printf("scan failed, err:%v", err)
// 	}

// 	named = wk.Name
// 	authors = wk.Author
// 	commits = wk.CommitID

// 	return named, authors, commits, err
// }
