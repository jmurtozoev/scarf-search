package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func ConnectDb(dbPath string) (err error) {
	if dbPath == "" {
		dbPath = "./cloth-shop.db"
	}
	db, err = sql.Open("sqlite3", dbPath)
	return
}

func CreateTable(db *sql.DB) {
	createStudentTableSQL := `create table if not exists scarves(
        id integer primary key,
        material VARCHAR(255),
        manufacturer VARCHAR(255), 
        price INTEGER,
        colour VARCHAR(100),
        width INTEGER,
        length INTEGER, 
        size DECIMAL(10,2));` // SQL Statement for Create Table

	log.Println("Create scarves table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatal(err.Error())
	}

}

func GetAll() ([]Scarf, error) {
	scarves := make([]Scarf, 0)
	rows, err := db.Query("SELECT id, material, manufacturer, price, colour, width, length FROM scarves")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var scarf Scarf
		err = rows.Scan(
			&scarf.Id,
			&scarf.Material,
			&scarf.Manufacturer,
			&scarf.Price,
			&scarf.Colour,
			&scarf.Width,
			&scarf.Length)

		if err != nil {
			return nil, err
		}
		scarves = append(scarves, scarf)
	}

	return scarves, err
}