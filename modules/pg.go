package modules

import (
	"database/sql"
	"fmt"
	models "orm_compare/models"
	"time"

	"github.com/go-pg/pg/v10"
)

// We have to declare all structures first

// ProductPG represents a product in database
type ProductPG struct {
	// if struct has different name than the table, it is necessary
	// to define the correct name.
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

	if err := pgReadOne(db); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while reading value %v", err)
	}

	if err := pgReadAll(db); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while reading all %v", err)
	}

	if err := pgComplexQuery(db); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while doing complex query %v", err)
	}

	return nil
}

func pgReadOne(db *pg.DB) error {
	fmt.Println("Read one element")

	product := &ProductPG{ID: 2}
	err := db.Select(product)
	if err != nil {
		return fmt.Errorf("error while selecting product %v", err)
	}

	fmt.Printf("Product with id 2 is: %v\n", product)

	return nil
}

func pgReadAll(db *pg.DB) error {
	fmt.Println("Read data")
	var products []ProductPG

	cpt, err := db.Model(&products).Count()
	if err != nil {
		return fmt.Errorf("error while prompting count %v", err)
	}
	fmt.Printf("Products in database: %d\n", cpt)

	return nil
}

func pgComplexQuery(db *pg.DB) error {
	var qpps []models.QuantityPerProduct
	_, err := db.Query(&qpps, `SELECT product_id, SUM(quantity) quantity FROM purchase_items GROUP BY product_id ORDER BY product_id`)

	if err != nil {
		return fmt.Errorf("error while fetching results %v", err)
	}

	for i := 0; i < len(qpps); i++ {
		fmt.Printf("%v\n", qpps[i])
	}

	return nil
}
