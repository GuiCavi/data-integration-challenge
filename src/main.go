package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type company struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Website string `json:"website,omitempty"`
}

/**
This function reads a csv file base on the name given
*/
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

/**
This function opens the database and create a new table called companies
with the fields id, name, zip
*/
func prepareDatabase() *sql.DB {
	database, _ := sql.Open("sqlite3", "./yawoen.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS companies (id INTEGER PRIMARY KEY, name TEXT, zip TEXT)")
	statement.Exec()

	return database
}

/**
This function adds a new company to companies table.
It takes in consideration only the `name` and `zip` values so it
can satisfact the first part of the test.
*/
func insertValue(database *sql.DB, companyName string, zipCode string) {
	statement, _ := database.Prepare("INSERT INTO companies (name, zip) VALUES (?, ?)")
	statement.Exec(companyName, zipCode)
}

/**
This function receives a slice of slices of strings and iterate over them
to insert the values into database
*/
func insertValues(database *sql.DB, values [][]string) {
	for _, s := range values {
		insertValue(database, strings.ToUpper(s[0]), s[1])
	}
}

/**
This function adds a column to the companies table.
It could generic but for this example it was chosen to be specific for `website` column
*/
func updateTable(database *sql.DB) {
	statement, _ := database.Prepare("ALTER TABLE companies ADD website TEXT")
	statement.Exec()
}

/**
This function updates website from companies based on the zip code.
It could be generic but for this example it was chosen to be specific for `website` and `zip` columns
*/
func updateValue(database *sql.DB, selectValue string, newValue string) {
	statement, _ := database.Prepare("UPDATE companies SET website=? WHERE zip=?")
	statement.Exec(newValue, selectValue)
}

/**
This function finds a company using name and zip values.
*/
func findCompany(database *sql.DB, name string, zip string) company {
	rows, _ := database.Query("SELECT * FROM companies WHERE zip=? AND name LIKE '%' || $1 || '%' LIMIT 1", zip, name)

	var company = company{}

	for rows.Next() {
		rows.Scan(&company.ID, &company.Name, &company.Zip, &company.Website)
	}

	return company
}

/**
Here is the start of the API.
It depends on a sql.DB variable.

This function declares the URIs that will be existing in the API and
serve it on port 3333.
*/
func startAPI(database *sql.DB) {
	router := mux.NewRouter()

	router.HandleFunc("/company_info/{name}/{zip}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		company := findCompany(database, params["name"], params["zip"])

		fmt.Println(company)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(company)
	}).Methods("GET")

	http.ListenAndServe(":3333", router)
}

func main() {
	database := prepareDatabase()

	startAPI(database)

	// var records [][]string
	// records = readFile("./files/q1_catalog.csv")
	// insertValues(database, records[1:])

	// rows, _ := database.Query("SELECT * FROM companies")

	// var id int
	// var name string
	// var zip string
	// for rows.Next() {
	// 	rows.Scan(&id, &name, &zip)
	// 	fmt.Println(strconv.Itoa(id) + ": " + name + " " + zip)
	// }

	// updateTable(database)

	// var newRecords = readFile("./files/q2_clientData.csv")
	// for _, s := range newRecords {
	// 	updateValue(database, s[1], s[2])
	// }

	// defer database.Close()
}
