package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bpeters-cmu/dns-threat-analyser/graph/model"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	os.Remove("./threat_analyser.db")

	var err error
	db, err = sql.Open("sqlite3", "./threat_analyser.db")
	if err != nil {
		log.Fatal(err)
	}
	createCmd := `
	create table ip (ip_address TEXT PRIMARY KEY,
					 uuid TEXT,
					 created_at DATETIME,
					 updated_at DATETIME,
					 response_code TEXT);
	`
	_, err = db.Exec(createCmd)
	if err != nil {
		log.Fatal("Error creating DB table", err)
		return
	}
}

type Database interface {
	SaveIp(ip *model.IP) error
	GetIp(ipAddr string) (*model.IP, error)
}

type SqliteDB struct {
}

func (sqlDb *SqliteDB) SaveIp(ip *model.IP) error {
	upsert, err := db.Prepare("INSERT OR REPLACE INTO ip (ip_address, uuid, created_at, updated_at, response_code) VALUES (?, ?, ?, ?, ?)")
	defer upsert.Close()
	if err != nil {
		return errors.New(fmt.Sprint("ERROR preparing db insert statement:", err.Error()))
	}
	_, err = upsert.Exec(ip.IPAddress, ip.UUID, ip.CreatedAt, ip.UpdatedAt, ip.ResponseCode)

	if err != nil {
		return errors.New(fmt.Sprint("ERROR executing DB insert:", err.Error()))
	}
	return nil
}

func (sqlDb *SqliteDB) GetIp(ipAddr string) (*model.IP, error) {
	row := db.QueryRow("SELECT * FROM ip WHERE ip_address = ?", ipAddr)
	ip := model.IP{}

	err := row.Scan(&ip.IPAddress, &ip.UUID, &ip.CreatedAt, &ip.UpdatedAt, &ip.ResponseCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}

	}
	return &ip, nil

}
