package main

import (
	"orm_compare/modules"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	boil.DebugMode = false

	modules.RunSQLBoiler()
	//modules.RunPG()

}
