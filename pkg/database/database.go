package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb() (dbc *DbClient, err error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}

	dbc = new(DbClient)
	dbc.conn = db

	query := "SELECT * FROM `sqlite_master` WHERE `type`='table' AND `name`='opengfw_record'"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		// 查询结果为空
		query = "CREATE TABLE `opengfw_record` (`id` bigint not null constraint opengfw_record_pk primary key,`rule_name` varchar(255) not null,`src` varchar(255) not null,`dst` varchar(255) not null,`time` varchar(255) not null, action varchar(255) not null);"
		_, err := db.Exec(query)
		if err != nil {
			return nil, err
		}
	}
	return dbc, nil
}

func (dbc *DbClient) QueryRecord(query string) (result []OpenGFWRecord, err error) {
	db := dbc.conn
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rs []OpenGFWRecord

	for rows.Next() {
		var data OpenGFWRecord
		err := rows.Scan(&data)
		if err != nil {
			return nil, err
		}
		rs = append(rs, data)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rs, nil
}

func (dbc *DbClient) InsertRecord(query string) (err error) {
	db := dbc.conn
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DbClient) UpdateRecord(query string) (err error) {
	db := dbc.conn
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DbClient) RemoveRecord(query string) (err error) {
	db := dbc.conn
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

type DbClient struct {
	conn *sql.DB
}

type OpenGFWRecord struct {
	Id       int
	RuleName string
	Src      string
	Dst      string
	Time     string
	Action   string
}
