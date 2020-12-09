package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func readFile(name string) [][]string {
	var records [][]string
	file, err := os.Open(name)

	if err != nil {
		log.Fatalln("Could not open the csv file!", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.Comment = '#'

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(records)
		records = append(records, row)
	}

	return records
}

func prepareDatabase() *sql.DB {
	database, _ := sql.Open("sqlite3", "./yawoen.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS companies (id INTEGER PRIMARY KEY, name TEXT, zip TEXT)")
	statement.Exec()

	return database
}

func insertValue(database *sql.DB, companyName string, zipCode string) {
	statement, _ := database.Prepare("INSERT INTO companies (name, zip) VALUES (?, ?)")
	statement.Exec(companyName, zipCode)
}

func insertValues(database *sql.DB, values [][]string) {
	for _, s := range values {
		log.Println(s)

		insertValue(database, strings.ToUpper(s[0]), s[1])
	}
}

func main() {
	database := prepareDatabase()

	// var records [][]string
	// records = readFile("./files/q1_catalog.csv")
	// insertValues(database, records[1:])

	rows, _ := database.Query("SELECT * FROM companies")

	var id int
	var name string
	var zip string
	for rows.Next() {
		rows.Scan(&id, &name, &zip)
		fmt.Println(strconv.Itoa(id) + ": " + name + " " + zip)
	}

	defer database.Close()
}
