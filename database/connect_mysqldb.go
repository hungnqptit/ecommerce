package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (db *sql.DB) {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:admin123456@tcp(127.0.0.1:3306)/golangdb")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("success connect mysql db")
	}

	// if there is an error inserting, handle it
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer db.Close()
	return db
}
