// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for sql usage in go
//
// See github.com/framps/golang_tutorial for latest code
//
// See github.com/framps/golang_tutorial for latest code

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type person struct {
	name       string
	department string
	created    string
}

// convenience helper to be able to pass one parm for exec
func (p *person) forSQLExec() (s []interface{}) {
	return []interface{}{p.name, p.department, p.created}
}

type sqlPerson struct {
	uid string
	person
}

func main() {

	peoples := []person{
		{"Walter", "Steinmeier", "1956-01-05"},
		{"Albert", "Einstein", "1879-03-14"},
		{"Herbert", "Gönemeier", "1956-04-12"},
	}

	// open database
	fmt.Println("... Open database")
	db, err := sql.Open("sqlite3", "./sample1.db")
	checkErr(err)
	defer db.Close() // close db at end of prog

	fmt.Println("... Delete all entries from database")
	stmt, err := db.Prepare("DELETE FROM userinfo")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	for _, person := range peoples {
		fmt.Printf("... Inserting %s\n", person.name)
		_, err = stmt.Exec(person.forSQLExec()...)
		checkErr(err)
	}

	// variables required for query
	var (
		uid                        int
		username, department, born string
	)

	// just to remember the ids of the db entries
	var idList map[string]int

	var rows *sql.Rows // query results

	query := func() { // annonymous function to query the db

		idList = make(map[string]int)

		// query
		rows, err = db.Query("SELECT * FROM userinfo")
		defer rows.Close() // close query at function end
		checkErr(err)

		fmt.Println("... Execute query")
		for rows.Next() {
			err = rows.Scan(&uid, &username, &department, &born)
			checkErr(err)
			idList[username] = uid
			fmt.Printf("... Retrieved: %s %s, born %s\n", username, department, born)
		}
	}

	query() // show current db contents

	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)
	fmt.Printf("... Updating %s to %s\n", "Walter", "Frank-Walter")
	_, err = stmt.Exec("Frank-Walter", idList["Walter"])
	checkErr(err)

	fmt.Println("... Deleting Herbert")
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)
	_, err = stmt.Exec(idList["Herbert"])
	checkErr(err)

	query() // show current db contents

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
