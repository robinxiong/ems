package main

import (
	"ems/l10n"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println(l10n.Global)

}
