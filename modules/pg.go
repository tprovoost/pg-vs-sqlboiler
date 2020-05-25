package modules

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

// We have to declare all structures first

// ProductPG represents a product in database
type ProductPG struct {
	tableName struct{} `pg:"select:products"`
	ID        int64
	Title     string
	Price     float64
	CreatedAt time.Time
	DeletedAt sql.NullTime
	Tags      []string `pg:",array"`
}

// RunPG executes all PG commands
func RunPG() error {
	fmt.Println("Starts PG")
	db := pg.Connect(&pg.Options{
		User:      "thomasprovoost",
		Database:  "pgguide",
		TLSConfig: nil,
	})
	defer db.Close()

	err := pgRead(db)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while reading value %v", err)
	}

	return nil
}

func pgRead(db *pg.DB) error {
	fmt.Println("Read data")
	var products []ProductPG

	cpt, err := db.Model(&products).Count()
	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}
	fmt.Printf("Products in database: %d\n", cpt)

	return nil
}
