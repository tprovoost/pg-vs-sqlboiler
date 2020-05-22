package main

import (
	"orm_compare/modules"

	_ "github.com/lib/pq"
)

func main() {

	/*var count int16
	row := db.QueryRow("SELECT COUNT(*) FROM products")
	err := row.Scan(&count)

	if err != nil {
		fmt.Printf("Error while prompting count %v\n", err)
	}

	fmt.Printf("Count 1 : %d\n", count)*/

	modules.RunSqlBoiler()

}
