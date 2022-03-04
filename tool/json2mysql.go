package tool

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const TABLE_CHINA_GEO_SQL = "create table if not exists china_geo (id int auto_increment comment 'id' primary key,code int default 0  not null comment 'code',parent_code int default 0  not null comment 'parent code',name varchar(100) default '' not null comment 'name',level tinyint default 0  not null comment 'level',center varchar(50) default ''  not null comment 'center point', `polyline` mediumtext NOT NULL, KEY `china_geo_code_index` (`code`), KEY `china_geo_parent_code_index` (`parent_code`))comment 'china_geo'"

var db *sql.DB

type mysqlConf struct {
	host     string
	port     string
	userName string
	pass     string
	dbName   string
}

func NewMysqlConf(host string, port string, userName string, pass string, dbName string) *mysqlConf {
	return &mysqlConf{
		host:     host,
		port:     port,
		userName: userName,
		pass:     pass,
		dbName:   dbName,
	}
}

func (sc *mysqlConf) CreateTable() error {
	err := sc.conn()
	if err != nil {
		return err
	}

	err = createTable(TABLE_CHINA_GEO_SQL)
	if err != nil {
		return err
	}
	fmt.Println("create table region finished")

	return nil
}

func (sc *mysqlConf) conn() (err error) {
	dataSourceName := sc.userName + ":" + sc.pass + "@tcp(" + sc.host + ":" + sc.port + ")/" + sc.dbName + "?autocommit=true"
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	return nil
}

func createTable(sqlStr string) error {
	_, err := db.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}

func InsertChinaGeo(code uint, parentCode uint, name string, level uint, center string, polyline string) error {
	sql := fmt.Sprintf("insert into china_geo (code, parent_code, name, level, center, polyline) values(%d, %d, '%s', %d, '%s', '%s')", code, parentCode, name, level, center, polyline)
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (sc *mysqlConf) QueryData() (*sql.Rows, error) {
	err := sc.conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := time.Now()
	rows, err := db.Query("select * from china_geo")
	endTime := time.Now()
	useTime := endTime.Sub(startTime).Milliseconds()
	fmt.Printf("query use %d ms\n", useTime)
	return rows, err
}
