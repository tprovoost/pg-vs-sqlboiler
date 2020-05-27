package main

import (
	_ "github.com/lib/pq"
	"github.com/tprovoost/pg-vs-sqlboiler/modules"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	boil.DebugMode = false

	modules.RunSQLBoiler()
	//modules.RunPG()

}
