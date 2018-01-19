package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)



func main() {
	db, err := sql.Open("mysql", "root:root@/test?multiStatements=true")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Println(err.Error())
	}





	rows, err := db.Query(`
	select username from users limit 0, 10;
	select uuid from users limit 0, 10;`)

	if err != nil {
		log.Print(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		name := ""
		rows.Scan(&name)
		fmt.Println(name)
	}

	if !rows.NextResultSet() {
		log.Fatal("expected more result sets", rows.Err())
	}
	for rows.Next() {
		uuid := ""
		rows.Scan(&uuid)
		fmt.Println(uuid)
	}

	rows, err = db.Query(`call id_users(?)`, 10)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		name := ""
		rows.Scan(&name)
		fmt.Println(name)
	}

}
